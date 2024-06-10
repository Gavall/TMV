package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"tmv/handlers"
	"tmv/storage"

	"github.com/gin-gonic/gin"
)

func main() {
	mongoStorage, err := storage.NewMongoStorage("mongodb://localhost:27017", "mydatabase", "employees")
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err = mongoStorage.Client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()

	router := gin.Default()

	memoryStorage := storage.NewMemoryStorage()
	handler2 := handlers.NewHandler(mongoStorage)
	handler := handlers.NewHandler(memoryStorage)

	router.POST("/employee", handler.CreateEmployee)
	router.GET("/employee/:id", handler.GetEmployee)
	router.GET("/employee", handler.GetAllEmployees)
	router.PUT("/employee/:id", handler.UpdateEmployee)
	router.DELETE("/employee/:id", handler.DeleteEmployee)

	router.POST("/remployee", handler2.CreateEmployee)
	router.GET("/remployee/:id", handler2.GetEmployee)
	router.GET("/remployee", handler2.GetAllEmployees)
	router.PUT("/remployee/:id", handler2.UpdateEmployee)
	router.DELETE("/remployee/:id", handler2.DeleteEmployee)

	srv := &http.Server{
		Addr:    ":8000",
		Handler: router,
	}

	// Канал для получения сигналов завершения
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// Запуск сервера в горутине
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("ListenAndServe(): %s", err)
		}
	}()
	log.Println("Server is running on port 8000")

	// Блокируемся до получения сигнала завершения
	<-quit
	log.Println("Shutting down server...")

	// Создаем контекст с тайм-аутом для завершения активных запросов
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	// Останавливаем сервер
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %s", err)
	}

	log.Println("Server exiting")
}
