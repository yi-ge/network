package network

import (
	"fmt"
	"testing"
)

func TestGetInterfacesList(t *testing.T) {
	networkInterface := New(nil)
	interfacesList, err := networkInterface.GetInterfacesList()

	if err != nil {
		fmt.Println(err)
	}

	for _, interfaces := range interfacesList {
		SmartPrint(interfaces)
	}
}
