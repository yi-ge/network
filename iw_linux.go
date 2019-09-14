package network

import (
	"regexp"
	"strings"
)

// InterfaceType .
type InterfaceType struct {
	Name string
	Type string
}

// ConnectWifi .
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

// DisconnectWifi .
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

func (runner *runner) getIwconfig() ([]InterfaceType, error) {
	out, err := runner.exec.Command("iwconfig").CombinedOutput()
	if err != nil {
		return []InterfaceType{}, err
	}

	output := string(out[:])

	repItem := regexp.MustCompile(`(?m)^\s*$[\r\n]*|[\r\n]+\s+\z`)

	itemOutIndex := repItem.FindStringSubmatchIndex(output)
	interfaceType := InterfaceType{}
	interfaceTypeList := []InterfaceType{}

	if len(itemOutIndex) != 0 {
		itemOut := output[:itemOutIndex[1]]
		output = output[itemOutIndex[1]:]

		if strings.Contains(itemOut, "no wireless") {
			fs := strings.Fields(itemOut)
			interfaceType.Name = fs[0]
			interfaceType.Type = "Wired"
			interfaceTypeList = append(interfaceTypeList, interfaceType)
			interfaceType = InterfaceType{}
		} else {
			fs := strings.Fields(itemOut)
			interfaceType.Name = fs[0]
			interfaceType.Type = "Wi-Fi"
			interfaceTypeList = append(interfaceTypeList, interfaceType)
			interfaceType = InterfaceType{}
		}

		itemOutIndex = repItem.FindStringSubmatchIndex(output)
	}

	return interfaceTypeList, nil
}
