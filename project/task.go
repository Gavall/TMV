package project

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Task struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"` // Уникальный идентификатор задачи
	ProjectID    primitive.ObjectID `bson:"projectId" json:"projectId"`        // Идентификатор проекта
	Name         string             `bson:"name" json:"name"`                  // Название задачи
	Description  string             `bson:"description" json:"description"`    // Описание задачи
	Priority     int                `bson:"priority" json:"priority"`          // Приоритет задачи (от 1 до 10)
	Author       string             `bson:"author" json:"author"`              // Автор
	Responsible  string             `bson:"responsible" json:"responsible"`    // Ответственный
	Performers   string             `bson:"performers" json:"performers"`      // Исполнители
	DateCreation time.Time          `bson:"dateCreation" json:"dateCreation"`  // Дата создания
	Deadline     time.Time          `bson:"deadline" json:"deadline"`          // Планируемая дата окончания
	Guests       string             `bson:"guests" json:"guests"`              // Гости
	Status       string             `bson:"status" json:"status"`              // Статус задачи
}

func NewTask(projectID primitive.ObjectID, name, description string, priority int, author, responsible, performers string, deadline time.Time, guests, status string) *Task {
	return &Task{
		ID:           primitive.NewObjectID(),
		ProjectID:    projectID,
		Name:         name,
		Description:  description,
		Priority:     priority,
		Author:       author,
		Responsible:  responsible,
		Performers:   performers,
		DateCreation: time.Now(),
		Deadline:     deadline,
		Guests:       guests,
		Status:       status,
	}
}
