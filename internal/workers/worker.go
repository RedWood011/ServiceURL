// Package workers Пакет для конкурентной работы с потоками
package workers

import (
	"context"
	"fmt"
	"log"
	"sync"
)

// WorkerPool Пул потоков c функцией исполнения
type WorkerPool struct {
	numWorkers int
	input      chan func(ctx context.Context) error
}

// New Создание пула
func New(numWorkers, buf int) *WorkerPool {
	return &WorkerPool{
		numWorkers: numWorkers,
		input:      make(chan func(ctx context.Context) error, buf),
	}
}

// Run Запуск потоков
func (wp *WorkerPool) Run(ctx context.Context) {
	wg := &sync.WaitGroup{}
	for i := 0; i < wp.numWorkers; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			fmt.Printf("Start worker %d \n", i)

		loop:
			for {
				select {

				case f := <-wp.input:
					err := f(ctx)
					if err != nil {
						fmt.Printf("Error in worker %d: %v\n", i, err.Error())
					}
				case <-ctx.Done():
					break loop
				}
			}
			log.Printf("Finish worker %d \n", i)

		}(i)
	}
	wg.Wait()
	close(wp.input)
}

// Add Добавление задачи
func (wp *WorkerPool) Add(job func(ctx context.Context) error) {
	wp.input <- job
}
