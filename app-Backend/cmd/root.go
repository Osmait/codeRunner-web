package cmd

import "github.com/spf13/cobra"

var rootCmd = &cobra.Command{
	Use:   "coderunner",
	Short: "CodeRunner is a tool to run code and start a web server.",
}

// Execute inicia la ejecución de los comandos.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	cobra.OnInitialize()
}
