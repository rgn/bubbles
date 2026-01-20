package main

import (
	"fmt"
	"log"
	"os"

	"github.com/charmbracelet/bubbles/textinput_autocomplete"
	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	textInput textinput_autocomplete.Model
	err       error
}

func initialModel() model {
	ti := textinput_autocomplete.New()
	ti.Placeholder = "Type to search..."
	ti.Focus()
	ti.CharLimit = 156
	ti.Width = 30

	// Enable suggestions and dropdown
	ti.ShowSuggestions = true
	ti.ShowDropdown = true

	// Set available suggestions
	suggestions := []string{
		"Apple",
		"Apricot",
		"Avocado",
		"Banana",
		"Blueberry",
		"Blackberry",
		"Cherry",
		"Coconut",
		"Cranberry",
		"Dragon Fruit",
		"Elderberry",
		"Fig",
		"Grape",
		"Grapefruit",
		"Guava",
		"Honeydew",
		"Kiwi",
		"Lemon",
		"Lime",
		"Mango",
		"Melon",
		"Orange",
		"Papaya",
		"Peach",
		"Pear",
		"Pineapple",
		"Plum",
		"Pomegranate",
		"Raspberry",
		"Strawberry",
		"Tangerine",
		"Watermelon",
		"Option 1",
		"Option 2",
		"Option 3",
		"Option 4",
		"Option 5",
		"Option 6",
		"Option 7",
		"Option 8",
	}
	ti.SetSuggestions(suggestions)

	return model{
		textInput: ti,
		err:       nil,
	}
}

func (m model) Init() tea.Cmd {
	return textinput_autocomplete.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter, tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		}
	case error:
		m.err = msg
		return m, nil
	}

	m.textInput, cmd = m.textInput.Update(msg)
	return m, cmd
}

func (m model) View() string {
	return fmt.Sprintf(
		"Type to filter fruits (up/down to navigate, tab to complete, enter to submit):\n\n%s\n\n%s",
		m.textInput.View(),
		"(esc to quit)",
	) + "\n"
}

func main() {
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}

	// After program exits, print the final value
	fmt.Printf("\nYou selected: %s\n", os.Args[0])
}
