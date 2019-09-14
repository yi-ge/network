package network

import (
	"net"
	"strings"
	"sync"

	utilexec "github.com/yi-ge/network/exec"
)

// Interfaces - Network interface struct.
type Interfaces struct {
	Name                  string // e.g., "en0", "lo0", "eth0.100"
	HardwareAddr          string // IEEE MAC-48, EUI-48 and EUI-64 form
	Type                  string // Wired, Wi-Fi
	DHCPEnabled           bool
	IPv4Address           string
	SubnetPrefix          string
	DefaultGatewayAddress string
	DNSPrimary            string
	DNSBack               string
	Description           string // In mac is HardwarePort, other is net card name
	Connected             bool
	Mode                  string
	AdminState            string
}

// runner implements Interface in terms of exec("ipconfig").
type runner struct {
	mu   sync.Mutex
	exec utilexec.Interface
}

// Interface - Network manage interface.
type Interface interface {
	GetInterfacesList() ([]Interfaces, error)
	EnabledNetworkInterface(interfaceName string) error
	DisableNetworkInterface(interfaceName string) error
	SetStaticIP(interfaceName string, addr string, mask string, gateway string) error
	SetInterfaceUseDHCP(interfaceName string) error
	SetDNS(interfaceName string, primaryAddr string, backAddr string) error
	SetDNSUseDHCP(interfaceName string) error
	ScanWIFI(wifiInterface ...string) (wifiList []Wifi, err error)
	ConnectWifi(interfaceName string, ssid string, password string, securityType string, broadcast bool) (string, error)
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
						interfacesOut.SubnetPrefix = hex2dot(ip.Mask.String())
					}
				}
			}

			interfacesListOut = append(interfacesListOut, interfacesOut)
		}
	}

	return runner.parseInterfacesList(interfacesListOut), nil
}
