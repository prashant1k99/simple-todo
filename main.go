/*
Copyright © 2024 Prashant Singh <prashantco111+github@gmail.com>
*/
package main

func main() {
	initDB()
	defer closeDB()

	rootCmd.Execute()
}
