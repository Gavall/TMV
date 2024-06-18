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
func (h *Handler) GetProjectsByUser(c *gin.Context) {
	// Получаем userId из параметров запроса
	userIdHex := c.Param("userId")
	userId, err := primitive.ObjectIDFromHex(userIdHex)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid userId format"})
		return
	}

	// Получаем проекты пользователя
	projects, err := h.Storage.GetProjectByUser(userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	// Возвращаем список проектов
	c.JSON(http.StatusOK, projects)
}
func (h *Handler) GetProject(c *gin.Context) {
	// Получаем userId и projectId из параметров запроса
	userIdHex := c.Param("userId")
	userId, err := primitive.ObjectIDFromHex(userIdHex)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid userId format"})
		return
	}

	projectIdHex := c.Param("projectId")
	projectId, err := primitive.ObjectIDFromHex(projectIdHex)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid projectId format"})
		return
	}

	// Получаем проект
	project, err := h.Storage.GetProject(userId, projectId)
	if err != nil {
		if err.Error() == "проект не найден" {
			c.JSON(http.StatusNotFound, gin.H{"message": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		}
		return
	}

	// Возвращаем проект
	c.JSON(http.StatusOK, project)
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

func (h *Handler) GetAlltasks(c *gin.Context) {
	storage := h.Storage.GetAllTasks()
	c.JSON(http.StatusOK, storage)
}
func (h *Handler) GetTasksByProject(c *gin.Context) {
	projectIdParam := c.Param("projectId")
	projectId, err := primitive.ObjectIDFromHex(projectIdParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid projectId format"})
		return
	}

	tasks, err := h.Storage.GetTasksByProject(projectId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, tasks)
}
func (h *Handler) CreateTask(c *gin.Context) {
	var task project.Task

	if err := c.BindJSON(&task); err != nil {
		fmt.Printf("failed to bind project: %s\n", err.Error())
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Message: err.Error(),
		})
		return
	}

	projectIDStr := c.Param("projectId")
	projectID, err := primitive.ObjectIDFromHex(projectIDStr)
	if err != nil {
		fmt.Printf("failed to convert userId to ObjectID: %s\n", err.Error())
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Message: "invalid userId format",
		})
		return
	}

	err = h.Storage.InsertTask(&task, projectID)
	if err != nil {
		fmt.Printf("failed to insert project: %s\n", err.Error())
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"taskId":    task.ID,
		"projectId": projectID.Hex(),
	})
}
func (h *Handler) GetTask(c *gin.Context) {
	projectIdParam := c.Param("projectId")
	taskIdParam := c.Param("taskId")

	projectId, err := primitive.ObjectIDFromHex(projectIdParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid projectId format"})
		return
	}

	taskId, err := primitive.ObjectIDFromHex(taskIdParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid taskId format"})
		return
	}

	task, err := h.Storage.GetTask(projectId, taskId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	if task == nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "task not found"})
		return
	}

	c.JSON(http.StatusOK, task)
}
func (h *Handler) DeleteTask(c *gin.Context) {
	projectIdParam := c.Param("projectId")
	taskIdParam := c.Param("taskId")

	projectId, err := primitive.ObjectIDFromHex(projectIdParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid projectId format"})
		return
	}

	taskId, err := primitive.ObjectIDFromHex(taskIdParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid taskId format"})
		return
	}

	err = h.Storage.DeleteTask(projectId, taskId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "task deleted successfully"})
}
func (h *Handler) DeleteTasks(c *gin.Context) {
	projectIdParam := c.Param("projectId")

	// Проверка правильности формата projectId
	projectId, err := primitive.ObjectIDFromHex(projectIdParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid projectId format"})
		return
	}

	var taskIds []string
	if err := c.ShouldBindJSON(&taskIds); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid taskIds format"})
		return
	}

	// Преобразование taskIds в []primitive.ObjectID
	var objectIDs []primitive.ObjectID
	for _, id := range taskIds {
		objectID, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "invalid taskId format"})
			return
		}
		objectIDs = append(objectIDs, objectID)
	}

	// Вызов метода DeleteTasks для удаления задач
	err = h.Storage.DeleteTasks(projectId, objectIDs)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "tasks deleted successfully"})
}
func (h *Handler) UpdateTask(c *gin.Context) {
	projectIdParam := c.Param("projectId")
	taskIdParam := c.Param("taskId")

	projectId, err := primitive.ObjectIDFromHex(projectIdParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid projectId format"})
		return
	}

	taskId, err := primitive.ObjectIDFromHex(taskIdParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid taskId format"})
		return
	}

	var updateFields map[string]interface{}
	if err := c.ShouldBindJSON(&updateFields); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid task format"})
		return
	}

	err = h.Storage.UpdateTask(projectId, taskId, updateFields)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "task updated successfully"})
}
