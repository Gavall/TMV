package storage

import (
	"errors"
	"sync"
	"tmv/employee"
)

type MemoryStorage struct {
	Counter int
	Data    map[int]employee.Employee
	sync.Mutex
}

func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{
		Counter: 1,
		Data:    make(map[int]employee.Employee),
	}
}

func (m *MemoryStorage) GetAll() map[int]employee.Employee {
	m.Lock()
	defer m.Unlock()

	return m.Data
}
func (m *MemoryStorage) Get(id int) (employee.Employee, error) {
	m.Lock()
	defer m.Unlock()
	employee, ok := m.Data[id]
	if !ok {
		return employee, errors.New("employee not found")
	}
	return employee, nil
}

func (m *MemoryStorage) Insert(e *employee.Employee) {
	m.Lock()
	defer m.Unlock()
	e.Id = m.Counter

	m.Data[e.Id] = *e
	m.Counter++
}
func (m *MemoryStorage) Update(id int, e *employee.Employee) {
	m.Lock()
	defer m.Unlock()

	m.Data[id] = *e
}
func (m *MemoryStorage) Delete(id int) {
	m.Lock()
	defer m.Unlock()

	delete(m.Data, id)
}
