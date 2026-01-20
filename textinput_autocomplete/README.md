# TextInput Autocomplete with Dropdown

An enhanced text input component for Bubble Tea with autocomplete suggestions displayed in a dropdown list.

## Features

- **Dropdown Display**: Shows matching suggestions in a dropdown below the input field
- **Limited Visible Items**: Displays up to 5 items at a time (configurable via `MaxDropdownItems`)
- **Overflow Indicator**: Shows "..." when more items exist beyond the visible range
- **Keyboard Navigation**: Navigate through suggestions with up/down arrow keys
- **Smart Scrolling**: Automatically scrolls to keep the selected item visible
- **Bold Matching Text**: The matching portion of each suggestion is displayed in bold
- **Flexible Modes**: Can use dropdown mode or inline completion mode

## Usage

```go
package main

import (
	"github.com/charmbracelet/bubbles/textinput_autocomplete"
	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	textInput textinput_autocomplete.Model
}

func initialModel() model {
	ti := textinput_autocomplete.New()
	ti.Placeholder = "Type to search..."
	ti.Focus()
	
	// Enable dropdown mode
	ti.ShowSuggestions = true
	ti.ShowDropdown = true
	ti.MaxDropdownItems = 5  // Show up to 5 items (default)
	
	// Set suggestions
	ti.SetSuggestions([]string{
		"Apple",
		"Apricot",
		"Banana",
		"Cherry",
		// ... more items
	})

	return model{textInput: ti}
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	m.textInput, cmd = m.textInput.Update(msg)
	return m, cmd
}

func (m model) View() string {
	return m.textInput.View()
}
```

## Configuration

### Dropdown Settings

- `ShowDropdown` (bool): Enable/disable dropdown display
- `MaxDropdownItems` (int): Maximum number of items to show at once (default: 5)
- `ShowSuggestions` (bool): Enable/disable autocomplete suggestions

### Styling

The component provides several styling options:

- `DropdownStyle`: Style for dropdown items
- `DropdownSelectedStyle`: Style for the currently selected item
- `DropdownMatchStyle`: Style for the matching text portion (bold by default)

Example custom styling:

```go
ti.DropdownStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("252"))
ti.DropdownSelectedStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("15")).
	Background(lipgloss.Color("240"))
ti.DropdownMatchStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("205"))
```

## Key Bindings

- **Up/Down Arrow**: Navigate through suggestions
- **Tab**: Accept current suggestion
- **Enter**: Submit value
- **Escape**: Cancel/blur

## How It Works

1. When the user types, the component filters suggestions that match the input (case-insensitive **substring** matching - matches anywhere in the suggestion)
2. Matching suggestions are displayed in a dropdown below the input
3. The matching portion of each suggestion is highlighted in bold
4. Users can navigate with arrow keys; the list scrolls automatically
5. If there are more than `MaxDropdownItems` suggestions, "..." is shown at the bottom
6. Pressing Tab accepts the currently selected suggestion

## Matching Behavior

The component uses **substring matching**, which means it will match your input anywhere within the suggestion text:

- Typing "apple" will match both "**Apple**" and "Pine**apple**"
- Typing "script" will match "Java**Script**", "Type**Script**", and "Coffee**Script**"
- Typing "berry" will match "Straw**berry**", "Blue**berry**", "Rasp**berry**", etc.

Matching is case-insensitive, so "APPLE", "apple", and "ApPlE" will all match "Pineapple".

## Example

Run the example:

```bash
cd textinput_autocomplete/example
go run main.go
```

## Modes

### Dropdown Mode (New)
- Set `ShowDropdown = true`
- Shows all matching suggestions in a list
- Better for longer lists with many options
- Visual feedback with selection highlighting

### Inline Mode (Original)
- Set `ShowDropdown = false`
- Shows completion inline with grayed text
- Better for single-line completions
- More subtle, less intrusive

You can use both modes simultaneously or choose one based on your needs.
