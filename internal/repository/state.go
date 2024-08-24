package repository

import (
	"errors"
	"github.com/google/uuid"
	"sync"
)

type ID string

func (id ID) String() string {
	return string(id)
}

type StateRepository interface {
	Create(value string) (id string)
	Get(id string) (value string, err error)
}

type stateRepository struct {
	cache     map[string]string
	cacheLock sync.RWMutex
}

func NewStateRepository() StateRepository {
	return &stateRepository{
		cache: make(map[string]string),
	}
}

func (repo *stateRepository) Create(value string) (id string) {
	id = uuid.NewString()

	repo.cacheLock.Lock()
	defer repo.cacheLock.Unlock()

	repo.cache[id] = value

	return
}

func (repo *stateRepository) Get(id string) (value string, err error) {
	repo.cacheLock.Lock()
	defer repo.cacheLock.Unlock()

	value, exists := repo.cache[id]
	if !exists {
		return "", errors.New("state not found")
	}

	return value, nil
}
