package temporal_storage

import (
	"errors"
	"fmt"
)

type TemporalStorage interface {
	Set(id int64, key string, value string) error
	Del(id int64) error
	Get(id int64, key string) (string, error)
	DelValue(id int64, key string) error
}

type MemoryStorage struct {
	data map[int64]map[string]string
}

func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{
		data: make(map[int64]map[string]string),
	}
}

func (m *MemoryStorage) Set(id int64, key string, value string) error {
	if _, ok := m.data[id]; !ok {
		m.data[id] = make(map[string]string)
	}

	m.data[id][key] = value
	return nil
}

func (m *MemoryStorage) Del(id int64) error {
	delete(m.data, id)
	return nil
}

func (m *MemoryStorage) Get(id int64, key string) (string, error) {
	if _, ok := m.data[id]; ok {
		if value, ok := m.data[id][key]; ok {
			return value, nil
		}
	}

	return "", errors.New(fmt.Sprintf("No info for chat by key = %s", key))
}

func (m *MemoryStorage) DelValue(id int64, key string) error {
	if _, ok := m.data[id]; ok {
		delete(m.data[id], key)
		return nil
	}

	return errors.New("Chat is not found")
}
