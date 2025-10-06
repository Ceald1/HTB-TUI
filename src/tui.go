package src

// contains all TUIs

import (
	"fmt"
	"os"
	"strconv"

	BloodModel "github.com/Ceald1/HTB-TUI/src/models/blood"
	BoxModel "github.com/Ceald1/HTB-TUI/src/models/boxes"
	ChallengeModel "github.com/Ceald1/HTB-TUI/src/models/challenges"
	Fortress "github.com/Ceald1/HTB-TUI/src/models/fortress"
	ProlabModel "github.com/Ceald1/HTB-TUI/src/models/prolabs"
	SherlockModel "github.com/Ceald1/HTB-TUI/src/models/sherlocks"

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
			huh.NewOption(lipgloss.NewStyle().Foreground(format.TextRed).Background(format.BaseBG).Render("Monitor Bloods"), "blood"), // Bloods
			huh.NewOption(lipgloss.NewStyle().Foreground(format.TextCyan).Background(format.BaseBG).Render("View All Boxes"), "boxes"), // Boxes
			huh.NewOption(lipgloss.NewStyle().Foreground(format.TextYellow).Background(format.BaseBG).Render("View All Challenges"), "challenge"), // Challenges
			huh.NewOption(lipgloss.NewStyle().Foreground(format.DarkPurple).Background(format.BaseBG).Render("View Fortresses"), "fortress"), // Fortresses
			huh.NewOption(lipgloss.NewStyle().Foreground(format.LightGreen).Background(format.BaseBG).Render("View Pro Labs"), "prolabs"), // Pro Labs
			huh.NewOption(lipgloss.NewStyle().Foreground(format.Blue).Background(format.BaseBG).Render("View Sherlocks"), "sherlocks"), // Sherlocks

			huh.NewOption(lipgloss.NewStyle().Foreground(format.TextDefault).Background(format.BaseBG).Render("Quit"), "quit"), // Quit
		).Value(&option).Run()
	
	
	switch option{
		case "blood":
			ClearTerminal()
			BoxBlood(HTBClient)
		case "boxes":
			var box_selected = "1"
			for box_selected != ""{
				ClearTerminal()
				BoxModel.Run(HTBClient)
				box_selected = BoxModel.SelectedBox
				if box_selected == ""{
					break
				}
				boxInfo, machineHandle := BoxModel.BoxInfo(box_selected, HTBClient)
				BoxModel.BoxInfoMenu(boxInfo, machineHandle)
			}
		case "fortress":
			var Fortress_selected = 1
			for Fortress_selected != 0{
				ClearTerminal()
				Fortress_selected = Fortress.SelectFortresses(HTBClient)
				if Fortress_selected == 0 {
					break
				}
				Fortress.ViewFort(HTBClient, Fortress_selected)
			}
		case "challenge":
			var challenge_selected = "1"
			for challenge_selected != ""{
				ClearTerminal()
				ChallengeModel.Run(HTBClient)
				challenge_selected = ChallengeModel.SelectedChallenge
				if challenge_selected == ""{
					break
				}
				challengeInfo, machineHandle := ChallengeModel.ChallengeInfo(challenge_selected, HTBClient)
				ChallengeModel.ChallengeInfoMenu(challengeInfo, machineHandle)
			}

		case "prolabs":
			var prolabSelected = 1
			for prolabSelected != 0{
				ClearTerminal()
				prolabSelected = ProlabModel.SelectProlabs(HTBClient)
				if prolabSelected == 0 {
					break
				}
				ProlabModel.ViewProLab(HTBClient, prolabSelected)
			}
		
		case "sherlocks":
			var sherlock_selected = "1"
			var sherlock_selected_i int
			for sherlock_selected != ""{
				ClearTerminal()
				SherlockModel.Run(HTBClient)
				sherlock_selected = SherlockModel.Selected
				if sherlock_selected == ""{
					break
				}
				sherlock_selected_i, _ = strconv.Atoi(sherlock_selected)
				SherlockModel.ViewSherlock(HTBClient, sherlock_selected_i)
			}

		
		case "quit":
			os.Exit(0)
			return
		default:
			return
	}
	MainMenu(HTBClient)
}