package network

func (runner *runner) parseInterfacesList(interfacesList []Interfaces) []Interfaces {
	output, err := runner.getIfconfigOutPut()
	if err != nil {
		return interfacesList
	}

	ipconfigInterfacesList := parseIfconfig(output)

	for inx, interfaces := range interfacesList {
		for _, ipconfigInterfaces := range ipconfigInterfacesList {
			if interfaces.HardwareAddr == ipconfigInterfaces.HardwareAddr {
				interfacesList[inx].DefaultGatewayAddress = ipconfigInterfaces.DefaultGatewayAddress
				interfacesList[inx].DNSPrimary = ipconfigInterfaces.DNSPrimary
				interfacesList[inx].DNSBack = ipconfigInterfaces.DNSBack
				interfacesList[inx].Connected = ipconfigInterfaces.Connected
			}
		}
	}

	return interfacesList
}

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
