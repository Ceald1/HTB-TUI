package yaml


// TODO: make outputs pretty
import (
	"strings"

	HTB "github.com/gubarz/gohtb"
	// "github.com/charmbracelet/lipgloss"
	// "github.com/Ceald1/HTB-TUI/src/format"
	"context"
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)


var (
	ctx = context.Background()
)

type Config struct {
	Global  GlobalSection `yaml:"global"`
	Actions []Action      `yaml:"action"`
}
type GlobalSection struct {
	Token string `yaml:"token"`
}

type Action struct {
	Type string
	Data any
}
func (a *Action) UnmarshalYAML(value *yaml.Node) error {
	if len(value.Content) != 2 {
		return fmt.Errorf("invalid action format")
	}

	keyNode := value.Content[0]
	valNode := value.Content[1]

	switch keyNode.Value {
	case "flagSubmit":
		var fs FlagSubmit
		if err := valNode.Decode(&fs); err != nil {
			return err
		}
		a.Type = "flagSubmit"
		a.Data = fs
	case "info":
		var info Info
		if err := valNode.Decode(&info); err != nil {
			return err
		}
		a.Type = "info"
		a.Data = info
	default:
		return fmt.Errorf("unknown action type: %s", keyNode.Value)
	}
	return nil
}



type FlagSubmit struct {
	BoxID int `yaml:"boxID"`
	BoxName string `yaml:"box"`
	Challenge string `yaml:"challenge"`
	ChallengeID int `yaml:"challengeID"`
	Fortress 	string 	`yaml:"fortress"`
	FortressID 	int 	`yaml:"fortressID"`
	Flag	string `yaml:"flag"`
}
type Info struct {
	BoxID int `yaml:"boxID"`
	BoxName string `yaml:"box"`
	Challenge string `yaml:"challenge"`
	ChallengeID int `yaml:"challengeID"`
	Fortress 	string 	`yaml:"fortress"`
	FortressID 	int 	`yaml:"fortressID"`	
}



