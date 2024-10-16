package cmd

import (
	"fmt"

	programinglanguages "github.com/Osmait/CodeRunner-web/internal/modules/programingLanguages"
	"github.com/spf13/cobra"
)

var listProgramingLanguage = &cobra.Command{
	Use:   "list",
	Short: "List Availabe Programming languages .",
	Run: func(cmd *cobra.Command, args []string) {
		availableLang := programinglanguages.NewAvailablePrograminLanguages()
		for _, lang := range availableLang.GetListOfAvailablesPrograminLanguages() {
			fmt.Println(lang.GetName())
		}
	},
}

func init() {
	rootCmd.AddCommand(listProgramingLanguage)
}
