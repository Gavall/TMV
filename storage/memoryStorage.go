package storage

// import (
// 	"errors"
// 	"sync"
// 	"tmv/user"
// )

// type MemoryStorage struct {
// 	Counter int
// 	Data    map[int]user.User
// 	sync.Mutex
// }

// func NewMemoryStorage() *MemoryStorage {
// 	return &MemoryStorage{
// 		Counter: 1,
// 		Data:    make(map[int]user.User),
// 	}
// }

// func (m *MemoryStorage) GetAll() map[int]user.User {
// 	m.Lock()
// 	defer m.Unlock()

// 	return m.Data
// }
// func (m *MemoryStorage) Get(id int) (user.User, error) {
// 	m.Lock()
// 	defer m.Unlock()
// 	user, ok := m.Data[id]
// 	if !ok {
// 		return user, errors.New("user not found")
// 	}
// 	return user, nil
// }

// func (m *MemoryStorage) Insert(e *user.User) {
// 	m.Lock()
// 	defer m.Unlock()
// 	e.UserId = m.Counter

// 	m.Data[e.UserId] = *e
// 	m.Counter++
// }
// func (m *MemoryStorage) Update(id int, e *user.User) {
// 	m.Lock()
// 	defer m.Unlock()

// 	m.Data[id] = *e
// }
// func (m *MemoryStorage) Delete(id int) {
// 	m.Lock()
// 	defer m.Unlock()

// 	delete(m.Data, id)
// }