func RunAutomation(yaml_file string){
	// parse yaml
	var HTBClient *HTB.Client
	f_data, err := os.ReadFile(yaml_file)
	if err != nil {
		panic("File not found!")
	}
	yaml_file = os.ExpandEnv(string(f_data))
	var cfg Config
	yaml.Unmarshal([]byte(yaml_file), &cfg)
	var token = cfg.Global.Token
	if token == "" {
		panic("Token is required! Use an environment variable")
	}
	HTBClient, err = HTB.New(token)
	if err != nil {
		panic(err)
	}
	actions := cfg.Actions
	for _, action := range actions {

		switch action.Type{
			case "flagSubmit":
				fs := action.Data.(FlagSubmit)
				flag := fs.Flag
				content_type := ""
				boxID := fs.BoxID
				boxName := fs.BoxName
				ChallengeID := fs.ChallengeID
				ChallengeName := fs.Challenge
				fortressID := fs.FortressID
				fortressName := fs.Fortress
				if ( boxID == 0 && ChallengeID == 0 && fortressID == 0){ // grab based on name
					if boxName != ""{
						boxID, err = GetMachineID(boxName, *HTBClient)
						if err != nil {
							fmt.Println(err.Error())
						}
						content_type = "box"
					}
					if ChallengeName != "" && content_type == ""{
						ChallengeID, err = GetChallengeID(ChallengeName, *HTBClient)
						if err != nil {
							fmt.Println(err.Error())
						}
						content_type = "challenge"
					}
					if fortressName != "" && content_type == ""{
						fortressID, err = GetFortressID(fortressName, *HTBClient)
						if err != nil {
							fmt.Println(err.Error())
						}
						content_type = "fortress"
					}
				} // end of grab based on name
				switch content_type{ // submit flag
					case "box":
						Handle := HTBClient.Machines.Machine(boxID)
						resp, err := Handle.Own(ctx, flag)
						if err != nil {
							fmt.Println(err.Error())
						}else{
							fmt.Println(resp)
						}
					case "challenge":
						Handle := HTBClient.Challenges.Challenge(ChallengeID)
						resp, err := Handle.Own(ctx, flag)
						if err != nil {
							fmt.Println(err.Error())
						}else{
							fmt.Println(resp)
						}
					case "fortress":
						Handle := HTBClient.Fortresses.Fortress(fortressID)
						resp, err := Handle.SubmitFlag(ctx, flag)
						if err != nil {
							fmt.Println(err.Error())
						}else{
							fmt.Println(resp)
						}
				} // end of flag submission
			case "info":
				fs := action.Data.(Info)
				content_type := ""
				boxID := fs.BoxID
				boxName := fs.BoxName
				ChallengeID := fs.ChallengeID
				ChallengeName := fs.Challenge
				fortressID := fs.FortressID
				fortressName := fs.Fortress
				if ( boxID == 0 && ChallengeID == 0 && fortressID == 0){ // grab based on name
					if boxName != ""{
						boxID, err = GetMachineID(boxName, *HTBClient)
						if err != nil {
							fmt.Println(err.Error())
						}
						content_type = "box"
					}
					if ChallengeName != "" && content_type == ""{
						ChallengeID, err = GetChallengeID(ChallengeName, *HTBClient)
						if err != nil {
							fmt.Println(err.Error())
						}
						content_type = "challenge"
					}
					if fortressName != "" && content_type == ""{
						fortressID, err = GetFortressID(fortressName, *HTBClient)
						if err != nil {
							fmt.Println(err.Error())
						}
						content_type = "fortress"
					}
				} // end of grab based on name
				switch content_type{ // submit flag
					case "box":
						Handle := HTBClient.Machines.Machine(boxID)
						resp, err := Handle.Info(ctx)
						if err != nil {
							fmt.Println(err.Error())
						}else{
							fmt.Println(resp)
						}
					case "challenge":
						Handle := HTBClient.Challenges.Challenge(ChallengeID)
						resp, err := Handle.Info(ctx)
						if err != nil {
							fmt.Println(err.Error())
						}else{
							fmt.Println(resp)
						}
					case "fortress":
						Handle := HTBClient.Fortresses.Fortress(fortressID)
						resp, err := Handle.Info(ctx)
						if err != nil {
							fmt.Println(err.Error())
						}else{
							fmt.Println(resp)
						}
				} // end of flag submission

		}

	}
}


func GetMachineID(name string, HTBClient HTB.Client) (id int, err error) {
	name = strings.ToLower(name)
	contents, err := HTBClient.Machines.ListActive().AllResults(ctx)
	if err != nil {
		return
	}
	Response := contents.Data
	for _, r := range Response {
		if strings.ToLower(r.Name) == name {
			id = r.Id
			return
		}
	}
	contents, err = HTBClient.Machines.ListRetired().AllResults(ctx)
	if err != nil {
		return
	}
	Response = contents.Data
	for _, r := range Response {
		if strings.ToLower(r.Name) == name {
			id = r.Id
			return
		}
	}

	contents1, err := HTBClient.Machines.ListUnreleased().AllResults(ctx)
	if err != nil {
		return
	}
	Response1 := contents1.Data
	for _, r := range Response1 {
		if strings.ToLower(r.Name) == name {
			id = r.Id
			return
		}
	}


	err = fmt.Errorf("unable to find machine")
	return

}




func GetChallengeID(name string, HTBClient HTB.Client) (id int, err error) {
	name = strings.ToLower(name)
	contents, err := HTBClient.Challenges.List().AllResults(ctx)
	if err != nil {
		return
	}
	Response := contents.Data
	for _, r := range Response {
		if strings.ToLower(r.Name) == name {
			id = r.Id
			return
		}
	}
	err = fmt.Errorf("unable to find challenge")
	return

}

func GetFortressID(name string, HTBClient HTB.Client) (id int, err error) {
	name = strings.ToLower(name)
	contents, err := HTBClient.Fortresses.List(ctx)
	if err != nil {
		return
	}
	Response := contents.Data
	for _, r := range Response {
		if strings.ToLower(r.Name) == name {
			id = r.Id
			return
		}
	}
	err = fmt.Errorf("unable to find challenge")
	return

}