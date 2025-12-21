package seasons

import (
	"fmt"
	"time"

	"context"
	"github.com/Ceald1/HTB-TUI/src/format"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
	HTB "github.com/gubarz/gohtb"
)

var (
	ctx = context.Background()
)


func SeasonalMachine(HTBClient *HTB.Client) {
	var err error
	machines, err := HTBClient.Machines.List().PerPage(1).Page(1).First(ctx)
	

	if err != nil {
		panic(err)
	}
	machineID := machines.Data[0].Id
	machine, _ := HTBClient.Machines.Machine(machineID).Info(ctx)
	boxStatus := lipgloss.NewStyle().Foreground(format.Pink).Render(machine.Data.InfoStatus)

	machineHandle := HTBClient.Machines.Machine(machine.Data.Id)
	ip := machine.Data.Ip
	var flagInputPlaceholder = lipgloss.NewStyle().Foreground(format.TextBlue).Faint(true).Blink(true).Render("enter flag.. > ")
	var FormInfo = lipgloss.NewStyle().Background(format.BaseBG).Render(fmt.Sprintf(
		"%s \nIP: %s \nOS: %s \nDifficulty: %s \nBreach Info: %s",
		format.LoadImage(machine.Data.Avatar),
		ip,
		format.CheckOS(machine.Data.Os),
		format.CheckDiff(machine.Data.DifficultyText),
		boxStatus,
	),
	)
	var boxAction string
	var flag string

	huh.NewForm(
		huh.NewGroup(
			huh.NewNote().Title(lipgloss.NewStyle().Foreground(format.TextDefault).Background(format.BaseBG).Padding(1).Render(machine.Data.Name)).
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
		if flag != "" {
			resp, err := machineHandle.Own(ctx, flag)
			if err != nil {
				fmt.Println("unable to submit flag! ", err.Error())
				time.Sleep(10 * time.Second)
				SeasonalMachine(HTBClient)
			} else {
				fmt.Println(resp.Data.Message)
				// time.Sleep( 10 * time.Second )
				SeasonalMachine(HTBClient)
			}
		} else {
			SeasonalMachine(HTBClient)
		}
	case "quit":
		return

	case "reset":
		resp, err := machineHandle.Reset(ctx)
		if err != nil {
			fmt.Println("unable to reset! ", err.Error())
			time.Sleep(10 * time.Second)
			SeasonalMachine(HTBClient)
		} else {
			fmt.Println(resp.Data.Message)
			// time.Sleep( 10 * time.Second )
			SeasonalMachine(HTBClient)
		}
	case "spawn":
		resp, err := machineHandle.Spawn(ctx)
		if err != nil {
			fmt.Println("unable to spawn! ", err.Error())
			time.Sleep(10 * time.Second)
			SeasonalMachine(HTBClient)
		} else {
			fmt.Println(resp.Data.Message)
			// time.Sleep( 10 * time.Second )
			SeasonalMachine(HTBClient)
		}
	case "terminate":
		resp, err := machineHandle.Terminate(ctx)
		if err != nil {
			fmt.Println("unable to terminate! ", err.Error())
			time.Sleep(10 * time.Second)
			SeasonalMachine(HTBClient)
		} else {
			fmt.Println(resp.Data.Message)
			return
			// time.Sleep( 10 * time.Second )
			// BoxInfoMenu(boxInfo, machineHandle)
		}
	case "extend":
		resp, err := machineHandle.Extend(ctx)
		if err != nil {
			fmt.Println("unable to extend! ", err.Error())
			time.Sleep(10 * time.Second)
			SeasonalMachine(HTBClient)
		} else {
			fmt.Println(resp.Data.Message)
			// time.Sleep( 10 * time.Second )
			SeasonalMachine(HTBClient)
		}
	}
}

