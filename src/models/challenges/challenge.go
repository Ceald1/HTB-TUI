package challenges

import (
	"fmt"
	"strconv"
	"time"

	"github.com/Ceald1/HTB-TUI/src/format"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
	HTB "github.com/gubarz/gohtb"
	"github.com/gubarz/gohtb/services/challenges"
)


func ChallengeInfo(box_id string,  HTBClient *HTB.Client) (ChallengeInfo challenges.InfoResponse, challengeHandle *challenges.Handle) {
	id, err := strconv.Atoi(box_id)
	if err != nil {
		panic(err)
	}
	challengeHandle = HTBClient.Challenges.Challenge(id)
	ChallengeInfo, err = challengeHandle.Info(ctx)
	
	if err != nil {
		panic(err)
	}
	return ChallengeInfo, challengeHandle
}



func ChallengeInfoMenu(ChallengeInfo challenges.InfoResponse, challengeHandle *challenges.Handle) {

	var action string
	var flag string
	ChallengeInfoData := ChallengeInfo.Data
	var FormInfo = lipgloss.NewStyle().Render(fmt.Sprintf(
    "Category: %s\nDifficulty: %s\nDescription: %s",
    ChallengeInfoData.CategoryName,
    format.CheckDiff(ChallengeInfoData.Difficulty),
	lipgloss.NewStyle().Foreground(format.Pink).Render(ChallengeInfoData.Description),
))

	var flagInputPlaceholder = lipgloss.NewStyle().Foreground(format.TextBlue).Faint(true).Blink(true).Render("enter flag.. > ")
	huh.NewForm(
		huh.NewGroup(
			huh.NewNote().Title(lipgloss.NewStyle().Foreground(format.TextTitle).Background(format.BaseBG).Padding(1,1).Render(ChallengeInfoData.Name)).
				Description(FormInfo),
			huh.NewInput().Prompt(flagInputPlaceholder).Title(lipgloss.NewStyle().Foreground(format.TextTitle).Background(format.BaseBG).Render("Submit Flag")).Value(&flag),
			huh.NewSelect[string]().Title(lipgloss.NewStyle().Foreground(format.TextTitle).Background(format.BaseBG).Render("Action")).
				Options(
					huh.NewOption(lipgloss.NewStyle().Background(format.BaseBG).Foreground(format.TextDefault).Render("continue"), "none"),
					huh.NewOption(lipgloss.NewStyle().Background(format.BaseBG).Foreground(format.TextCyan).Render("Spawn"), "spawn"),
					huh.NewOption(lipgloss.NewStyle().Background(format.BaseBG).Foreground(format.TextPink).Render("Terminate"), "terminate"),
					huh.NewOption(lipgloss.NewStyle().Background(format.BaseBG).Foreground(format.TextDefault).Render("quit"), "quit"),

				).Value(&action),

		),
	).Run()
	switch action {
		default:
			return
		case "none":
			if(flag != "") {
				resp,  err  := challengeHandle.Own(ctx, flag)
				if err != nil {
					fmt.Println("unable to submit flag! ", err.Error())
					time.Sleep( 10 * time.Second)
					ChallengeInfoMenu(ChallengeInfo, challengeHandle)
				}else {
					fmt.Println(resp.Data.Message)
					// time.Sleep( 10 * time.Second )
					ChallengeInfoMenu(ChallengeInfo, challengeHandle)
				}
			}else{
				ChallengeInfoMenu(ChallengeInfo, challengeHandle)
			}
		case "quit":
			return
		
		

		case "spawn":
			resp, err := challengeHandle.Start(ctx)
			if err != nil {
				fmt.Println("unable to spawn! ", err.Error())
				time.Sleep( 10 * time.Second)
				ChallengeInfoMenu(ChallengeInfo, challengeHandle)
			}else {
				fmt.Println(resp.Data.Message)
				time.Sleep( 10 * time.Second )
				ChallengeInfoMenu(ChallengeInfo, challengeHandle)
			}
		case "terminate":
			resp, err := challengeHandle.Stop(ctx)
			if err != nil {
				fmt.Println("unable to terminate! ", err.Error())
				time.Sleep( 10 * time.Second)
				ChallengeInfoMenu(ChallengeInfo, challengeHandle)
			}else {
				fmt.Println(resp.Data.Message)
				return
				// time.Sleep( 10 * time.Second )
				// ChallengeInfoMenu(ChallengeInfo, challengeHandle)
			}
	}
}