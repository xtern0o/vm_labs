package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"vm_lab2/handler"
	"vm_lab2/internal/catalog"
	"vm_lab2/internal/solver"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

func main() {
	fmt.Println(solver.SolveSystem2(catalog.GetSystem(1), catalog.GetPhiSystem(1), solver.NewSystemSimpleIterSolver(), 1, 4, 0.001, 200))

	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer) // recover от паник
	r.Use(middleware.RequestID)
	r.Use(middleware.Timeout(60 * time.Second)) // таймаут на запрос

	// CORS
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost:5173", "http://localhost:5174"},
		AllowedMethods:   []string{"GET", "POST"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	r.Route("/api", func(r chi.Router) {
		r.Use(middleware.AllowContentType("application/json"))

		r.Post("/system/{solver_id}", handler.SystemHandler)
		r.Post("/equation/{solver_id}", handler.EquationHandler)
	})

	srv := &http.Server{
		Addr:         ":8080",
		Handler:      r,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		log.Println("сервер запущен!!")
		if err := srv.ListenAndServe(); err != nil {
			log.Fatalf("это конец... %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	log.Println("выключаемся...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("ВЫКЛЮЧАЕМСЯ...")
	}

	log.Println("выкличились")
}
