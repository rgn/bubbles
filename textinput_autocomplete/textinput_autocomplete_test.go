package textinput_autocomplete

import (
	"fmt"
	"strconv"
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

func Test_CurrentSuggestion(t *testing.T) {
	textinput := New()
	textinput.ShowSuggestions = true

	suggestion := textinput.CurrentSuggestion()
	expected := ""
	if suggestion != expected {
		t.Fatalf("Error: expected no current suggestion but was %s", suggestion)
	}

	textinput.SetSuggestions([]string{"test1", "test2", "test3"})
	suggestion = textinput.CurrentSuggestion()
	expected = ""
	if suggestion != expected {
		t.Fatalf("Error: expected no current suggestion but was %s", suggestion)
	}

	textinput.SetValue("test")
	textinput.updateSuggestions()
	textinput.nextSuggestion()
	suggestion = textinput.CurrentSuggestion()
	expected = "test2"
	if suggestion != expected {
		t.Fatalf("Error: expected first suggestion but was %s", suggestion)
	}

	textinput.Blur()
	if strings.HasSuffix(textinput.View(), "test2") {
		t.Fatalf("Error: suggestions should not be rendered when input isn't focused. expected \"> test\" but got \"%s\"", textinput.View())
	}
}

func Test_SlicingOutsideCap(t *testing.T) {
	textinput := New()
	textinput.Placeholder = "作業ディレクトリを指定してください"
	textinput.Width = 32
	textinput.View()
}

func TestChinesePlaceholder(t *testing.T) {
	textinput := New()
	textinput.Placeholder = "输入消息..."
	textinput.Width = 20

	got := textinput.View()
	expected := "> 输入消息...       "
	if got != expected {
		t.Fatalf("expected %q but got %q", expected, got)
	}
}

func TestPlaceholderTruncate(t *testing.T) {
	textinput := New()
	textinput.Placeholder = "A very long placeholder, or maybe not so much"
	textinput.Width = 10

	got := textinput.View()
	expected := "> A very …"
	if got != expected {
		t.Fatalf("expected %q but got %q", expected, got)
	}
}

func ExampleValidateFunc() {
	creditCardNumber := New()
	creditCardNumber.Placeholder = "4505 **** **** 1234"
	creditCardNumber.Focus()
	creditCardNumber.CharLimit = 20
	creditCardNumber.Width = 30
	creditCardNumber.Prompt = ""
	// This anonymous function is a valid function for ValidateFunc.
	creditCardNumber.Validate = func(s string) error {
		// Credit Card Number should a string less than 20 digits
		// It should include 16 integers and 3 spaces
		if len(s) > 16+3 {
			return fmt.Errorf("CCN is too long")
		}

		if len(s) == 0 || len(s)%5 != 0 && (s[len(s)-1] < '0' || s[len(s)-1] > '9') {
			return fmt.Errorf("CCN is invalid")
		}

		// The last digit should be a number unless it is a multiple of 4 in which
		// case it should be a space
		if len(s)%5 == 0 && s[len(s)-1] != ' ' {
			return fmt.Errorf("CCN must separate groups with spaces")
		}

		// The remaining digits should be integers
		c := strings.ReplaceAll(s, " ", "")
		_, err := strconv.ParseInt(c, 10, 64)

		return err
	}
}

func keyPress(key rune) tea.Msg {
	return tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{key}, Alt: false}
}

func sendString(m Model, str string) Model {
	for _, k := range str {
		m, _ = m.Update(keyPress(k))
	}

	return m
}

// Tests for dropdown functionality

func TestNew_DropdownDefaults(t *testing.T) {
	m := New()

	if m.ShowDropdown != false {
		t.Errorf("Expected ShowDropdown to be false, got %v", m.ShowDropdown)
	}

	if m.MaxDropdownItems != 5 {
		t.Errorf("Expected MaxDropdownItems to be 5, got %d", m.MaxDropdownItems)
	}

	if m.dropdownScrollOffset != 0 {
		t.Errorf("Expected dropdownScrollOffset to be 0, got %d", m.dropdownScrollOffset)
	}
}

