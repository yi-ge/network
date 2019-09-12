package network

func (runner *runner) ConnectWifi(interfaceName string, ssid string, password string, securityType string, broadcast bool) (string, error) {
	args := []string{
		"-setairportnetwork",
		interfaceName,
		ssid,
		password,
	}

	out, err := runner.exec.Command("networksetup", args...).CombinedOutput()
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

	out, err := runner.exec.Command("networksetup", args...).CombinedOutput()
	if err != nil {
		return "", err
	}

	return string(out[:]), nil
}
