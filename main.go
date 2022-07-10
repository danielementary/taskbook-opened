package main

import "fmt"

type task struct {
	id          uint
	description string
	completed   bool
}

type note struct {
	id          uint
	description string
}

type taskbook struct {
	boards []*board
}

type board struct {
	name            string
	tasksPending    []*task
	tasksInProgress []*task
	tasksCompleted  []*task
	notes           []*note
}

func (b *board) display() {
	numberOfTasksCompleted := len(b.tasksCompleted)
	numberOfTasks := len(b.tasksPending) + len(b.tasksInProgress) + numberOfTasksCompleted

	fmt.Printf("  #%s [%d/%d]", b.name, numberOfTasksCompleted, numberOfTasks)
}

func (tb *taskbook) newBoard(name string) {
	tasksPending := []*task{}
	tasksInProgress := []*task{}
	tasksCompleted := []*task{}
	notes := []*note{}

	board := board{name, tasksPending, tasksInProgress, tasksCompleted, notes}

	tb.boards = append(tb.boards, &board)
}

func (tb *taskbook) display() {
	for _, board := range tb.boards {
		board.display()
		fmt.Println()
	}
}

func main() {
	tb := taskbook{boards: []*board{}}

	fmt.Println("Taskbook opened!")

	tb.newBoard("Coding")
	tb.newBoard("Chill")

	tb.display()
}
