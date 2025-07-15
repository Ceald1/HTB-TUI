package fortress

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
	"github.com/gubarz/gohtb/services/fortresses"
)

var (
	ctx = context.Background()
	Fortress_selected int
)


func SelectFortresses(HTBClient *HTB.Client) (int) {
	// returns the fortress ID
	task := format.Task(func(a any) any {
		if client, ok := a.(*HTB.Client);ok {
			Fortresses, err := client.Fortresses.List(ctx)
			if err != nil {
				panic(err)
			}
			return Fortresses
		}
		panic("other error occurred")
	})
	err := format.RunLoading(task, HTBClient)
	if err != nil {
		panic(err)
	}
	Fortresses, ok := format.TaskResult.(fortresses.ListResponse)
	if !ok {
		panic("error occurred when checking fortress task result")
	}

	
	var options []huh.Option[int]
	var quit_value = 9999999999999
	quit_op := huh.NewOption(lipgloss.NewStyle().Foreground(format.Red).Background(format.BaseBG).Render("Quit"), quit_value)
	options = append(options, quit_op)
	format.ColorIndex = 0
	for _, f := range Fortresses.Data{
		info := lipgloss.NewStyle().Foreground(format.NextColor()).Background(format.BaseBG).Render(f.Name)
		op := huh.NewOption(info, f.Id)
		options = append(options, op)
	}
	huh.NewSelect[int]().
		Title(lipgloss.NewStyle().Foreground(format.TextTitle).Background(format.BaseBG).Render("Fortresses")).
		Options(options...).Value(&Fortress_selected).Run()
	
	switch Fortress_selected{
		case quit_value:
			Fortress_selected = 0
			return 0
		default:
			return Fortress_selected

	}
}

func ViewFort(HTBClient *HTB.Client, fortressID int){
	fortData := HTBClient.Fortresses.Fortress(fortressID)
	info, _ := fortData.Info(ctx)
	var action string
	var flagInputPlaceholder = lipgloss.NewStyle().Foreground(format.TextBlue).Faint(true).Blink(true).Render("enter flag.. > ")
	var flag string
	
	flags := strconv.Itoa(len(info.Data.Flags))
	var title = lipgloss.NewStyle().Foreground(format.TextTitle).Background(format.BaseBG).Padding(1).Render(info.Data.Name)
	
	description := lipgloss.NewStyle().Render(fmt.Sprintf(
		"Company: %s \nCompleted: %.0f%% \nDescription: %s \nPoints: %s \nFlags: %s \n ", 
			lipgloss.NewStyle().Foreground(format.LightBlue).Render(info.Data.Company.Name), 
			info.Data.ProgressPercent, 
			lipgloss.NewStyle().Foreground(format.Pink).Render(strings.TrimSuffix(info.Data.Description, "\n")), 
			lipgloss.NewStyle().Foreground(format.TextCyan).Render(strconv.Itoa(info.Data.Points)), 
			lipgloss.NewStyle().Foreground(format.DarkPurple).Render(flags),
	))
	huh.NewForm(
		huh.NewGroup(
			huh.NewNote().Title(title).Description(description),
			huh.NewInput().Prompt(flagInputPlaceholder).Title(lipgloss.NewStyle().Foreground(format.TextTitle).Background(format.BaseBG).Render("Submit Flag")).Value(&flag),
			huh.NewSelect[string]().Title(lipgloss.NewStyle().Foreground(format.TextTitle).Background(format.BaseBG).Render("Box Action")).
				Options(
					huh.NewOption(lipgloss.NewStyle().Background(format.BaseBG).Foreground(format.TextDefault).Render("continue"), "none"),
					huh.NewOption(lipgloss.NewStyle().Background(format.BaseBG).Foreground(format.TextYellow).Render("Reset"), "reset"),
					huh.NewOption(lipgloss.NewStyle().Background(format.BaseBG).Foreground(format.TextDefault).Render("quit"), "quit"),
				).Value(&action),
		),

	).Run()
	switch action {
		default:
			return
		case "none":
			if(flag != "") {
				resp,  err  := fortData.SubmitFlag(ctx, flag)
				if err != nil {
					fmt.Println("unable to submit flag! ", err.Error())
					time.Sleep( 10 * time.Second)
					ViewFort(HTBClient, fortressID)
				}else {
					fmt.Println(resp.Data.Message)
					// time.Sleep( 10 * time.Second )
					ViewFort(HTBClient, fortressID)
				}
			}else{
				ViewFort(HTBClient, fortressID)
			}
		case "quit":
			return
		
		
		case "reset":
			resp, err := fortData.Reset(ctx)
			if err != nil {
				fmt.Println("unable to reset! ", err.Error())
				time.Sleep( 10 * time.Second)
				ViewFort(HTBClient, fortressID)
			}else {
				fmt.Println(resp.Data.Message)
					// time.Sleep( 10 * time.Second )
				ViewFort(HTBClient, fortressID)
			}
	}
}