package main

import (
	"nobozo/config"
	"nobozo/handle"
)

func main() {
	connectToDatabase()
	initSystemObjects()
	config.LoadConfigFromDatabase()
}

func connectToDatabase() {
	handle.Connect()
}

func initSystemObjects() {
}
