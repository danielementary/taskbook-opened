package main

import (
	"errors"
	"flag"
	"fmt"
	"strings"
)

type taskStatus int

const (
	InProgress taskStatus = iota
	Pending
	Completed
	Note
)

type task struct {
	id          uint
	description string
	status      taskStatus
}

func (t *task) display() {
	fmt.Printf("    %d. ", t.id)

	switch t.status {
	case InProgress:
		fmt.Print("[/]")
	case Pending:
		fmt.Print("[ ]")
	case Completed:
		fmt.Print("[X]")
	case Note:
		fmt.Print("-N-")
	}

	fmt.Printf(" %s", t.description)
}

type board struct {
	name    string
	counter uint
	tasks   []*task
}

func (b *board) getNumberOfTasksCompleted() (numberOfTasksCompleted int) {
	numberOfTasksCompleted = 0

	for _, t := range b.tasks {
		if t.status == Completed {
			numberOfTasksCompleted++
		}
	}

	return
}

func (b *board) display() {
	numberOfTasksCompleted := b.getNumberOfTasksCompleted()
	numberOfTasks := len(b.tasks) + numberOfTasksCompleted

	fmt.Printf("  #%s [%d/%d]\n", b.name, numberOfTasksCompleted, numberOfTasks)

	if numberOfTasks <= 0 {
		fmt.Print("    This board is empty.")
	} else {
		for _, t := range b.tasks {
			t.display()
		}
	}
	fmt.Print("\n\n")
}

type taskbook struct {
	boards []*board
}

func (tb *taskbook) newBoard(name string) {
	counter := uint(1)
	tasks := []*task{}

	board := board{name, counter, tasks}

	tb.boards = append(tb.boards, &board)
}

func (tb *taskbook) display() {
	for _, board := range tb.boards {
		board.display()
	}
}

func parseBoardName(s *string) (*string, error) {
	boardName := strings.Split(*s, " ")[0]

	if boardName[0] != '#' || len(boardName) <= 1 {
		return nil, errors.New("invalid board name")
	}

	return &boardName, nil
}

func parseDescription(s *string) (*string, error) {
	description := strings.Join(strings.Split(*s, " ")[1:], " ")

	if len(description) <= 1 {
		return nil, errors.New("invalid description")
	}

	return &description, nil
}

func main() {
	fmt.Print(" Taskbook opened!\n\n")

	tb := taskbook{boards: []*board{}}
	defer tb.display()

	tb.newBoard("Coding")
	tb.boards[0].tasks = append(tb.boards[0].tasks, &task{id: 1, description: "implement taskbook opened!", status: Pending})

	tb.newBoard("Chill")

	taskPtr := flag.String("task", "", "the description of the new task to add preceded by the corresponding #board")
	notePtr := flag.String("note", "", "the description of the new note to add preceded by the corresponding #board")

	flag.Parse()

	if len(*taskPtr) > 0 {
		// add a new task
	}

	if len(*notePtr) > 0 {
		// add a new note
	}
}
