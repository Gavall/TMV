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
	router.GET("/project/:userId", handlerMongo.GetProject)
	router.GET("/projects/", handlerMongo.GetAllProjects)
	router.DELETE("/project/:id", handlerMongo.DeleteProject)
	router.DELETE("/user/:userId/projects", handlerMongo.DeleteProjects)
	router.PATCH("/project/:projectId", handlerMongo.UpdateProject)

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

// 6666df7fe46ca0eec261ce5d

// "667046cad72cfdacf52c6bf3",
// "667046d2d72cfdacf52c6bf4",
//
// "name": "New Project",
// "description": "Description of the new project",
// "priority": 5,
// "author": "Author Name",
// "responsible": "Responsible Person",
// "performers": "Performer1, Performer2",
// "deadline": "2024-12-31T23:59:59Z",
// "guests": "Guest1, Guest2",
// "status": "open"
//
//
//    "projectIDs": [
// 	"60b8d6e4f1a2c3b9d567e98a",
// 	"60b8d6e4f1a2c3b9d567e98b",
// 	"60b8d6e4f1a2c3b9d567e98c"
// ]
//
//
//
//
//
//
//
