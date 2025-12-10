package users

import (
	"context"
	"fmt"

	"github.com/Ceald1/HTB-TUI/src/format"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
	HTB "github.com/gubarz/gohtb"
	users "github.com/gubarz/gohtb/services/users"
)
var (
	ctx = context.Background()
)

func SearchUser(userName string, HTBClient *HTB.Client) (userID int, err error) {
	results, err := HTBClient.Search.Query(userName).Users(ctx)
	if err != nil {
		return
	}
	for _, result := range results.Data.Users {
		if result.Value == userName {
			userID = result.Id
			return
		}
	}
	return 0, fmt.Errorf("user does not exist!")
}

func getUser(userId int, HTBClient  *HTB.Client) (profile users.ProfileBasicResponse, err error) {
	profile, err = HTBClient.Users.User(userId).ProfileBasic(ctx)
	return
}

func UserForm(userId int, HTBClient *HTB.Client) (err error) { // display basic user info
	profile, err := getUser(userId, HTBClient)
	if err != nil {
		return
	}
	var FormInfo = lipgloss.NewStyle().Background(format.BaseBG).Render(fmt.Sprintf(
    "Country: %s\nPoints: %d\nBloods: %d\nOwns: %d\n",
    profile.Data.CountryName,
    profile.Data.Points,
    profile.Data.ChallengeBloods + profile.Data.UserBloods + profile.Data.SystemBloods,
	profile.Data.SystemOwns + profile.Data.UserOwns,
		),
	)
	huh.NewForm(
		huh.NewGroup(
			huh.NewNote().Title(lipgloss.NewStyle().Foreground(format.TextDefault).Background(format.BaseBG).Padding(1).Render(profile.Data.Name)).Description(FormInfo),
			huh.NewSelect[string]().Options(huh.NewOption(lipgloss.NewStyle().Background(format.BaseBG).Foreground(format.TextDefault).Render("quit"),"")),
		),
	).Run()


	return
}