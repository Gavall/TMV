package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"tmv/employee"
	"tmv/storage"

	"github.com/gin-gonic/gin"
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

func (h *Handler) CreateEmployee(c *gin.Context) {
	var employee employee.Employee

	if err := c.BindJSON(&employee); err != nil {
		fmt.Printf("failer to bind employee: %s\n", err.Error())
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Message: err.Error(),
		})
		return
	}

	h.Storage.Insert(&employee)

	c.JSON(http.StatusOK, map[string]interface{}{
		"id": employee.Id,
	})
}

func (h *Handler) UpdateEmployee(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		fmt.Printf("failed to convert params id to int: %s\n", err.Error())
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Message: err.Error(),
		})
		return
	}

	// Получение существующего сотрудника из хранилища
	existingEmployee, err := h.Storage.Get(id)
	if err != nil {
		fmt.Printf("failed to get employee: %s\n", err.Error())
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Message: err.Error(),
		})
		return
	}

	// Чтение новых данных сотрудника из тела запроса
	var newEmployee employee.Employee
	if err := c.BindJSON(&newEmployee); err != nil {
		fmt.Printf("failed to bind employee: %s\n", err.Error())
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Message: err.Error(),
		})
		return
	}

	// Обновление полей сотрудника только для тех, которые были переданы в теле запроса
	if newEmployee.Name != "" {
		existingEmployee.Name = newEmployee.Name
	}
	if newEmployee.Work != "" {
		existingEmployee.Work = newEmployee.Work
	}
	if newEmployee.Age != 0 {
		existingEmployee.Age = newEmployee.Age
	}
	if newEmployee.Salary != 0 {
		existingEmployee.Salary = newEmployee.Salary
	}
	// Добавьте дополнительные условия для других полей, если необходимо

	// Обновление сотрудника в хранилище
	h.Storage.Update(id, &existingEmployee)

	c.JSON(http.StatusOK, map[string]interface{}{
		"id": existingEmployee.Id,
	})
}
func (h *Handler) GetAllEmployees(c *gin.Context) {
	storage := h.Storage.GetAll()
	c.JSON(http.StatusOK, storage)
}
func (h *Handler) GetEmployee(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		fmt.Printf("failer convert params id in int: %s\n", err.Error())
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Message: err.Error(),
		})
		return
	}

	employee, err := h.Storage.Get(id)
	if err != nil {
		fmt.Printf("failed to get employee %s\n ", err.Error())
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, employee)
}

func (h *Handler) DeleteEmployee(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		fmt.Printf("failed to convert id param to int: %s\n", err.Error())
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Message: err.Error(),
		})
		return
	}

	h.Storage.Delete(id)
	c.String(http.StatusOK, "employee deleted")
}
