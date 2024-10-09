package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		PrintHelp()
		return
	}
	subCommand := os.Args[1]
	switch subCommand {
	case "update-list":
		UpdateList()
	case "update-paper":
		UpdatePaper()
	case "search":
		Search()
	default:
		PrintHelp()
	}
}
func PrintHelp() {
	fmt.Println("Usage:\n" +
		"- update-list: Update the paper list from selected subs\n" +
		"- search: Search papers by keywords")
}
