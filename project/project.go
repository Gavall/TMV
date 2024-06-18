package project

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Project struct {
	Id           primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID       primitive.ObjectID `bson:"userId" json:"userId"`
	Name         string             `bson:"name" json:"name"`                 // Название проекта
	Descript     string             `bson:"description" json:"description"`   // Описание проекта
	Priority     int                `bson:"priority" json:"priority"`         // Приоритет проекта (от 1 до 10)
	Author       string             `bson:"author" json:"author"`             // Автор
	Responsible  string             `bson:"responsible" json:"responsible"`   // Ответственный
	Performers   string             `bson:"performers" json:"performers"`     // Исполнители
	DateCreation time.Time          `bson:"dateCreation" json:"dateCreation"` // Дата создания
	Deadline     time.Time          `bson:"deadline" json:"deadline"`         // Планируемая дата окончания
	Guests       string             `bson:"guests" json:"guests"`             // Гости
	Tasks        []int              `bson:"tasks" json:"tasks"`               // Задачи
	Status       string             `bson:"status" json:"status"`
}

func NewProject(userId primitive.ObjectID, name, desc string, priority int, author, responsible, performers string, deadline time.Time, guests string, tasks []int, status string) *Project {
	return &Project{
		Id:           primitive.NewObjectID(),
		UserID:       userId,
		Name:         name,
		Descript:     desc,
		Priority:     priority,
		Author:       author,
		Responsible:  responsible,
		Performers:   performers,
		DateCreation: time.Now(),
		Deadline:     time.Time{},
		Guests:       guests,
		Tasks:        tasks,
		Status:       status,
	}
}
