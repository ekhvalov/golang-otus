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
	doneCh := make(chan struct{})
	for i := 0; i < workersLimit; i++ {
		wg.Add(1)
		go handleTasks(tasksCh, errorsCh, &wg)
	}
	go countErrors(errorsCh, errorsLimit, doneCh)
	result := produceTasks(doneCh, tasks, tasksCh)
	close(tasksCh)
	wg.Wait()
	close(errorsCh)
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

func countErrors(errorsCh <-chan error, errorsLimit int, doneCh chan<- struct{}) {
	errorsCount := 0
	for range errorsCh {
		errorsCount++
		if errorsCount == errorsLimit {
			doneCh <- struct{}{}
			// don't break to drain remain errors
		}
	}
}

func produceTasks(doneCh <-chan struct{}, tasks []Task, tasksCh chan<- Task) error {
	for _, task := range tasks {
		select {
		case <-doneCh:
			return ErrErrorsLimitExceeded
		default:
		}
		select {
		case tasksCh <- task:
		case <-doneCh:
			return ErrErrorsLimitExceeded
		}

	}
	return nil
}
