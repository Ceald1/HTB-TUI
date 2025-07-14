package challenges

import (
	"context"
	"fmt"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/evertras/bubble-table/table"
	HTB "github.com/gubarz/gohtb"
	"github.com/gubarz/gohtb/services/challenges"
	"github.com/Ceald1/HTB-TUI/src/format"
)

var (
	ctx         = context.Background()
	SelectedChallenge = ""
)

const (
	ColumnName     = "name"
	ColumnDifficulty = "difficulty"
	ColumnCategory   = "category"
	ColumnStatus   	 = "status"
	ColumnChallengeID= "id"

	minWidth            = 30
	minHeight           = 8
	fixedVerticalMargin = 4
)

type model struct {
	Challenges             table.Model
	Selected          string
	totalWidth        int
	totalHeight       int
	horizontalMargin  int
	verticalMargin    int
}

type challengeListData struct {
	Challenges []challenges.ChallengeList
}

func getChallenges(HTBClient *HTB.Client) (machineList challengeListData) {
	
	challenges, err := HTBClient.Challenges.List().AllResults(ctx)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	machineList.Challenges = challenges.Data
	return machineList
}

func NewModel(Challenges challengeListData) model {
	columns := []table.Column{
		table.NewColumn(ColumnName, lipgloss.NewStyle().Foreground(format.TextTitle).Render("Name"), 30).WithFiltered(true),
		table.NewFlexColumn(ColumnCategory, lipgloss.NewStyle().Foreground(format.TextTitle).Render("Category"), 1).WithFiltered(true),
		table.NewFlexColumn(ColumnDifficulty, lipgloss.NewStyle().Foreground(format.TextTitle).Render("Difficulty"), 1).WithFiltered(true),
		table.NewFlexColumn(ColumnStatus, lipgloss.NewStyle().Foreground(format.TextTitle).Render("Status"), 1).WithFiltered(true),
		table.NewFlexColumn(ColumnChallengeID, lipgloss.NewStyle().Foreground(format.TextTitle).Render("ID"), 1).WithFiltered(true),
	}

	rows := genRows(Challenges)

	keys := table.DefaultKeyMap()
	keys.RowDown.SetKeys("j", "down", "s")
	keys.RowUp.SetKeys("k", "up", "w")

	m := model{
		Challenges: table.New(columns).
			WithRows(rows).
			SelectableRows(true).
			Filtered(true).
			Focused(true).
			WithKeyMap(keys).
			WithPageSize(10).
			WithBaseStyle(lipgloss.NewStyle().Background(format.BaseBG)).
			WithSelectedText(" ", "✓"),
	}

	m.updateFooter()

	return m
}

func (m *model) updateFooter() {
	highlightedRow := m.Challenges.HighlightedRow()

	footerText := fmt.Sprintf(
		"Pg. %d/%d - Currently looking at: %s",
		m.Challenges.CurrentPage(),
		m.Challenges.MaxPages(),
		highlightedRow.Data[ColumnName],
	)

	m.Challenges = m.Challenges.WithStaticFooter(footerText)
}

func genRows(Challenges challengeListData) (rows []table.Row) {
	for _, Challenge := range Challenges.Challenges {
		c := strings.ToLower(Challenge.State)
		var state string
		switch c {
			case "active":
				state = lipgloss.NewStyle().Foreground(format.TextLightGreen).Render("Active")
			case "retired":
				state = lipgloss.NewStyle().Foreground(format.TextYellow).Render("Retired")
			case "unreleased":
				state = lipgloss.NewStyle().Foreground(format.Pink).Render("Unreleased")

		}
		rows = append(rows, table.NewRow(table.RowData{
			ColumnName:     Challenge.Name,
			ColumnCategory:       Challenge.CategoryName,
			ColumnDifficulty: format.CheckDiff(Challenge.Difficulty),
			ColumnStatus:   lipgloss.NewStyle().Foreground(format.TextPink).Render(state),
			ColumnChallengeID:    fmt.Sprintf("%d", Challenge.Id),
		}))
	}
	return rows
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	m.Challenges, cmd = m.Challenges.Update(msg)
	cmds = append(cmds, cmd)

	m.updateFooter()

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			cmds = append(cmds, tea.Quit)

		case "enter", "return":
			SelectedChallenge = m.Challenges.HighlightedRow().Data[ColumnChallengeID].(string)
			cmds = append(cmds, tea.Quit)
		case "esc":
			m.Challenges.Filtered(false)
		case "/":
			m.Challenges.Filtered(true)

		}

	case tea.WindowSizeMsg:
		m.totalWidth = msg.Width
		m.totalHeight = msg.Height
		m.recalculateTable()
	}

	return m, tea.Batch(cmds...)
}

func (m *model) recalculateTable() {
	m.Challenges = m.Challenges.
		WithTargetWidth(m.calculateWidth()).
		WithMinimumHeight(m.calculateHeight()).
		WithPageSize(m.calculateHeight() - 7)
}

func (m model) calculateWidth() int {
	return m.totalWidth - m.horizontalMargin
}

func (m model) calculateHeight() int {
	return m.totalHeight - m.verticalMargin - fixedVerticalMargin 
}

func (m model) View() string {
	body := strings.Builder{}

	body.WriteString(fmt.Sprintf("Target size: %d W x %d H\n", m.calculateWidth(), m.calculateHeight()))
	body.WriteString("Use the arrow keys to navigate, 'esc' and '/' to toggle filtering, Enter to select, q to quit\n\n")
	body.WriteString(m.Challenges.View())
	body.WriteString("\n")

	return body.String()
}

func Run(HTBClient *HTB.Client) string {
	SelectedChallenge = ""
	fmt.Println(`fetching challenges.... If loaded before type "q" to skip loading`)
	// Challenges := getChallenges(HTBClient)
	task := format.Task(func(a any) any {
		if client, ok := a.(*HTB.Client); ok {
			Challenges := getChallenges(client)
			return Challenges
		}
		return nil
	})
	err := format.RunLoading(task, HTBClient)
	if err != nil {
		panic(err)
	}
	Challenges, ok := format.TaskResult.(challengeListData)
	if !ok {
		panic("error checking task result for machines!")
	}
	fmt.Print("\033[H\033[2J")
	

	
	p := tea.NewProgram(NewModel(Challenges), tea.WithAltScreen())

	// fmt.Println("done fetching Challenges!")

	p.Run()
	return SelectedChallenge
}
