package leaderboards
import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/evertras/bubble-table/table"
	"github.com/Ceald1/HTB-TUI/src/format"
)
// Make a table for ranks
const (
	ColumnNum = "#"
	ColumnName = "name"
	ColumnBloods = "bloods"
	ColumnPoints = "points"
	fixedVerticalMargin = 4

)
type TableDataRow struct {
	Name string
	Num string
	Bloods string
	Points string
}

type model struct {
	RankingData 	table.Model
	Selected    	string
	totalWidth		int
	totalHeight 	int
	horizontalMargin  int
	verticalMargin    int
}

func genRows(rankItems []TableDataRow) (rows []table.Row) {
	for _, rankItem := range rankItems {
		rows = append(rows, table.NewRow(table.RowData{
			ColumnNum: rankItem.Num,
			ColumnName: rankItem.Name,
			ColumnBloods: rankItem.Bloods,
			ColumnPoints: rankItem.Points,
		}))
	}
	return rows
}
func NewModel(rankItems []TableDataRow) model {
	columns := []table.Column{
		table.NewFlexColumn(ColumnNum, lipgloss.NewStyle().Foreground(format.TextYellow).Render("#"), 1).WithFiltered(true),
		table.NewFlexColumn(ColumnName, lipgloss.NewStyle().Foreground(format.TextDefault).Render("Name"), 10).WithFiltered(true),
		table.NewFlexColumn(ColumnBloods, lipgloss.NewStyle().Foreground(format.TextRed).Render("Bloods"), 10).WithFiltered(true),
		table.NewFlexColumn(ColumnPoints, lipgloss.NewStyle().Foreground(format.TextLightGreen).Render("Points"), 10).WithFiltered(true),
	}
	rows := genRows(rankItems)

	keys := table.DefaultKeyMap()
	keys.RowDown.SetKeys("j", "down", "s")
	keys.RowUp.SetKeys("k", "up", "w")

	m := model{
		RankingData: table.New(columns).
			WithRows(rows).
			// SelectableRows(true).
			Filtered(true).
			Focused(true).
			WithKeyMap(keys).
			WithPageSize(10).
			WithBaseStyle(lipgloss.NewStyle().Background(format.BaseBG)),
			// WithSelectedText(" ", "✓"),

	}

	m.updateFooter()
	return m
}



func (m model) Init() tea.Cmd {
	return tea.WindowSize()
}


func (m *model) updateFooter() {
	highlightedRow := m.RankingData.HighlightedRow()

	footerText := fmt.Sprintf(
		"Pg. %d/%d - Currently looking at: %s",
		m.RankingData.CurrentPage(),
		m.RankingData.MaxPages(),
		highlightedRow.Data[ColumnName],
	)

	m.RankingData = m.RankingData.WithStaticFooter(footerText)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	m.RankingData, cmd = m.RankingData.Update(msg)
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
			m.RankingData.Filtered(false)
		case "/":
			m.RankingData.Filtered(true)

		}

	case tea.WindowSizeMsg:
		m.totalWidth = msg.Width
		m.totalHeight = msg.Height
		m.recalculateTable()
	}

	return m, tea.Batch(cmds...)
}

func (m *model) recalculateTable() {
	m.RankingData = m.RankingData.
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
	body.WriteString(m.RankingData.View())
	body.WriteString("\n")

	return body.String()
}

func RunRankTable(rankData []TableDataRow) (err error){
	p := tea.NewProgram(NewModel(rankData), tea.WithAltScreen())

	_, err = p.Run()
	return
}