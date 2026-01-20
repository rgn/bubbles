# Quick Start Guide

## Enabling the Dropdown

To enable the dropdown feature in your textinput_autocomplete component:

```go
import "github.com/charmbracelet/bubbles/textinput_autocomplete"

// Create and configure the text input
ti := textinput_autocomplete.New()
ti.ShowSuggestions = true   // Required: Enable autocomplete
ti.ShowDropdown = true       // Required: Enable dropdown display
ti.SetSuggestions(yourSuggestionsList)
```

## Configuration Options

### Basic Setup
```go
ti.MaxDropdownItems = 5      // Show up to 5 items (default)
ti.Width = 30                // Input field width
ti.Placeholder = "Type to search..."
```

### Custom Styling
```go
// Normal dropdown items
ti.DropdownStyle = lipgloss.NewStyle().
    Foreground(lipgloss.Color("252"))

// Selected/highlighted item
ti.DropdownSelectedStyle = lipgloss.NewStyle().
    Foreground(lipgloss.Color("15")).
    Background(lipgloss.Color("240"))

// Matching text portion
ti.DropdownMatchStyle = lipgloss.NewStyle().
    Bold(true).
    Foreground(lipgloss.Color("205"))
```

## Keyboard Controls

| Key | Action |
|-----|--------|
| Type characters | Filter suggestions |
| ↑ (Up) | Previous suggestion |
| ↓ (Down) | Next suggestion |
| Tab | Accept current suggestion |
| Enter | Submit value |

## Behavior

1. **Auto-filtering**: Suggestions are filtered as you type using **substring matching** (matches anywhere in the text, case-insensitive)
2. **Auto-scrolling**: List scrolls automatically to keep selected item visible
3. **Overflow indicator**: "..." appears when more than 5 items match
4. **Bold matching**: The matched text is highlighted in bold within each suggestion

### Substring Matching Examples

- Type "apple" → matches "**Apple**", "Pine**apple**"
- Type "berry" → matches "Straw**berry**", "Blue**berry**", "Rasp**berry**"
- Type "script" → matches "Java**Script**", "Type**Script**"

## Example: Fruit Selector

```go
type model struct {
    input textinput_autocomplete.Model
}

func (m model) Init() tea.Cmd {
    return textinput_autocomplete.Blink
}

func initialModel() model {
    ti := textinput_autocomplete.New()
    ti.Placeholder = "Search fruits..."
    ti.Focus()
    ti.ShowSuggestions = true
    ti.ShowDropdown = true
    
    ti.SetSuggestions([]string{
        "Apple", "Apricot", "Avocado",
        "Banana", "Blueberry",
        "Cherry", "Coconut",
        // ... more fruits
    })
    
    return model{input: ti}
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    var cmd tea.Cmd
    
    if msg, ok := msg.(tea.KeyMsg); ok {
        if msg.Type == tea.KeyEnter {
            // User submitted - get the value
            value := m.input.Value()
            // Handle the selected value
        }
    }
    
    m.input, cmd = m.input.Update(msg)
    return m, cmd
}

func (m model) View() string {
    return "Select a fruit:\n\n" + m.input.View()
}
```

## Tips

- Set `ShowDropdown = false` to use inline completion mode instead
- Adjust `MaxDropdownItems` based on your terminal size
- Customize styles to match your application theme
- The dropdown automatically positions below the input field
