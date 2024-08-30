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
		fmt.Printf("ID: %d, Name: %s, Description: %s, IsComplete: %t, CreatedAt: %s\n", todo.ID, todo.Name, todo.Description, todo.IsClosed, todo.CreatedAt)
	}
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
}
