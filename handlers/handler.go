package handlers

import (
	"fmt"
	"net/http"
	"tmv/project"
	"tmv/storage"
	"tmv/user"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ErrorResponse struct {
	Message string `json:"message"`
}
type Handler struct {
	Storage storage.Storage
}

func NewHandler(st storage.Storage) *Handler {
	return &Handler{Storage: st}
}

func (h *Handler) CreateUser(c *gin.Context) {
	var newUser user.User

	if err := c.BindJSON(&newUser); err != nil {
		fmt.Printf("failer to bind user: %s\n", err.Error())
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Message: err.Error(),
		})
		return
	}

	h.Storage.InsertUser(&newUser)

	c.JSON(http.StatusOK, map[string]interface{}{
		"userId": newUser.Id.Hex(),
	})
}

func (h *Handler) CreateProject(c *gin.Context) {
	var proj project.Project

	if err := c.BindJSON(&proj); err != nil {
		fmt.Printf("failed to bind project: %s\n", err.Error())
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Message: err.Error(),
		})
		return
	}

	userIDStr := c.Param("userId")
	userID, err := primitive.ObjectIDFromHex(userIDStr)
	if err != nil {
		fmt.Printf("failed to convert userId to ObjectID: %s\n", err.Error())
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Message: "invalid userId format",
		})
		return
	}

	err = h.Storage.InsertProject(&proj, userID)
	if err != nil {
		fmt.Printf("failed to insert project: %s\n", err.Error())
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"projectId": proj.Id,
		"userId":    userID.Hex(),
	})
}
func (h *Handler) UpdateUser(c *gin.Context) {
	userId, err := primitive.ObjectIDFromHex(c.Param("userId"))
	if err != nil {
		fmt.Printf("failed to convert params userId to ObjectID: %s\n", err.Error())
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Message: err.Error(),
		})
		return
	}

	existingUser, err := h.Storage.GetUser(userId)
	if err != nil {
		fmt.Printf("failed to get user: %s\n", err.Error())
		c.JSON(http.StatusNotFound, ErrorResponse{
			Message: err.Error(),
		})
		return
	}

	var newUser user.User
	if err := c.BindJSON(&newUser); err != nil {
		fmt.Printf("failed to bind user: %s\n", err.Error())
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Message: err.Error(),
		})
		return
	}

	if newUser.Name != "" {
		existingUser.Name = newUser.Name
	}
	if newUser.Work != "" {
		existingUser.Work = newUser.Work
	}
	if newUser.Age != 0 {
		existingUser.Age = newUser.Age
	}
	if newUser.Salary != 0 {
		existingUser.Salary = newUser.Salary
	}
	if newUser.Email != "" {
		existingUser.Email = newUser.Email
	}

	h.Storage.UpdateUser(userId, &existingUser)

	c.JSON(http.StatusOK, map[string]interface{}{
		"userId": existingUser.Id.Hex(),
	})
}
func (h *Handler) GetAllUsers(c *gin.Context) {
	storage := h.Storage.GetAllUsers()
	c.JSON(http.StatusOK, storage)
}
func (h *Handler) GetUser(c *gin.Context) {
	userId, err := primitive.ObjectIDFromHex(c.Param("userId"))
	if err != nil {
		fmt.Printf("failer convert params userId to ObjectID: %s\n", err.Error())
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Message: err.Error(),
		})
		return
	}

	user, err := h.Storage.GetUser(userId)
	if err != nil {
		fmt.Printf("failed to get user %s\n", err.Error())
		c.JSON(http.StatusNotFound, ErrorResponse{
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, user)
}
func (h *Handler) DeleteUser(c *gin.Context) {
	userId, err := primitive.ObjectIDFromHex(c.Param("userId"))
	if err != nil {
		fmt.Printf("failed to convert userId param to ObjectID: %s\n", err.Error())
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Message: err.Error(),
		})
	}

	h.Storage.DeleteUser(userId)
	c.String(http.StatusOK, "user deleted")
}

func (h *Handler) GetProject(c *gin.Context) {
	userId := c.Param("userId")

	// Проверяем, является ли userId корректным ObjectID
	objID, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		fmt.Printf("failed to convert params userId to ObjectID: %s\n", err.Error())
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Message: "invalid userId format",
		})
		return
	}

	// Вызываем метод для получения проекта
	projects, err := h.Storage.GetProject(objID)
	if err != nil {
		fmt.Printf("failed to get project: %s\n", err.Error())
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, projects)
}
func (h *Handler) GetAllProjects(c *gin.Context) {
	storage := h.Storage.GetAllProjects()
	c.JSON(http.StatusOK, storage)
}
func (h *Handler) UpdateProject(c *gin.Context) {
	projectID := c.Param("projectId")
	if projectID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid projectId"})
		return
	}

	// Преобразуем строку projectID в ObjectID
	projectObjectID, err := primitive.ObjectIDFromHex(projectID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid projectId format"})
		return
	}

	var updateFields map[string]interface{}
	if err := c.BindJSON(&updateFields); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid request body"})
		return
	}

	err = h.Storage.UpdateProject(projectObjectID, updateFields)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "failed to update project", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "project updated successfully"})
}

func (h *Handler) DeleteProject(c *gin.Context) {
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		fmt.Printf("failed to convert id param to ProjectID: %s\n", err.Error())
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Message: err.Error(),
		})
	}

	h.Storage.DeleteProject(id)
	c.String(http.StatusOK, "project deleted")
}
func (h *Handler) DeleteProjects(c *gin.Context) {
	userID := c.Param("userId")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid userId"})
		return
	}

	// Преобразуем строку userID в ObjectID
	userObjectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid userId format"})
		return
	}

	var requestBody struct {
		ProjectIDs []string `json:"projectIDs"`
	}

	if err := c.BindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid request body"})
		return
	}

	// Преобразуем строковые идентификаторы проектов в ObjectID
	projectObjectIDs := make([]primitive.ObjectID, len(requestBody.ProjectIDs))
	for i, projectID := range requestBody.ProjectIDs {
		projectObjectID, err := primitive.ObjectIDFromHex(projectID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "invalid projectID format"})
			return
		}
		projectObjectIDs[i] = projectObjectID
	}

	// Вызовем метод для удаления проектов
	err = h.Storage.DeleteProjects(userObjectID, projectObjectIDs)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "failed to delete projects", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "projects deleted successfully"})
}
