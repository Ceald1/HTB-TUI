package prolabs

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/Ceald1/HTB-TUI/src/format"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
	HTB "github.com/gubarz/gohtb"
	"github.com/gubarz/gohtb/services/problabs"
)


var (
	ctx = context.Background()
	LabSelected int
)

func SelectProlabs(HTBClient *HTB.Client) (selectedLab int) {
	// get prolabs and do fancy ass loading
	var labresp prolabs.ProlabDataItems
	task := format.Task(func(a any) any {
		if client, ok := a.(*HTB.Client);ok {
			labs, err := client.Prolabs.List(ctx)
			if err != nil {
				panic(err)
			}
			return labs
		}
		panic("other error occurred")
	})
	err := format.RunLoading(task, HTBClient)
	if err != nil {
		panic("Error fetching prolabs")
	}
	labs, ok := format.TaskResult.(prolabs.ListResponse)
	if !ok {
		panic("Error checking typing for prolabs request")
	}
	labresp = labs.Data.Labs

	var options []huh.Option[int]
	var quit_value = 9999999999999
	quit_op := huh.NewOption(lipgloss.NewStyle().Foreground(format.Red).Background(format.BaseBG).Render("Quit"), quit_value)
	options = append(options, quit_op)

	for _, lab := range labresp {
		info := lipgloss.NewStyle().Foreground(format.NextColor()).Background(format.BaseBG).Render(lab.Name)
		op := huh.NewOption(info, lab.Id)
		options = append(options, op)

	}
	huh.NewSelect[int]().
		Title(lipgloss.NewStyle().Foreground(format.TextTitle).Background(format.BaseBG).Render("Pro Labs")).
		Options(options...).Value(&selectedLab).Run()
	
	switch selectedLab{
		case quit_value:
			return 0
		default:
			return selectedLab

	}
}


func ViewProLab(HTBClient *HTB.Client, selectedProlab int) {
	labData := HTBClient.Prolabs.Prolab(selectedProlab)
	info, _ := labData.Info(ctx)
	var action string
	var flagInputPlaceholder = lipgloss.NewStyle().Foreground(format.TextBlue).Faint(true).Blink(true).Render("enter flag.. > ")
	var flag string
	
	flags := strconv.Itoa(info.Data.ProFlagsCount)
	name := format.Sanitize(info.Data.Name)
	var title = lipgloss.NewStyle().Foreground(format.TextTitle).Background(format.BaseBG).Padding(1).Render(name)
	progressResp, _ := labData.Progress(ctx)
	progress := progressResp.Data.Ownership

	description := lipgloss.NewStyle().Render(fmt.Sprintf(
		"Completed: %.0f%% \nDescription: %s \nFlags: %s \n ", 
			progress, 
			lipgloss.NewStyle().Foreground(format.Pink).Render(strings.TrimSuffix(strings.TrimSpace(format.Sanitize(info.Data.Description)), "\n")), 
			lipgloss.NewStyle().Foreground(format.DarkPurple).Render(flags),
	))
	huh.NewForm(
		huh.NewGroup(
			huh.NewNote().Title(title).Description(description),
			huh.NewInput().Prompt(flagInputPlaceholder).Title(lipgloss.NewStyle().Foreground(format.TextTitle).Background(format.BaseBG).Render("Submit Flag")).Value(&flag),
			huh.NewSelect[string]().Title(lipgloss.NewStyle().Foreground(format.TextTitle).Background(format.BaseBG).Render("Action")).
				Options(
					huh.NewOption(lipgloss.NewStyle().Background(format.BaseBG).Foreground(format.TextDefault).Render("continue"), "none"),
					huh.NewOption(lipgloss.NewStyle().Background(format.BaseBG).Foreground(format.TextYellow).Render("View Machines"), "machines"),
					huh.NewOption(lipgloss.NewStyle().Background(format.BaseBG).Foreground(format.TextDefault).Render("quit"), "quit"),
				).Value(&action),
		),

	).Run()
		switch action {
		default:
			return
		case "none":
			if(flag != "") {
				resp,  err  := labData.SubmitFlag(ctx, flag)
				if err != nil {
					fmt.Println("unable to submit flag! ", err.Error())
					time.Sleep( 10 * time.Second)
					ViewProLab(HTBClient, selectedProlab)
				}else {
					fmt.Println(resp.Data.Message)
					time.Sleep( 10 * time.Second )
					ViewProLab(HTBClient, selectedProlab)
				}
			}else{
				ViewProLab(HTBClient, selectedProlab)
			}
		case "quit":
			return
		
		case "machines":
			LabTable(labData)
			ViewProLab(HTBClient, selectedProlab)
		
	}
}

