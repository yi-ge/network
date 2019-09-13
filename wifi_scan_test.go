package network

import (
	"fmt"
	"testing"
)

func TestScanWIFI(t *testing.T) {
	networkInterface := New(nil)
	// wifiList, err := networkInterface.ScanWIFI("WLAN 2")
	wifiList, err := networkInterface.ScanWIFI("en0")
	if err != nil {
		fmt.Println(err)
	}

	t.Log(len(wifiList))

	for _, wifi := range wifiList {
		SmartPrint(wifi)
	}
}