func TestDropdownView_NoSuggestions(t *testing.T) {
	m := New()
	m.ShowDropdown = true
	m.SetSuggestions([]string{})

	view := m.dropdownView()
	if view != "" {
		t.Errorf("Expected empty view with no suggestions, got %q", view)
	}
}

func TestDropdownView_NoInput(t *testing.T) {
	m := New()
	m.ShowDropdown = true
	m.SetSuggestions([]string{"Apple", "Banana"})

	view := m.dropdownView()
	if view != "" {
		t.Errorf("Expected empty view with no input, got %q", view)
	}
}

func TestDropdownView_WithMatches(t *testing.T) {
	m := New()
	m.ShowDropdown = true
	m.ShowSuggestions = true
	m.SetSuggestions([]string{"Apple", "Apricot", "Banana"})
	m.SetValue("ap")
	m.updateSuggestions()

	view := m.dropdownView()

	if !strings.Contains(view, "Apple") {
		t.Error("Expected view to contain 'Apple'")
	}

	if !strings.Contains(view, "Apricot") {
		t.Error("Expected view to contain 'Apricot'")
	}

	if strings.Contains(view, "Banana") {
		t.Error("Expected view to not contain 'Banana' (doesn't match)")
	}
}

func TestDropdownView_MaxItems(t *testing.T) {
	m := New()
	m.ShowDropdown = true
	m.ShowSuggestions = true
	m.MaxDropdownItems = 3

	suggestions := []string{
		"Apple", "Apricot", "Avocado", "Artichoke", "Asparagus",
	}
	m.SetSuggestions(suggestions)
	m.SetValue("a")
	m.updateSuggestions()

	view := m.dropdownView()
	lines := strings.Split(strings.TrimSpace(view), "\n")

	// Should show 3 items + "..." = 4 lines
	if len(lines) != 4 {
		t.Errorf("Expected 4 lines (3 items + ...), got %d", len(lines))
	}

	if !strings.Contains(view, "...") {
		t.Error("Expected view to contain '...' overflow indicator")
	}
}

func TestDropdownView_NoOverflow(t *testing.T) {
	m := New()
	m.ShowDropdown = true
	m.ShowSuggestions = true
	m.MaxDropdownItems = 5

	suggestions := []string{"Apple", "Apricot"}
	m.SetSuggestions(suggestions)
	m.SetValue("a")
	m.updateSuggestions()

	view := m.dropdownView()

	if strings.Contains(view, "...") {
		t.Error("Expected view to not contain '...' when items fit")
	}
}

func TestUpdateDropdownScroll_ScrollDown(t *testing.T) {
	m := New()
	m.ShowSuggestions = true
	m.MaxDropdownItems = 3

	suggestions := []string{
		"Item1", "Item2", "Item3", "Item4", "Item5", "Item6",
	}
	m.SetSuggestions(suggestions)
	m.SetValue("item")
	m.updateSuggestions()

	if m.dropdownScrollOffset != 0 {
		t.Errorf("Expected initial offset 0, got %d", m.dropdownScrollOffset)
	}

	// Select item beyond visible range
	m.currentSuggestionIndex = 4
	m.updateDropdownScroll()

	expectedOffset := 4 - m.MaxDropdownItems + 1 // = 2
	if m.dropdownScrollOffset != expectedOffset {
		t.Errorf("Expected scroll offset %d, got %d", expectedOffset, m.dropdownScrollOffset)
	}
}

func TestUpdateDropdownScroll_ScrollUp(t *testing.T) {
	m := New()
	m.ShowSuggestions = true
	m.MaxDropdownItems = 3

	suggestions := []string{
		"Item1", "Item2", "Item3", "Item4", "Item5", "Item6",
	}
	m.SetSuggestions(suggestions)
	m.SetValue("item")
	m.updateSuggestions()

	// Start scrolled down
	m.currentSuggestionIndex = 5
	m.updateDropdownScroll()

	// Select item above visible range
	m.currentSuggestionIndex = 1
	m.updateDropdownScroll()

	if m.dropdownScrollOffset != 1 {
		t.Errorf("Expected scroll offset 1, got %d", m.dropdownScrollOffset)
	}
}

