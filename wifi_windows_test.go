package network

import (
	"testing"
)

func TestSetWifiProfile(t *testing.T) {
	networkInterface := New(nil)
	msg, err := networkInterface.ConnectWifi("WLAN 2", "ProjectX", "nideshengri", "WPA", true)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(msg)
}
