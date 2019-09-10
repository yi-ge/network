package network

import (
	"testing"
)

func TestParseIpconfig(t *testing.T) {
	XPOutput := "Windows IP Configuration\r\n        Host Name . . . . . . . . . . . . : win-5hDxCp2Miw8\r\n        Primary Dns Suffix  . . . . . . . :\r\n        Node Type . . . . . . . . . . . . : Unknown\r\n        IP Routing Enabled. . . . . . . . : No\r\n        WINS Proxy Enabled. . . . . . . . : No\r\nEthernet adapter 本地连接:\r\n        Connection-specific DNS Suffix  . :\r\n        Description . . . . . . . . . . . : Intel(R) PRO/1000 MT Network Connection\r\n        Physical Address. . . . . . . . . : 02-00-2D-98-0A-08\r\n        Dhcp Enabled. . . . . . . . . . . : No\r\n        IP Address. . . . . . . . . . . . : 10.53.10.8\r\n        Subnet Mask . . . . . . . . . . . : 255.255.0.0\r\n        Default Gateway . . . . . . . . . : 10.53.0.1\r\n        DNS Servers . . . . . . . . . . . : 61.147.37.1\r\n                                            8.8.8.8\r\n        NetBIOS over Tcpip. . . . . . . . : Disabled\r\n"

	ipconfigInterfaces := parseIpconfig(XPOutput)

	for _, ipconfigInterface := range ipconfigInterfaces {
		SmartPrint(ipconfigInterface)
	}
}
