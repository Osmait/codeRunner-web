package runner

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os/exec"
	"path/filepath"

	"github.com/Osmait/CodeRunner-web/internal/modules/dispacher"
)

type RunnerInterface interface {
	Execute(code string, lang string, output *dispacher.Dispacher) error
}

type Runner struct{}

func NewRunner() *Runner {
	return &Runner{}
}

func (r *Runner) Execute(code string, lang string, output *dispacher.Dispacher) error {
	// Ejemplo de comando que se ejecutará (en este caso, ls)
	cmd := exec.Command("docker", "run", "--rm", "-i", "-v", fmt.Sprintf("%s:/app", filepath.Dir("temp.js")), "-w", "/app", "node:latest", "node", filepath.Base("temp.js"))

	// Capturar salida estándar del comando
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}

	// Capturar salida de error estándar
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return err
	}

	// Iniciar la ejecución del comando
	if err := cmd.Start(); err != nil {
		return err
	}

	go func() {
		sendLogs(stdout, output) // Usar la función sendLogs para enviar los logs
	}()

	go func() {
		sendLogs(stderr, output) // Captura y envía errores al canal
	}()

	// Esperar a que el comando termine su ejecución
	if err := cmd.Wait(); err != nil {
		return err
	}

	log.Println("Comando finalizado y eliminado del canal.")
	return nil
}

func sendLogs(pipe io.ReadCloser, output *dispacher.Dispacher) {
	scanner := bufio.NewScanner(pipe)
	for scanner.Scan() {
		message := scanner.Text()
		output.Notifique([]byte(message))
	}
	if err := scanner.Err(); err != nil {
		log.Println("Error leyendo el log:", err)
	}
}
