package storage

import (
	"tmv/project"
	"tmv/user"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Storage interface {
	GetAllUsers() map[primitive.ObjectID]user.User
	GetUser(userId primitive.ObjectID) (user.User, error)
	InsertUser(u *user.User) error
	UpdateUser(userId primitive.ObjectID, e *user.User) error
	DeleteUser(userId primitive.ObjectID) error

	GetAllProjects() map[primitive.ObjectID]project.Project
	GetProject(userId, projectId primitive.ObjectID) (*project.Project, error)
	GetProjectByUser(userId primitive.ObjectID) ([]project.Project, error)
	InsertProject(p *project.Project, userId primitive.ObjectID) error
	UpdateProject(projectID primitive.ObjectID, updateFields bson.M) error
	DeleteProject(projectId primitive.ObjectID) error
	DeleteProjects(userID primitive.ObjectID, projectIDs []primitive.ObjectID) error

	GetAllTasks() map[primitive.ObjectID]project.Task
	InsertTask(t *project.Task, projectId primitive.ObjectID) error
	GetTasksByProject(projectId primitive.ObjectID) ([]project.Task, error)
	GetTask(projectId, taskId primitive.ObjectID) (*project.Task, error)
	DeleteTasks(projectId primitive.ObjectID, taskIds []primitive.ObjectID) error
	UpdateTask(projectId, taskId primitive.ObjectID, updateFields bson.M) error
	DeleteTask(projectId, taskId primitive.ObjectID) error
}
