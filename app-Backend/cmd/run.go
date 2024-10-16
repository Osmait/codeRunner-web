package cmd

import (
	"fmt"
	"log"

	coderunner "github.com/Osmait/CodeRunner-web/internal/app/CodeRunner"
	"github.com/Osmait/CodeRunner-web/internal/modules/dispacher"
	programinglanguages "github.com/Osmait/CodeRunner-web/internal/modules/programingLanguages"
	"github.com/Osmait/CodeRunner-web/internal/modules/runner"
	"github.com/spf13/cobra"
)

var (
	language string
	file     string
)

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run the application as a terminal application.",
	Run: func(cmd *cobra.Command, args []string) {
		if language == "" || file == "" {
			log.Fatal("Both --language and --file must be provided")
		}

		outputs := make(chan []byte)
		output := dispacher.NewNotifier(outputs)
		runner := runner.NewRunner() // Aquí debes inicializar el runner según tu implementación.
		availableLang := programinglanguages.NewAvailablePrograminLanguages()
		if !availableLang.IsAvaliable(language) {
			log.Fatal("languages not suported")
		}
		coderun := coderunner.NewCodeRunner(runner, output, availableLang)
		go func() {
			for v := range output.Consumer() {
				fmt.Println((string(v)))
			}
		}()

		fmt.Println("Running the code...")
		codeRequest := coderunner.CodeRequest{Code: file, Lang: language}
		coderun.RunCode(codeRequest)
		fmt.Println("Code executed successfully.")
	},
}

func init() {
	runCmd.Flags().StringVarP(&language, "language", "l", "", "Programming language to use (e.g., python, javascript, etc.)")
	runCmd.Flags().StringVarP(&file, "file", "f", "", "File path of the code to execute")
	rootCmd.AddCommand(runCmd)
}
