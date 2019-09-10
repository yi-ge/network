package network

import (
	"testing"
)

func TestSetWifiProfile(t *testing.T) {
	networkInterface := New(nil)
	msg, err := networkInterface.SetWifiProfile("ProjectX", "WPA", "nideshengri", true)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(msg)

	networkInterface.ConnectWifi("WLAN 2", "ProjectX", "ProjectX")
}
