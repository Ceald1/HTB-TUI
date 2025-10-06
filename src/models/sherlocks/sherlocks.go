package sherlocks

import (
	"context"
	"fmt"
	// "strconv"
	"strings"
	// "time"

	"github.com/Ceald1/HTB-TUI/src/format"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
	HTB "github.com/gubarz/gohtb"
	"github.com/gubarz/gohtb/services/sherlocks"
)


var (
	ctx = context.Background()
	LabSelected int
)

func SelectSherlock(HTBClient *HTB.Client) (selected int) {
	var sherlock_resp []sherlocks.SherlockItem
	task := format.Task(func(a any) any {
		if client, ok := a.(*HTB.Client);ok {
			sherlocks, err := client.Sherlocks.List().AllResults(ctx)
			if err != nil {
				panic(err)
			}
			return sherlocks
		}
		panic("other error occurred")
	})
	err := format.RunLoading(task, HTBClient)
	if err != nil {
		panic("error while getting sherlocks!")
	}
	sherlocks_, ok := format.TaskResult.(sherlocks.SherlockListResponse)
	if !ok {
		panic("error while checking sherlock request")
	}
	sherlock_resp = sherlocks_.Data
	var options []huh.Option[int]
	var quit_value = 9999999999999
	quit_op := huh.NewOption(lipgloss.NewStyle().Foreground(format.Red).Background(format.BaseBG).Render("Quit"), quit_value)
	options = append(options, quit_op)
	

	for _, sherlock := range sherlock_resp {
		info := lipgloss.NewStyle().Foreground(format.NextColor()).Background(format.BaseBG).Render(sherlock.Name)
		op := huh.NewOption(info, sherlock.Id)
		options = append(options, op)
	}
	huh.NewSelect[int]().
		Title(lipgloss.NewStyle().Foreground(format.TextTitle).Background(format.BaseBG).Render("Sherlocks")).
		Options(options...).Value(&selected).Run()
	
	switch selected{
		case quit_value:
			return 0
		default:
			return selected

	}
}

func ViewSherlock(HTBClient *HTB.Client, selected int) {
	labData := HTBClient.Sherlocks.Sherlock(selected)
	info, _ := labData.Info(ctx)
	var action string
	var title = lipgloss.NewStyle().Foreground(format.TextTitle).Background(format.BaseBG).Padding(1).Render(info.Data.Name)
	progressResp, _ := labData.Progress(ctx)
	progress := float64(progressResp.Data.Progress)
	description := lipgloss.NewStyle().Render(fmt.Sprintf(
		"Completed: %.0f%% \nCategory: %s \n ", 
			progress, 
			lipgloss.NewStyle().Foreground(format.Pink).Render(strings.TrimSuffix(strings.TrimSpace(format.Sanitize(info.Data.CategoryName)), "\n")), 
	))
	huh.NewForm(
		huh.NewGroup(
			huh.NewNote().Title(title).Description(description),
			huh.NewSelect[string]().Title(lipgloss.NewStyle().Foreground(format.TextTitle).Background(format.BaseBG).Render("Action")).
				Options(
					// huh.NewOption(lipgloss.NewStyle().Background(format.BaseBG).Foreground(format.TextDefault).Render("continue"), "none"),
					// huh.NewOption(lipgloss.NewStyle().Background(format.BaseBG).Foreground(format.TextYellow).Render("View Sherlocks"), "sherlocks"),
					huh.NewOption(lipgloss.NewStyle().Background(format.BaseBG).Foreground(format.TextDefault).Render("quit"), "quit"),
				).Value(&action),
		),
	).Run()
		switch action {
		default:
			return
		case "quit":
			return
		
	}

}