package main

import (
	"encoding/csv"
	"os"
	"path/filepath"
	"strconv"
	"time"

	gap "github.com/muesli/go-app-paths"
)

type Todo struct {
	ID          int
	Name        string
	Description string
	IsClosed    bool
	CreatedAt   string
}

var dbFilePath string

func initDB() {
	// Initialize the database path
	scope := gap.NewScope(gap.User, "simple-todo")

	// Get the data directory for the application
	dbDir, err := scope.DataDirs()
	if err != nil {
		panic(err)
	}

	// Create the database path
	dbFilePath = filepath.Join(dbDir[0], "todos.csv")

	// Check if the directory exists, otherwise create directory
	if _, err := os.Stat(filepath.Dir(dbFilePath)); os.IsNotExist(err) {
		if err := os.MkdirAll(filepath.Dir(dbFilePath), 0755); err != nil {
			panic(err)
		}
	}

	// Create the CSV file if it doesn't exist
	if _, err := os.Stat(dbFilePath); os.IsNotExist(err) {
		file, err := os.Create(dbFilePath)
		if err != nil {
			panic(err)
		}
		defer file.Close()

		writer := csv.NewWriter(file)
		defer writer.Flush()

		// Write the header
		writer.Write([]string{"ID", "Name", "Description", "IsClosed", "CreatedAt"})
	}
}

func readTodos() ([]Todo, error) {
	file, err := os.Open(dbFilePath)
	if err != nil {
		return nil, err
	}

	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	var todos []Todo

	for i, record := range records {
		if i == 0 {
			continue
		}
		id, _ := strconv.Atoi(record[0])
		isClosed, _ := strconv.ParseBool(record[3])
		todos = append(todos, Todo{
			ID:          id,
			Name:        record[1],
			Description: record[2],
			IsClosed:    isClosed,
			CreatedAt:   record[4],
		})
	}

	return todos, nil
}

func getTodoById(id int) (Todo, error) {
	todos, err := readTodos()
	if err != nil {
		return Todo{}, err
	}

	for _, todo := range todos {
		if todo.ID == id {
			return todo, nil
		}
	}

	return Todo{}, nil
}

func writeTodos(todos []Todo) error {
	file, err := os.Create(dbFilePath)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write the header
	writer.Write([]string{"ID", "Name", "Description", "IsClosed", "CreatedAt"})

	for _, todo := range todos {
		writer.Write([]string{
			strconv.Itoa(todo.ID),
			todo.Name,
			todo.Description,
			strconv.FormatBool(todo.IsClosed),
			todo.CreatedAt,
		})
	}

	return nil
}

func addTodo(name, description string) error {
	todos, err := readTodos()
	if err != nil {
		return err
	}

	id := 1
	if len(todos) > 0 {
		id = todos[len(todos)-1].ID + 1
	}

	todo := Todo{
		ID:          id,
		Name:        name,
		Description: description,
		IsClosed:    false,
		CreatedAt:   time.Now().Format(time.RFC3339),
	}

	todos = append(todos, todo)
	return writeTodos(todos)
}

func deleteTodo(id int) error {
	todos, err := readTodos()
	if err != nil {
		return err
	}

	var updatedTodos []Todo
	for _, todo := range todos {
		if todo.ID != id {
			updatedTodos = append(updatedTodos, todo)
		}
	}

	return writeTodos(updatedTodos)
}

func updateTodo(id int, name, description string) error {
	todos, err := readTodos()
	if err != nil {
		return err
	}

	for i, todo := range todos {
		if todo.ID == id {
			todos[i].Name = name
			todos[i].Description = description
			break
		}
	}

	return writeTodos(todos)
}

func updateTodoStatus(id int, isClosed bool) error {
	todos, err := readTodos()
	if err != nil {
		return err
	}

	for i, todo := range todos {
		if todo.ID == id {
			todos[i].IsClosed = isClosed
			break
		}
	}

	return writeTodos(todos)
}
