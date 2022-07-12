package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

const storageDirectory = ".taskbook-opened"
const storageFile = "storage.json"

type taskStatus int

const (
	InProgress taskStatus = iota
	Pending
	Completed
	Note
)

type task struct {
	Id          uint
	Description string
	Status      taskStatus
}

func (t *task) display() {
	fmt.Printf("    %d. ", t.Id)

	switch t.Status {
	case InProgress:
		fmt.Print("[/]")
	case Pending:
		fmt.Print("[ ]")
	case Completed:
		fmt.Print("[X]")
	case Note:
		fmt.Print("-N-")
	}

	fmt.Printf(" %s\n", t.Description)
}

type board struct {
	Name          string
	Counter       uint
	Tasks         []*task
	NumberOfTasks map[taskStatus]uint
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
	totalNumberOfTasks := b.NumberOfTasks[InProgress] + b.NumberOfTasks[Pending] + b.NumberOfTasks[Completed]

	fmt.Printf("  %s [%d/%d]\n", b.Name, b.NumberOfTasks[Completed], totalNumberOfTasks)

	if totalNumberOfTasks <= 0 {
		fmt.Print("    This board is empty.")
	} else {
		for _, t := range b.Tasks {
			t.display()
		}
	}
	fmt.Print("\n")
}

type taskbook struct {
	Boards map[string]*board
}

func newTaskbook() *taskbook {
	boards := make(map[string]*board)

	return &taskbook{boards}
}

func (tb *taskbook) display() {
	for _, board := range tb.Boards {
		board.display()
	}
}

func (tb *taskbook) saveToFile() {
	json, err := json.MarshalIndent(tb, "", " ")
	if err != nil {
		fmt.Println(" Failed to save to file:", err)
		return
	}

	userHomeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Println(" Failed to save to file:", err)
		return
	}

	storageDirectorypath := filepath.Join(userHomeDir, storageDirectory)
	err = os.MkdirAll(storageDirectorypath, os.ModePerm)
	if err != nil {
		fmt.Println(" Failed to save to file:", err)
		return
	}

	storageFilepath := filepath.Join(userHomeDir, storageDirectory, storageFile)
	err = ioutil.WriteFile(storageFilepath, json, 0644)
	if err != nil {
		fmt.Println(" Failed to save to file:", err)
		return
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
	board, found := tb.Boards[boardName]

	if !found {
		tb.Boards[boardName] = newBoard(boardName)
		board = tb.Boards[boardName]
	}

	id := board.Counter
	board.Counter++
	description := taskDescription
	status := taskStatus

	board.Tasks = append(board.Tasks, &task{id, description, status})
	board.NumberOfTasks[status]++
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
	defer tb.saveToFile()

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
