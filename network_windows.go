package network

func (runner *runner) parseInterfacesList(interfacesList []Interfaces) []Interfaces {
	output, err := runner.getIpconfigOutPut()
	if err != nil {
		return interfacesList
	}

	ipconfigInterfacesList := parseIpconfig(output)

	ipconfigInterfacesList = runner.parseInterfaceStatus(ipconfigInterfacesList)

	for _, ipconfigInterfaces := range ipconfigInterfacesList {
		have := false
		for inx, interfaces := range interfacesList {
			if interfaces.HardwareAddr == ipconfigInterfaces.HardwareAddr {
				interfacesList[inx].Type = ipconfigInterfaces.Type
				interfacesList[inx].DefaultGatewayAddress = ipconfigInterfaces.DefaultGatewayAddress
				interfacesList[inx].DNSPrimary = ipconfigInterfaces.DNSPrimary
				interfacesList[inx].DNSBack = ipconfigInterfaces.DNSBack
				interfacesList[inx].Mode = ipconfigInterfaces.Mode
				interfacesList[inx].Description = ipconfigInterfaces.Description
				interfacesList[inx].Connected = ipconfigInterfaces.Connected
				interfacesList[inx].AdminState = ipconfigInterfaces.AdminState
				have = true
			}
		}

		if !have {
			currentInterface := Interfaces{
				Name:       ipconfigInterfaces.Name,
				Mode:       ipconfigInterfaces.Mode,
				AdminState: ipconfigInterfaces.AdminState,
				Connected:  ipconfigInterfaces.Connected,
			}
			interfacesList = append(interfacesList, currentInterface)
		}
	}

	return interfacesList
}

func (runner *runner) parseInterfaceStatus(ipconfigInterfacesList []IpconfigInterfaces) []IpconfigInterfaces {
	maps, err := runner.getNetworkInterfaceStatus()
	if err != nil {
		return ipconfigInterfacesList
	}

	for inx, interfaces := range ipconfigInterfacesList {
		if val, ok := maps[interfaces.Name]; ok {
			ipconfigInterfacesList[inx].Mode = val.Mode
			ipconfigInterfacesList[inx].AdminState = val.AdminState
			ipconfigInterfacesList[inx].Connected = val.Connected
			delete(maps, interfaces.Name)
		}
	}

	for _, item := range maps {
		newIpconfigInterface := IpconfigInterfaces{
			Name:       item.Name,
			Mode:       item.Mode,
			AdminState: item.AdminState,
			Connected:  item.Connected,
		}
		ipconfigInterfacesList = append(ipconfigInterfacesList, newIpconfigInterface)
	}

	return ipconfigInterfacesList
}

func (runner *runner) EnableNetworkInterface(interfaceName string) error {
	return runner.enableNetworkInterfaceByNetsh(interfaceName)
}

func (runner *runner) DisabledNetworkInterface(interfaceName string) error {
	return runner.disabledNetworkInterfaceByNetsh(interfaceName)
}

func (runner *runner) SetStaticIP(interfaceName string, addr string, mask string, gateway string) error {
	return runner.setStaticIPByNetsh(interfaceName, addr, mask, gateway)
}

func (runner *runner) SetInterfaceUseDHCP(interfaceName string) error {
	return runner.setInterfaceUseDHCPByNetsh(interfaceName)
}

func (runner *runner) SetDNS(interfaceName string, primaryAddr string, backAddr string) error {
	return runner.setDNSByNetsh(interfaceName, primaryAddr, backAddr)
}

func (runner *runner) SetDNSUseDHCP(interfaceName string) error {
	return runner.setDNSUseDHCPByNetsh(interfaceName)
}
