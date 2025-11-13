package main

import (
	"sync"
	"sync/atomic"
)

var (
	tasks       = sync.Map{}
	taskCounter int32
)

func NextTaskID() int {
	return int(atomic.AddInt32(&taskCounter, 1))
}
func SaveTask(task *Task) {
	tasks.Store(task.ID, task)
}

func GetTask(id int) (*Task, bool) {
	t, ok := tasks.Load(id)
	if !ok {
		return nil, false
	}
	return t.(*Task), true
}

func GetTasks(ids []int) []*Task {
	var result []*Task
	for _, id := range ids {
		if t, ok := GetTask(id); ok {
			result = append(result, t)
		}
	}
	return result
}
