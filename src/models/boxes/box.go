package boxes

import (
	"fmt"
	"strconv"
	"time"

	"github.com/Ceald1/HTB-TUI/src/format"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
	HTB "github.com/gubarz/gohtb"
	"github.com/gubarz/gohtb/services/machines"
)


func BoxInfo(box_id string,  HTBClient *HTB.Client) (boxInfo machines.InfoResponse, machineHandle *machines.Handle) {
	id, err := strconv.Atoi(box_id)
	if err != nil {
		panic(err)
	}
	machineHandle = HTBClient.Machines.Machine(id)
	boxInfo, err = machineHandle.Info(ctx)
	
	if err != nil {
		panic(err)
	}
	return boxInfo, machineHandle
}



func BoxInfoMenu(boxInfo machines.InfoResponse, machineHandle *machines.Handle) {

	var boxAction string
	var flag string
	boxInfoData := boxInfo.Data
	boxStatus := lipgloss.NewStyle().Foreground(format.Pink).Render(boxInfoData.InfoStatus)
	var FormInfo = lipgloss.NewStyle().Background(format.BaseBG).Render(fmt.Sprintf(
    "OS: %s\nDifficulty: %s\nBreach Info: %s",
    format.CheckOS(boxInfoData.Os),
    format.CheckDiff(boxInfoData.DifficultyText),
    boxStatus,
))

	var flagInputPlaceholder = lipgloss.NewStyle().Foreground(format.TextBlue).Faint(true).Blink(true).Render("enter flag.. > ")
	huh.NewForm(
		huh.NewGroup(
			huh.NewNote().Title(lipgloss.NewStyle().Foreground(format.TextDefault).Background(format.BaseBG).Padding(1).Render(boxInfoData.Name)).
				Description(FormInfo),
			huh.NewInput().Prompt(flagInputPlaceholder).Title(lipgloss.NewStyle().Foreground(format.TextTitle).Background(format.BaseBG).Render("Submit Flag")).Value(&flag),
			huh.NewSelect[string]().Title(lipgloss.NewStyle().Foreground(format.TextTitle).Background(format.BaseBG).Render("Box Action")).
				Options(
					huh.NewOption(lipgloss.NewStyle().Background(format.BaseBG).Foreground(format.TextDefault).Render("continue"), "none"),
					huh.NewOption(lipgloss.NewStyle().Background(format.BaseBG).Foreground(format.TextYellow).Render("Reset"), "reset"),
					huh.NewOption(lipgloss.NewStyle().Background(format.BaseBG).Foreground(format.TextCyan).Render("Spawn"), "spawn"),
					huh.NewOption(lipgloss.NewStyle().Background(format.BaseBG).Foreground(format.TextPink).Render("Terminate"), "terminate"),
					huh.NewOption(lipgloss.NewStyle().Background(format.BaseBG).Foreground(format.DarkPurple).Render("Extend"), "extend"),
					huh.NewOption(lipgloss.NewStyle().Background(format.BaseBG).Foreground(format.TextDefault).Render("quit"), "quit"),

				).Value(&boxAction),

		),
	).Run()
	switch boxAction {
		default:
			return
		case "none":
			if(flag != "") {
				resp,  err  := machineHandle.Own(ctx, flag)
				if err != nil {
					fmt.Println("unable to submit flag! ", err.Error())
					time.Sleep( 10 * time.Second)
					BoxInfoMenu(boxInfo, machineHandle)
				}else {
					fmt.Println(resp.Data.Message)
					// time.Sleep( 10 * time.Second )
					BoxInfoMenu(boxInfo, machineHandle)
				}
			}else{
				BoxInfoMenu(boxInfo, machineHandle)
			}
		case "quit":
			return
		
		
		case "reset":
			resp, err := machineHandle.Reset(ctx)
			if err != nil {
				fmt.Println("unable to reset! ", err.Error())
				time.Sleep( 10 * time.Second)
				BoxInfoMenu(boxInfo, machineHandle)
			}else {
				fmt.Println(resp.Data.Message)
					// time.Sleep( 10 * time.Second )
				BoxInfoMenu(boxInfo, machineHandle)
			}
		case "spawn":
			resp, err := machineHandle.Spawn(ctx)
			if err != nil {
				fmt.Println("unable to spawn! ", err.Error())
				time.Sleep( 10 * time.Second)
				BoxInfoMenu(boxInfo, machineHandle)
			}else {
				fmt.Println(resp.Data.Message)
				time.Sleep( 10 * time.Second )
				BoxInfoMenu(boxInfo, machineHandle)
			}
		case "terminate":
			resp, err := machineHandle.Terminate(ctx)
			if err != nil {
				fmt.Println("unable to terminate! ", err.Error())
				time.Sleep( 10 * time.Second)
				BoxInfoMenu(boxInfo, machineHandle)
			}else {
				fmt.Println(resp.Data.Message)
				return
				// time.Sleep( 10 * time.Second )
				// BoxInfoMenu(boxInfo, machineHandle)
			}
		case "extend":
			resp, err := machineHandle.Extend(ctx)
			if err != nil {
				fmt.Println("unable to extend! ", err.Error())
				time.Sleep( 10 * time.Second)
				BoxInfoMenu(boxInfo, machineHandle)
			}else {
				fmt.Println(resp.Data.Message)
				time.Sleep( 10 * time.Second )
				BoxInfoMenu(boxInfo, machineHandle)
			}
	}
}