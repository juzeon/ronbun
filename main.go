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
	case "update-abstract":
		UpdateAbstract()
	case "update-embedding":
		UpdateEmbedding()
	case "search":
		Search()
	case "search-vec":
		SearchVec()
	case "translate":
		Translate()
	default:
		PrintHelp()
	}
}
func PrintHelp() {
	fmt.Println("Usage:\n" +
		"- update-list: Update the paper list from selected subs\n" +
		"- update-abstract: Update papers' abstracts\n" +
		"- update-embedding: Update papers' embeddings\n" +
		"- search: Search papers by keywords\n" +
		"- search-vec: Search papers by document\n" +
		"- translate: Translate a pdf document to Chinese")
}
