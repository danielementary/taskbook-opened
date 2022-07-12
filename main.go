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

	fmt.Printf(" %s\n", t.description)
}

type board struct {
	name    string
	counter uint
	tasks   []*task
}

func newBoard(name string) *board {
	counter := uint(1)
	tasks := []*task{}

	return &board{name, counter, tasks}
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

	fmt.Printf("  %s [%d/%d]\n", b.name, numberOfTasksCompleted, numberOfTasks)

	if numberOfTasks <= 0 {
		fmt.Print("    This board is empty.")
	} else {
		for _, t := range b.tasks {
			t.display()
		}
	}
	fmt.Print("\n")
}

type taskbook struct {
	boards map[string]*board
}

func newTaskbook() *taskbook {
	boards := make(map[string]*board)

	return &taskbook{boards}
}

func (tb *taskbook) display() {
	for _, board := range tb.boards {
		board.display()
	}
}

func (tb *taskbook) addTask(s string, taskStatus taskStatus) {
	boardName, taskDescription, err := parseBoardNameAndTaskDescription(s)

	if err != nil {
		fmt.Println(" Failed to add task:", err)
	} else {
		tb.addTaskToBoard(boardName, taskDescription, taskStatus)
	}
}

func (tb *taskbook) addTaskToBoard(boardName, taskDescription string, taskStatus taskStatus) {
	board, found := tb.boards[boardName]

	if !found {
		tb.boards[boardName] = newBoard(boardName)
		board = tb.boards[boardName]
	}

	id := board.counter
	board.counter++
	description := taskDescription
	status := taskStatus

	board.tasks = append(board.tasks, &task{id, description, status})
}

func parseBoardNameAndTaskDescription(s string) (string, string, error) {
	splitS := strings.Split(s, " ")

	boardName := splitS[0]
	description := strings.Join(splitS[1:], " ")

	if boardName[0] != '#' || len(boardName) <= 1 {
		return "", "", errors.New("invalid board name")
	}

	if len(description) <= 1 {
		return "", "", errors.New("invalid description")
	}

	return boardName, description, nil
}

func main() {
	fmt.Print(" Taskbook opened!\n\n")

	tb := newTaskbook()
	defer tb.display()

	taskPtr := flag.String("task", "", "the description of the new task to add preceded by the corresponding #board")
	notePtr := flag.String("note", "", "the description of the new note to add preceded by the corresponding #board")

	flag.Parse()

	if len(*taskPtr) > 0 {
		tb.addTask(*taskPtr, Pending)
	}

	if len(*notePtr) > 0 {
		tb.addTask(*notePtr, Note)
	}
}
