# Network

[![GoDoc](https://godoc.org/github.com/yi-ge/network?status.svg)](https://godoc.org/github.com/yi-ge/network)

Cross-Platform Network manage, golang interface. Develop for IoT.

```golang
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
```

## Chinese Guide

注意：因此库尚未完整完成开发，并且大部分功能是调用各个平台上的命令实现，不保证该库的稳定及使用效果。

将来会完善测试及控制命令版本尽可能的实现稳定高效。

```golang
  // 获取接口列表，返回Interfaces结构体数组
  GetInterfacesList() ([]Interfaces, error)

  // 启用指定的网络接口
  EnabledNetworkInterface(interfaceName string) error

  // 禁用指定的网络接口
  DisableNetworkInterface(interfaceName string) error

  // 设置指定接口的IP地址
	SetStaticIP(interfaceName string, addr string, mask string, gateway string) error

  // 设置指定接口使用DHCP获取IP等信息
  SetInterfaceUseDHCP(interfaceName string) error

  // 设置指定接口的DNS
  SetDNS(interfaceName string, primaryAddr string, backAddr string) error

  // 设置指定接口使用DHCP获取DNS服务器信息
  SetDNSUseDHCP(interfaceName string) error

  // 扫描指定WIFI接口的网络信号
	ScanWIFI(wifiInterface ...string) (wifiList []Wifi, err error)

  // 连接指定的WIFI接口的指定信号
  ConnectWifi(interfaceName string, ssid string, password string, securityType string, broadcast bool) (string, error)

  // 断开指定WIFI接口的信号
  DisconnectWifi(interfaceName string) (string, error)
```
