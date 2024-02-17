package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/adamelfsborg-code/food/user/config"
	"github.com/adamelfsborg-code/food/user/data"
	"github.com/adamelfsborg-code/food/user/db"
	"github.com/go-pg/pg/v10"
)

type Server struct {
	router http.Handler
	data   data.DataConn
}

func New(config config.Environments) *Server {
	dataCon := data.DataConn{
		Env: config,
	}

	d := pg.Connect(&pg.Options{
		Addr:     config.DatabaseAddr,
		Database: config.DatabaseName,
		User:     config.DatabaseUser,
		Password: config.DatabasePassword,
	})

	nats, _ := ConnectNats(&Nats{
		host: config.NatsAddr,
	})
	dataCon.Nats = nats

	jetstream, _ := ConnectJetstream(nats)
	dataCon.JS = jetstream

	dataCon.DB = *d

	server := &Server{
		data: dataCon,
	}

	server.loadRoutes()

	return server
}

func (a *Server) Start(ctx context.Context) error {
	server := &http.Server{
		Addr:    a.data.Env.ServerAddr,
		Handler: a.router,
	}

	err := a.data.DB.Ping(ctx)
	if err != nil {
		log.Fatal("failed to connect to repo: %w", err)
	}

	defer func() {
		err := a.data.DB.Close()
		if err != nil {
			log.Fatal("failed to close Repo")
		}
	}()

	defer func() {
		a.data.Nats.Close()
	}()

	err = a.data.Nats.Publish("your.subject", []byte("your message"))
	if err != nil {
		log.Fatal(err)
	}

	a.data.DB.AddQueryHook(&db.QueryLogger{})

	go func() {
		ticker := time.NewTicker(time.Minute)
		defer ticker.Stop()
		for range ticker.C {
			a.data.DB.Ping(ctx)

			if err != nil {
				log.Fatal("Database connection lost:", err)
			}
		}
	}()

	ch := make(chan error, 1)

	go func() {
		err := server.ListenAndServe()
		if err != nil {
			ch <- fmt.Errorf("failed to start server: %w", err)
		}

		close(ch)
	}()

	fmt.Println("Server started")

	select {
	case err := <-ch:
		return err
	case <-ctx.Done():
		timeout, cancel := context.WithTimeout(context.Background(), time.Second*10)
		defer cancel()
		return server.Shutdown(timeout)
	}
}
