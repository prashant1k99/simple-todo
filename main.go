package main

func main() {
	initDB()
	defer closeDB()

	testQuery()
}
