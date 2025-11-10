package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/dinesh00509/gitease/internals"
)

func main() {
	if len(os.Args) > 1 && (os.Args[1] == "-h" || os.Args[1] == "--help") {
		fmt.Println("GitEase - Interactive Git Assistant")
		fmt.Println("Usage: gitease")
		fmt.Println("Navigate and manage Git tasks interactively in your terminal.")
		os.Exit(0)
	}
	p := tea.NewProgram(internals.IntialModel())

	if err := p.Start(); err != nil {
		fmt.Println("Error running GitFlow:", err)
		os.Exit(1)
	}
}
