package network

// https://github.com/rakelkar/gonetsh/blob/master/netsh/netsh.go
import (
	"errors"
	"regexp"
	"strings"
)

// InterfaceStatus .
type InterfaceStatus struct {
	Name       string
	Mode       string
	AdminState string
	Connected  bool
}

func (runner *runner) getNetworkInterfaceStatus() (map[string]InterfaceStatus, error) {
	args := []string{
		"interface", "show", "interface",
	}

	output, err := runner.exec.Command("netsh", args...).CombinedOutput()
	outputString := string(output[:])
	// outputString := ConvertToString(string(output[:]), "gbk", "utf8")

	if err != nil {
		return nil, err
	}

	// Split output by line
	outputString = strings.TrimSpace(outputString)
	var outputLines = strings.Split(outputString, "\n")

	if len(outputLines) < 2 {
		return nil, errors.New("unexpected netsh output:\n" + outputString)
	}

	// Remove first two lines of header text
	outputLines = outputLines[2:]

	indexMap := make(map[string]InterfaceStatus)

	reg := regexp.MustCompile("\\s{2,}")

	for _, line := range outputLines {

		line = strings.TrimSpace(line)

		// Split the line by two or more whitespace characters, returning all substrings (n < 0)
		splitLine := reg.Split(line, -1)

		name := strings.Join(splitLine[3:], " ")
		currentInterfaceStatus := InterfaceStatus{
			Name:       name,
			Mode:       splitLine[2],
			AdminState: splitLine[0],
			Connected:  splitLine[1] == "Connected",
		}
		indexMap[name] = currentInterfaceStatus
	}

	return indexMap, nil
}

func (runner *runner) enableNetworkInterfaceByNetsh(interfaceName string) error {
	args := []string{
		"interface", "set", "interface", "name=\"" + interfaceName + "\"", "enabled",
	}

	_, err := runner.exec.Command("netsh", args...).CombinedOutput()
	if err != nil {
		return err
	}

	return nil
}

func (runner *runner) disabledNetworkInterfaceByNetsh(interfaceName string) error {
	args := []string{
		"interface", "set", "interface", "name=\"" + interfaceName + "\"", "disabled",
	}

	_, err := runner.exec.Command("netsh", args...).CombinedOutput()
	if err != nil {
		return err
	}

	return nil
}

func (runner *runner) setStaticIPByNetsh(interfaceName string, addr string, mask string, gateway string) error {
	args := []string{
		"interface", "ip", "set", "address", "name=\"" + interfaceName + "\"", "source=static", "addr=" + addr, "mask=" + mask, "gateway=" + gateway, "gwmetric=auto",
	}

	_, err := runner.exec.Command("netsh", args...).CombinedOutput()
	if err != nil {
		return err
	}

	return nil
}

func (runner *runner) setInterfaceUseDHCPByNetsh(interfaceName string) error {
	args := []string{
		"interface", "ip", "set", "address", "name=\"" + interfaceName + "\"", "source=dhcp",
	}

	_, err := runner.exec.Command("netsh", args...).CombinedOutput()
	if err != nil {
		return err
	}

	return nil
}

func (runner *runner) setDNSByNetsh(interfaceName string, primaryAddr string, backAddr string) error {
	primaryArgs := []string{
		"interface", "ip", "set", "dns", "name=\"" + interfaceName + "\"", "source=static", "addr=" + primaryAddr, "register=primary",
	}

	backArgs := []string{
		"interface", "ip", "add", "dns", "name=\"" + interfaceName + "\"", "source=static", "addr=" + backAddr, "index=2",
	}

	_, err := runner.exec.Command("netsh", primaryArgs...).CombinedOutput()
	if err != nil {
		return err
	}
	_, err = runner.exec.Command("netsh", backArgs...).CombinedOutput()

	if err != nil {
		return err
	}

	return nil
}

func (runner *runner) setDNSUseDHCPByNetsh(interfaceName string) error {
	args := []string{
		"interface",
		"ip",
		"set",
		"dns",
		"name=\"" + interfaceName + "\"",
		"source=dhcp",
	}

	_, err := runner.exec.Command("netsh", args...).CombinedOutput()
	if err != nil {
		return err
	}

	return nil
}

func (runner *runner) getNetworkInterfaceDriversInfo(interfaceName string) (string, error) {
	args := []string{
		"wlan show drivers interface=\"" + interfaceName + "\"",
	}

	out, err := runner.exec.Command("netsh", args...).CombinedOutput()
	if err != nil {
		return "", err
	}

	return string(out[:]), nil
}

func (runner *runner) ConnectWifi(interfaceName string, ssid string, password string, securityType string, broadcast bool) (string, error) {
	msg, err := runner.SetWifiProfile(ssid, securityType, password, broadcast)
	if err != nil {
		return "", err
	}

	if strings.Contains(msg, "err") {
		return "", errors.New(msg)
	}

	name := ssid
	args := []string{
		"wlan",
		"connect",
		"name=\"" + name + "\"",
		"ssid=\"" + ssid + "\"",
		"interface=\"" + interfaceName + "\"",
	}

	out, err := runner.exec.Command("netsh", args...).CombinedOutput()
	if err != nil {
		return "", err
	}

	return string(out[:]), nil
}

func (runner *runner) DisconnectWifi(interfaceName string) (string, error) {
	args := []string{
		"wlan",
		"disconnect",
		"interface=\"" + interfaceName + "\"",
	}

	out, err := runner.exec.Command("netsh", args...).CombinedOutput()
	if err != nil {
		return "", err
	}

	return string(out[:]), nil
}
