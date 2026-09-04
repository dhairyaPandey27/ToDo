package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/dhairyaPandey27/ToDo/internal/config"
)

func main() {

	//load config
	cfg := config.MustLoad()

	//load database

	//setup router
	router := http.NewServeMux()

	router.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Welcome to students api"))
	})

	//setup server
	server := http.Server{
		Addr:    cfg.Addr,
		Handler: router,
	}

	fmt.Printf("Server started %s: ", cfg.Addr)

	done := make(chan os.Signal,1)
	signal.Notify(done,os.Interrupt,syscall.SIGTERM, syscall.SIGINT)

	go func() {
		err := server.ListenAndServe()
		if err != nil {
			slog.Error("failed to start the server",slog.String("error" ,err.Error()))
		}
	}()

	<- done

	// Gracefull Shutdown

	ctx,cancel := context.WithTimeout(context.Background(),5 * time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err!=nil {
		slog.Error("failed to shutdown the server",slog.String("error", err.Error()))
	}

	slog.Info("successfully shutdown the server")

}
