package main

import (
	"fmt"
	"time"
)

type Task struct {
	Name     string
	Consumes map[string]int
	Produces map[string]int
}

type Worker struct {
	Name       string
	ActiveTask *Task
}

type Game struct {
	AvailableTasks []Task
	Workers        []Worker
	State          map[string]int
}

func main() {
	game := Game{
		AvailableTasks: generateAvailableTasks(),
		Workers: []Worker{
			{
				Name: "Bob",
			},
		},
		State: map[string]int{
			"Wood":  0,
			"Stone": 0,
			"Ore":   0,
			"Pipe":  0,
			"House": 0,
		},
	}

	ticker := time.NewTicker(time.Second)

	for {
		<-ticker.C

		fmt.Printf(
			"** Resources: W %d, S %d, O %d, P %d, H %d\n",
			game.State["Wood"],
			game.State["Stone"], 
			game.State["Ore"],
			game.State["Pipe"],
			game.State["House"],
		)

		if game.State["House"] == 1 {
			fmt.Println("3 houses built")
			return
		}

		if game.Workers[0].ActiveTask == nil {
			// Don't have 3 houses, build a house
			task := game.generateActiveTask("House")
			game.Workers[0].ActiveTask = task
			fmt.Println("   starting new task:", task.Name)
			continue
		}

		// Already had a task, complete it
		fmt.Println("   completing task:", game.Workers[0].ActiveTask.Name)
		for produce, produceAmount := range game.Workers[0].ActiveTask.Produces {
			game.State[produce] += produceAmount
			for consumes, consumesAmount := range game.Workers[0].ActiveTask.Consumes {
				game.State[consumes] -= consumesAmount

			}
		}

		game.Workers[0].ActiveTask = nil
	}
}

func (g *Game) generateActiveTask(target string) *Task {
	for _, task := range g.AvailableTasks {
		if _, ok := task.Produces[target]; ok {
			// Check requirements
			for consume, consumeAmount := range task.Consumes {
				fmt.Printf("   %s requires %d x %s, got %d\n", target, consumeAmount, consume, g.State[consume])
				if g.State[consume] < consumeAmount {
					return g.generateActiveTask(consume)
				}
			}
			return &task
		}
	}

	return nil
}

func generateAvailableTasks() []Task {
	tasks := []Task{
		{
			Name:     "Collect Wood",
			Consumes: map[string]int{},
			Produces: map[string]int{
				"Wood": 1,
			},
		},
		{
			Name:     "Collect Stone",
			Consumes: map[string]int{},
			Produces: map[string]int{
				"Stone": 1,
			},
		},
		{
			Name:     "Collect Ore",
			Consumes: map[string]int{},
			Produces: map[string]int{
				"Ore": 1,
			},
		},
		{
			Name: "Make Pipes",
			Consumes: map[string]int{
				"Ore": 3,
			},
			Produces: map[string]int{
				"Pipe": 1,
			},
		},
		{
			Name: "Build House",
			Consumes: map[string]int{
				"Wood":  4,
				"Stone": 2,
				"Pipe":  2,
			},
			Produces: map[string]int{
				"House": 1,
			},
		},
	}

	return tasks
}
