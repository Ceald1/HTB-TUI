package vpn

import (
	"context"

	"github.com/Ceald1/HTB-TUI/src/format"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
	HTB "github.com/gubarz/gohtb"
	"github.com/gubarz/gohtb/services/prolabs"
	"github.com/gubarz/gohtb/services/vpn"
)

var (
	ctx = context.Background()
	VPNSelected string
)
const (
	quit_value = 9999999999999
)


func SelectVPNLabs(HTBClient *HTB.Client) (vpn_data []byte) {
	var selected string

	var options []huh.Option[string]
	var quit_value = "9999999999999"
	quit_op := huh.NewOption(lipgloss.NewStyle().Foreground(format.Red).Background(format.BaseBG).Render("Quit"), quit_value)
	options = append(options, quit_op)
	var vpn_options = []string{"labs","prolabs", "starting_point", "fortresses"}
	for _, option := range vpn_options {
		info := lipgloss.NewStyle().Foreground(format.NextColor()).Background(format.BaseBG).Render(option)
		op := huh.NewOption(info, option)
		options = append(options, op)
	}
	huh.NewSelect[string]().
		Title(lipgloss.NewStyle().Foreground(format.TextTitle).Background(format.BaseBG).Render("VPN Lab Options")).
		Options(options...).Value(&selected).Run()
	switch selected {
		case quit_value:
			return
		case "labs", "starting_point", "fortresses":
			return DownloadLabVPN(HTBClient, selected)
		case "prolabs":
			return DownloadProlabVPN(HTBClient)
		default:
			return
	}
	
}

func DownloadLabVPN(HTBClient *HTB.Client, product string) (vpn_data []byte) {
	servers, err := HTBClient.VPN.Servers(product).Results(ctx)
	var selected int
	var proto bool
	if err != nil {
		panic(err)
	}
	var theme = format.HTBTheme()

	var options []huh.Option[int]
	for _, server := range servers.Data.Options{
		info := lipgloss.NewStyle().Foreground(format.NextColor()).Background(format.BaseBG).Render(server.FriendlyName)
		op := huh.NewOption(info, server.Id)
		options = append(options, op)
	}
	huh.NewForm(
		huh.NewGroup(
			huh.NewConfirm().
				Title(lipgloss.NewStyle().Foreground(format.TextTitle).Background(format.BaseBG).Render("Servers")).
				Value(&proto).
				Affirmative("TCP").
				Negative("UDP"),
			huh.NewSelect[int]().
				Title(lipgloss.NewStyle().Foreground(format.TextTitle).Background(format.BaseBG).Render("Servers")).
				Options(options...).Value(&selected),
		),
	).WithTheme(theme).Run()
	
	switch selected{
		case quit_value:
			return
		default:
			task := format.Task(func(a any) any {
				var vpnConfig vpn.VPNFileResponse
				if proto{
					vpnConfig, err = HTBClient.VPN.VPN(selected).SwitchAndDownloadTCP(ctx)
					if err != nil {
						panic(err)
					}
				}else{
					vpnConfig, err = HTBClient.VPN.VPN(selected).SwitchAndDownloadUDP(ctx)
					if err != nil {
						panic(err)
					}
				}
				return vpnConfig.Data
			})
			// var vpnConfig vpn.VPNFileResponse
			// if proto{
			// 	vpnConfig, err = HTBClient.VPN.VPN(selected).SwitchAndDownloadTCP(ctx)
			// 	if err != nil {
			// 		panic(err)
			// 	}
			// }else{
			// 	vpnConfig, err = HTBClient.VPN.VPN(selected).SwitchAndDownloadUDP(ctx)
			// 	if err != nil {
			// 		panic(err)
			// 	}
			// }
			// vpn_data = vpnConfig.Data

			err := format.RunLoading(task, HTBClient)
			if err != nil {
				panic(err)
			}
			var ok bool
			vpn_data, ok = format.TaskResult.([]byte)
			if !ok {
				panic("error occurred")
			}
			return 

	}
}


func DownloadProlabVPN(HTBClient *HTB.Client) (vpn_data []byte) {
	selected_lab := prolabList(HTBClient)
	var proto bool
	var selected int
	if selected_lab == quit_value {
		return
	}

	servers, err := HTBClient.VPN.ProlabServers(selected_lab).Results(ctx)
	if err != nil {
		panic(err)
	}
	var theme = format.HTBTheme()

	var options []huh.Option[int]
	for _, server := range servers.Data.Options{
		info := lipgloss.NewStyle().Foreground(format.NextColor()).Background(format.BaseBG).Render(server.FriendlyName)
		op := huh.NewOption(info, server.Id)
		options = append(options, op)
	}
	huh.NewForm(
		huh.NewGroup(
			huh.NewConfirm().
				Title(lipgloss.NewStyle().Foreground(format.TextTitle).Background(format.BaseBG).Render("Servers")).
				Value(&proto).
				Affirmative("TCP").
				Negative("UDP"),
			huh.NewSelect[int]().
				Title(lipgloss.NewStyle().Foreground(format.TextTitle).Background(format.BaseBG).Render("Servers")).
				Options(options...).Value(&selected),
		),
	).WithTheme(theme).Run()
	
	switch selected{
		case quit_value:
			return
		default:
			task := format.Task(func(a any) any {
				var vpnConfig vpn.VPNFileResponse
				if proto{
					vpnConfig, err = HTBClient.VPN.VPN(selected).SwitchAndDownloadTCP(ctx)
					if err != nil {
						panic(err)
					}
				}else{
					vpnConfig, err = HTBClient.VPN.VPN(selected).SwitchAndDownloadUDP(ctx)
					if err != nil {
						panic(err)
					}
				}
				vpn_data = vpnConfig.Data
				return vpn_data
			})
			// var vpnConfig vpn.VPNFileResponse
			// if proto{
			// 	vpnConfig, err = HTBClient.VPN.VPN(selected).SwitchAndDownloadTCP(ctx)
			// 	if err != nil {
			// 		panic(err)
			// 	}
			// }else{
			// 	vpnConfig, err = HTBClient.VPN.VPN(selected).SwitchAndDownloadUDP(ctx)
			// 	if err != nil {
			// 		panic(err)
			// 	}
			// }
			// vpn_data = vpnConfig.Data
			err := format.RunLoading(task, HTBClient)
			if err != nil {
				panic(err)
			}
			var ok bool
			vpn_data, ok = format.TaskResult.([]byte)
			if !ok {
				panic("error occurred")
			}
			return 

	}



}


func prolabList(HTBClient *HTB.Client) (selectedLab int) {
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
		return quit_value
	default:
		return
	}
}