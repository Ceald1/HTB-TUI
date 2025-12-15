
package searchengine

import (
	"fmt"
	"strings"

	"time"

	// "github.com/charmbracelet/glamour"

	"context"

	"github.com/Ceald1/HTB-TUI/src/format"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	HTB "github.com/gubarz/gohtb"
	v4 "github.com/gubarz/gohtb/httpclient/v4"
	S "github.com/gubarz/gohtb/services/search"
)
var ctx = context.Background()
var SELECTED_ITEM any
var lastInput string

type timeoutMsg struct{}

// Command that waits 1 seconds then sends timeoutMsg
func waitForTimeout() tea.Cmd {
	return tea.Tick(1*time.Second, func(t time.Time) tea.Msg {
		return timeoutMsg{}
	})
}


func NewSearch(HTBClient *HTB.Client, keyword string) (result S.SearchResponse, err error) {
	result, err = HTBClient.Search.Query(keyword).All(ctx)
	return result,err
}



func parseResults(r S.SearchResponse) (result SearchResult) {
	challenges := r.Data.Challenges
	machines := r.Data.Machines
	sherlocks := r.Data.Sherlocks
	teams := r.Data.Teams
	users := r.Data.Users
	for _, challenge := range challenges {
		challenge.Value = fmt.Sprintf("%s: %s", lipgloss.NewStyle().Foreground(format.TextYellow).Render("challenge") ,lipgloss.NewStyle().Foreground(format.NextColor()).Render(challenge.Value))
		
		result = append(result, challenge)
	}
	for _, machine := range machines {
		machine.Value = fmt.Sprintf("%s: %s", lipgloss.NewStyle().Foreground(format.LightGreen).Render("box"), lipgloss.NewStyle().Foreground(format.NextColor()).Render(machine.Value))
		result = append(result, machine)
	}
	for _, sherlock := range sherlocks {
		sherlock.Value = fmt.Sprintf("%s: %s", lipgloss.NewStyle().Foreground(format.TextCyan).Render("sherlock"),lipgloss.NewStyle().Foreground(format.NextColor()).Render(sherlock.Value))
		result = append(result, sherlock)
	}
	for _, team := range teams {
		team.Value = fmt.Sprintf("%s: %s", lipgloss.NewStyle().Foreground(format.TextPink).Render("team"),lipgloss.NewStyle().Foreground(format.NextColor()).Render(team.Value))
		result = append(result, team)
	}
	for _, user := range users {
		user.Value = fmt.Sprintf("%s: %s", lipgloss.NewStyle().Foreground(format.TextPurple).Render("user"),lipgloss.NewStyle().Foreground(format.NextColor()).Render(user.Value))
		result = append(result, user)
	}
	return
}

type model struct {
	choices SearchResult
	cursor int
	selected any
	input string
	timeoutActive bool
	HTBClient *HTB.Client
}


func initialModel(choices SearchResult, HTBClient *HTB.Client) model {
	return model{
		choices: choices,
		selected: nil,
		input: "",
		timeoutActive: true,
		HTBClient: HTBClient,
	}	
}

func (m model) Init() tea.Cmd {
	// Just return `nil`, which means "no I/O right now, please."
	return waitForTimeout()
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case timeoutMsg:
		// Handle timeout - run your function here
		if m.timeoutActive {
			if len(m.input) > 2 && m.input != lastInput {
				lastInput = m.input
				newSearch, err := NewSearch(m.HTBClient, m.input)
				if err != nil {
					fmt.Println(err.Error())
					return m, nil
				}
				m.choices = parseResults(newSearch)
				// Reset cursor when new results come in
				m.cursor = 0
			}
		}
		return m, nil

	case tea.KeyMsg:
		// Reset timeout on any keypress
		var cmd tea.Cmd = waitForTimeout()
		m.timeoutActive = true

		switch msg.String() {

		case "ctrl+c", "esc":
			return m, tea.Quit
		
		case "up":
			// Get filtered choices first
			filteredChoices := m.getFilteredChoices()
			if len(filteredChoices) > 0 {
				if m.cursor > 0 {
					m.cursor--
				} else {
					m.cursor = len(filteredChoices) - 1
				}
			}
		
		case "down":
			// Get filtered choices first
			filteredChoices := m.getFilteredChoices()
			if len(filteredChoices) > 0 {
				if m.cursor < len(filteredChoices)-1 {
					m.cursor++
				} else {
					m.cursor = 0 // loop over to the cursor top
				}
			}
		
		case "enter":
			// Set selected item from filtered choices
			filteredChoices := m.getFilteredChoices()
			if len(filteredChoices) > 0 && m.cursor < len(filteredChoices) {
				m.selected = filteredChoices[m.cursor]
				SELECTED_ITEM = m.selected
			}
			return m, tea.Quit

		case "backspace":
			if m.input != "" {
				m.input = m.input[:len(m.input)-1]
				// Reset cursor when input changes
				m.cursor = 0
			}
		
		default:
			// Only add printable characters
			if len(msg.String()) == 1 {
				m.input += msg.String()
				// Reset cursor when input changes
				m.cursor = 0
			}
		}

		return m, cmd
	}
	return m, nil
}


