package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"rvdl/internal/server"
	"rvdl/pkg/rvdl"
	"rvdl/pkg/util"
	"syscall"
	"time"
)

var Address = util.EnvString("RVDL_ADDRESS", "")

func main() {
	err := rvdl.Init()
	if err != nil {
		fmt.Println(err)
		return
	}

	srv := &http.Server{
		Addr: Address,
		Handler: http.HandlerFunc(server.ServeHTTP),
	}

	go func() {
		if err := srv.ListenAndServeTLS("certs/rvdl.crt", "certs/rvdl.key"); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	fmt.Println("Server Started")

	done := make(chan os.Signal)
	signal.Notify(done, os.Interrupt, syscall.SIGTERM)
	<-done

	fmt.Println("Server Stopped")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer func() {
		// extra handling here
		cancel()
	}()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server Shutdown Failed:%+v", err)
	}
	log.Print("Server Exited Properly")

	rvdl.Close()
}
