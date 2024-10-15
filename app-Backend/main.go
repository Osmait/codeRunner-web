package main

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"sync"
	"syscall"
	"time"

	"github.com/Osmait/CodeRunner-web/internal/modules/dispacher"
	programinglanguages "github.com/Osmait/CodeRunner-web/internal/modules/programingLanguages"
	"github.com/Osmait/CodeRunner-web/internal/server"

	"github.com/gorilla/websocket"
)

type CodeRequest struct {
	Lang   string `json:"lang"`
	Code   string `json:"code"`
	Action string `json:"action,omitempty"`
}

type Client struct {
	conn     *websocket.Conn
	cmdChan  chan *exec.Cmd
	stopChan chan struct{}
	stopOnce sync.Once
}

func enableCors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func executeCode(lang string, code string, client *Client) {
	availablePrograminLanguages := programinglanguages.NewAvailablePrograminLanguages()
	languages, err := availablePrograminLanguages.SearchLanguage(lang)
	if err != nil {
		client.conn.WriteMessage(websocket.TextMessage, []byte("Error: lenguaje no soportado"))
	}

	filename := fmt.Sprintf("temp.%s", languages.GetExtension())
	runner := fmt.Sprintf("%s:latest", languages.GetRunner())
	fmt.Println(filename)
	os.WriteFile(filename, []byte(code), 0644)
	runDockerCommand(runner, languages.GetRunner(), filename, client)
}

func runDockerCommand(image, command, filename string, client *Client) {
	cmd := exec.Command("docker", "run", "--rm", "-i", "-v", fmt.Sprintf("%s:/app", filepath.Dir(filename)), "-w", "/app", image, command, filepath.Base(filename))

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		client.conn.WriteMessage(websocket.TextMessage, []byte("Error al obtener stdout: "+err.Error()))
		return
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		client.conn.WriteMessage(websocket.TextMessage, []byte("Error al obtener stderr: "+err.Error()))
		return
	}

	// Iniciar el comando en un goroutine
	go func() {
		// Enviar los logs al WebSocket
		go sendLogs(stdout, client)
		go sendLogs(stderr, client)

		if err := cmd.Start(); err != nil {
			client.conn.WriteMessage(websocket.TextMessage, []byte("Error al iniciar el comando: "+err.Error()))
			return
		}

		log.Println("Proceso Docker iniciado con PID:", cmd.Process.Pid)

		// Aquí se envía el comando al canal
		client.cmdChan <- cmd
		log.Println("Comando enviado al canal para posible detención...")

		if err := cmd.Wait(); err != nil {
			client.conn.WriteMessage(websocket.TextMessage, []byte("Error al ejecutar el comando: "+err.Error()))
		}

		// Eliminar el comando cuando termine
		client.cmdChan <- nil
		log.Println("Comando finalizado y eliminado del canal.")
	}()
}

func sendLogs(pipe io.ReadCloser, client *Client) {
	scanner := bufio.NewScanner(pipe)
	for scanner.Scan() {
		message := scanner.Text()
		client.conn.WriteMessage(websocket.TextMessage, []byte(message))
	}
	if err := scanner.Err(); err != nil {
		log.Println("Error leyendo el log:", err)
	}
}

func handleConnection(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Error al establecer conexión:", err)
		return
	}
	defer conn.Close()

	client := &Client{
		conn:     conn,
		cmdChan:  make(chan *exec.Cmd),
		stopChan: make(chan struct{}),
	}

	go handleCommands(client)

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Println("Error al leer mensaje:", err)
			break
		}

		var request CodeRequest
		if err := json.Unmarshal(msg, &request); err != nil {
			conn.WriteMessage(websocket.TextMessage, []byte("Error de deserialización: "+err.Error()))
			continue
		}

		// Verificar si el mensaje es una acción de "stop"
		if request.Action == "stop" {
			log.Println("Recibida solicitud de parada")
			client.stopChan <- struct{}{} // Enviar señal de stop, no cerrar el canal
			continue
		}

		// Ejecutar el código si no es una acción de parada
		executeCode(request.Lang, request.Code, client)
	}
}

func handleCommands(client *Client) {
	for {
		select {
		case cmd := <-client.cmdChan:
			if cmd == nil {
				log.Println("No hay ningún comando en ejecución en el canal")
				continue
			}
			log.Println("Comando recibido en cmdChan, esperando señal de parada...")

			// Iniciar un nuevo goroutine para ejecutar el comando
			go func() {
				defer func() {
					// Limpieza: Cerramos el canal o hacemos otras tareas necesarias
				}()

				// Esperar a que termine el comando
				if err := cmd.Wait(); err != nil {
					log.Println("Error al esperar el comando:", err)
					return
				}
			}()

			// Esperar la señal de parada en otra gorutina
			go func() {
				select {
				case <-client.stopChan: // Esperar a que llegue la señal de stop
					log.Println("Intentando detener el proceso...")
					if cmd.Process != nil {
						// Intentar primero con SIGINT
						if err := cmd.Process.Signal(os.Interrupt); err != nil {
							log.Println("Error al enviar señal de interrupción:", err)
						} else {
							log.Println("Señal de interrupción enviada.")
						}

						// Esperar un tiempo para permitir que el proceso maneje la señal
						time.Sleep(2 * time.Second)

						// Verificar si sigue ejecutándose
						if err := cmd.Process.Signal(syscall.SIGTERM); err != nil {
							log.Println("Error al enviar señal de SIGTERM:", err)
						} else {
							log.Println("Señal SIGTERM enviada.")
						}

						// Si no responde, enviar kill
						if err := cmd.Process.Kill(); err != nil {
							log.Println("Error al enviar señal de kill:", err)
						} else {
							log.Println("Señal de kill enviada.")
						}
					}
				}
			}()
		}
	}
}

func main() {
	outputs := make(chan []byte)
	dispacher := dispacher.NewNotifier(outputs)
	if dispacher == nil || dispacher.NotificationChan == nil {
		log.Fatal("Notifier no se inicializó correctamente")
	}

	ctx, server := server.New(context.Background(), "127.0.0.1", 8080, dispacher)
	server.Run(ctx)

	//
	// http.HandleFunc("/ws", handleConnection)
	//
	// fmt.Println("Servidor WebSocket escuchando en ws://localhost:8080/ws")
	// if err := http.ListenAndServe(":8080", enableCors(http.DefaultServeMux)); err != nil {
	// 	log.Fatal("Error al iniciar servidor:", err)
	// }
}
