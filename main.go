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
	number_of_boards uint
	boards           []*board
}

type board struct {
	id    uint
	name  string
	tasks []*task
	notes []*note
}

func (tb *taskbook) newBoard(name string) {
	var id = tb.number_of_boards
	tb.number_of_boards++

	var tasks = []*task{}
	var notes = []*note{}

	var board = board{id, name, tasks, notes}

	tb.boards = append(tb.boards, &board)
}

func main() {
	var tb = taskbook{number_of_boards: 0, boards: []*board{}}
	fmt.Println("Taskbook opened!")

	tb.newBoard("Coding")
}
