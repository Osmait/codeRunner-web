package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/gorilla/websocket"
)

type CodeRequest struct {
	Lang string `json:"lang"`
	Code string `json:"code"`
}
type CodeResponse struct {
	Result string `json:"result"`
}

func enableCors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*") // Permite todas las solicitudes
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		// Si es una solicitud OPTIONS (preflight), respondemos con un 200 y terminamos
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Permite todas las solicitudes de origen
	},
}

func executeCode(lang string, code string) (string, error) {
	// Crear un archivo temporal según el lenguaje
	var filename string
	fmt.Println(lang)
	switch lang {
	case "javascript":
		filename = "temp.js"
		os.WriteFile(filename, []byte(code), 0644)
		return runDockerCommand("node:latest", "node", filename)
	case "python":
		filename = "temp.py"
		os.WriteFile(filename, []byte(code), 0644)
		return runDockerCommand("python:latest", "python", filename)
	default:
		return "", fmt.Errorf("language not supported")
	}
}

func runDockerCommand(image, command, filename string) (string, error) {
	// Ejecutar el código en un contenedor Docker
	cmd := exec.Command("docker", "run", "--rm", "-v", fmt.Sprintf("%s:/app", filepath.Dir(filename)), "-w", "/app", image, command, filepath.Base(filename))
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", err
	}
	fmt.Println(string(output))
	return string(output), nil
}

func handleConnection(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Error al establecer conexión:", err)
		return
	}
	defer conn.Close()

	for {
		// Leer el mensaje del cliente
		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Println("Error al leer mensaje:", err)
			break
		}
		// Deserializar el mensaje JSON
		var request CodeRequest
		if err := json.Unmarshal(msg, &request); err != nil {
			conn.WriteMessage(websocket.TextMessage, []byte("Error de deserialización: "+err.Error()))
			continue
		}

		// Ejecutar el código
		result, err := executeCode(request.Lang, request.Code)
		ver := CodeResponse{Result: result}
		res, _ := json.Marshal(ver)
		if err != nil {
			conn.WriteMessage(websocket.TextMessage, []byte("Error: "+err.Error()))
		} else {
			conn.WriteMessage(websocket.TextMessage, res)
		}
	}
}

func main() {
	http.HandleFunc("/ws", handleConnection)

	fmt.Println("Servidor WebSocket escuchando en ws://localhost:8080/ws")
	if err := http.ListenAndServe(":8080", enableCors(http.DefaultServeMux)); err != nil {
		log.Fatal("Error al iniciar servidor:", err)
	}
}
