package main

import (
	"fmt"
	"strings"
	"time"
)

type Task struct {
	Name     string
	Consumes map[string]int
	Produces map[string]int
	Requires string
}

type Worker struct {
	Name       string
	ActiveTask *Task
}

type Game struct {
	AvailableTasks     []Task
	AvailableBuildings []Building
	Workers            []Worker
	State              map[string]int
	Buildings          map[string]int
}

func main() {
	game := Game{
		AvailableTasks:     generateAvailableTasks(),
		AvailableBuildings: generateBuildingList(),
		Workers: []Worker{
			{
				Name: "Bob",
			},
		},
		State: map[string]int{
			"Wood":     0,
			"Stone":    0,
			"Ore":      0,
			"Pipe":     0,
			"Resident": 0,
		},
		Buildings: map[string]int{
			"Forest":  1,
			"Quarry":  0,
			"Mine":    0,
			"Factory": 0,
		},
	}

	ticker := time.NewTicker(time.Second)

	for {
		<-ticker.C

		fmt.Printf(
			"\n** Resources: W %d, S %d, O %d, P %d -- F %d, Q %d, M %d, Fc %d, H %d\n",
			game.State["Wood"],
			game.State["Stone"],
			game.State["Ore"],
			game.State["Pipe"],

			game.Buildings["Forest"],
			game.Buildings["Quarry"],
			game.Buildings["Mine"],
			game.Buildings["Factory"],
			game.Buildings["House"],
		)

		if game.State["Resident"] == 1 {
			fmt.Println("resident moved in!")
			return
		}

		if game.Workers[0].ActiveTask == nil {
			// Final target, get a Resident
			task := game.generateActiveTask("Resident")
			game.Workers[0].ActiveTask = task
			if task != nil {
				//fmt.Println("   starting new task:", task.Name)
			}
			continue
		}

		// Already had a task, complete it
		// fmt.Println("   completing task:", game.Workers[0].ActiveTask.Name)
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
	// Check for a build task first
	if strings.HasPrefix(target, "build:") {
		buildingName := strings.TrimPrefix(target, "build:")

		for _, building := range g.AvailableBuildings {
			if building.Name == buildingName {
				// Is there enough resources for the build
				for consume, consumeAmount := range building.Requires {
					// fmt.Printf("   %s requires %d x %s, got %d\n", target, consumeAmount, consume, g.State[consume])
					if g.State[consume] < consumeAmount {
						return g.generateActiveTask(consume)
					}
				}

				// Enough resources
				fmt.Println("building", buildingName)
				g.Buildings[buildingName]++
				for consumes, consumesAmount := range building.Requires {
					g.State[consumes] -= consumesAmount
				}
				return nil
			}
		}
	}

	// Do a resource task
	for _, task := range g.AvailableTasks {
		if _, ok := task.Produces[target]; ok {
			// Check requirements

			// Do we have the required building
			if task.Requires != "" && g.Buildings[task.Requires] == 0 {
				return g.generateActiveTask(fmt.Sprintf("build:%s", task.Requires))
			}

			// Is there enough resources for the task
			for consume, consumeAmount := range task.Consumes {
				//fmt.Printf("   %s requires %d x %s, got %d\n", target, consumeAmount, consume, g.State[consume])
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
			Requires: "Forest",
		},
		{
			Name:     "Collect Stone",
			Consumes: map[string]int{},
			Produces: map[string]int{
				"Stone": 1,
			},
			Requires: "Quarry",
		},
		{
			Name:     "Collect Ore",
			Consumes: map[string]int{},
			Produces: map[string]int{
				"Ore": 1,
			},
			Requires: "Mine",
		},
		{
			Name: "Make Pipes",
			Consumes: map[string]int{
				"Ore": 3,
			},
			Produces: map[string]int{
				"Pipe": 1,
			},
			Requires: "Factory",
		},
		{
			Name:     "Move in",
			Consumes: map[string]int{},
			Produces: map[string]int{
				"Resident": 1,
			},
			Requires: "House",
		},
	}

	return tasks
}

type Building struct {
	Name     string
	Requires map[string]int
}

func generateBuildingList() []Building {
	buildings := []Building{
		{
			Name: "Forest",
			Requires: map[string]int{
				"Wood": 10,
			},
		},
		{
			Name: "Quarry",
			Requires: map[string]int{
				"Wood": 10,
			},
		},
		{
			Name: "Mine",
			Requires: map[string]int{
				"Stone": 10,
			},
		},
		{
			Name: "Factory",
			Requires: map[string]int{
				"Wood":  5,
				"Stone": 5,
			},
		},
		{
			Name: "House",
			Requires: map[string]int{
				"Wood":  4,
				"Stone": 2,
				"Pipe":  2,
			},
		},
	}

	return buildings
}
