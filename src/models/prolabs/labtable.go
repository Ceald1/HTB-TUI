package prolabs

import (
	"fmt"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/evertras/bubble-table/table"
	"github.com/gubarz/gohtb/services/prolabs"
	"github.com/Ceald1/HTB-TUI/src/format"
)

const (
	ColumnName     = "name"
	ColumnOS       = "os"
	ColumnBoxID	   = "id"
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

func getBoxes(labData *prolabs.Handle) (labs []prolabs.Machine){
	labResp, err := labData.Machines(ctx)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	labs = labResp.Data
	return
}

func genRows(boxes []prolabs.Machine) (rows []table.Row) {
	for _, box := range boxes {
		rows = append(rows, table.NewRow(table.RowData{
			ColumnName:     box.Name,
			ColumnOS:       format.CheckOS(box.Os),
			ColumnBoxID:    fmt.Sprintf("%d", box.Id),
		}))
	}
	return rows
}
func NewModel(Boxes []prolabs.Machine) model {
	columns := []table.Column{
		table.NewColumn(ColumnName, lipgloss.NewStyle().Foreground(format.TextTitle).Render("Name"), 30).WithFiltered(true),
		table.NewFlexColumn(ColumnOS, lipgloss.NewStyle().Foreground(format.TextTitle).Render("OS"), 1).WithFiltered(true),
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

		// case "enter", "return":
		// 	SelectedBox = m.Boxes.HighlightedRow().Data[ColumnBoxID].(string)
		// 	cmds = append(cmds, tea.Quit)
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
	body.WriteString("Use the arrow keys to navigate, 'esc' and '/' to toggle filtering, q or esc to quit\n\n")
	body.WriteString(m.Boxes.View())
	body.WriteString("\n")

	return body.String()
}
func LabTable(labData *prolabs.Handle) {
	fmt.Println(`fetching boxes.... `)
	// boxes := getBoxes(HTBClient)
	task := format.Task(func(a any) any {
		if client, ok := a.(*prolabs.Handle); ok {
			boxes := getBoxes(client)
			return boxes
		}
		return nil
	})
	err := format.RunLoading(task, labData)
	if err != nil {
		panic(err)
	}
	boxes, ok := format.TaskResult.([]prolabs.Machine)
	if !ok {
		panic("error checking task result for machines!")
	}
	fmt.Print("\033[H\033[2J")
	

	
	p := tea.NewProgram(NewModel(boxes), tea.WithAltScreen())


	p.Run()
}