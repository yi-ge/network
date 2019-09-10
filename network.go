package network

import (
	utilexec "network/exec"
	"net"
	"strings"
	"sync"
)

// Interfaces .
type Interfaces struct {
	Name                  string // e.g., "en0", "lo0", "eth0.100"
	HardwareAddr          string // IEEE MAC-48, EUI-48 and EUI-64 form
	Type                  string // Ethernet, Wireless LAN
	DHCPEnabled           bool
	IPv4Address           string
	SubnetPrefix          string
	DefaultGatewayAddress string
	DNSPrimary            string
	DNSBack               string
	Description           string
	Connected             bool
	Mode                  string
	AdminState            string
}

// runner implements Interface in terms of exec("ipconfig").
type runner struct {
	mu   sync.Mutex
	exec utilexec.Interface
}

// Interface .
type Interface interface {
	GetInterfacesList() ([]Interfaces, error)
	EnableNetworkInterface(interfaceName string) error
	DisabledNetworkInterface(interfaceName string) error
	ScanWIFI(wifiInterface ...string) (wifiList []Wifi, err error)
	SetWifiProfile(ssid string, securityType string, wifiKey string, ssidBroadcast bool) (msg string, err error)
	ConnectWifi(interfaceName string, name string, ssid string) (string, error)
	DisconnectWifi(interfaceName string) (string, error)
}

// New returns a new Interface which will exec ipconfig.
func New(exec utilexec.Interface) Interface {
	if exec == nil {
		exec = utilexec.New()
	}
	runner := &runner{
		exec: exec,
	}
	return runner
}

// GetInterfacesList .
func (runner *runner) GetInterfacesList() ([]Interfaces, error) {
	interfacesList, err := net.Interfaces()
	if err != nil {
		return []Interfaces{}, err
	}

	interfacesListOut := []Interfaces{}
	for _, interfaces := range interfacesList {
		if interfaces.HardwareAddr != nil { // 排除无用网卡 (interfaces.Flags&net.FlagUp) != 0 &&
			addrs, err := interfaces.Addrs()
			if err != nil {
				return []Interfaces{}, err
			}

			interfacesOut := Interfaces{}
			interfacesOut.Name = interfaces.Name
			interfacesOut.HardwareAddr = strings.ToUpper(interfaces.HardwareAddr.String())

			for _, addr := range addrs {
				if ip, ok := addr.(*net.IPNet); ok && !ip.IP.IsLoopback() {
					if ip.IP.To4() != nil {
						interfacesOut.IPv4Address = ip.IP.String()
						interfacesOut.SubnetPrefix = ip.Mask.String()
					}
				}
			}

			interfacesListOut = append(interfacesListOut, interfacesOut)
		}
	}

	return runner.parseInterfacesList(interfacesListOut), nil
}
