package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"

	coderunner "github.com/Osmait/CodeRunner-web/internal/app/CodeRunner"
	"github.com/Osmait/CodeRunner-web/internal/modules/dispacher"
	"github.com/Osmait/CodeRunner-web/internal/modules/runner"
)

type Server struct {
	Engine   *http.ServeMux
	Notifer  *dispacher.Notifier
	httpAddr string
}

func New(ctx context.Context, host string, port uint, notifer *dispacher.Notifier) (context.Context, *Server) {
	srv := Server{
		Engine:   http.DefaultServeMux,
		httpAddr: fmt.Sprintf("%s:%d", host, port),
		Notifer:  notifer,
	}
	srv.Routes()
	return serverContext(ctx), &srv
}

func (s *Server) Routes() {
	http.HandleFunc("/test", func(http.ResponseWriter, *http.Request) {
		runner := runner.NewRunner()
		app := coderunner.NewCodeRunner(runner, s.Notifer)
		app.RunCode()
	})
	http.HandleFunc("/ws", handlerWebSocket(s.Notifer))
}

func (s *Server) Run(ctx context.Context) error {
	log.Println("Server running on", s.httpAddr)

	srv := &http.Server{
		Addr:    s.httpAddr,
		Handler: s.Engine,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("server shut down", err)
		}
	}()

	<-ctx.Done()
	ctxShutDown, cancel := context.WithTimeout(context.Background(), 5)
	defer cancel()

	return srv.Shutdown(ctxShutDown)
}

func serverContext(ctx context.Context) context.Context {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	ctx, cancel := context.WithCancel(ctx)
	go func() {
		<-c
		cancel()
	}()

	return ctx
}
