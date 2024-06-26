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
	mongoStorage, err := storage.NewMongoStorage("mongodb://localhost:27017", "tmv", "users", "projects", "tasks")
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err = mongoStorage.Client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()

	router := gin.Default()

	handlerMongo := handlers.NewHandler(mongoStorage)

	router.POST("/user", handlerMongo.CreateUser)
	router.GET("/user/:userId", handlerMongo.GetUser)
	router.GET("/users", handlerMongo.GetAllUsers)
	router.PUT("/user/:userId", handlerMongo.UpdateUser)
	router.DELETE("/user/:userId", handlerMongo.DeleteUser)

	router.POST("/project/:userId", handlerMongo.CreateProject)
	router.GET("/projects/", handlerMongo.GetAllProjects)

	router.GET("/projects/:userId", handlerMongo.GetProjectsByUser)
	router.GET("/project/:userId/:projectId", handlerMongo.GetProject)

	router.DELETE("/project/:id", handlerMongo.DeleteProject)
	router.DELETE("/user/:userId/projects", handlerMongo.DeleteProjects)
	router.PATCH("/project/:projectId", handlerMongo.UpdateProject)

	router.GET("/tasks/", handlerMongo.GetAlltasks)
	router.GET("/tasks/:projectId", handlerMongo.GetTasksByProject)
	router.GET("/task/:projectId/:taskId", handlerMongo.GetTask)
	router.POST("/task/:projectId", handlerMongo.CreateTask)
	router.DELETE("/task/:projectId/:taskId", handlerMongo.DeleteTask)
	router.DELETE("/tasks/:projectId", handlerMongo.DeleteTasks)
	router.PUT("/projects/:projectId/task/:taskId", handlerMongo.UpdateTask)

	srv := &http.Server{
		Addr:    ":8080",
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
	log.Println("Server is running on port 8080")

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
