package network

import (
	"bufio"
	"strings"
)

// InterfaceConnected .
type InterfaceConnected struct {
	Name      string
	Connected bool
}

func (runner *runner) getInterfaceConnected() ([]InterfaceConnected, error) {
	out, err := runner.exec.Command("ip", "link").CombinedOutput()
	if err != nil {
		return []InterfaceConnected{}, err
	}

	output := string(out[:])
	interfaceConnected := InterfaceConnected{}
	interfaceConnectedList := []InterfaceConnected{}

	scanner := bufio.NewScanner(strings.NewReader(output))

	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, "mtu") {
			fs := strings.Fields(line)
			interfaceConnected.Name = strings.Replace(fs[1], ":", "", -1)
			interfaceConnected.Connected = strings.Contains(line, "state UP")
			interfaceConnectedList = append(interfaceConnectedList, interfaceConnected)
			interfaceConnected = InterfaceConnected{}
		}
	}

	return interfaceConnectedList, nil
}
