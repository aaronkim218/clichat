package commands

import "fmt"

// print help message
func Help() {
	fmt.Println("Usage: clichat <command> [arguments]")
	fmt.Println("Commands:")
	fmt.Println("\tregister")
	fmt.Println("\tunregister")
	fmt.Println("\tlogin")
	fmt.Println("\tlogout")
	fmt.Println("\tjoin <room_id>")
	fmt.Println("\tleave <room_id>")
	fmt.Println("\tenter <room_id>")
	fmt.Println("\tcreate")
	fmt.Println("\tdelete <room_id>")
	fmt.Println("\tlist")
	fmt.Println("\tinvite <room_id> <user_id>")
	fmt.Println("\tkick <room_id> <user_id>")
	fmt.Println("\thelp")
	fmt.Println("\tversion")
}
