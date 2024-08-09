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
		commands.Register()
	case "delete":
		fmt.Println("Deleting user...")
		commands.Delete()
	case "login":
		fmt.Println("Logging in user...")
		commands.Login()
	case "logout":
		fmt.Println("Logging out user...")
		commands.Logout()
	case "join":
		// TODO
		rid := args[1]
		fmt.Printf("Joining room %s...\n", rid)
	case "leave":
		// TODO
		rid := args[1]
		fmt.Printf("Leaving room %s...\n", rid)
	case "enter":
		rid := args[1]
		fmt.Printf("Entering room %s...\n", rid)
		commands.Enter(rid)
	case "create":
		rid := args[1]
		fmt.Println("Creating room...")
		commands.Create(rid)
	case "destroy":
		rid := args[1]
		fmt.Printf("Destroying room %s...\n", rid)
		commands.Destroy(rid)
	case "list":
		// TODO
		fmt.Println("Listing rooms...")
	case "invite":
		// TODO
		rid, uid := args[1], args[2]
		fmt.Printf("Inviting user %s to room %s...\n", uid, rid)
	case "kick":
		// TODO
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
