// Package views collects the various app state views in the Termsheet app
package views

import "fmt"

func RenderMenu(cursor int, choices []string) string {
	s := "Please select a view\n\n"

	for i, choice := range choices {
		cursorStr := " "
		if cursor == i {
			cursorStr = ">"
		}
		s += fmt.Sprintf("%s %s\n", cursorStr, choice)
	}

	s += "\nPress q to quit."
	return s
}
