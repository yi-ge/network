package network

func getDefaultGateway(interfaceName string) string {
	command := "route -n | grep " + interfaceName + " | grep UG | awk '{print $2}'"

	return command
}

// EnabledNetworkInterface .
func (runner *runner) EnabledNetworkInterface(interfaceName string) error {
	return runner.EnabledNetworkInterfaceByIfconfig(interfaceName)
}

// DisableNetworkInterface .
func (runner *runner) DisableNetworkInterface(interfaceName string) error {
	return runner.DisableNetworkInterfaceByIfconfig(interfaceName)
}
