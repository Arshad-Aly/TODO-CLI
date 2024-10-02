package storage

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/Arshad-Aly/TODO-CLI/internal/model"
)

const filename = "tasks.csv"

func SaveTasks(tasks []model.Task) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}

	defer file.Close()

	writer := csv.NewWriter(file)

	defer writer.Flush()

	for _, task := range tasks {
		dueDate := ""
		if task.DueDate != nil {
			dueDate = task.DueDate.Format(time.RFC3339)
		}


		row := []string{
			strconv.Itoa(task.ID),
			task.Title,
			strconv.FormatBool(task.Done),
			task.CreatedAt.Format(time.RFC3339),
			dueDate,
		}
		if err := writer.Write(row); err != nil {
			return err
		}
	}
	return nil
}

func GetTask(id int) (model.Task, error) {
	tasks, err := LoadTasks()
	if err != nil {
		return model.Task{}, err
	}

	for _, task := range tasks {
		if task.ID == id {
			return task, nil
		}
	}
	return model.Task{}, fmt.Errorf("task with ID %d not found", id)
}

func UpdateTask(updateTask model.Task) error {
	tasks, err := LoadTasks()
	if err != nil {
		return err
	}

	found := false
	for i, task := range tasks {
		if task.ID == updateTask.ID {
			tasks[i] = updateTask
			found = true
			break
		}
	}

	if !found {
		return fmt.Errorf("task with ID %d not found", updateTask.ID)
	}

	return SaveTasks(tasks)
}

func DeleteTask(taskID int) error {
	tasks, err := LoadTasks()
	if err != nil {
		return err
	}

	index := -1
	for i, task := range tasks {
		if task.ID == taskID {
			index = i
			break
		}
	}

	if index == -1 {
		return fmt.Errorf("task with ID %d not found", taskID)
	}

	tasks = append(tasks[:index], tasks[index+1:]...)

	return SaveTasks(tasks)
}

func LoadTasks() ([]model.Task, error) {
	file, err := os.Open(filename)
	if err != nil {
		if os.IsNotExist(err) {
			return []model.Task{}, nil
		}
		return nil, err
	}

	defer file.Close()

	reader := csv.NewReader(file)
	rows, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	var tasks []model.Task

	for _, row := range rows {
		id, _ := strconv.Atoi(row[0])
		done, _ := strconv.ParseBool(row[2])
		createdAt, _ := time.Parse(time.RFC3339, row[3])

		
		task := model.Task{
			ID:        id,
			Title:     row[1],
			Done:      done,
			CreatedAt: createdAt,
		}
		tasks = append(tasks, task)
	}

	return tasks, nil

}
