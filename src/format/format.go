package format

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
)

var (
	// based on https://github.com/silofy/hackthebox/tree/master
	// Base background color
	BaseBG  = lipgloss.Color("#141D2B")

	// Primary palette
	Purple     = lipgloss.Color("#A4B1CD")
	Red        = lipgloss.Color("#FF8484")
	Cyan       = lipgloss.Color("#5CECC6")
	Pink       = lipgloss.Color("#FFB3DE") // replaced
	Yellow     = lipgloss.Color("#FFCC5C")
	LightBlue  = lipgloss.Color("#C5D1EB")
	Blue       = lipgloss.Color("#5CB2FF")
	LightGreen = lipgloss.Color("#C5F467")
	DarkPurple = lipgloss.Color("#A000FF")

	// Optional: Text color for dark backgrounds (use primary palette for accents)
	TextDefault = Purple
	TextRed     = Red
	TextCyan    = Cyan
	TextPink    = Pink
	TextYellow  = Yellow
	TextLightBlue = LightBlue
	TextBlue    = Blue
	TextLightGreen = LightGreen
	TextTitle = LightGreen
	TextPurple = DarkPurple

)
var ColorsBrightToDark = []lipgloss.Color{
	Yellow,     // #FFCC5C
	LightGreen, // #C5F467
	Cyan,       // #5CECC6
	LightBlue,  // #C5D1EB
	Pink,       // #FFB3DE
	Blue,       // #5CB2FF
	Red,        // #FF8484
	Purple,     // #A4B1CD
	DarkPurple, // #A000FF
}
var ColorIndex = 0

func NextColor() lipgloss.Color {
	if ColorIndex >= len(ColorsBrightToDark) {
		ColorIndex = 0
	}
	color := ColorsBrightToDark[ColorIndex]
	ColorIndex++
	return color
}
func CheckOS(BoxOS string) (color string) {
	BoxOS = strings.ToLower(BoxOS)
	switch BoxOS{
		case "linux":
			color = lipgloss.NewStyle().Foreground(TextLightGreen).Render(BoxOS)
		case "windows":
			color = lipgloss.NewStyle().Foreground(TextBlue).Render(BoxOS)
		case "freebsd":
			color = lipgloss.NewStyle().Foreground(TextRed).Render(BoxOS)
		case "openbsd":
			color = lipgloss.NewStyle().Foreground(TextYellow).Render(BoxOS)
		case "other":
			color = lipgloss.NewStyle().Foreground(TextDefault).Render(BoxOS)
	}
	return
}
func CheckDiff(difficulty string) (color string) {
	difficulty = strings.ToLower(difficulty)
	switch difficulty{
		case "easy":
			color = lipgloss.NewStyle().Foreground(TextLightGreen).Render(difficulty)
		case "medium":
			color = lipgloss.NewStyle().Foreground(TextYellow).Render(difficulty)
		case "hard":
			color = lipgloss.NewStyle().Foreground(TextRed).Render(difficulty)
		case "insane":
			color = lipgloss.NewStyle().Foreground(TextPurple).Render(difficulty)
		default:
			color = lipgloss.NewStyle().Foreground(TextLightBlue).Render(difficulty)
		
		
	}
	return
}
// func CheckCategory(categoryName string) (color string){
// 	categoryName = strings.ToLower(categoryName)
// 	switch categoryName{
// 		case 
// 	}

// 	return
// }




type Task func(any) any

type loading struct {
	model spinner.Model
	task  tea.Cmd
}

type doneMSG struct {
	result any
}

var TaskResult any

func HelpTask(task Task, args any) tea.Cmd {
	return func() tea.Msg {
		result := task(args)
		return doneMSG{result: result}
	}
}

func InitialLoadingModel(task Task, args any) loading {
	s := spinner.New()
	s.Spinner = spinner.Meter
	return loading{
		model: s,
		task:  HelpTask(task, args),
	}
}

func (m loading) Init() tea.Cmd {
	return tea.Batch(
		m.model.Tick,
		m.task,
	)
}

func (m loading) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "esc", "ctrl+c":
			return m, tea.Quit
		}
	case doneMSG:
		TaskResult = msg.result
		return m, tea.Quit
	default:
		var cmd tea.Cmd
		m.model, cmd = m.model.Update(msg)
		return m, cmd
	}
	return  m, nil
}

func (m loading) View() string {
	str := lipgloss.NewStyle().Foreground(TextTitle).Render(m.model.View())
	return str
}

func RunLoading(task Task, args any) (err error) {
	p := tea.NewProgram(InitialLoadingModel(task, args))
	if _, err = p.Run(); err != nil {
		return err
	}
	return
}