package network

import (
	"bufio"
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"
)

// Wifi is the data structure containing the basic
// elements
type Wifi struct {
	SSID           string `json:"ssid"`
	NetworkType    string `json:"networkType"`
	Authentication string
	Encryption     string
	ChannelList    []ChannelInfo
}

// ChannelInfo .
type ChannelInfo struct {
	BSSID     string `json:"bssid"`
	Signal    string `json:"signal"`
	RSSI      int    `json:"rssi"`
	RadioType string
	Channel   string
}

// Parse will parse wifi output and extract the access point
// information.
func Parse(output, os string) (wifiList []Wifi, err error) {
	switch os {
	case "windows":
		wifiList, err = parseWindows(output)
	case "darwin":
		wifiList, err = parseDarwin(output)
	case "linux":
		wifiList, err = parseLinux(output)
	default:
		err = fmt.Errorf("%s is not a recognized OS", os)
	}
	return
}

func parseWindows(output string) (wifiList []Wifi, err error) {
	scanner := bufio.NewScanner(strings.NewReader(output))
	wifiList = []Wifi{}
	w := Wifi{}
	channelInfo := ChannelInfo{}
	channelInfoList := []ChannelInfo{}

	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, "SSID") && !strings.Contains(line, "BSSID") {
			if w.SSID != "" {
				if channelInfo.BSSID != "" && channelInfo.RSSI != 0 {
					channelInfoList = append(channelInfoList, channelInfo)
					channelInfo = ChannelInfo{}
				}
				w.ChannelList = channelInfoList
				wifiList = append(wifiList, w)
				w = Wifi{}
				channelInfo = ChannelInfo{}
				channelInfoList = []ChannelInfo{}
			}

			fs := strings.Fields(line)
			if len(fs) >= 3 {
				ssid, errOut := hex.DecodeString(strings.Join(fs[3:], " "))
				if errOut != nil {
					w.SSID = strings.Join(fs[3:], " ")
				} else {
					w.SSID = string(ssid[:])
				}
			}
		} else if strings.Contains(line, "Network type") {
			fs := strings.Fields(line)
			if len(fs) == 4 {
				w.NetworkType = fs[3]
			}
		} else if strings.Contains(line, "Authentication") {
			fs := strings.Fields(line)
			if len(fs) == 3 {
				w.Authentication = fs[2]
			}
		} else if strings.Contains(line, "Encryption") {
			fs := strings.Fields(line)
			if len(fs) == 3 {
				w.Encryption = fs[2]
			}
		} else if strings.Contains(line, "BSSID") {
			if channelInfo.BSSID != "" && channelInfo.RSSI != 0 {
				channelInfoList = append(channelInfoList, channelInfo)
				channelInfo = ChannelInfo{}
			}
			fs := strings.Fields(line)
			if len(fs) == 4 {
				channelInfo.BSSID = fs[3]
			}
		} else if strings.Contains(line, "Signal") {
			if strings.Contains(line, "%") {
				fs := strings.Fields(line)
				if len(fs) == 3 {
					channelInfo.Signal = fs[2]
					channelInfo.RSSI, err = strconv.Atoi(strings.Replace(fs[2], "%", "", 1))
					if err != nil {
						continue
					}
					channelInfo.RSSI = (channelInfo.RSSI / 2) - 100
				}
			}
		} else if strings.Contains(line, "Radio type") {
			fs := strings.Fields(line)
			if len(fs) == 4 {
				channelInfo.RadioType = fs[3]
			}
		} else if strings.Contains(line, "Channel") {
			fs := strings.Fields(line)
			if len(fs) == 3 {
				channelInfo.Channel = fs[2]
			}
		}
	}

	if channelInfo.BSSID != "" && channelInfo.RSSI != 0 {
		channelInfoList = append(channelInfoList, channelInfo)
		channelInfo = ChannelInfo{}
	}

	if w.SSID != "" {
		w.ChannelList = channelInfoList
		wifiList = append(wifiList, w)
		w = Wifi{}
		channelInfoList = []ChannelInfo{}
	}

	return
}

func ssidInArray(wifiList []Wifi, ssid string) (bool, int) {
	for index, wifi := range wifiList {
		if wifi.SSID == ssid {
			return true, index
		}
	}

	return false, -1
}

