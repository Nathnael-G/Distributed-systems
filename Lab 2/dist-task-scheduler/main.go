// Nathnael Getachew
// UGR/8932/13
package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// Task represents a unit of work.
type Task struct {
	id      int
	data    string
	retries int
}

// Worker function that processes tasks. If a worker fails, the task will be sent to failChan.
func worker(id int, taskChan <-chan Task, wg *sync.WaitGroup, failChan chan<- Task, maxRetries int, failedTasks map[int]bool, mu *sync.Mutex) {
	defer wg.Done()

	for task := range taskChan {
		fmt.Printf("Worker %d started task %d: %s (Retry: %d)\n", id, task.id, task.data, task.retries)

		// Simulate random failure (30% chance of failure)
		if rand.Float32() < 0.3 {
			fmt.Printf("Worker %d failed on task %d\n", id, task.id)

			// Increment retries and check if it exceeds maxRetries
			task.retries++
			if task.retries > maxRetries {
				mu.Lock()
				if !failedTasks[task.id] {
					fmt.Printf("Task %d has exceeded maximum retries and will not be retried.\n", task.id)
					failedTasks[task.id] = true
				}
				mu.Unlock()
				return
			}

			failChan <- task // Send the failed task for reassignment
			return
		}

		// Simulate task processing time
		time.Sleep(time.Duration(rand.Intn(3)+1) * time.Second)
		fmt.Printf("Worker %d completed task %d\n", id, task.id)
	}
}

func main() {
	rand.Seed(time.Now().UnixNano())

	// Define a set of tasks to be executed
	tasks := []Task{
		{id: 1, data: "Task 1"},
		{id: 2, data: "Task 2"},
		{id: 3, data: "Task 3"},
		{id: 4, data: "Task 4"},
		{id: 5, data: "Task 5"},
	}

	// Channels for task distribution and failure handling
	taskChan := make(chan Task, len(tasks))
	failChan := make(chan Task, len(tasks))

	// WaitGroup to ensure all workers finish their tasks
	var wg sync.WaitGroup

	// Number of workers (simulating processors)
	numWorkers := 3

	// Maximum number of retries for each task
	maxRetries := 3

	// Map to track which tasks have failed
	failedTasks := make(map[int]bool)
	var mu sync.Mutex // Mutex to protect access to failedTasks

	// Start worker goroutines
	for i := 1; i <= numWorkers; i++ {
		wg.Add(1)
		go worker(i, taskChan, &wg, failChan, maxRetries, failedTasks, &mu)
	}

	// Distribute tasks to workers in a round-robin manner
	counter := 0
	for _, task := range tasks {
		workerId := (counter % numWorkers) + 1
		fmt.Printf("Assigning task %d to worker %d\n", task.id, workerId)
		taskChan <- task
		counter++
	}
	close(taskChan)

	// Handle failed tasks by redistributing them
	go func() {
		for failedTask := range failChan {
			mu.Lock()
			// Find the next worker in round-robin fashion for reassignment
			nextWorkerId := (len(failedTasks) % numWorkers) + 1
			fmt.Printf("Reassigning failed task %d to worker %d\n", failedTask.id, nextWorkerId)
			wg.Add(1)
			go worker(nextWorkerId, taskChan, &wg, failChan, maxRetries, failedTasks, &mu)
			mu.Unlock()
		}
	}()

	// Wait for all workers to finish
	wg.Wait()
	close(failChan)

	fmt.Println("All tasks completed.")
}
