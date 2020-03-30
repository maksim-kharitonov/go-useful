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
	// очередь на получение ресурса
	subscribers []chan int
	// мьютекс для пула ожидания
	smu *sync.Mutex
}

// Создание пустого пула ресурсов
func NewFixedPool(capacity int, factory PooledObjectFactory) *FixedPool {
	return &FixedPool{
		available: make([]PooledObject, 0),
		inUse:     make([]PooledObject, 0),
		capacity:  capacity,
		mu:        new(sync.Mutex),
		factory:   factory,
		subscribers: make([]chan int, 0),
		smu: 			 new(sync.Mutex)
	}
}

func (p *FixedPool) Get() (PooledObject, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	var obj PooledObject
	var err error

	if len(p.available) == 0 {
		if len(p.inUse) == p.capacity {
			p.smu.Lock()
			c := make(chan int)
			p.subscribers = append(p.subscribers,c)
			p.smu.Unlock()

			resTimeout := time.After(3000 * time.Millisecond)

			p.mu.Unlock()

			select {
			case <-c:
				p.mu.Lock()
				obj, p.available = p.available[0], p.available[1:]
				err = nil
				break
			case <-resTimeout:
				err = errors.New("fixed Pool reached maximum capacity")
				break
			}

		} else {
			obj, err = p.factory.Create()
			p.inUse = append(p.inUse, obj)
		}
	} else {
		obj, p.available = p.available[0], p.available[1:]
		err = nil
		p.inUse = append(p.inUse, obj)
	}

	//p.mu.Unlock()

	return obj, err
}

func (p *FixedPool) Return(obj PooledObject) error {
	obj.Reset()

	var err error

	p.mu.Lock()
	if idx := findIndex(obj, p.inUse); idx != -1 {
		p.inUse = append(p.inUse[:idx], p.inUse[idx+1:]...)
		p.available = append(p.available, obj)
		if (len(p.subscribers) != 0) {
			p.smu.Lock()
			p.subscribers[0]<-1
			p.subscribers = p.subscribers[1:]
			p.smu.Unlock()
		}
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