func TestUpdateDropdownScroll_Bounds(t *testing.T) {
	m := New()
	m.ShowSuggestions = true
	m.MaxDropdownItems = 3

	suggestions := []string{"Item1", "Item2", "Item3", "Item4"}
	m.SetSuggestions(suggestions)
	m.SetValue("item")
	m.updateSuggestions()

	// Try to scroll beyond max
	m.currentSuggestionIndex = 3
	m.dropdownScrollOffset = 10
	m.updateDropdownScroll()

	maxOffset := len(m.matchedSuggestions) - m.MaxDropdownItems
	if m.dropdownScrollOffset > maxOffset {
		t.Errorf("Expected offset <= %d, got %d", maxOffset, m.dropdownScrollOffset)
	}

	// Try negative scroll
	m.currentSuggestionIndex = 0
	m.dropdownScrollOffset = -5
	m.updateDropdownScroll()

	if m.dropdownScrollOffset < 0 {
		t.Errorf("Expected offset >= 0, got %d", m.dropdownScrollOffset)
	}
}

func TestNextSuggestion_WithScrolling(t *testing.T) {
	m := New()
	m.ShowSuggestions = true
	m.MaxDropdownItems = 3

	suggestions := []string{
		"Item1", "Item2", "Item3", "Item4", "Item5",
	}
	m.SetSuggestions(suggestions)
	m.SetValue("item")
	m.updateSuggestions()

	if m.currentSuggestionIndex != 0 {
		t.Errorf("Expected initial index 0, got %d", m.currentSuggestionIndex)
	}

	// Move to next (index 1) - should not scroll yet
	m.nextSuggestion()
	if m.currentSuggestionIndex != 1 {
		t.Errorf("Expected index 1, got %d", m.currentSuggestionIndex)
	}
	if m.dropdownScrollOffset != 0 {
		t.Errorf("Expected offset 0, got %d", m.dropdownScrollOffset)
	}

	// Move to index 2 - still visible
	m.nextSuggestion()
	if m.currentSuggestionIndex != 2 {
		t.Errorf("Expected index 2, got %d", m.currentSuggestionIndex)
	}

	// Move to index 3 - should scroll
	m.nextSuggestion()
	if m.currentSuggestionIndex != 3 {
		t.Errorf("Expected index 3, got %d", m.currentSuggestionIndex)
	}
	if m.dropdownScrollOffset != 1 {
		t.Errorf("Expected offset 1, got %d", m.dropdownScrollOffset)
	}
}

func TestPreviousSuggestion_WithScrolling(t *testing.T) {
	m := New()
	m.ShowSuggestions = true
	m.MaxDropdownItems = 3

	suggestions := []string{
		"Item1", "Item2", "Item3", "Item4", "Item5",
	}
	m.SetSuggestions(suggestions)
	m.SetValue("item")
	m.updateSuggestions()

	// Start at end
	m.currentSuggestionIndex = 4
	m.updateDropdownScroll()

	// Move to previous
	m.previousSuggestion()
	if m.currentSuggestionIndex != 3 {
		t.Errorf("Expected index 3, got %d", m.currentSuggestionIndex)
	}

	// Continue moving back - should scroll
	m.previousSuggestion() // index 2
	m.previousSuggestion() // index 1 - should scroll up

	if m.currentSuggestionIndex != 1 {
		t.Errorf("Expected index 1, got %d", m.currentSuggestionIndex)
	}
	if m.dropdownScrollOffset != 1 {
		t.Errorf("Expected offset 1, got %d", m.dropdownScrollOffset)
	}
}

func TestNextSuggestion_Wrapping(t *testing.T) {
	m := New()
	m.ShowSuggestions = true

	suggestions := []string{"Item1", "Item2", "Item3"}
	m.SetSuggestions(suggestions)
	m.SetValue("item")
	m.updateSuggestions()

	m.currentSuggestionIndex = 2

	m.nextSuggestion()
	if m.currentSuggestionIndex != 0 {
		t.Errorf("Expected index to wrap to 0, got %d", m.currentSuggestionIndex)
	}

	if m.dropdownScrollOffset != 0 {
		t.Errorf("Expected offset to reset to 0, got %d", m.dropdownScrollOffset)
	}
}

