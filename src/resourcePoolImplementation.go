package go_object_pool

import (
	"errors"
	"sync"
)

type PooledObject interface {
	Reset()
}

type PooledObjectFactory interface {
	Create() (PooledObject, error)
}

type Pool interface {
	Get() (PooledObject, error)
	Return(obj PooledObject) error
}

type FixedPool struct {
	// Доступные ресурсы
	available []PooledObject
	// Занятые ресурсы
	inUse []PooledObject
	// Предельный размер
	capacity int
	// мьютекс для безопасной работы с пулом
	mu *sync.Mutex
	// фабрика создающая ресурс
	factory PooledObjectFactory
}

// Создание пустого пула ресурсов
func NewFixedPool(capacity int, factory PooledObjectFactory) *FixedPool {
	return &FixedPool{
		available: make([]PooledObject, 0),
		inUse:     make([]PooledObject, 0),
		capacity:  capacity,
		mu:        new(sync.Mutex),
		factory:   factory,
	}
}

func (p *FixedPool) Get() (PooledObject, error) {
	p.mu.Lock()

	var obj PooledObject
	var err error

	if len(p.available) == 0 {
		if len(p.inUse) == p.capacity {
			err = errors.New("fixed Pool reached maximum capacity")
		} else {
			obj, err = p.factory.Create()
			p.inUse = append(p.inUse, obj)
		}
	} else {
		obj, p.available = p.available[0], p.available[1:]
		err = nil
		p.inUse = append(p.inUse, obj)
	}

	p.mu.Unlock()

	return obj, err
}

func (p *FixedPool) Return(obj PooledObject) error {
	obj.Reset()

	var err error

	p.mu.Lock()
	if idx := findIndex(obj, p.inUse); idx != -1 {
		p.inUse = append(p.inUse[:idx], p.inUse[idx+1:]...)
		p.available = append(p.available, obj)
		err = nil
	} else {
		err = errors.New("unrecognized pooled object returned")
	}
	p.mu.Unlock()

	return err
}

func findIndex(target PooledObject, slice []PooledObject) int {
	for idx, obj := range slice {
		if target == obj {
			return idx
		}
	}

	return -1
}
