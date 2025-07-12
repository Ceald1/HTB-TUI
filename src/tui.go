package src

// contains all TUIs

import (
	"os"
	BloodModel "github.com/Ceald1/HTB-TUI/src/modesl/blood"

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

	title := lipgloss.NewStyle().Foreground(TextTitle).Padding(1,1,1,1).Background(BaseBG).Render(`Please set the "HTB_TOKEN" environment variable to avoid going through this form`)
	prompt := lipgloss.NewStyle().Foreground(TextDefault).Render("Enter your HTB Token > ")
	huh.NewInput().
		Title(title).Prompt(prompt).
		Value(&token).Run()
	instance, err = HTB.New(token)
	return
}




func BoxBlood() {
	BloodModel.Run()
}