func TestPreviousSuggestion_Wrapping(t *testing.T) {
	m := New()
	m.ShowSuggestions = true

	suggestions := []string{"Item1", "Item2", "Item3"}
	m.SetSuggestions(suggestions)
	m.SetValue("item")
	m.updateSuggestions()

	m.currentSuggestionIndex = 0

	m.previousSuggestion()
	if m.currentSuggestionIndex != 2 {
		t.Errorf("Expected index to wrap to 2, got %d", m.currentSuggestionIndex)
	}
}

func TestNextSuggestion_EmptySuggestions(t *testing.T) {
	m := New()
	m.ShowSuggestions = true
	m.SetSuggestions([]string{})

	m.nextSuggestion()

	if m.currentSuggestionIndex != 0 {
		t.Errorf("Expected index to remain 0, got %d", m.currentSuggestionIndex)
	}
}

func TestPreviousSuggestion_EmptySuggestions(t *testing.T) {
	m := New()
	m.ShowSuggestions = true
	m.SetSuggestions([]string{})

	m.previousSuggestion()

	if m.currentSuggestionIndex != 0 {
		t.Errorf("Expected index to remain 0, got %d", m.currentSuggestionIndex)
	}
}

func TestUpdateSuggestions_ResetsScroll(t *testing.T) {
	m := New()
	m.ShowSuggestions = true
	m.MaxDropdownItems = 3

	suggestions := []string{
		"Apple", "Apricot", "Avocado", "Banana", "Blueberry",
	}
	m.SetSuggestions(suggestions)
	m.SetValue("a")
	m.updateSuggestions()

	// Navigate down
	m.currentSuggestionIndex = 2
	m.updateDropdownScroll()

	// Change input - should reset
	m.SetValue("b")
	m.updateSuggestions()

	if m.currentSuggestionIndex != 0 {
		t.Errorf("Expected index to reset to 0, got %d", m.currentSuggestionIndex)
	}

	if m.dropdownScrollOffset != 0 {
		t.Errorf("Expected scroll offset to reset to 0, got %d", m.dropdownScrollOffset)
	}
}

func TestView_DropdownEnabled(t *testing.T) {
	m := New()
	m.ShowSuggestions = true
	m.ShowDropdown = true
	m.Focus()

	suggestions := []string{"Apple", "Apricot"}
	m.SetSuggestions(suggestions)
	m.SetValue("ap")
	m.updateSuggestions()

	view := m.View()

	if !strings.Contains(view, "Apple") {
		t.Error("Expected view to contain dropdown with 'Apple'")
	}

	if !strings.Contains(view, "\n") {
		t.Error("Expected view to contain newlines for dropdown")
	}
}

func TestView_DropdownDisabled(t *testing.T) {
	m := New()
	m.ShowSuggestions = true
	m.ShowDropdown = false
	m.Focus()

	suggestions := []string{"Apple", "Apricot"}
	m.SetSuggestions(suggestions)
	m.SetValue("ap")

	view := m.View()

	// In inline mode, there may still be completion text
	// but dropdown items should not appear as separate lines
	lines := strings.Split(view, "\n")
	if len(lines) > 1 {
		t.Errorf("Expected single line view (inline mode), got %d lines", len(lines))
	}
}

func TestView_DropdownNotFocused(t *testing.T) {
	m := New()
	m.ShowSuggestions = true
	m.ShowDropdown = true

	suggestions := []string{"Apple", "Apricot"}
	m.SetSuggestions(suggestions)
	m.SetValue("ap")

	view := m.View()

	// Should not show dropdown when not focused
	lines := strings.Split(view, "\n")
	if len(lines) > 1 {
		t.Error("Expected no dropdown when not focused")
	}
}

func TestUpdate_NavigationKeys_Dropdown(t *testing.T) {
	m := New()
	m.ShowSuggestions = true
	m.ShowDropdown = true
	m.Focus()

	suggestions := []string{"Item1", "Item2", "Item3"}
	m.SetSuggestions(suggestions)
	m.SetValue("item")
	m.updateSuggestions()

	// Test down key
	m, _ = m.Update(tea.KeyMsg{Type: tea.KeyDown})
	if m.currentSuggestionIndex != 1 {
		t.Errorf("Expected index 1 after down key, got %d", m.currentSuggestionIndex)
	}

	// Test up key
	m, _ = m.Update(tea.KeyMsg{Type: tea.KeyUp})
	if m.currentSuggestionIndex != 0 {
		t.Errorf("Expected index 0 after up key, got %d", m.currentSuggestionIndex)
	}
}

