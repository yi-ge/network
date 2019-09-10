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

// EnableNetworkInterface .
func (runner *runner) EnableNetworkInterface(interfaceName string) error {
	return runner.EnableNetworkInterfaceByIfconfig(interfaceName)
}

// DisabledNetworkInterface .
func (runner *runner) DisabledNetworkInterface(interfaceName string) error {
	return runner.DisabledNetworkInterfaceByIfconfig(interfaceName)
}
