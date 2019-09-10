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
	Channel   int
}

// Parse will parse wifi output and extract the access point
// information.
func Parse(output, os string) (wifiList []Wifi, err error) {
	switch os {
	case "windows":
		wifiList, err = parseWindows(output)
	case "darwin":
		// wifiList, err = parseDarwin(output)
	case "linux":
		// wifiList, err = parseLinux(output)
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
				channel, err := strconv.Atoi(fs[2])
				if err != nil {
					continue
				}
				channelInfo.Channel = channel
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

// func parseDarwin(output string) (wifiList []Wifi, err error) {
// 	scanner := bufio.NewScanner(strings.NewReader(output))
// 	wifiList = []Wifi{}
// 	for scanner.Scan() {
// 		line := scanner.Text()
// 		fs := strings.Fields(line)
// 		if len(fs) < 6 {
// 			continue
// 		}
// 		rssi, errParse := strconv.Atoi(fs[2])
// 		if errParse != nil {
// 			continue
// 		}
// 		if rssi > 0 {
// 			continue
// 		}
// 		wifiList = append(wifiList, Wifi{SSID: strings.ToLower(fs[1]), RSSI: rssi})
// 	}
// 	return
// }

// func parseLinux(output string) (wifiList []Wifi, err error) {
// 	scanner := bufio.NewScanner(strings.NewReader(output))
// 	w := Wifi{}
// 	wifiList = []Wifi{}
// 	for scanner.Scan() {
// 		line := scanner.Text()
// 		if w.SSID == "" {
// 			if strings.Contains(line, "Address") {
// 				fs := strings.Fields(line)
// 				if len(fs) == 5 {
// 					w.SSID = strings.ToLower(fs[4])
// 				}
// 			} else {
// 				continue
// 			}
// 		} else {
// 			if strings.Contains(line, "Signal level=") {
// 				level, errParse := strconv.Atoi(strings.Split(strings.Split(strings.Split(line, "level=")[1], "/")[0], " dB")[0])
// 				if errParse != nil {
// 					continue
// 				}
// 				if level > 0 {
// 					level = (level / 2) - 100
// 				}
// 				w.RSSI = level
// 			}
// 		}
// 		if w.SSID != "" && w.RSSI != 0 {
// 			wifiList = append(wifiList, w)
// 			w = Wifi{}
// 		}
// 	}
// 	return
// }
