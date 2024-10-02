package task

import (
	"fmt"
	"os"
	"text/tabwriter"
	"time"
	"github.com/dustin/go-humanize"

	"github.com/Arshad-Aly/TODO-CLI/internal/model"
	"github.com/Arshad-Aly/TODO-CLI/internal/storage"
)

// type Task struct {
// 	ID        int
// 	Title     string
// 	Done      bool
// 	CreatedAt time.Time
// }

var tasks []model.Task
var lastId int

func init() {
	loadedTasks, err := storage.LoadTasks()
	if err != nil {
		fmt.Printf("Error loading tasks: %v\n", err)
		return
	}

	tasks = loadedTasks
	for _, task := range tasks {
		if task.ID > lastId {
			lastId = task.ID
		}
	}

}

func Add(title string, dueDate *time.Time) int {
	lastId++

	task := model.Task{
		ID:        lastId,
		Title:     title,
		Done:      false,
		CreatedAt: time.Now(),
		DueDate: dueDate,
	}
	
	tasks = append(tasks, task)
	saveTasks()
	return lastId
}

func List(flag bool) []model.Task {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)

	fmt.Fprintln(w, "ID\tTitle\tDone\tCreated At")
	fmt.Fprintln(w, "--\t----\t----\t----------")
	
	// fmt.Println(flag)
	for _, task := range tasks {
		if !task.Done && !flag {
			fmt.Fprintf(w, "%d\t%s\t%v\t%s\n",
				task.ID,
				task.Title,
				task.Done,
				humanize.Time(task.CreatedAt),
			)
		} 
		if flag {
			fmt.Fprintf(w, "%d\t%s\t%v\t%s\n",
				task.ID,
				task.Title,
				task.Done,
				humanize.Time(task.CreatedAt),
			)
		}

	}

	w.Flush()
	return tasks
}


func Complete(id int) error {
	for i, task := range tasks {
		if task.ID == id {
			tasks[i].Done = true
			saveTasks()
			return nil
		}
	}
	return fmt.Errorf("task with ID %d not found", id)
}

func saveTasks() {
	err := storage.SaveTasks(tasks)
	if err != nil {
		fmt.Printf("Error saving tasks: %v\n", err)
	}
}