func TestMaxDropdownItems_Configuration(t *testing.T) {
	tests := []struct {
		name             string
		maxItems         int
		totalSuggestions int
		expectedLines    int
	}{
		{"Default 5", 5, 10, 6},  // 5 items + "..."
		{"Custom 3", 3, 10, 4},   // 3 items + "..."
		{"Custom 7", 7, 10, 8},   // 7 items + "..."
		{"No overflow", 5, 3, 3}, // 3 items, no "..."
		{"Exact fit", 5, 5, 5},   // 5 items, no "..."
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := New()
			m.ShowSuggestions = true
			m.ShowDropdown = true
			m.MaxDropdownItems = tt.maxItems

			suggestions := make([]string, tt.totalSuggestions)
			for i := 0; i < tt.totalSuggestions; i++ {
				suggestions[i] = strings.Repeat("a", i+1)
			}
			m.SetSuggestions(suggestions)
			m.SetValue("a")
			m.updateSuggestions()

			view := m.dropdownView()
			lines := strings.Split(strings.TrimSpace(view), "\n")

			if len(lines) != tt.expectedLines {
				t.Errorf("Expected %d lines, got %d", tt.expectedLines, len(lines))
			}
		})
	}
}

func TestDropdownView_CaseInsensitiveMatching(t *testing.T) {
	m := New()
	m.ShowSuggestions = true
	m.ShowDropdown = true
	m.SetSuggestions([]string{"Apple", "APRICOT", "aPpLe"})
	m.SetValue("AP")
	m.updateSuggestions()

	view := m.dropdownView()

	if !strings.Contains(view, "Apple") {
		t.Error("Expected lowercase 'Apple' to match")
	}

	if !strings.Contains(view, "APRICOT") {
		t.Error("Expected uppercase 'APRICOT' to match")
	}

	if !strings.Contains(view, "aPpLe") {
		t.Error("Expected mixed case 'aPpLe' to match")
	}
}

func TestDropdownView_MatchHighlighting(t *testing.T) {
	m := New()
	m.ShowSuggestions = true
	m.ShowDropdown = true
	m.SetSuggestions([]string{"Apple", "Apricot"})
	m.SetValue("ap")
	m.updateSuggestions()

	view := m.dropdownView()

	// View should contain the suggestions
	if !strings.Contains(view, "ple") {
		t.Error("Expected 'ple' to be in view (part of Apple)")
	}

	if !strings.Contains(view, "ricot") {
		t.Error("Expected 'ricot' to be in view (part of Apricot)")
	}
}

func TestDropdownScrolling_LargeList(t *testing.T) {
	m := New()
	m.ShowSuggestions = true
	m.ShowDropdown = true
	m.MaxDropdownItems = 5

	// Create 20 suggestions
	suggestions := make([]string, 20)
	for i := 0; i < 20; i++ {
		suggestions[i] = fmt.Sprintf("Item%02d", i+1)
	}
	m.SetSuggestions(suggestions)
	m.SetValue("item")
	m.updateSuggestions()

	// Navigate to the last item
	for i := 0; i < 19; i++ {
		m.nextSuggestion()
	}

	if m.currentSuggestionIndex != 19 {
		t.Errorf("Expected index 19, got %d", m.currentSuggestionIndex)
	}

	// Scroll should be at bottom
	expectedOffset := 19 - m.MaxDropdownItems + 1
	if m.dropdownScrollOffset != expectedOffset {
		t.Errorf("Expected offset %d, got %d", expectedOffset, m.dropdownScrollOffset)
	}

	// Navigate back to top
	for i := 0; i < 19; i++ {
		m.previousSuggestion()
	}

	if m.currentSuggestionIndex != 0 {
		t.Errorf("Expected index 0, got %d", m.currentSuggestionIndex)
	}

	if m.dropdownScrollOffset != 0 {
		t.Errorf("Expected offset 0, got %d", m.dropdownScrollOffset)
	}
}

