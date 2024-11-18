package main

import (
	"fmt"
	"time"

	"math/rand"
)

type TaskList struct {
	Tasks []Task
}

type Task struct {
	Name     string
	Duration time.Duration
}

type Worker struct {
	Name             string
	ActiveTask       *Task
	ActiveTaskDoneAt time.Time
	t                *time.Ticker
}

func main() {
	taskList := TaskList{
		Tasks: []Task{
			{
				Name:     "Gather wood",
				Duration: time.Second * 2,
			},
			{
				Name:     "Gather stone",
				Duration: time.Second * 2,
			},
			{
				Name:     "Gather food",
				Duration: time.Second * 2,
			},
		},
	}

	w := Worker{
		Name: "Bob",
		t:    time.NewTicker(time.Second),
	}

	for {
		<-w.t.C
		if w.ActiveTask == nil {
			newTask, ok := taskList.GetTask()
			if !ok {
				fmt.Println("no more tasks")
				return
			}

			w.ActiveTask = &newTask
			w.ActiveTaskDoneAt = time.Now().Add(newTask.Duration)
			fmt.Println("new task:", w.ActiveTask.Name)
		}

		if time.Now().After(w.ActiveTaskDoneAt) {
			fmt.Println("completing task:", w.ActiveTask.Name)
			w.ActiveTask = nil
		} else {
			fmt.Println("working on task:", w.ActiveTask.Name)
		}
	}
}

func (tl *TaskList) GetTask() (Task, bool) {
	if len(tl.Tasks) == 0 {
		return Task{}, false
	}

	i := rand.Intn(len(tl.Tasks))
	t := tl.Tasks[i]

	tl.Tasks = append(tl.Tasks[:i], tl.Tasks[i+1:]...)

	return t, true
}
