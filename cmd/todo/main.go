package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/araddon/dateparse"

	"github.com/Arshad-Aly/TODO-CLI/internal/storage"
	"github.com/Arshad-Aly/TODO-CLI/internal/task"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "todo",
	Short: "Todo is a task manager",
	Long:  `A Fast and Flexible Task Manager.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Welcome to your TODO list manager")
	},
}

var dueDate string

var addCmd = &cobra.Command {
	Use:   "add",
	Short: "Add a new task",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			fmt.Println("You must provide a task description")
			return
		}
		taskTitle := args[0]
		fmt.Printf("Added task: %s\n", taskTitle)

		var dueDateParsed *time.Time
		if dueDate != "" {
			parsed, err := dateparse.ParseAny(dueDate)
			if err != nil {
				fmt.Printf("Error parsing due date: %v\n", err)
				return
			}
			dueDateParsed = &parsed
		}
		task.Add(taskTitle, dueDateParsed)
		if dueDateParsed != nil {
			fmt.Printf("Due date set to: %s\n", dueDateParsed.Format("2006-01-02"))
		}
	},
}

var showAll bool

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all tasks",
	Run: func(cmd *cobra.Command, args []string) {

		task.List(showAll)
		// tasks := task.List(showAll)
		// 	if len(tasks) == 0 {
		// 		fmt.Println("No tasks to display.")
		// 		return
		// 	}
		// 	fmt.Println("Your tasks:")
		// 	for _, t := range tasks {
		// 		status := " "
		// 		if t.Done {
		// 			status = "✓"
		// 		}
		// 		dueString := "No due date"
		// 		if t.DueDate != nil {
		// 			dueString = t.DueDate.Format("2006-01-02")
		// 		}
		// 		fmt.Printf("[%s] %d: %s (Due: %s)\n", status, t.ID, t.Title, dueString)
		// 	}

	},
}

func updateCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "update [task ID]",
		Short: "update a task",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			id, err := strconv.Atoi(args[0])
			if err != nil {
				fmt.Println("Invalid task ID")
				return
			}

			task, err := storage.GetTask(id)
			if err != nil {
				fmt.Printf("Task with ID %d not found\n", id)
				return
			}

			fmt.Println("Current task:")
			fmt.Printf("Title: %s\nStatus: %v\n", task.Title, task.Done)

			var choice string
			fmt.Print("Do you want to update the title (t) or status (s)?: ")
			fmt.Scanln(&choice)

			switch choice {
			case "t":
				reader := bufio.NewReader(os.Stdin)
				fmt.Print("Enter new Title: ")
				newTitle, _ := reader.ReadString('\n')
				newTitle = strings.TrimSpace(newTitle)
				task.Title = newTitle
			case "s":
				task.Done = !task.Done
				fmt.Printf("Status updated to: %v\n", task.Done)
			default:
				fmt.Println("Invalid choice. No updates made.")
				return
			}

			err = storage.UpdateTask(task)
			if err != nil {
				fmt.Println("Error updating task: ", err)
				return
			}

			fmt.Println("Task update successfully!")

		},
	}
}

func completeCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "complete",
		Short: "Mark a task as Complete",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) < 1 {
				fmt.Println("You must provide a task ID")
				return
			}
			taskId := args[0]
			fmt.Printf("Marked task %s as complete\n", taskId)
			ID, err := strconv.Atoi(taskId)
			if err != nil {
				panic(err)
			}
			err = task.Complete(ID)
			if err != nil {
				fmt.Println(err)
			}
		},
	}
}

func deleteCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "delete [task ID]",
		Short: "Delete a task",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			id, err := strconv.Atoi(args[0])
			if err != nil {
				fmt.Println("Invalid task ID")
				return
			}

			err = storage.DeleteTask(id)
			if err != nil {
				fmt.Printf("Error deleting task: %v\n", err)
				return
			}

			fmt.Printf("Task with ID %d has been deleted successfully!\n", id)
		},
	}

}


func Execute() {

	rootCmd.AddCommand(addCmd)
	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(completeCmd())
	rootCmd.AddCommand(updateCmd())
	rootCmd.AddCommand(deleteCmd())

	// addCmd.Flags().StringVarP(&dueDate, "due", "d", "", "Due date for the task")

	listCmd.Flags().BoolVar(&showAll, "all", false, "Show all tasks, including completed ones.")

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func main() {
	Execute()
}
