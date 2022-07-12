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
	name          string
	counter       uint
	tasks         []*task
	numberOfTasks map[taskStatus]uint
}

func newBoard(name string) *board {
	counter := uint(1)
	tasks := []*task{}
	numberOfTasks := make(map[taskStatus]uint)

	numberOfTasks[InProgress] = uint(0)
	numberOfTasks[Pending] = uint(0)
	numberOfTasks[Completed] = uint(0)
	numberOfTasks[Note] = uint(0)

	return &board{name, counter, tasks, numberOfTasks}
}

func (b *board) display() {
	totalNumberOfTasks := b.numberOfTasks[InProgress] + b.numberOfTasks[Pending] + b.numberOfTasks[Completed]

	fmt.Printf("  %s [%d/%d]\n", b.name, b.numberOfTasks[Completed], totalNumberOfTasks)

	if totalNumberOfTasks <= 0 {
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
	board.numberOfTasks[status]++
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

	tb.addTask("#coding work on taskbook-opened", Pending)
	tb.addTask("#coding do not forget about testing...", Note)

	tb.addTask("#chill go to bed", Pending)

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
