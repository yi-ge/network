package network

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
