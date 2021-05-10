package main

import (
	"context"
	"errors"
	"log"
	"net"
	"net/http"
	"os"
	"runtime"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"golang.org/x/sync/errgroup"
)

type Server struct {
	server        *http.Server
	timeout       time.Duration
	router        *mux.Router
	lis           net.Listener
	network       string
	address       string
	useGoShutdown bool
}

func (s *Server) init() {
	if runtime.Version() >= "go1.8" {
		s.useGoShutdown = true
	} else {
		s.useGoShutdown = false
	}
}

func (s *Server) Handle(path string, h http.Handler) {
	s.router.Handle(path, h)
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), s.timeout)
	defer cancel()
	s.router.ServeHTTP(w, r.WithContext(ctx))
}

func (s *Server) Start() error {
	lis, err := net.Listen(s.network, s.address)
	if err != nil {
		return err
	}
	s.lis = lis
	if err := s.server.Serve(lis); !errors.Is(err, http.ErrServerClosed) {
		return err
	}
	return nil
}

func (s *Server) Stop() error {
	return s.server.Shutdown(context.Background())
}

func (s *Server) Signal(sig os.Signal) {
	if sig == syscall.SIGTERM {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := s.server.Shutdown(ctx); err != nil {

		}
	}
}

func main() {
	server1 := &Server{}
	server2 := &Server{}
	var g errgroup.Group
	g.Go(func() error {
		return server1.server.ListenAndServe()
	})
	g.Go(func() error {
		return server2.server.ListenAndServe()
	})
	if err := g.Wait(); err != nil {
		log.Fatal(err)
	}
}
