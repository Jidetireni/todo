package task

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"text/tabwriter"
	"time"
)

type Manager interface {
	ReadTasksToFIle(paste interface{}) error
	WriteTaskToFile(filedata interface{}) error
}

type task struct {
	id         int       `json:"id"`
	actions    string    `json:"task"`
	status     bool      `json:"status"`
	createdAt  time.Time `json:"created"`
	finishedAt time.Time `json:"finished"`
}

type todo struct {
	tasks map[int]task
	// date  time.Time
	Io Manager
}

func (t *todo) Add(action string) (int, error) {

	newTasks := task{
		id:        len(t.tasks) + 1,
		actions:   action,
		status:    false,
		createdAt: time.Now(),
	}

	t.tasks[newTasks.id] = newTasks
	fmt.Println("Tasks before writing to file:", t.tasks)

	err := t.Io.WriteTaskToFile(t.tasks)
	if err != nil {
		return 0, err
	}

	return newTasks.id, nil
}

func (t *todo) List() {

	writer := tabwriter.NewWriter(os.Stdout, 0, 8, 2, '\t', 0)
	fmt.Fprintln(writer, "ID\tTask\tStatus\tCreated\tFinished")

	for id, task := range t.tasks {
		finishedTime := "N/A"
		if task.status {
			finishedTime = task.finishedAt.Format("2006-01-02 15:04:05")
		}
		fmt.Fprintf(writer, "%d\t%s\t%s\t%s\t%s\n",
			id,
			task.actions, map[bool]string{true: "Done", false: "Pending"}[task.status],
			task.createdAt.Format("2006-01-02 15:04:05"),
			finishedTime,
		)
	}
	writer.Flush()
}

func (t *todo) Done(id int) (bool, error) {

	task, exists := t.tasks[id]
	if !exists {
		return false, fmt.Errorf("task with id %d not found", id)
	}
	task.status = true
	task.finishedAt = time.Now()

	t.tasks[id] = task
	err := t.Io.WriteTaskToFile(t.tasks)
	if err != nil {
		return false, err
	}

	return false, nil
}

func (t *todo) Delete(id int) error {

	if task, exists := t.tasks[id]; exists {

		reader := bufio.NewScanner(os.Stdin)

		fmt.Printf("Are you sure you want to delete task '%v'? \n(yes/no) \n", task.actions)

		if reader.Scan() {
			choice := strings.TrimSpace(strings.ToLower(reader.Text()))

			switch choice {
			case "yes", "y":
				delete(t.tasks, id)
				fmt.Println("Task deleted successfully.")
				return nil
			case "no", "n":
				fmt.Println("Task not deleted.")
				return nil
			default:
				return fmt.Errorf("invalid response: %s. Task not deleted", choice)
			}
		}

	}

	return fmt.Errorf("task with id %d not found", id)
}

func (t *todo) GetTasks() map[int]task {
	return t.tasks
}

func (t *todo) LoadTasks() error {
	return t.Io.ReadTasksToFIle(&t.tasks)
}

func (t *todo) PrintTasksFromFile() error {
	err := t.LoadTasks()
	if err != nil {
		return fmt.Errorf("error loading tasks from file: %v", err)
	}
	t.List()
	return nil
}

func New(io Manager) *todo {
	todo := &todo{
		tasks: make(map[int]task),
		Io:    io,
	}

	if err := todo.LoadTasks(); err != nil {
		fmt.Println("Error loading tasks:", err)
	}

	return todo
}