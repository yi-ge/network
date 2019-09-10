package network

func getDefaultGateway(interfaceName string) string {
	command := "route -n | grep " + interfaceName + " | grep UG | awk '{print $2}'"

	return command
}

// EnableNetworkInterface .
func (runner *runner) EnableNetworkInterface(interfaceName string) error {
	return runner.EnableNetworkInterfaceByIfconfig(interfaceName)
}

// DisabledNetworkInterface .
func (runner *runner) DisabledNetworkInterface(interfaceName string) error {
	return runner.DisabledNetworkInterfaceByIfconfig(interfaceName)
}
