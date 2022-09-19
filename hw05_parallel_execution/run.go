package hw05parallelexecution

import (
	"errors"
	"sync"
)

var (
	ErrErrorsLimitExceeded = errors.New("errors limit exceeded")
	ErrInsufficientWorkers = errors.New("insufficient workers")
)

type Task func() error

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, workersLimit, errorsLimit int) error {
	if workersLimit <= 0 {
		return ErrInsufficientWorkers
	}
	if len(tasks) < workersLimit {
		workersLimit = len(tasks)
	}
	wg := sync.WaitGroup{}
	tasksCh := make(chan Task)
	errorsCh := make(chan error, workersLimit)
	for i := 0; i < workersLimit; i++ {
		wg.Add(1)
		go handleTasks(tasksCh, errorsCh, &wg)
	}
	result := produce(&tasks, tasksCh, errorsCh, errorsLimit)
	close(tasksCh)
	wg.Wait()
	return result
}

func handleTasks(tasksCh <-chan Task, errorsCh chan<- error, wg *sync.WaitGroup) {
	for task := range tasksCh {
		err := task()
		if err != nil {
			errorsCh <- err
		}
	}
	wg.Done()
}

func produce(tasks *[]Task, tasksCh chan<- Task, errorsCh <-chan error, errorsLimit int) error {
	errorsCount := 0
	for _, task := range *tasks {
		isTaskSent := false
		for !isTaskSent {
			if errorsLimit > 0 && errorsCount == errorsLimit {
				return ErrErrorsLimitExceeded
			}
			select {
			case <-errorsCh:
				errorsCount++
				continue
			default:
			}
			select {
			case <-errorsCh:
				errorsCount++
				continue
			case tasksCh <- task:
				isTaskSent = true
			}
		}
	}
	return nil
}
