package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os/exec"
	"sync"

	"github.com/Osmait/CodeRunner-web/internal/modules/dispacher"
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

func handlerWebSocket(notifier *dispacher.Notifier) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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
		fmt.Println(notifier)
		go func() {
			for msg := range notifier.Consumer() {
				client.conn.WriteMessage(websocket.TextMessage, msg)
			}
		}()

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
		}
	}
}
