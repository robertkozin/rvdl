package server

import (
	"context"
	"fmt"
	"github.com/robertkozin/rvdl/internal/server"
	"github.com/robertkozin/rvdl/pkg/rvdl"
	"github.com/robertkozin/rvdl/pkg/util"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var Address = util.EnvString("RVDL_ADDRESS", "")

func main() {
	err := rvdl.Init()
	if err != nil {
		log.Println(err)
		return
	}

	srv := &http.Server{
		Addr:    Address,
		Handler: http.HandlerFunc(server.ServeHTTP),
	}

	go func() {
		if err := srv.ListenAndServeTLS("certs/rvdl.crt", "certs/rvdl.key"); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	log.Println("Server Started")

	done := make(chan os.Signal)
	signal.Notify(done, os.Interrupt, syscall.SIGTERM)
	<-done

	fmt.Println("Server Stopped")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("Server Shutdown Failed:%+v", err)
	}

	rvdl.Close()
}
