package format

import (
	"strings"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
	"github.com/microcosm-cc/bluemonday"
	"golang.org/x/term"
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

func HTBTheme()*huh.Theme{
	t := huh.ThemeBase()

	var (
		background = lipgloss.AdaptiveColor{Dark: "#141D2B"}
		selection  = lipgloss.AdaptiveColor{Dark: "#C5D1EB"}
		foreground = lipgloss.AdaptiveColor{Dark: "#A4B1CD"}
		comment    = lipgloss.AdaptiveColor{Dark: "#5CB2FF"}
		green      = lipgloss.AdaptiveColor{Dark: "#C5F467"}
		purple     = lipgloss.AdaptiveColor{Dark: "#A000FF"}
		red        = lipgloss.AdaptiveColor{Dark: "#FF8484"}
		yellow     = lipgloss.AdaptiveColor{Dark: "#FFCC5C"}
	)

	t.Focused.Base = t.Focused.Base.BorderForeground(selection)
	t.Focused.Card = t.Focused.Base
	t.Focused.Title = t.Focused.Title.Foreground(purple)
	t.Focused.NoteTitle = t.Focused.NoteTitle.Foreground(purple)
	t.Focused.Description = t.Focused.Description.Foreground(comment)
	t.Focused.ErrorIndicator = t.Focused.ErrorIndicator.Foreground(red)
	t.Focused.Directory = t.Focused.Directory.Foreground(purple)
	t.Focused.File = t.Focused.File.Foreground(foreground)
	t.Focused.ErrorMessage = t.Focused.ErrorMessage.Foreground(red)
	t.Focused.SelectSelector = t.Focused.SelectSelector.Foreground(yellow)
	t.Focused.NextIndicator = t.Focused.NextIndicator.Foreground(yellow)
	t.Focused.PrevIndicator = t.Focused.PrevIndicator.Foreground(yellow)
	t.Focused.Option = t.Focused.Option.Foreground(foreground)
	t.Focused.MultiSelectSelector = t.Focused.MultiSelectSelector.Foreground(yellow)
	t.Focused.SelectedOption = t.Focused.SelectedOption.Foreground(green)
	t.Focused.SelectedPrefix = t.Focused.SelectedPrefix.Foreground(green)
	t.Focused.UnselectedOption = t.Focused.UnselectedOption.Foreground(foreground)
	t.Focused.UnselectedPrefix = t.Focused.UnselectedPrefix.Foreground(comment)
	t.Focused.FocusedButton = t.Focused.FocusedButton.Foreground(green).Background(purple).Bold(true)
	t.Focused.BlurredButton = t.Focused.BlurredButton.Foreground(foreground).Background(background)

	t.Focused.TextInput.Cursor = t.Focused.TextInput.Cursor.Foreground(yellow)
	t.Focused.TextInput.Placeholder = t.Focused.TextInput.Placeholder.Foreground(comment)
	t.Focused.TextInput.Prompt = t.Focused.TextInput.Prompt.Foreground(yellow)

	t.Blurred = t.Focused
	t.Blurred.Base = t.Blurred.Base.BorderStyle(lipgloss.HiddenBorder())
	t.Blurred.Card = t.Blurred.Base
	t.Blurred.NextIndicator = lipgloss.NewStyle()
	t.Blurred.PrevIndicator = lipgloss.NewStyle()

	t.Group.Title = t.Focused.Title
	t.Group.Description = t.Focused.Description

	return t
}

func Sanitize(inputStr string) (string) {
	p := bluemonday.StripTagsPolicy()
	return p.Sanitize(inputStr)
}

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

func BoxState(state string) (color string) {
	state = strings.ToLower(state)
	switch state {
		case "free","active":
			color = lipgloss.NewStyle().Foreground(TextLightGreen).Render(state)
		case "retired_free":
			color = lipgloss.NewStyle().Foreground(TextYellow).Render(state)
		case "retired":
			color = lipgloss.NewStyle().Foreground(TextPink).Render(state)
		case "unreleased":
			color = lipgloss.NewStyle().Foreground(TextPurple).Render(state)
		default:
			color = lipgloss.NewStyle().Foreground(TextDefault).Render(state)
	}
	return
}


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

func SplitResp() (resp string) {
	width, _, _ := term.GetSize(0)
	resp = strings.Repeat("-", width - 1)
	return resp
}