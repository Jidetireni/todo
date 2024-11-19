package main

import (
	"fmt"

	"github.com/Jidetireni/todo/filemanger"
	"github.com/Jidetireni/todo/task"
)

func main() {
	fm := filemanger.New("todolist.json", "todolist.json")
	todoList := task.New(fm)

	id1, err := todoList.Add("Learn Golang")
	if err != nil {
		fmt.Println("Error adding task:", err)
		return
	}
	id2, err := todoList.Add("Build a Go Project")
	if err != nil {
		fmt.Println("Error adding task:", err)
		return
	}

	fmt.Println("Tasks after adding:")
	todoList.List()

	_, err = todoList.Done(id1)
	if err != nil {
		fmt.Println("Error marking task as done:", err)
		return
	}
	fmt.Println("\nTasks after marking the first one as done:")
	todoList.List()
	err = todoList.Delete(id2)

	if err != nil {
		fmt.Println("Error deleting task:", err)
		return
	}

	fmt.Println("\nTasks after deleting the second one:")
	todoList.List()

	fmt.Println("Tasks loaded from file:")
	err = todoList.PrintTasksFromFile()
	if err != nil {
		fmt.Println("Error printing tasks from file:", err)
	}
}
