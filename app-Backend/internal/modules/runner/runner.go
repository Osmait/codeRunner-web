package runner

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os/exec"
	"path/filepath"

	programinglanguages "github.com/Osmait/CodeRunner-web/internal/modules/programingLanguages"
)

type RunnerInterface interface {
	Execute(code string, lang programinglanguages.ProgramingLanguages)
}

type Runner struct{}

func NewRunner() *Runner {
	return &Runner{}
}

func (r *Runner) Execute(code string, lang string, output chan []byte) error {
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
	close(output) // Cerrar el canal al finalizar
	return nil
}

func sendLogs(pipe io.ReadCloser, output chan []byte) {
	scanner := bufio.NewScanner(pipe)
	for scanner.Scan() {
		message := scanner.Text()
		output <- []byte(message) // Enviar la salida al canal
	}
	if err := scanner.Err(); err != nil {
		log.Println("Error leyendo el log:", err)
	}
}