func parseDarwin(output string) (wifiList []Wifi, err error) {
	scanner := bufio.NewScanner(strings.NewReader(output))
	wifiList = []Wifi{}
	for scanner.Scan() {
		line := scanner.Text()
		fs := strings.Fields(line)
		if len(fs) < 6 {
			continue
		}
		rssi, errParse := strconv.Atoi(fs[2])
		if errParse != nil {
			continue
		}
		if rssi > 0 {
			continue
		}

		radioType := "802.11n"
		channel, err := strconv.Atoi(fs[3])
		if err == nil && channel > 14 {
			radioType = "802.11ac"
		}

		signal := "0%"
		signalPer := (rssi + 100) * 2
		if signalPer > 100 {
			signal = "100%"
		} else if signalPer < 0 {
			signal = "0%"
		} else {
			signal = strconv.Itoa(signalPer) + "%"
		}

		channelInfo := ChannelInfo{
			BSSID:     strings.ToLower(fs[1]),
			Signal:    signal,
			RSSI:      rssi,
			RadioType: radioType,
			Channel:   fs[3],
		}

		ssid := fs[0]

		inArray, index := ssidInArray(wifiList, ssid)

		if inArray {
			wifiList[index].ChannelList = append(wifiList[index].ChannelList, channelInfo)
		} else {
			currentWifi := Wifi{}
			currentWifi.SSID = ssid
			currentWifi.NetworkType = "Infrastructure" // or Ad-hoc/AP/Station
			currentWifi.Authentication = strings.Join(fs[6:], " ")
			if strings.Contains(strings.Join(fs[6:], " "), "AES") {
				currentWifi.Encryption = "CCMP"
			} else if strings.Contains(strings.Join(fs[6:], " "), "TKIP") {
				currentWifi.Encryption = "TKIP"
			} else {
				currentWifi.Encryption = "" // WEP/TKIP/Other
			}
			currentWifi.ChannelList = append(currentWifi.ChannelList, channelInfo)

			wifiList = append(wifiList, currentWifi)
		}
	}
	return
}

func parseLinux(output string) (wifiList []Wifi, err error) {
	scanner := bufio.NewScanner(strings.NewReader(output))
	w := Wifi{}
	channelInfo := ChannelInfo{}
	wifiList = []Wifi{}
	for scanner.Scan() {
		line := scanner.Text()

		if strings.Contains(line, "Address") {
			fs := strings.Split(line, ":")
			bssid := strings.TrimSpace(strings.ToLower(fs[2]))
			line = scanner.Text()
			if strings.Contains(line, "ESSID") {
				if channelInfo.BSSID != "" {
					inArray, index := ssidInArray(wifiList, line)
					if inArray {
						wifiList[index].ChannelList = append(wifiList[index].ChannelList, channelInfo)
					} else {
						w.ChannelList = append(w.ChannelList, channelInfo)
						wifiList = append(wifiList, w)
						w = Wifi{}
					}

					channelInfo = ChannelInfo{}
				}

				fs := strings.Split(line, ":")
				w.SSID = fs[2]
				channelInfo.BSSID = bssid
			} else {
				continue
			}
		} else if strings.Contains(line, "Mode") {
			fs := strings.Split(line, ":")
			if fs[2] == "Managed" {
				w.NetworkType = "Infrastructure"
			} else {
				w.NetworkType = fs[2]
			}
		} else if strings.Contains(line, "IE: IEEE") {
			if strings.Contains(line, "WPA2") {
				w.Authentication = "WPA2"
			}
		} else if strings.Contains(line, "Pairwise Ciphers") {
			if strings.Contains(line, "CCMP") {
				w.Encryption = "CCMP"
			} else if strings.Contains(line, "TKIP") {
				w.Encryption = "TKIP"
			} else {
				w.Encryption = ""
			}
		} else if strings.Contains(line, "Encryption key") {
			fs := strings.Split(line, ":")
			if fs[2] == "off" {
				w.Encryption = "None"
				w.Authentication = "None"
			}
		} else if strings.Contains(line, "Signal level=") {
			signal := strings.Split(strings.Split(strings.Split(line, "level=")[1], "/")[0], " dB")[0]
			channelInfo.Signal = signal
			level, errParse := strconv.Atoi(signal)
			if errParse != nil {
				continue
			}
			if level > 0 {
				level = (level / 2) - 100
			}
			channelInfo.RSSI = level
		} else if strings.Contains(line, "Frequency=5") {
			channelInfo.RadioType = "802.11ac"
		} else if strings.Contains(line, "Frequency=2") {
			channelInfo.RadioType = "802.11n"
		} else if strings.Contains(line, "(Channel") {
			channelInfo.Channel = strings.TrimSpace(strings.Replace(strings.Split(line, "(Channel")[2], ")", "", -1))
		}
	}

	if channelInfo.BSSID != "" {
		inArray, index := ssidInArray(wifiList, w.SSID)
		if inArray {
			wifiList[index].ChannelList = append(wifiList[index].ChannelList, channelInfo)
		} else {
			w.ChannelList = append(w.ChannelList, channelInfo)
			wifiList = append(wifiList, w)
			w = Wifi{}
		}

		channelInfo = ChannelInfo{}
	}
	return wifiList, nil
}
