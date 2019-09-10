package network

import (
	"os/exec"
	"regexp"
	"strings"
)

// IfconfigInterfaces .
type IfconfigInterfaces struct {
	Name                  string // e.g., "en0", "lo0", "eth0.100"
	HardwareAddr          string // IEEE MAC-48, EUI-48 and EUI-64 form
	Type                  string // Ethernet, Wireless LAN
	DHCPEnabled           bool
	IPv4Address           string
	SubnetPrefix          string
	DefaultGatewayAddress string
	DNSPrimary            string
	DNSBack               string
	Description           string
	Connected             bool
	Mode                  string
	AdminState            string
}

// IsInstalled checks if the program ifconfig exists using PATH environment variable
func IsInstalled() bool {
	_, err := exec.LookPath("ifconfig")
	if err != nil {
		return false
	}
	return true
}

// getIfconfigOutPut .
func (runner *runner) getIfconfigOutPut() (string, error) {
	out, err := runner.exec.Command("ifconfig -a").CombinedOutput()

	if err != nil {
		return "", err
	}

	return string(out[:]), nil
}

// minIndexAndCardType .
func minIndexAndCardType(x []int, xType *regexp.Regexp, y []int, yType *regexp.Regexp) (int, *regexp.Regexp, string) {
	if len(x) != 0 && len(y) != 0 && x[1] < y[1] {
		return x[1], xType, "Ethernet"
	} else if len(x) != 0 && len(y) != 0 && x[1] > y[1] {
		return y[1], yType, "Wireless LAN"
	} else if len(x) != 0 && len(y) == 0 {
		return x[1], xType, "Ethernet"
	} else if len(x) == 0 && len(y) != 0 {
		return y[1], yType, "Wireless LAN"
	}

	return 0, nil, ""
}

func parseIfconfig(str string) []IfconfigInterfaces {
	repEthernet := regexp.MustCompile(`\bEthernet adapter ([^:\r\n]+):`)     // 判断有线网卡
	repWireless := regexp.MustCompile(`\bWireless LAN adapter ([^:\r\n]+):`) // 判断无线网卡
	repItem := regexp.MustCompile(`(?m)^\s*$[\r\n]*|[\r\n]+\s+\z`)           // 判断空行

	var (
		IfconfigInterfacesList []IfconfigInterfaces
		currentInterface       IfconfigInterfaces
		inDNSPrimary           = false
	)

	output := str

	cardIndex, cardType, typeName := minIndexAndCardType(repEthernet.FindStringSubmatchIndex(output), repEthernet, repWireless.FindStringSubmatchIndex(output), repWireless)
	for cardIndex != 0 {
		card := cardType.FindStringSubmatch(output)

		currentInterface = IfconfigInterfaces{
			Name: card[1],
		}
		output = output[cardIndex+4:]
		itemOutIndex := repItem.FindStringSubmatchIndex(output)
		itemOut := output[:itemOutIndex[1]]
		output = output[itemOutIndex[1]:]

		outputLines := strings.Split(itemOut, "\r\n")

		for _, outputLine := range outputLines {
			parts := strings.SplitN(outputLine, ":", 2)
			if len(parts) != 2 {
				if inDNSPrimary {
					currentInterface.DNSBack = strings.TrimSpace(outputLine)
					inDNSPrimary = false
				}
				continue
			}
			if inDNSPrimary {
				inDNSPrimary = false
			}
			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])
			if strings.HasPrefix(key, "Physical Address") {
				currentInterface.HardwareAddr = strings.ToUpper(strings.Replace(value, "-", ":", -1))
			} else if strings.HasPrefix(key, "DHCP enabled") {
				if value == "Yes" {
					currentInterface.DHCPEnabled = true
				}
			} else if strings.HasPrefix(key, "IPv4 Address") || strings.HasPrefix(key, "IP Address") {
				currentInterface.IPv4Address = strings.Replace(value, "(Preferred)", "", -1)
			} else if strings.HasPrefix(key, "Subnet Prefix") || strings.HasPrefix(key, "Subnet Mask") {
				currentInterface.SubnetPrefix = value
			} else if strings.HasPrefix(key, "Default Gateway") {
				currentInterface.DefaultGatewayAddress = value
			} else if strings.HasPrefix(key, "DNS Servers") {
				currentInterface.DNSPrimary = value
				inDNSPrimary = true
			} else if strings.HasPrefix(key, "Description") {
				currentInterface.Description = value
			} else if strings.HasPrefix(key, "Media State") {
				if strings.Contains(value, "disconnected") {
					currentInterface.Connected = false
				} else {
					currentInterface.Connected = true
				}
			}
		}

		if !strings.Contains(itemOut, "Media State") {
			currentInterface.Connected = true
		}

		currentInterface.Type = typeName

		if currentInterface != (IfconfigInterfaces{}) {
			IfconfigInterfacesList = append(IfconfigInterfacesList, currentInterface)
		}
		cardIndex, cardType, typeName = minIndexAndCardType(repEthernet.FindStringSubmatchIndex(output), repEthernet, repWireless.FindStringSubmatchIndex(output), repWireless)
	}

	return IfconfigInterfacesList
}

// EnableNetworkInterfaceByIfconfig .
func (runner *runner) EnableNetworkInterfaceByIfconfig(interfaceName string) error {
	command := "ifconfig " + interfaceName + " up"
	_, err := runner.exec.Command(command).CombinedOutput()
	if err != nil {
		return err
	}
	return nil
}

// DisabledNetworkInterfaceByIfconfig .
func (runner *runner) DisabledNetworkInterfaceByIfconfig(interfaceName string) error {
	command := "ifconfig " + interfaceName + " down"
	_, err := runner.exec.Command(command).CombinedOutput()
	if err != nil {
		return err
	}
	return nil
}