type SearchResult []any

func ExtractValue(choice any) (value string) {
	var data string
	chal, ok := choice.(v4.SearchChallengeItem)
	if ok {
		// data = fmt.Sprintf("challenge: %s", chal.Value)
		data = chal.Value
	}
	box, ok := choice.(v4.SearchFetchMachinesItem)
	if ok {
		// data = fmt.Sprintf("box: %s", box.Value)
		data = box.Value
	}
	sherlock, ok := choice.(v4.SearchSherlockItem)
	if ok {
		// data= fmt.Sprintf("sherlock: %s", sherlock.Value)
		data = sherlock.Value
	}

	team, ok := choice.(v4.SearchTeamItem)
	if ok {
		// data= fmt.Sprintf("team: %s", team.Value)
		data = team.Value
	}
	user, ok := choice.(v4.SearchUserItem)
	if ok {
		// data = fmt.Sprintf("user: %s", user.Value)
		data = user.Value
	}
	// colored := lipgloss.NewStyle().Foreground(format.NextColor()).Render(data)
	return data
}

// Helper method to get filtered choices
func (m model) getFilteredChoices() SearchResult {
	var choices SearchResult
	if len(m.input) > 0 {
		for _, choice := range m.choices {
			c := ExtractValue(choice)
			if c != "" {
				if strings.Contains(strings.ToLower(c), strings.ToLower(m.input)) {
					choices = append(choices, choice)
				}
			}
		}
	} else {
		choices = m.choices[:]
	}
	return choices
}


func (m model) View() string {
	// Input field with better styling
	inputStyle := lipgloss.NewStyle().
		Foreground(format.TextYellow).
		Bold(true)
	
	s := inputStyle.Render("Search: ") + m.input + "\n"
	s += strings.Repeat("─", 50) + "\n\n"

	// Get filtered choices
	choices := m.getFilteredChoices()

	if len(choices) == 0 {
		s += "No results found.\n"
	} else {
		// Iterate over filtered choices
		for i, choice := range choices {
			// Is the cursor pointing at this choice?
			cursor := " " // no cursor
			if m.cursor == i {
				cursor = ">" // cursor!
			}

			// Render the row
			value := ExtractValue(choice)
			if m.cursor == i {
				// Highlight selected row
				s += fmt.Sprintf("%s %s\n", 
					lipgloss.NewStyle().Foreground(format.TextYellow).Render(cursor),
					lipgloss.NewStyle().Bold(true).Render(value))
			} else {
				s += fmt.Sprintf("%s %s\n", cursor, value)
			}
		}
	}

	// The footer
	s += "\n" + strings.Repeat("─", 50) + "\n"
	s += "↑/↓: navigate • enter: select • esc: quit\n"

	// Send the UI for rendering
	return s
}

func NewModel(HTBClient *HTB.Client) model{
	m := model{
		HTBClient: HTBClient,
		cursor: 0,
		selected: nil,
		input: "",
		timeoutActive: true,
	}
	return m
}

type SelectedChoice struct {
	Product string
	Id int
	Avatar string
}


func ExtractSearchValue() (result SelectedChoice) {
	choice := SELECTED_ITEM
	chal, ok := choice.(v4.SearchChallengeItem)
	if ok {
		result = SelectedChoice{
			Product: "challenge",
			Id: chal.Id,
		}
	}
	box, ok := choice.(v4.SearchFetchMachinesItem)
	if ok {
		result = SelectedChoice{
			Product: "box",
			Id: box.Id,
			Avatar: box.Avatar,
		}
	}
	sherlock, ok := choice.(v4.SearchSherlockItem)
	if ok {
		result = SelectedChoice{
			Product: "sherlock",
			Id: sherlock.Id,
			Avatar: sherlock.Avatar,
		}
	}
	team, ok := choice.(v4.SearchTeamItem)
	if ok {
		result = SelectedChoice{
			Product: "team",
			Id: team.Id,
			Avatar: team.Avatar,
		}
	}
	user, ok := choice.(v4.SearchUserItem)
	if ok {
		result = SelectedChoice{
			Product: "user",
			Id: user.Id,
			Avatar: user.Avatar,
		}
	}
	return result
}



func Run(HTBClient *HTB.Client){
	p := tea.NewProgram(NewModel(HTBClient), tea.WithAltScreen())
	p.Run()
}