package yaml

import (
	"fmt"
	"os"

	HTB "github.com/gubarz/gohtb"
	"github.com/gubarz/gohtb/services/vpn"
)


type VPNDownload struct{
	Region 		string `yaml:"region"`
	Tier   		string `yaml:"tier"`
	Type   		string `yaml:"type"`
	Outfile 	string `yaml:"outfile"`
	Protocol	string `yaml:"protocol"`
}

func DownloadVPN(HTBClient *HTB.Client, region, tier, vpnType, outFile,  protocol string) (err error) {
	// Download HTB VPN
	servers, err := HTBClient.VPN.Servers(vpnType).ByTier(tier).ByLocation(region).Results(ctx)
	// Servers("labs").ByTier("free").ByLocation("US")
	if err != nil {
		return err
	}
	best := servers.Data.Options.SortByCurrentClients().First()
	if best.Id == 0 {
		err = fmt.Errorf("unable to find VPN server")
		return err
	}
	var vpnConfig vpn.VPNFileResponse
	if protocol == "udp" {
		vpnConfig, err = HTBClient.VPN.VPN(best.Id).SwitchAndDownloadUDP(ctx)
		if err != nil {
			return err
		}
	}else{
		vpnConfig, err = HTBClient.VPN.VPN(best.Id).SwitchAndDownloadTCP(ctx)
		if err != nil {
			return err
		}
	}
	outFile = fmt.Sprintf("%s.ovpn", outFile)
	err = os.WriteFile(outFile, vpnConfig.Data, 0644)
	return

}