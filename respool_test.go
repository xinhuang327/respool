package respool

import (
	"fmt"
	"math/rand"
	"sync"
	"testing"
	"time"
)

type ResObj struct {
	Value string
}

func (self *ResObj) DoWork() {
	rand.Seed(time.Now().UnixNano())
	timeToExecute := time.Duration(rand.Float64()*100) * time.Millisecond
	<-time.After(timeToExecute)
	fmt.Println("### My value is", self.Value, " ### Used time:", timeToExecute)
}

func (self *ResObj) Dispose() {
	fmt.Println("Disposing...")
}

func Test_ResourcePool(t *testing.T) {
	mutex := &sync.Mutex{}
	total := 0

	pool := NewResourcePool(10)

	pool.NewResourceFunc = func() interface{} {
		mutex.Lock()
		defer mutex.Unlock()
		total += 1
		return &ResObj{
			Value: time.Now().String(),
		}
	}

	jobCount := 500
	var waitGroup sync.WaitGroup
	waitGroup.Add(jobCount)

	for i := 0; i < jobCount; i++ {
		go func() {
			if res, ok := pool.GetResource().(*ResObj); ok {
				UseResource("A", pool, res)
			} else {
				t.Error("GetResource returns nil")
			}
			waitGroup.Done()
		}()
		<-time.After(5 * time.Millisecond)
	}

	waitGroup.Wait()

	fmt.Println(pool.Size())
	fmt.Println("Allocated:", total)

}

func UseResource(prefix string, pool *ResourcePool, res *ResObj) {
	defer pool.PutResource(res)
	fmt.Println(prefix, pool.Size())
	res.DoWork()
}
