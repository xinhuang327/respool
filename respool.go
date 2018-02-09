package respool

import "sync"

type IDisposable interface {
	Dispose()
}

type ResourcePool struct {
	MaxSize         int
	buffer          []interface{}
	mutex           *sync.Mutex
	NewResourceFunc func() interface{}
}

func NewResourcePool(maxSize int) *ResourcePool {
	pool := ResourcePool{}
	pool.buffer = []interface{}{}
	pool.mutex = &sync.Mutex{}
	pool.MaxSize = maxSize
	return &pool
}

func (pool *ResourcePool) GetResource() interface{} {
	pool.mutex.Lock()
	defer pool.mutex.Unlock()

	if len(pool.buffer) > 0 {
		res, updatedBuffer := pool.buffer[0], pool.buffer[1:]
		pool.buffer = updatedBuffer
		return res
	} else if pool.NewResourceFunc != nil {
		return pool.NewResourceFunc()
	}

	return nil
}

func (pool *ResourcePool) PutResource(res interface{}) {
	pool.mutex.Lock()
	defer pool.mutex.Unlock()

	if len(pool.buffer) < pool.MaxSize {
		pool.buffer = append(pool.buffer, res)
	} else {
		if disposable, ok := res.(IDisposable); ok {
			disposable.Dispose()
		}
	}
}

func (pool *ResourcePool) Size() int {
	pool.mutex.Lock()
	defer pool.mutex.Unlock()

	return len(pool.buffer)
}
