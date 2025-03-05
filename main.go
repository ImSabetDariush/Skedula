package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3" // SQLite driver
)

var db *sql.DB

func main() {
	// Connect to the SQLite database (creates it if it doesn't exist)
	var err error
	db, err = sql.Open("sqlite3", "./tasks.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Create the tasks table if it doesn't exist
	createTable()

	// Show the main menu
	mainMenu()
}

// Create the tasks table in the database
func createTable() {
	createTableSQL := `
	CREATE TABLE IF NOT EXISTS tasks (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		description TEXT NOT NULL
	);`
	_, err := db.Exec(createTableSQL)
	if err != nil {
		log.Fatal(err)
	}
}

// Display the main menu
func mainMenu() {
	for {
		fmt.Println("\nTimeWise - Time and Study Manager")
		fmt.Println("1. View Tasks")
		fmt.Println("2. Add Task")
		fmt.Println("3. Remove Task")
		fmt.Println("4. Exit")
		fmt.Print("Please choose an option: ")

		var choice string
		fmt.Scanln(&choice)

		switch choice {
		case "1":
			viewTasks()
		case "2":
			addTask()
		case "3":
			removeTask()
		case "4":
			fmt.Println("Exiting the program...")
			os.Exit(0)
		default:
			fmt.Println("Invalid choice! Please enter a number between 1 and 4.")
		}
	}
}

// View all tasks
func viewTasks() {
	// Query all tasks from the database
	rows, err := db.Query("SELECT id, description FROM tasks")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	fmt.Println("\nYour Tasks:")
	// Loop through the rows and print each task
	for rows.Next() {
		var id int
		var description string
		err = rows.Scan(&id, &description)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%d. %s\n", id, description)
	}
}

// Add a new task
func addTask() {
	var task string
	fmt.Print("Enter a new task: ")
	fmt.Scanln(&task)

	// Check if the task is empty
	if task == "" {
		fmt.Println("Task cannot be empty!")
		return
	}

	// Insert the new task into the database
	insertSQL := `INSERT INTO tasks (description) VALUES (?);`
	_, err := db.Exec(insertSQL, task)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Task added successfully!")
}

// Remove a task by ID
func removeTask() {
	// Show all tasks to the user
	viewTasks()

	var id int
	fmt.Print("Enter the ID of the task to remove: ")
	_, err := fmt.Scanln(&id)
	if err != nil {
		fmt.Println("Invalid input! Please enter a valid ID.")
		return
	}

	// Delete the task with the given ID
	deleteSQL := `DELETE FROM tasks WHERE id = ?;`
	result, err := db.Exec(deleteSQL, id)
	if err != nil {
		log.Fatal(err)
	}

	// Check if any rows were affected
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		fmt.Println("No task found with the given ID.")
	} else {
		fmt.Println("Task removed successfully!")
	}
}
