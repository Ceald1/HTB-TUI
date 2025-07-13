package src

// contains all TUIs

import (
	"fmt"
	"os"

	BloodModel "github.com/Ceald1/HTB-TUI/src/models/blood"
	BoxModel "github.com/Ceald1/HTB-TUI/src/models/boxes"

	"github.com/Ceald1/HTB-TUI/src/format"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
	HTB "github.com/gubarz/gohtb"
)

// Initial TUI Prompt
func API_Token() (instance *HTB.Client, err error){
	
	// creating HTB instance and set token
	token := os.Getenv("HTB_TOKEN")
	if len(token) > 1 {
		instance, err = HTB.New(token)
		return
	}

	title := lipgloss.NewStyle().Foreground(format.TextTitle).Padding(1,1,1,1).Background(format.BaseBG).Render(`Please set the "HTB_TOKEN" environment variable to avoid going through this form`)
	prompt := lipgloss.NewStyle().Foreground(format.TextDefault).Render("Enter your HTB Token > ")
	huh.NewInput().
		Title(title).Prompt(prompt).
		Value(&token).Run()
	instance, err = HTB.New(token)
	return
}




func BoxBlood(HTBClient *HTB.Client) {
	BloodModel.Run(HTBClient)
}
func ClearTerminal() {
	fmt.Fprint(os.Stdout, "\033[H\033[2J")
}

func MainMenu(HTBClient *HTB.Client) {
	title := lipgloss.NewStyle().Foreground(format.TextTitle).Padding(1,1,1,1).Background(format.BaseBG).Render(`Main Menu`)
	var option string

	

	huh.NewSelect[string]().
		Title(title).
		Options(
			huh.NewOption(lipgloss.NewStyle().Foreground(format.TextRed).Render("Monitor Bloods"), "blood"), // Bloods
			huh.NewOption(lipgloss.NewStyle().Foreground(format.TextCyan).Render("View All Boxes"), "boxes"), // Bloods


			huh.NewOption(lipgloss.NewStyle().Foreground(format.TextDefault).Render("Quit"), "quit"), // Quit
		).Value(&option).Run()
	
	
	switch option{
		case "blood":
			ClearTerminal()
			BoxBlood(HTBClient)
		case "boxes":
			ClearTerminal()
			BoxModel.Run(HTBClient)
		case "quit":
			os.Exit(0)
			return
		default:
			return
	}
	MainMenu(HTBClient)
}