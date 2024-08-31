/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

type ToDo struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	IsClosed    bool   `json:"isClosed"`
	CreatedAt   string `json:"created_at"`
}

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

		createTODO(&ToDo{Name: name, Description: description})
	},
}

func createTODO(todo *ToDo) {
	if todo.Name == "" {
		fmt.Println("Name is required to create a TODO")
		return
	}
	fmt.Println("Creating todo", todo.Name)
	result, err := dbQueries.Exec("INSERT INTO todos (name, description) VALUES (?, ?)", todo.Name, todo.Description)
	if err != nil {
		fmt.Println(err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("Created todo with ID: %d\n", id)
}

var listTodoCmd = &cobra.Command{
	Use:   "list",
	Short: "List all todos",
	Run: func(cmd *cobra.Command, args []string) {
		listTODO()
	},
}

func listTODO() {
	rows, err := dbQueries.Query("SELECT * FROM todos")
	if err != nil {
		fmt.Println(err)
	}
	defer rows.Close()

	fmt.Println(rows.Columns())

	var todos []ToDo
	for rows.Next() {
		var todo ToDo
		err := rows.Scan(&todo.ID, &todo.Name, &todo.Description, &todo.IsClosed, &todo.CreatedAt)
		if err != nil {
			fmt.Println(err)
			return
		}
		todos = append(todos, todo)
	}

	for _, todo := range todos {
		fmt.Printf("ID: %d, Name: %s, Description: %s, IsClosed: %t, CreatedAt: %s\n", todo.ID, todo.Name, todo.Description, todo.IsClosed, todo.CreatedAt)
	}
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

	_, err = dbQueries.Exec("DELETE FROM todos WHERE id = ?", id)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("Deleted todo with ID: %d\n", id)
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

	isClosed, err := cmd.Flags().GetBool("is-closed")
	if err != nil {
		fmt.Println(err)
		return
	}

	_, err = dbQueries.Exec("UPDATE todos SET is_closed = ? WHERE id = ?", isClosed, id)
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

	var todo ToDo
	err = dbQueries.QueryRow("SELECT id, name, description, is_closed, created_at FROM todos WHERE id = ?", id).Scan(&todo.ID, &todo.Name, &todo.Description, &todo.IsClosed, &todo.CreatedAt)
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

	if name != "" {
		todo.Name = name
	}
	if description != "" {
		todo.Description = description
	}

	_, err = dbQueries.Exec("UPDATE todos SET name = ?, description = ? WHERE id = ?", todo.Name, todo.Description, todo.ID)
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
	setToDoStatusCmd.Flags().BoolP("is-closed", "c", false, "Status of the todo")
	rootCmd.AddCommand(setToDoStatusCmd)

	updateTodoCmd.Flags().IntP("id", "i", 0, "ID of the todo")
	updateTodoCmd.Flags().StringP("name", "n", "", "Name of the todo")
	updateTodoCmd.Flags().StringP("description", "d", "", "Description of the todo")
	rootCmd.AddCommand(updateTodoCmd)
}
