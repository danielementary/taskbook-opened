package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

const storageDirectory = ".taskbook-opened"
const storageFile = "storage.json"
const storageFilePermission = 0644

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
	storageFile, err := os.OpenFile(storageFilepath, os.O_CREATE, storageFilePermission)
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
	Pending taskStatus = iota
	InProgress
	Completed
	Note
)

func highestPriorityTaskStatus(a, b taskStatus) taskStatus {
	if a >= b {
		return a
	} else {
		return b
	}
}

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
	Tag           string
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

	fmt.Printf("  %s [%d/%d]\n", b.Tag, b.NumberOfTasks[Completed], totalNumberOfTasks)

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
	Tasks   []*task
}

func newTaskbook() *taskbook {
	counter := uint(1)
	boards := make(map[string]*board)
	tasks := make([]*task, 0)

	return &taskbook{counter, boards, tasks}
}

func (tb *taskbook) newTask(description string, status taskStatus) *task {
	id := tb.Counter
	tb.Counter++

	task := &task{id, description, status}

	tb.Tasks = append(tb.Tasks, task)

	return task
}

func (tb *taskbook) display() {
	for _, board := range tb.Boards {
		board.display()
	}
}

func readTaskbookFromStorageFile(storageFilepath string) (tb *taskbook) {
	tbJson, err := os.ReadFile(storageFilepath)
	if err != nil {
		log.Fatalf("Failed to read taskbook from storage file: %s", err)
	}

	err = json.Unmarshal(tbJson, &tb)
	if err != nil {
		fmt.Println("Instantiating a new taskbook...")
		tb = newTaskbook()
	}

	return
}

func (tb *taskbook) saveToStorageFile(storageFilepath string) {
	tbJson, err := json.MarshalIndent(tb, "", " ")
	if err != nil {
		log.Fatalf("Failed to marshal taskbook: %s", err)
	}

	err = os.WriteFile(storageFilepath, tbJson, storageFilePermission)
	if err != nil {
		fmt.Println(" Failed to save to file:", err)
		return
	}
}

// func (tb *taskbook) addTaskToBoard(boardName, taskDescription string, taskStatus taskStatus) {
// 	board, found := tb.Boards[boardName]

// 	if !found {
// 		tb.Boards[boardName] = newBoard(boardName)
// 		board = tb.Boards[boardName]
// 	}

// 	id := tb.Counter
// 	tb.Counter++
// 	description := taskDescription
// 	status := taskStatus

// 	board.Tasks = append(board.Tasks, &task{id, description, status})
// 	board.NumberOfTasks[status]++
// }

func parseBoardTags(taskTextSlice []string) (boardTags []string, taskStatus taskStatus) {
	boardTags = make([]string, 0)
	taskStatus = Pending

	for _, s := range taskTextSlice {
		if s[0] == '#' {
			switch s {
			case "#pending":
				taskStatus = highestPriorityTaskStatus(Pending, taskStatus)
			case "#inProgress":
				taskStatus = highestPriorityTaskStatus(InProgress, taskStatus)
			case "#completed":
				taskStatus = highestPriorityTaskStatus(Completed, taskStatus)
			case "#note":
				taskStatus = highestPriorityTaskStatus(Note, taskStatus)
			default:
				boardTags = append(boardTags, s)
			}
		}
	}

	return boardTags, taskStatus
}

func main() {
	fmt.Print(" Taskbook opened!\n\n")

	storageFilepath := initStorageFile()

	tb := readTaskbookFromStorageFile(storageFilepath)
	defer tb.saveToStorageFile(storageFilepath)
	defer tb.display()

	if len(os.Args) <= 1 {
		return
	}

	taskText := os.Args[1:]
	boardTags, _ := parseBoardTags(taskText)

	fmt.Printf("DEBUG task text:  %s\n", taskText)
	fmt.Printf("DEBUG board tags: %s\n", boardTags)
}
