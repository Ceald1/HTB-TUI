package yaml

// TODO: make outputs pretty
import (
	"context"
	"strings"

	"github.com/Ceald1/HTB-TUI/src/format"
	"github.com/charmbracelet/lipgloss"
	HTB "github.com/gubarz/gohtb"

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
// TODO: Add VPN downloading support, and pro lab support 




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
		a.Type = "vpn"
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


func RunAutomation(yaml_file string) {
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
		switch action.Type {
		case "vpn":
			fs := action.Data.(VPNDownload)
			region := strings.ToUpper(fs.Region)
			if region == ""{
				fmt.Println(ErrorText(fmt.Errorf("Region is required!")))
				break
			}
			tier := strings.ToLower(fs.Tier)
			if tier == ""{
				fmt.Println(ErrorText(fmt.Errorf("Tier is required (free, vip, vip+, or unknown)!")))
				break
			}
			typeVpn := strings.ToLower(fs.Type)
			if typeVpn == ""{
				fmt.Println(ErrorText(fmt.Errorf("Type is required (labs, release_arena, endgame, prolab, or fortress)!")))
				break
			}
			outFile := fs.Outfile
			if outFile == ""{
				fmt.Println(ErrorText(fmt.Errorf("Outfile name is required (don't put .ovpn)!")))
				break
			}
			protocol := strings.ToLower(fs.Protocol)
			if protocol != "udp" && protocol != "tcp"{
				fmt.Println("invalid protocol found, defaulting to TCP")
				protocol = "tcp"
			}
			err = DownloadVPN(HTBClient, region, tier, typeVpn, outFile, protocol)
			if err != nil {
				fmt.Println(ErrorText(err))
				return
			}


		
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

			// resolve ID by name
			if boxID == 0 && ChallengeID == 0 && fortressID == 0 {
				if boxName != "" {
					boxID, err = GetMachineID(boxName, *HTBClient)
					if err != nil {
						fmt.Println(ErrorText(err))
						break
					}
					content_type = "box"
				}
				if ChallengeName != "" && content_type == "" {
					ChallengeID, err = GetChallengeID(ChallengeName, *HTBClient)
					if err != nil {
						fmt.Println(ErrorText(err))
						break
					}
					content_type = "challenge"
				}
				if fortressName != "" && content_type == "" {
					fortressID, err = GetFortressID(fortressName, *HTBClient)
					if err != nil {
						fmt.Println(ErrorText(err))
						break
					}
					content_type = "fortress"
				}
			}

			switch content_type {
			case "box":
				Handle := HTBClient.Machines.Machine(boxID)
				resp, err := Handle.Own(ctx, flag)
				if err != nil {
					fmt.Println(ErrorText(err))
				} else {
					fmt.Println(SubmissionText(resp.Data.Message))
				}

			case "challenge":
				Handle := HTBClient.Challenges.Challenge(ChallengeID)
				resp, err := Handle.Own(ctx, flag)
				if err != nil {
					fmt.Println(ErrorText(err))
				} else {
					fmt.Println(SubmissionText(resp.Data.Message))
				}

			case "fortress":
				Handle := HTBClient.Fortresses.Fortress(fortressID)
				resp, err := Handle.SubmitFlag(ctx, flag)
				if err != nil {
					fmt.Println(ErrorText(err))
				} else {
					fmt.Println(SubmissionText(resp.Data.Message))
				}
			}

			fmt.Println(format.SplitResp()) // Print once at the end

		case "info":
			fs := action.Data.(Info)
			content_type := ""
			boxID := fs.BoxID
			boxName := fs.BoxName
			ChallengeID := fs.ChallengeID
			ChallengeName := fs.Challenge
			fortressID := fs.FortressID
			fortressName := fs.Fortress

			// resolve ID by name
			if boxID == 0 && ChallengeID == 0 && fortressID == 0 {
				if boxName != "" {
					boxID, err = GetMachineID(boxName, *HTBClient)
					if err != nil {
						fmt.Println(ErrorText(err))
						break
					}
					content_type = "box"
				}
				if ChallengeName != "" && content_type == "" {
					ChallengeID, err = GetChallengeID(ChallengeName, *HTBClient)
					if err != nil {
						fmt.Println(ErrorText(err))
						break
					}
					content_type = "challenge"
				}
				if fortressName != "" && content_type == "" {
					fortressID, err = GetFortressID(fortressName, *HTBClient)
					if err != nil {
						fmt.Println(ErrorText(err))
						break
					}
					content_type = "fortress"
				}
			}

			switch content_type {
			case "box":
				Handle := HTBClient.Machines.Machine(boxID)
				resp, err := Handle.Info(ctx)
				if err != nil {
					fmt.Println(ErrorText(err))
				} else {
					data := lipgloss.NewStyle().Render(
						lipgloss.NewStyle().Background(format.BaseBG).Foreground(format.TextTitle).Padding(1,1).Render(resp.Data.Name),
						fmt.Sprintf("\nDifficulty: %s", lipgloss.NewStyle().Render(format.CheckDiff(resp.Data.DifficultyText))),
						fmt.Sprintf("\nOS: %s", lipgloss.NewStyle().Foreground(format.TextDefault).Render(format.CheckOS(resp.Data.Os))),
						fmt.Sprintf("\nCredentials: %s:%s", lipgloss.NewStyle().Foreground(format.Pink).Render(resp.Data.Credentials.Username), lipgloss.NewStyle().Foreground(format.Pink).Render(resp.Data.Credentials.Password)),
					)
					
					
					fmt.Println(data)
				}

			case "challenge":
				Handle := HTBClient.Challenges.Challenge(ChallengeID)
				resp, err := Handle.Info(ctx)
				if err != nil {
					fmt.Println(ErrorText(err))
				} else {
					data := lipgloss.NewStyle().Render(
						lipgloss.NewStyle().Background(format.BaseBG).Foreground(format.TextTitle).Padding(1, 1).Render(resp.Data.Name),
						fmt.Sprintf("\nDifficulty: %s", lipgloss.NewStyle().Render(format.CheckDiff(resp.Data.Difficulty))),
						fmt.Sprintf("\nCategory: %s", lipgloss.NewStyle().Foreground(format.TextDefault).Render(resp.Data.CategoryName)),
						fmt.Sprintf("\nDescription: %s\n", lipgloss.NewStyle().Foreground(format.Pink).Render(strings.Replace(resp.Data.Description, "\n", "", -1))),
					)
					fmt.Println(data)
				}

			case "fortress":
				Handle := HTBClient.Fortresses.Fortress(fortressID)
				resp, err := Handle.Info(ctx)
				if err != nil {
					fmt.Println(ErrorText(err))
				} else {
					data := lipgloss.NewStyle().Render(
						lipgloss.NewStyle().Background(format.BaseBG).Foreground(format.TextTitle).Padding(1, 1).Render(resp.Data.Name),
						fmt.Sprintf("\n# Of Flags: %s", lipgloss.NewStyle().Render(string(len(resp.Data.Flags)))),
						fmt.Sprintf("\nCompany: %s", lipgloss.NewStyle().Foreground(format.TextDefault).Render(resp.Data.Company.Name)),
						fmt.Sprintf("\nDescription: %s\n", lipgloss.NewStyle().Foreground(format.Pink).Render(strings.Replace(resp.Data.Description, "\n", "", -1))),
					)
					fmt.Println(data)
				}
			}

			fmt.Println(format.SplitResp()) // Print once at the end
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








