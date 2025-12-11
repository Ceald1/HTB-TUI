package leaderboards

import (
	"context"
	"fmt"

	"github.com/Ceald1/HTB-TUI/src/format"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
	HTB "github.com/gubarz/gohtb"
	rank "github.com/gubarz/gohtb/services/rankings"
	users "github.com/Ceald1/HTB-TUI/src/models/social"
)

var (
	ctx = context.Background()
)

func userLeaderBoards(HTBClient *HTB.Client) {
	Selected_item = 0
	var RankData []TableDataRow
	task := format.Task(func(a any) any {
		if client, ok := a.(*HTB.Client); ok {
			rankData, err := client.Rankings.Users(ctx)
			if err != nil {
				panic(err)
			}
			return rankData
		}
		return nil
	})
	err := format.RunLoading(task, HTBClient)
	if err != nil {
		panic(err)
	}
	UserResponseItems := format.TaskResult.(rank.UserRankingsResponse).Data
	for _, user := range UserResponseItems {
		
		row := TableDataRow{
			Name: user.Name,
			Num: lipgloss.NewStyle().Foreground(format.LightGreen).Render(fmt.Sprintf("%d",user.Rank)),
			Bloods: lipgloss.NewStyle().Foreground(format.TextRed).Render(fmt.Sprintf("%d",user.RootBloods + user.UserBloods + user.ChallengeBloods)),
			Points: lipgloss.NewStyle().Foreground(format.TextYellow).Render(fmt.Sprintf("%d",user.Points)),
			ID: user.Id,
		}
		RankData = append(RankData, row)
	}
	err = RunRankTable(RankData)
	if err != nil {
		panic(err)
	}
	userId := Selected_item
	// userId, err := users.SearchUser(Selected_item, HTBClient)
	// if err != nil {
	// 		panic(err)
	// }
	if userId != 0 {
		err = users.UserForm(userId, HTBClient)
		if err != nil {
			panic(err)
		}
	}
}


func teamLeaderBoards(HTBClient *HTB.Client) {
	Selected_item = 0
	var RankData []TableDataRow
	var teamId = 0
	task := format.Task(func(a any) any {
		if client, ok := a.(*HTB.Client); ok {
			rankData, err := client.Rankings.Teams(ctx)
			if err != nil {
				panic(err)
			}
			return rankData
		}
		return nil
	})
	err := format.RunLoading(task, HTBClient)
	if err != nil {
		panic(err)
	}
	TeamResponseItems := format.TaskResult.(rank.TeamRankingsResponse).Data
	for _, Team := range TeamResponseItems {
		row := TableDataRow{
			Name: Team.Name,
			Num: lipgloss.NewStyle().Foreground(format.LightGreen).Render(fmt.Sprintf("%d",Team.Rank)),
			Bloods: lipgloss.NewStyle().Foreground(format.TextRed).Render(fmt.Sprintf("%d",Team.RootBloods + Team.UserBloods + Team.ChallengeBloods)),
			Points: lipgloss.NewStyle().Foreground(format.TextYellow).Render(fmt.Sprintf("%d",Team.Points)),
			ID: Team.Id,
		}
		RankData = append(RankData, row)
	}
	err = RunRankTable(RankData)
	if err != nil {
		panic(err)
	}
	teamId = Selected_item
	if teamId != 0 {
		err = users.TeamForm(teamId, HTBClient)
		if err != nil {
			panic(err)
		}
	}

}

func countryLeaderBoards(HTBClient *HTB.Client) {
	Selected_item = 0
	var RankData []TableDataRow
	task := format.Task(func(a any) any {
		if client, ok := a.(*HTB.Client); ok {
			rankData, err := client.Rankings.Countries(ctx)
			if err != nil {
				panic(err)
			}
			return rankData
		}
		return nil
	})
	err := format.RunLoading(task, HTBClient)
	if err != nil {
		panic(err)
	}
	CountryResponseItems := format.TaskResult.(rank.CountryRankingsResponse).Data
	for _, Country := range CountryResponseItems {
		row := TableDataRow{
			Name: Country.Name,
			Num: lipgloss.NewStyle().Foreground(format.LightGreen).Render(fmt.Sprintf("%d",Country.Rank)),
			Bloods: lipgloss.NewStyle().Foreground(format.TextRed).Render(fmt.Sprintf("%d",Country.RootBloods + Country.UserBloods)),
			Points: lipgloss.NewStyle().Foreground(format.TextYellow).Render(fmt.Sprintf("%d",Country.Points)),
			ID: 0,
		}
		RankData = append(RankData, row)
	}
	err = RunRankTable(RankData)
	if err != nil {
		panic(err)
	}
}


func Rankings(HTBClient *HTB.Client) int{
	var options []huh.Option[string]
	var selected string
	var quit_value = "9999999999999"
	quit_op := huh.NewOption(lipgloss.NewStyle().Foreground(format.Red).Background(format.BaseBG).Render("Quit"), quit_value)
	options = append(options, quit_op)
	string_options := []string{"Countries", "Users", "Teams"}

	for _, str := range string_options {
		info := lipgloss.NewStyle().Foreground(format.NextColor()).Background(format.BaseBG).Render(str)
		option :=  huh.NewOption(info, str)
		options = append(options, option)
	}
	huh.NewSelect[string]().
		Title(lipgloss.NewStyle().Foreground(format.TextTitle).Background(format.BaseBG).Render("Leaderboards")).
		Options(options...).Value(&selected).Run()

	switch selected{
		case quit_value:
			return 9999999999999
		case "Countries":
			countryLeaderBoards(HTBClient)
			return 0
		case "Users":
			userLeaderBoards(HTBClient)
			return 0
		case "Teams":
			teamLeaderBoards(HTBClient)
			return 0
		default:
			return 0
	}
}