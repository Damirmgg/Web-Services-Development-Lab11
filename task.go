package main

import (
	"errors"
	"strings"
	"sync"
)

type Task struct {
	ID        int
	Title     string
	Completed bool
}

type PageData struct {
	Tasks          []Task
	TotalTasks     int
	CompletedTasks int
	PendingTasks   int
	Error          string
	Title          string
}

var (
	tasks      []Task
	nextTaskID = 1
	tasksMu    sync.RWMutex
)

func addTask(title string) error {
	title = strings.TrimSpace(title)
	if title == "" {
		return errors.New("Введите название задачи: оно не может быть пустым.")
	}
	tasksMu.Lock()
	defer tasksMu.Unlock()
	tasks = append(tasks, Task{ID: nextTaskID, Title: title})
	nextTaskID++
	return nil
}

// Вызывается только при удержании tasksMu.
func findTaskIndex(id int) int {
	for i := range tasks {
		if tasks[i].ID == id {
			return i
		}
	}
	return -1
}

func toggleTask(id int) bool {
	tasksMu.Lock()
	defer tasksMu.Unlock()
	i := findTaskIndex(id)
	if i < 0 {
		return false
	}
	tasks[i].Completed = !tasks[i].Completed
	return true
}

func deleteTask(id int) bool {
	tasksMu.Lock()
	defer tasksMu.Unlock()
	i := findTaskIndex(id)
	if i < 0 {
		return false
	}
	copy(tasks[i:], tasks[i+1:])
	tasks[len(tasks)-1] = Task{}
	tasks = tasks[:len(tasks)-1]
	return true
}

func taskPageData() PageData {
	tasksMu.RLock()
	defer tasksMu.RUnlock()
	// Копия позволяет отрисовать страницу после освобождения блокировки.
	data := PageData{Tasks: append([]Task(nil), tasks...), TotalTasks: len(tasks)}
	for _, task := range tasks {
		if task.Completed {
			data.CompletedTasks++
		}
	}
	data.PendingTasks = data.TotalTasks - data.CompletedTasks
	return data
}
