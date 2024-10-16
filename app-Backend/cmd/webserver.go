package cmd

import (
	"context"
	"fmt"
	"log"

	"github.com/Osmait/CodeRunner-web/internal/modules/dispacher"
	"github.com/Osmait/CodeRunner-web/internal/server"
	"github.com/spf13/cobra"
)

var webserverCmd = &cobra.Command{
	Use:   "webserver",
	Short: "Start the web server.",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Starting web server on port 8080...")
		outputs := make(chan []byte)
		dispacher := dispacher.NewNotifier(outputs)
		if dispacher == nil || dispacher.NotificationChan == nil {
			log.Fatal("Notifier no se inicializó correctamente")
		}

		ctx, server := server.New(context.Background(), "127.0.0.1", 8080, dispacher)
		server.Run(ctx)
	},
}

func init() {
	rootCmd.AddCommand(webserverCmd)
}
