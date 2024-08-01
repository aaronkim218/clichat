package main

import (
	"cli/commands"
	"fmt"
	"os"
)

func main() {
	args := os.Args[1:]

	switch args[0] {
	case "register":
		fmt.Println("Registering user...")
	case "unregister":
		fmt.Println("Unregistering user...")
	case "login":
		fmt.Println("Logging in user...")
	case "logout":
		fmt.Println("Logging out user...")
	case "join":
		rid := args[1]
		fmt.Printf("Joining room %s...\n", rid)
	case "leave":
		rid := args[1]
		fmt.Printf("Leaving room %s...\n", rid)
	case "enter":
		rid := args[1]
		fmt.Printf("Entering room %s...\n", rid)
		commands.Enter(rid)
	case "create":
		fmt.Println("Creating room...")
	case "delete":
		rid := args[1]
		fmt.Printf("Deleting room %s...\n", rid)
	case "list":
		fmt.Println("Listing rooms...")
	case "invite":
		rid, uid := args[1], args[2]
		fmt.Printf("Inviting user %s to room %s...\n", uid, rid)
	case "kick":
		rid, uid := args[1], args[2]
		fmt.Printf("Kicking user %s from room %s...\n", uid, rid)
	case "help":
		commands.Help()
	case "version":
		commands.Version()
	default:
		fmt.Println("Invalid command")
		commands.Help()
		os.Exit(1)
	}
}
