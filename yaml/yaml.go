package yaml

// TODO: make outputs pretty
import (
	"strings"

	"github.com/Ceald1/HTB-TUI/src/format"
	"github.com/charmbracelet/lipgloss"
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
// TODO: Add VPN switching and downloading support
type VPNDownload struct{} //Future implementation

type VPNSwitch struct{} // Future implementation


func (a *Action) UnmarshalYAML(value *yaml.Node) error {
	if len(value.Content) != 2 {
		return fmt.Errorf("invalid action format")
	}

	keyNode := value.Content[0]
	valNode := value.Content[1]

	switch keyNode.Value {
	// Labs
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
	
	// VPN 
	case "vpnDownload":
		var vpn VPNDownload
		if err := valNode.Decode(&vpn); err != nil {
			return err
		}
		a.Type = "vpnDownload"
		a.Data = vpn
	case "vpnSwitch":
		var vpn VPNSwitch
		if err := valNode.Decode(&vpn); err != nil {
			return err
		}
		a.Type = "vpnSwitch"
		a.Data = vpn
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
							fmt.Println(ErrorText(err))
						}
						content_type = "box"
					}
					if ChallengeName != "" && content_type == ""{
						ChallengeID, err = GetChallengeID(ChallengeName, *HTBClient)
						if err != nil {
							fmt.Println(ErrorText(err))
						}
						content_type = "challenge"
					}
					if fortressName != "" && content_type == ""{
						fortressID, err = GetFortressID(fortressName, *HTBClient)
						if err != nil {
							fmt.Println(ErrorText(err))
						}
						content_type = "fortress"
					}
				} // end of grab based on name
				switch content_type{ // submit flag
					case "box":
						Handle := HTBClient.Machines.Machine(boxID)
						resp, err := Handle.Own(ctx, flag)
						if err != nil {
							fmt.Println(ErrorText(err))
						}else{
							fmt.Println(SubmissionText(resp.Data.Message))
						}
					case "challenge":
						Handle := HTBClient.Challenges.Challenge(ChallengeID)
						resp, err := Handle.Own(ctx, flag)
						if err != nil {
							fmt.Println(ErrorText(err))
						}else{
							fmt.Println(SubmissionText(resp.Data.Message))
						}
					case "fortress":
						Handle := HTBClient.Fortresses.Fortress(fortressID)
						resp, err := Handle.SubmitFlag(ctx, flag)
						if err != nil {
							fmt.Println(ErrorText(err))
						}else{
							fmt.Println(SubmissionText(resp.Data.Message))
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
							fmt.Println(ErrorText(err))
						}
						content_type = "box"
					}
					if ChallengeName != "" && content_type == ""{
						ChallengeID, err = GetChallengeID(ChallengeName, *HTBClient)
						if err != nil {
							fmt.Println(ErrorText(err))
						}
						content_type = "challenge"
					}
					if fortressName != "" && content_type == ""{
						fortressID, err = GetFortressID(fortressName, *HTBClient)
						if err != nil {
							fmt.Println(ErrorText(err))
						}
						content_type = "fortress"
					}
				} // end of grab based on name
				switch content_type{ // Lab info
					case "box":
						Handle := HTBClient.Machines.Machine(boxID)
						resp, err := Handle.Info(ctx)
						if err != nil {
							fmt.Println(ErrorText(err))
						}else{
							fmt.Println(resp)
						}
					case "challenge":
						Handle := HTBClient.Challenges.Challenge(ChallengeID)
						resp, err := Handle.Info(ctx)
						if err != nil {
							fmt.Println(ErrorText(err))
						}else{
							fmt.Println(resp)
						}
					case "fortress":
						Handle := HTBClient.Fortresses.Fortress(fortressID)
						resp, err := Handle.Info(ctx)
						if err != nil {
							fmt.Println(ErrorText(err))
						}else{
							fmt.Println(resp)
						}
				} // end of lab Info

		}

	}
}

func ErrorText(err error) (out string) {
	out = lipgloss.NewStyle().Foreground(format.Red).Render(err.Error())
	return
}

func SubmissionText(message string) (out string){
	out = lipgloss.NewStyle().Foreground(format.LightGreen).Render(message)
	return
}


func GetMachineID(name string, HTBClient HTB.Client) (id int, err error) {
	name = strings.ToLower(name)
	format.TaskResult = 0
	task := format.Task(func(a any) any {
		HTBClient, _ := a.(*HTB.Client)
		contents, _ := HTBClient.Machines.ListActive().AllResults(ctx)
		Response := contents.Data
		for _, r := range Response {
			if strings.ToLower(r.Name) == name {
				id = r.Id
				return id
			}
		}
		contents, err = HTBClient.Machines.ListRetired().AllResults(ctx)

		Response = contents.Data
		for _, r := range Response {
			if strings.ToLower(r.Name) == name {
				id = r.Id
				return id
			}
		}

		contents1, _ := HTBClient.Machines.ListUnreleased().AllResults(ctx)

		Response1 := contents1.Data
		for _, r := range Response1 {
			if strings.ToLower(r.Name) == name {
				id = r.Id
				return id
			}
		}
		return 0
	})
	err = format.RunLoading(task, &HTBClient)
	if err != nil {
		return
	}

	id, _ = format.TaskResult.(int)
	if id == 0{
		err = fmt.Errorf("unable to find machine")
	}
	return

}




func GetChallengeID(name string, HTBClient HTB.Client) (id int, err error) {
	name = strings.ToLower(name)
	format.TaskResult = 0
	task := format.Task(func(a any) any {
		HTBClient, _ := a.(*HTB.Client)
		contents, _ := HTBClient.Challenges.List().AllResults(ctx)
		Response := contents.Data
		for _, r := range Response {
			if strings.ToLower(r.Name) == name {
				id = r.Id
				return id
			}
		}
		return 0 
	})
	err = format.RunLoading(task, &HTBClient)
	if format.TaskResult == 0 {
		err = fmt.Errorf("unable to find challenge")
	}
	id, _ = format.TaskResult.(int)
	return 

}

func GetFortressID(name string, HTBClient HTB.Client) (id int, err error) {
	name = strings.ToLower(name)
	format.TaskResult = 0
	task := format.Task(func(a any) any {
		contents, _ := HTBClient.Fortresses.List(ctx)

		Response := contents.Data
		for _, r := range Response {
			if strings.ToLower(r.Name) == name {
				id = r.Id
				return id
			}
		}
	return 0
	})
	err = format.RunLoading(task, &HTBClient)
	if err != nil {
		return
	}
	if format.TaskResult == 0 {
		err = fmt.Errorf("unable to find challenge")
		return
	}
	id, _ = format.TaskResult.(int)
	return

}