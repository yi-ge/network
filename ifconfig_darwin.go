package network

import (
	"bufio"
	"os/exec"
	"strings"
)

// IfconfigInterfaces .
type IfconfigInterfaces struct {
	Name                  string // e.g., "en0", "lo0", "eth0.100"
	HardwareAddr          string // IEEE MAC-48, EUI-48 and EUI-64 form
	Type                  string // Wired, Wi-Fi
	DHCPEnabled           bool
	IPv4Address           string
	SubnetPrefix          string
	DefaultGatewayAddress string
	DNSPrimary            string
	DNSBack               string
	Description           string // In mac is HardwarePort, other is net card name
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
	out, err := runner.exec.Command("ifconfig", "-a").CombinedOutput()

	if err != nil {
		return "", err
	}

	return string(out[:]), nil
}

func (runner *runner) parseIfconfig(str string) []IfconfigInterfaces {
	var (
		IfconfigInterfacesList []IfconfigInterfaces
		currentInterface       IfconfigInterfaces
	)

	output := str
	currentInterface = IfconfigInterfaces{}
	scanner := bufio.NewScanner(strings.NewReader(output))

	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, "flags=") {
			if currentInterface.Name != "" {
				IfconfigInterfacesList = append(IfconfigInterfacesList, currentInterface)
				currentInterface = IfconfigInterfaces{}
			}

			fs := strings.Split(line, ":")
			currentInterface.Name = fs[0]
		} else if strings.Contains(line, "ether") {
			fs := strings.Fields(line)
			value := fs[1]
			currentInterface.HardwareAddr = strings.ToUpper(strings.Replace(value, "-", ":", -1))
		} else if strings.Contains(line, "inet") && !strings.Contains(line, "inet6") {
			fs := strings.Fields(line)
			currentInterface.IPv4Address = fs[1]
			if len(fs) > 3 {
				currentInterface.SubnetPrefix = hex2dot(fs[3])
			}
			if len(fs) > 5 {
				currentInterface.DefaultGatewayAddress = fs[5]
			}
		} else if strings.Contains(line, "status:") {
			if strings.Contains(line, "inactive") {
				currentInterface.Connected = false
				currentInterface.AdminState = "Disabled"
			} else {
				currentInterface.Connected = true
				currentInterface.AdminState = "Enable"
			}
		}
	}

	currentInterface.Type = "Wired"

	if currentInterface != (IfconfigInterfaces{}) {
		IfconfigInterfacesList = append(IfconfigInterfacesList, currentInterface)
	}

	hardwarePortList, err := runner.getAllHardwarePortList()

	if err != nil {
		return IfconfigInterfacesList
	}

	for index, ifconfigInterfaces := range IfconfigInterfacesList {
		for _, hardwarePort := range hardwarePortList {
			if ifconfigInterfaces.Name == hardwarePort.Device {
				IfconfigInterfacesList[index].Description = hardwarePort.HardwarePort
				isDHCP, primary, back, err := runner.getDNSServer(hardwarePort.HardwarePort)
				if err != nil {
					continue
				}
				if isDHCP {
					IfconfigInterfacesList[index].DHCPEnabled = true
				} else {
					IfconfigInterfacesList[index].DNSPrimary = primary
					IfconfigInterfacesList[index].DNSBack = back
				}
			}
		}
	}

	return IfconfigInterfacesList
}

// EnabledNetworkInterfaceByIfconfig .
func (runner *runner) EnabledNetworkInterfaceByIfconfig(interfaceName string) error {
	command := "ifconfig " + interfaceName + " up"
	_, err := runner.exec.Command(command).CombinedOutput()
	if err != nil {
		return err
	}
	return nil
}

// DisableNetworkInterfaceByIfconfig .
func (runner *runner) DisableNetworkInterfaceByIfconfig(interfaceName string) error {
	command := "ifconfig " + interfaceName + " down"
	_, err := runner.exec.Command(command).CombinedOutput()
	if err != nil {
		return err
	}
	return nil
}

// 查看wifi对应设备名
// networksetup -listallhardwareports
