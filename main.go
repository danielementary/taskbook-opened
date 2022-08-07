package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

const storageDirectory = ".taskbook-opened"
const storageFile = "storage.json"

func initStorageFile() string {
	userHomeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatalf("Failed to retrieve the user home directory: %s", err)
	}

	storageDirectorypath := filepath.Join(userHomeDir, storageDirectory)
	err = os.MkdirAll(storageDirectorypath, os.ModePerm)
	if err != nil {
		log.Fatalf("Failed to retrieve or create the storage directorypath: %s", err)
	}

	storageFilepath := filepath.Join(storageDirectorypath, storageFile)

	storageFile := openOrCreateStorageFile(storageFilepath)
	closeStorageFile(storageFile)

	return storageFilepath
}

func openOrCreateStorageFile(storageFilepath string) *os.File {
	storageFile, err := os.OpenFile(storageFilepath, os.O_CREATE, 0644)
	if err != nil {
		log.Fatalf("Failed to open or create the storage file: %s", err)
	}

	return storageFile
}

func closeStorageFile(storageFile *os.File) {
	err := storageFile.Close()
	if err != nil {
		log.Fatalf("Failed to close storage file: %s", err)
	}
}

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
	Tasks         []*task
	NumberOfTasks map[taskStatus]uint
}

func newBoard(name string) *board {
	tasks := []*task{}
	numberOfTasks := make(map[taskStatus]uint)

	numberOfTasks[InProgress] = uint(0)
	numberOfTasks[Pending] = uint(0)
	numberOfTasks[Completed] = uint(0)
	numberOfTasks[Note] = uint(0)

	return &board{name, tasks, numberOfTasks}
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
	Counter uint
	Boards  map[string]*board
}

func newTaskbook() *taskbook {
	counter := uint(1)
	boards := make(map[string]*board)

	return &taskbook{counter, boards}
}

func (tb *taskbook) display() {
	for _, board := range tb.Boards {
		board.display()
	}
}

func readFromFileOrCreate() (tb *taskbook) {
	tb, err := readFromFile()

	if err != nil {
		tb = newTaskbook()
	}

	return
}

func readFromStorageFile() (tb *taskbook, err error) {
	userHomeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	storageFilepath := filepath.Join(userHomeDir, storageDirectory, storageFile)

	tbJson, err := ioutil.ReadFile(storageFilepath)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(tbJson, &tb)

	return
}

func (tb *taskbook) saveToStorageFile(storageFile *os.File) {
	tbJson, err := json.MarshalIndent(tb, "", " ")
	if err != nil {
		fmt.Println(" Failed to save to file:", err)
		return
	}

	err = ioutil.WriteFile(storageFilepath, tbJson, 0644)
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

	id := tb.Counter
	tb.Counter++
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

	// storageFile := openStorageFile()
	// defer closeStorageFile(storageFile)

	// tb := readFromFileOrCreate()
	// defer tb.display()
	// defer tb.saveToFile()

	// taskPtr := flag.String("task", "", "the description of the new task to add preceded by the corresponding #board")
	// notePtr := flag.String("note", "", "the description of the new note to add preceded by the corresponding #board")

	// flag.Parse()

	// if len(*taskPtr) > 0 {
	// 	tb.addTask(*taskPtr, Pending)
	// }

	// if len(*notePtr) > 0 {
	// 	tb.addTask(*notePtr, Note)
	// }
}
