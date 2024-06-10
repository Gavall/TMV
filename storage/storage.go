package storage

import (
	"tmv/employee"
)

type Storage interface {
	GetAll() map[int]employee.Employee
	Get(id int) (employee.Employee, error)
	Insert(e *employee.Employee)
	Update(id int, e *employee.Employee)
	Delete(id int)

	// GetAll() ([]employee.Employee, error)
	// Get(id int) (employee.Employee, error)
	// Insert(e *employee.Employee) error
	// Update(id int, e *employee.Employee) error
	// Delete(id int) error
}
