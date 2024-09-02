/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"fmt"
	"os"

	"github.com/prashant1k99/simple-todo/form"
	"github.com/prashant1k99/simple-todo/list"
	"github.com/prashant1k99/simple-todo/table"
	"github.com/spf13/cobra"
)

var createTodoCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new todo",
	Run: func(cmd *cobra.Command, args []string) {
		name, err := cmd.Flags().GetString("name")
		if err != nil {
			fmt.Println(err)
		}
		description, err := cmd.Flags().GetString("description")
		if err != nil {
			fmt.Println(err)
		}

		if name == "" && description == "" {
			msg, err := form.RenderCreateForm(&form.SubmissionMsg{
				Name:        "",
				Description: "",
			})
			if err != nil {
				fmt.Printf("Error: %v\n", err)
				return
			}
			if !msg.Submitted {
				return
			}
			if msg.Name == "" {
				fmt.Println("Name is required to create a TODO")
				return
			}
			createTODO(&Todo{Name: msg.Name, Description: msg.Description})
		} else {
			createTODO(&Todo{Name: name, Description: description})
		}
	},
}

func createTODO(todo *Todo) {
	if todo.Name == "" {
		fmt.Println("Name is required to create a TODO")
		return
	}
	err := addTodo(todo.Name, todo.Description)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("Created todo")
}

var listTodoCmd = &cobra.Command{
	Use:   "list",
	Short: "List all todos",
	Run: func(cmd *cobra.Command, args []string) {
		listTODO()
	},
}

func getAllTODOs() ([]Todo, error) {

	var todos []Todo
	todos, err := readTodos()
	if err != nil {
		return nil, err
	}
	return todos, nil
}

func listTODO() {
	todos, err := getAllTODOs()
	if err != nil {
		fmt.Println(err)
		return
	}
	if len(todos) == 0 {
		fmt.Println("No TODOs found")
		return
	}
	tableTodos := []table.Item{}
	for _, todo := range todos {
		tableTodos = append(tableTodos, table.Item{
			Name:     todo.Name,
			Desc:     todo.Description,
			IsClosed: todo.IsClosed,
		})
	}
	table.RenderTable(tableTodos)
}

var deleteTodoCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete a todo",
	Run:   deleteTODO,
}

func deleteTODO(cmd *cobra.Command, args []string) {
	id, err := cmd.Flags().GetInt("id")
	if err != nil {
		fmt.Println(err)
		return
	}

	if id == 0 {
		id, err = selectTODO()
		if err != nil {
			fmt.Println(err)
			return
		}
	}

	err = deleteTodo(id)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("Deleted todo with ID: %d\n", id)
}

func selectTODO() (int, error) {
	todos, err := getAllTODOs()
	if err != nil {
		return 0, err
	}
	if len(todos) == 0 {
		return 0, fmt.Errorf("No TODOs found")
	}
	todoItems := make([]list.Item, 0)
	for _, todo := range todos {
		todoItems = append(todoItems, list.Item{
			ID:   todo.ID,
			Name: todo.Name,
			Desc: todo.Description,
		})
	}
	selectionResponse := list.RenderListItem(todoItems)
	if selectionResponse.Err != nil {
		return 0, selectionResponse.Err
	}
	if !selectionResponse.Selected {
		return 0, fmt.Errorf("No TODO selected")
	}
	return selectionResponse.Item.ID, nil
}

var setToDoStatusCmd = &cobra.Command{
	Use:   "set-status",
	Short: "Set the status of a todo",
	Run:   setToDoStatus,
}

func setToDoStatus(cmd *cobra.Command, args []string) {
	id, err := cmd.Flags().GetInt("id")

	if err != nil {
		fmt.Println(err)
		return
	}

	if id == 0 {
		id, err = selectTODO()
		if err != nil {
			fmt.Println(err)
			return
		}
	}

	isClosed, err := cmd.Flags().GetBool("is-closed")
	if err != nil {
		fmt.Println(err)
		return
	}
	if !cmd.Flags().Changed("is-closed") {
		selectionResponse := list.RenderListItem([]list.Item{
			{
				ID:   1,
				Name: "Close ToDo",
				Desc: "This will close the selected ToDo",
			},
			{
				ID:   2,
				Name: "Reopen/Keep ToDo",
				Desc: "This will reopen or keep the selected ToDo open",
			},
		})
		if selectionResponse.Err != nil {
			fmt.Println(selectionResponse.Err)
			return
		}
		if !selectionResponse.Selected {
			return
		}
		isClosed = selectionResponse.Item.ID == 1
	}

	err = updateTodoStatus(id, isClosed)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("Updated todo with ID: %d\n", id)
}

var updateTodoCmd = &cobra.Command{
	Use:   "update",
	Short: "Update a todo",
	Run:   updateTODO,
}

func updateTODO(cmd *cobra.Command, args []string) {
	id, err := cmd.Flags().GetInt("id")
	if err != nil {
		fmt.Println(err)
		return
	}
	if id == 0 {
		id, err = selectTODO()
		if err != nil {
			fmt.Println(err)
			return
		}
	}

	var todo Todo
	todo, err = getTodoById(id)
	if err != nil {
		fmt.Println(err)
		return
	}

	name, err := cmd.Flags().GetString("name")
	if err != nil {
		fmt.Println(err)
		return
	}

	description, err := cmd.Flags().GetString("description")
	if err != nil {
		fmt.Println(err)
		return
	}

	if name == "" && description == "" {
		msg, err := form.RenderCreateForm(&form.SubmissionMsg{
			Name:        todo.Name,
			Description: todo.Description,
		})
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}
		if !msg.Submitted {
			return
		}
		if msg.Name == "" {
			fmt.Println("Name is required to create a TODO")
			return
		}
		todo.Name = msg.Name
		todo.Description = msg.Description
	}

	err = updateTodo(id, todo.Name, todo.Description)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("Updated todo with ID: %d\n", id)
}

var rootCmd = &cobra.Command{
	Use:   "simple-todo",
	Short: "A CLI based simple todo application",
	Long:  `A CLI based simple todo application that allows you to add, list, and remove tasks.`,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	createTodoCmd.Flags().StringP("name", "n", "", "Name of the todo")
	createTodoCmd.Flags().StringP("description", "d", "", "Description of the todo")

	rootCmd.AddCommand(createTodoCmd)

	rootCmd.AddCommand(listTodoCmd)

	deleteTodoCmd.Flags().IntP("id", "i", 0, "ID of the todo")
	rootCmd.AddCommand(deleteTodoCmd)

	setToDoStatusCmd.Flags().IntP("id", "i", 0, "ID of the todo")
	setToDoStatusCmd.Flags().Bool("is-closed", true, "Status of the todo")
	rootCmd.AddCommand(setToDoStatusCmd)

	updateTodoCmd.Flags().IntP("id", "i", 0, "ID of the todo")
	updateTodoCmd.Flags().StringP("name", "n", "", "Name of the todo")
	updateTodoCmd.Flags().StringP("description", "d", "", "Description of the todo")
	rootCmd.AddCommand(updateTodoCmd)
}