// Tests for substring matching (not just prefix)

func TestSubstringMatching_MiddleOfString(t *testing.T) {
	m := New()
	m.ShowSuggestions = true
	m.ShowDropdown = true

	suggestions := []string{
		"Apple", "Pineapple", "Banana", "Orange",
	}
	m.SetSuggestions(suggestions)
	m.SetValue("apple")
	m.updateSuggestions()

	// Should match both "Apple" and "Pineapple"
	matched := m.MatchedSuggestions()
	if len(matched) != 2 {
		t.Errorf("Expected 2 matches, got %d: %v", len(matched), matched)
	}

	matchMap := make(map[string]bool)
	for _, match := range matched {
		matchMap[match] = true
	}

	if !matchMap["Apple"] {
		t.Error("Expected 'Apple' to match")
	}

	if !matchMap["Pineapple"] {
		t.Error("Expected 'Pineapple' to match")
	}

	if matchMap["Banana"] {
		t.Error("Expected 'Banana' to not match")
	}
}

func TestSubstringMatching_EndOfString(t *testing.T) {
	m := New()
	m.ShowSuggestions = true

	suggestions := []string{
		"JavaScript", "TypeScript", "CoffeeScript", "Python",
	}
	m.SetSuggestions(suggestions)
	m.SetValue("script")
	m.updateSuggestions()

	matched := m.MatchedSuggestions()
	if len(matched) != 3 {
		t.Errorf("Expected 3 matches, got %d: %v", len(matched), matched)
	}

	matchMap := make(map[string]bool)
	for _, match := range matched {
		matchMap[match] = true
	}

	if !matchMap["JavaScript"] {
		t.Error("Expected 'JavaScript' to match")
	}

	if !matchMap["TypeScript"] {
		t.Error("Expected 'TypeScript' to match")
	}

	if !matchMap["CoffeeScript"] {
		t.Error("Expected 'CoffeeScript' to match")
	}

	if matchMap["Python"] {
		t.Error("Expected 'Python' to not match")
	}
}

func TestSubstringMatching_AnyPosition(t *testing.T) {
	m := New()
	m.ShowSuggestions = true

	suggestions := []string{
		"before", "after", "afternoon", "beforehand", "therefor",
	}
	m.SetSuggestions(suggestions)
	m.SetValue("for")
	m.updateSuggestions()

	matched := m.MatchedSuggestions()
	if len(matched) != 3 {
		t.Errorf("Expected 3 matches, got %d: %v", len(matched), matched)
	}

	matchMap := make(map[string]bool)
	for _, match := range matched {
		matchMap[match] = true
	}

	if !matchMap["before"] {
		t.Error("Expected 'before' to match (contains 'for')")
	}

	if !matchMap["therefor"] {
		t.Error("Expected 'therefor' to match")
	}

	if !matchMap["beforehand"] {
		t.Error("Expected 'beforehand' to match")
	}
}

func TestSubstringMatching_CaseInsensitive(t *testing.T) {
	m := New()
	m.ShowSuggestions = true

	suggestions := []string{
		"GameObject", "ScriptableObject", "MonoBehaviour", "Component",
	}
	m.SetSuggestions(suggestions)
	m.SetValue("OBJECT")
	m.updateSuggestions()

	matched := m.MatchedSuggestions()
	if len(matched) != 2 {
		t.Errorf("Expected 2 matches, got %d: %v", len(matched), matched)
	}

	matchMap := make(map[string]bool)
	for _, match := range matched {
		matchMap[match] = true
	}

	if !matchMap["GameObject"] {
		t.Error("Expected 'GameObject' to match (case insensitive)")
	}

	if !matchMap["ScriptableObject"] {
		t.Error("Expected 'ScriptableObject' to match (case insensitive)")
	}
}

func TestSubstringMatching_NoMatches(t *testing.T) {
	m := New()
	m.ShowSuggestions = true

	suggestions := []string{
		"Apple", "Banana", "Cherry",
	}
	m.SetSuggestions(suggestions)
	m.SetValue("xyz")
	m.updateSuggestions()

	matched := m.MatchedSuggestions()
	if len(matched) != 0 {
		t.Errorf("Expected 0 matches, got %d: %v", len(matched), matched)
	}
}

