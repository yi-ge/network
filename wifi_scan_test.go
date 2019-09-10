package network

import (
	"fmt"
	"testing"
)

func TestScanWIFI(t *testing.T) {
	networkInterface := New(nil)
	wifiList, err := networkInterface.ScanWIFI("WLAN 2")
	if err != nil {
		fmt.Println(err)
	}

	for _, wifi := range wifiList {
		SmartPrint(wifi)
	}
}
