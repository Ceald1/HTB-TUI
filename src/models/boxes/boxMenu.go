package boxes

import (
	"context"
	"fmt"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/evertras/bubble-table/table"
	HTB "github.com/gubarz/gohtb"
	"github.com/gubarz/gohtb/services/machines"
	"github.com/Ceald1/HTB-TUI/src/format"
)

var (
	ctx         = context.Background()
	SelectedBox = ""
)

const (
	ColumnName     = "name"
	ColumnDifficulty = "difficulty"
	ColumnOS       = "os"
	ColumnStatus   = "status"
	ColumnBoxID    = "id"

	minWidth            = 30
	minHeight           = 8
	fixedVerticalMargin = 4
)

type model struct {
	Boxes             table.Model
	Selected          string
	totalWidth        int
	totalHeight       int
	horizontalMargin  int
	verticalMargin    int
}

type machineListData struct {
	Active     machines.MachineDataItems
	Retired    machines.MachineDataItems
	Unreleased machines.UnreleasedDataItems
}

func getBoxes(HTBClient *HTB.Client) (machineList machineListData) {
	activeMachines, err := HTBClient.Machines.ListActive().AllResults(ctx)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	retiredMachines, err := HTBClient.Machines.ListRetired().AllResults(ctx)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	unreleased, err := HTBClient.Machines.ListUnreleased().AllResults(ctx)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	machineList.Active = activeMachines.Data
	machineList.Retired = retiredMachines.Data
	machineList.Unreleased = unreleased.Data
	return machineList
}

func NewModel(Boxes machineListData) model {
	columns := []table.Column{
		table.NewColumn(ColumnName, lipgloss.NewStyle().Foreground(format.TextTitle).Render("Name"), 10).WithFiltered(true),
		table.NewFlexColumn(ColumnOS, lipgloss.NewStyle().Foreground(format.TextTitle).Render("OS"), 1).WithFiltered(true),
		table.NewFlexColumn(ColumnDifficulty, lipgloss.NewStyle().Foreground(format.TextTitle).Render("Difficulty"), 1).WithFiltered(true),
		table.NewFlexColumn(ColumnStatus, lipgloss.NewStyle().Foreground(format.TextTitle).Render("Box Status"), 1).WithFiltered(true),
		table.NewFlexColumn(ColumnBoxID, lipgloss.NewStyle().Foreground(format.TextTitle).Render("Box ID"), 1).WithFiltered(true),
	}

	rows := genRows(Boxes)

	keys := table.DefaultKeyMap()
	keys.RowDown.SetKeys("j", "down", "s")
	keys.RowUp.SetKeys("k", "up", "w")

	m := model{
		Boxes: table.New(columns).
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
	highlightedRow := m.Boxes.HighlightedRow()

	footerText := fmt.Sprintf(
		"Pg. %d/%d - Currently looking at: %s",
		m.Boxes.CurrentPage(),
		m.Boxes.MaxPages(),
		highlightedRow.Data[ColumnName],
	)

	m.Boxes = m.Boxes.WithStaticFooter(footerText)
}

func genRows(boxes machineListData) (rows []table.Row) {
	for _, box := range boxes.Active {
		rows = append(rows, table.NewRow(table.RowData{
			ColumnName:     box.Name,
			ColumnOS:       format.CheckOS(box.Os),
			ColumnDifficulty: format.CheckDiff(box.DifficultyText),
			ColumnStatus:   lipgloss.NewStyle().Foreground(format.TextLightGreen).Render("active"),
			ColumnBoxID:    fmt.Sprintf("%d", box.Id),
		}))
	}
	for _, box := range boxes.Retired {
		rows = append(rows, table.NewRow(table.RowData{
			ColumnName:     box.Name,
			ColumnOS:       format.CheckOS(box.Os),
			ColumnDifficulty: format.CheckDiff(box.DifficultyText),
			ColumnStatus:   lipgloss.NewStyle().Foreground(format.Purple).Render("retired"),
			ColumnBoxID:    fmt.Sprintf("%d", box.Id),
		}))
	}
	for _, box := range boxes.Unreleased {
		rows = append(rows, table.NewRow(table.RowData{
			ColumnName:     box.Name,
			ColumnOS:       format.CheckOS(box.Os),
			ColumnDifficulty: format.CheckDiff(box.DifficultyText),
			ColumnStatus:   lipgloss.NewStyle().Foreground(format.TextPink).Render("unreleased"),
			ColumnBoxID:    fmt.Sprintf("%d", box.Id),
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

	m.Boxes, cmd = m.Boxes.Update(msg)
	cmds = append(cmds, cmd)

	m.updateFooter()

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			cmds = append(cmds, tea.Quit)

		case "enter", "return":
			SelectedBox = m.Boxes.HighlightedRow().Data[ColumnBoxID].(string)
			cmds = append(cmds, tea.Quit)
		case "esc":
			m.Boxes.Filtered(false)
		case "/":
			m.Boxes.Filtered(true)

		}

	case tea.WindowSizeMsg:
		m.totalWidth = msg.Width
		m.totalHeight = msg.Height
		m.recalculateTable()
	}

	return m, tea.Batch(cmds...)
}

func (m *model) recalculateTable() {
	m.Boxes = m.Boxes.
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
	body.WriteString(m.Boxes.View())
	body.WriteString("\n")

	return body.String()
}

func Run(HTBClient *HTB.Client) string {
	SelectedBox = ""
	fmt.Printf("fetching boxes....\n")
	// boxes := getBoxes(HTBClient)
	task := format.Task(func(a any) any {
		if client, ok := a.(*HTB.Client); ok {
			boxes := getBoxes(client)
			return boxes
		}
		return nil
	})
	err := format.RunLoading(task, HTBClient)
	if err != nil {
		panic(err)
	}
	boxes, _ := format.TaskResult.(machineListData)
	fmt.Print("\033[H\033[2J")
	

	
	
	p := tea.NewProgram(NewModel(boxes), tea.WithAltScreen())

	// fmt.Println("done fetching boxes!")

	p.Run()
	return SelectedBox
}