func TestSubstringMatching_DropdownView(t *testing.T) {
	m := New()
	m.ShowSuggestions = true
	m.ShowDropdown = true

	suggestions := []string{
		"Strawberry", "Blueberry", "Raspberry", "Blackberry", "Cranberry",
	}
	m.SetSuggestions(suggestions)
	m.SetValue("berry")
	m.updateSuggestions()

	view := m.dropdownView()

	// All should match
	if !strings.Contains(view, "Strawberry") {
		t.Error("Expected 'Strawberry' in view")
	}

	if !strings.Contains(view, "Blueberry") {
		t.Error("Expected 'Blueberry' in view")
	}

	if !strings.Contains(view, "Raspberry") {
		t.Error("Expected 'Raspberry' in view")
	}

	if !strings.Contains(view, "Blackberry") {
		t.Error("Expected 'Blackberry' in view")
	}

	if !strings.Contains(view, "Cranberry") {
		t.Error("Expected 'Cranberry' in view")
	}
}

func TestSubstringMatching_PrefixStillWorks(t *testing.T) {
	m := New()
	m.ShowSuggestions = true

	suggestions := []string{
		"Apple", "Application", "Apricot", "Banana",
	}
	m.SetSuggestions(suggestions)
	m.SetValue("app")
	m.updateSuggestions()

	matched := m.MatchedSuggestions()
	if len(matched) != 2 {
		t.Errorf("Expected 2 matches, got %d: %v", len(matched), matched)
	}

	matchMap := make(map[string]bool)
	for _, match := range matched {
		matchMap[match] = true
	}

	if !matchMap["Apple"] {
		t.Error("Expected 'Apple' to match (prefix)")
	}

	if !matchMap["Application"] {
		t.Error("Expected 'Application' to match (prefix)")
	}
}

func TestSubstringMatching_WithDropdownScroll(t *testing.T) {
	m := New()
	m.ShowSuggestions = true
	m.ShowDropdown = true
	m.MaxDropdownItems = 3

	// Create items with common substring
	suggestions := []string{
		"before", "reform", "perform", "information", "formation", "foremost",
	}
	m.SetSuggestions(suggestions)
	m.SetValue("for")
	m.updateSuggestions()

	// Should match: before, reform, perform, information, formation, foremost (6 matches)
	matched := m.MatchedSuggestions()
	if len(matched) != 6 {
		t.Errorf("Expected 6 matches, got %d: %v", len(matched), matched)
	}

	// Test scrolling through matches
	view := m.dropdownView()
	lines := strings.Split(strings.TrimSpace(view), "\n")

	// Should show 3 items + "..."
	if len(lines) != 4 {
		t.Errorf("Expected 4 lines (3 items + ...), got %d", len(lines))
	}

	// Navigate to last item
	for i := 0; i < 4; i++ {
		m.nextSuggestion()
	}

	// Scroll should have adjusted
	if m.dropdownScrollOffset == 0 {
		t.Error("Expected scroll offset to have changed")
	}
}

func TestSubstringMatching_HighlightCorrectPart(t *testing.T) {
	m := New()
	m.ShowSuggestions = true
	m.ShowDropdown = true

	suggestions := []string{"Pineapple"}
	m.SetSuggestions(suggestions)
	m.SetValue("apple")
	m.updateSuggestions()

	view := m.dropdownView()

	// Should contain "Pineapple" with "apple" highlighted
	if !strings.Contains(view, "Pine") {
		t.Error("Expected 'Pine' to be in view")
	}

	if !strings.Contains(view, "apple") {
		t.Error("Expected 'apple' to be in view (highlighted)")
	}
}

func TestSubstringMatching_MultipleOccurrences(t *testing.T) {
	m := New()
	m.ShowSuggestions = true

	// Test that it matches even when substring appears multiple times
	suggestions := []string{
		"testtest", "contest", "fastest", "testing",
	}
	m.SetSuggestions(suggestions)
	m.SetValue("test")
	m.updateSuggestions()

	matched := m.MatchedSuggestions()
	if len(matched) != 4 {
		t.Errorf("Expected 4 matches, got %d: %v", len(matched), matched)
	}
}
