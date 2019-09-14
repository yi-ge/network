package network

import (
	"bufio"
	"fmt"
	"strings"
)

// HardwarePort .
type HardwarePort struct {
	HardwarePort    string
	Device          string
	EthernetAddress string
}

// Info .
type Info struct {
	IPAddress     string
	SubnetMask    string
	Router        string
	ClientID      string
	IPv6          string
	IPv6IPAddress string
	IPv6Router    string
	WiFiID        string
}

func (runner *runner) ConnectWifi(interfaceName string, ssid string, password string, securityType string, broadcast bool) (string, error) {
	args := []string{
		"-setairportpower",
		"\"" + interfaceName + "\"",
		"on",
	}

	_, err := runner.exec.Command("networksetup", args...).CombinedOutput()
	if err != nil {
		return "", err
	}

	argsConnect := []string{
		"-setairportpower",
		"\"" + interfaceName + "\"",
		ssid,
		password,
	}

	out, err := runner.exec.Command("networksetup", argsConnect...).CombinedOutput()
	if err != nil {
		return "", err
	}

	return string(out[:]), nil
}

func (runner *runner) DisconnectWifi(interfaceName string) (string, error) {
	args := []string{
		"-setairportpower",
		"\"" + interfaceName + "\"",
		"off",
	}

	out, err := runner.exec.Command("networksetup", args...).CombinedOutput()
	if err != nil {
		return "", err
	}

	return string(out[:]), nil
}

func (runner *runner) SetDNS(interfaceName string, primaryAddr string, backAddr string) error {
	args := []string{
		"-setdnsservers",
		"\"" + interfaceName + "\"",
		primaryAddr,
		backAddr,
	}

	_, err := runner.exec.Command("networksetup", args...).CombinedOutput()
	if err != nil {
		return err
	}

	return nil
}

func (runner *runner) SetDNSUseDHCP(interfaceName string) error {
	args := []string{
		"-setairportpower",
		"\"" + interfaceName + "\"",
		"empty",
	}

	_, err := runner.exec.Command("networksetup", args...).CombinedOutput()
	if err != nil {
		return err
	}

	return nil
}

func (runner *runner) SetInterfaceUseDHCP(interfaceName string) error {
	args := []string{
		"-setdhcp",
		"\"" + interfaceName + "\"",
	}

	_, err := runner.exec.Command("networksetup", args...).CombinedOutput()
	if err != nil {
		return err
	}

	return nil
}

func (runner *runner) SetStaticIP(interfaceName string, addr string, mask string, gateway string) error {
	args := []string{
		"-setmanual",
		"\"" + interfaceName + "\"",
		addr,
		mask,
		gateway,
	}

	_, err := runner.exec.Command("networksetup", args...).CombinedOutput()
	if err != nil {
		return err
	}

	return nil
}

func (runner *runner) getAllHardwarePortList() ([]HardwarePort, error) {
	args := []string{
		"-listallhardwareports",
	}

	out, err := runner.exec.Command("networksetup", args...).CombinedOutput()
	if err != nil {
		return []HardwarePort{}, err
	}

	scanner := bufio.NewScanner(strings.NewReader(string(out[:])))

	hardwarePortList := []HardwarePort{}
	hardwarePort := HardwarePort{}
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, "Hardware Port") {
			if hardwarePort.HardwarePort != "" {
				hardwarePortList = append(hardwarePortList, hardwarePort)
				hardwarePort = HardwarePort{}
			}
			hardwarePort.HardwarePort = strings.TrimSpace(strings.Split(line, ":")[1])
		} else if strings.Contains(line, "Device") {
			hardwarePort.Device = strings.TrimSpace(strings.Split(line, ":")[1])
		} else if strings.Contains(line, "Ethernet Address") {
			hardwarePort.EthernetAddress = strings.TrimSpace(strings.Split(line, ":")[1])
		}
	}

	if hardwarePort.HardwarePort != "" {
		hardwarePortList = append(hardwarePortList, hardwarePort)
		hardwarePort = HardwarePort{}
	}

	return hardwarePortList, nil
}

func (runner *runner) getDNSServer(hardwarePort string) (DHCP bool, primary string, back string, err error) {
	args := []string{
		"-getdnsservers",
		hardwarePort,
	}

	out, err := runner.exec.Command("networksetup", args...).CombinedOutput()
	if err != nil {
		return false, "", "", err
	}

	output := string(out[:])

	if strings.Contains(output, "t any DNS Servers") {
		return true, "", "", nil
	}

	outputLines := strings.Split(output, "\n")

	primary = outputLines[0]
	back = outputLines[1]

	return false, primary, back, nil
}

func (runner *runner) getNetworkServiceEnabled(hardwarePort string) (string, error) {
	args := []string{
		"-getnetworkserviceenabled",
		hardwarePort,
	}

	out, err := runner.exec.Command("networksetup", args...).CombinedOutput()
	if err != nil {
		return "", err
	}

	output := string(out[:])

	return output, nil
}

func (runner *runner) getInfo(hardwarePort string) (info Info, err error) {
	args := []string{
		"-getinfo",
		hardwarePort,
	}

	out, err := runner.exec.Command("networksetup", args...).CombinedOutput()
	if err != nil {
		return Info{}, err
	}

	output := string(out[:])
	scanner := bufio.NewScanner(strings.NewReader(output))

	// IP address: 192.168.10.9
	// Subnet mask: 255.255.255.0
	// Router: 192.168.10.1
	// Client ID:
	// IPv6: Automatic
	// IPv6 IP address: none
	// IPv6 Router: none
	// Wi-Fi ID: 34:36:3b:cc:51:50

	for scanner.Scan() {
		line := scanner.Text()
		fmt.Println(line)
		if strings.Contains(line, "IP address") && !strings.Contains(line, "IPv6 IP") {
			info.IPAddress = strings.TrimSpace(strings.Split(line, ":")[1])
		} else if strings.Contains(line, "Subnet mask") {
			info.SubnetMask = strings.TrimSpace(strings.Split(line, ":")[1])
		} else if strings.Contains(line, "Router") && !strings.Contains(line, "IPv6 Router") {
			info.Router = strings.TrimSpace(strings.Split(line, ":")[1])
		} else if strings.Contains(line, "Client ID") {
			info.ClientID = strings.TrimSpace(strings.Split(line, ":")[1])
		} else if strings.Contains(line, "IPv6") {
			info.IPv6 = strings.TrimSpace(strings.Split(line, ":")[1])
		} else if strings.Contains(line, "IPv6 IP") {
			info.IPv6IPAddress = strings.TrimSpace(strings.Split(line, ":")[1])
		} else if strings.Contains(line, "IPv6 Router") {
			info.IPv6Router = strings.TrimSpace(strings.Split(line, ":")[1])
		} else if strings.Contains(line, "Wi-Fi ID") {
			info.WiFiID = strings.TrimSpace(strings.Join(strings.Split(line, ":")[1:], ":"))
		}
	}

	return info, nil
}
