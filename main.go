package main

import (
	"context"
	"fmt"
	"github.com/hubaxis/mclight/internal/cache"
	"github.com/hubaxis/mclight/internal/config"
	"github.com/hubaxis/mclight/internal/server"
	"github.com/hubaxis/mclight/protocol/mclight"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"net"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	cfg, err := config.New()
	if err != nil {
		log.Fatalln(err)
	}

	// global context for background worker not necessary here
	_, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGTERM)

	cl, err := cache.New(cfg.MemcachedEndpoint)
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.Port))
	if err != nil {
		log.Fatal(err)
	}
	srv := server.New(cl)
	s := grpc.NewServer()
	mclight.RegisterMCLightServiceServer(s, srv)
	reflection.Register(s)

	log.Info("gRPC server started on ", cfg.Port)
	go func() {
		<-sigChan
		cancel()
		s.GracefulStop()
		err := cl.Close()
		if err != nil {
			log.Errorf("can't close memcached connection %v", err)
		}
	}()
	err = s.Serve(listener)
	if err != nil {
		log.Fatal(err)
	}
}
