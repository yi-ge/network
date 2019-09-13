// https://github.com/schollz/wifiscan

package network

import (
	"bytes"
	"runtime"
	"strconv"
	"strings"
	"time"
)

// TimeLimit .
var TimeLimit = 10 * time.Second

// ScanWIFI will scan the optional interface for wifi access points
func (runner *runner) ScanWIFI(wifiInterface ...string) (wifiList []Wifi, err error) {
	if runtime.GOOS == "linux" && (len(wifiInterface) == 0 || wifiInterface[0] == "") {
		var interfaces []string
		interfaces, err = runner.getInterfacesLinux()
		if err != nil {
			return []Wifi{}, err
		}
		for _, in := range interfaces {
			var w []Wifi
			w, err = runner.scan(in)
			if len(w) > 0 {
				wifiList = append(wifiList, w...)
			}
		}
		if len(wifiList) > 0 {
			err = nil
			wifiMap := make(map[string]Wifi)
			for _, w := range wifiList {
				wifiMap[w.SSID] = w
			}
			i := 0
			for _, w := range wifiMap {
				wifiList[i] = w
				i++
			}
			wifiList = wifiList[:i]
		}

		return
	}
	return runner.scan(wifiInterface...)
}

func (runner *runner) scan(wifiInterface ...string) (wifiList []Wifi, err error) {
	command := ""
	os := ""
	switch runtime.GOOS {
	case "windows":
		os = "windows"
		command = "netsh wlan show networks interface=\"" + wifiInterface[0] + "\" mode=Bssid"
	case "darwin":
		os = "darwin"
		command = "/System/Library/PrivateFrameworks/Apple80211.framework/Versions/Current/Resources/airport \"" + wifiInterface[0] + "\" -s"
	default:
		os = "linux"
		command = "iwlist \"" + wifiInterface[0] + "\" scan"
		// if len(wifiInterface) > 0 && len(wifiInterface[0]) > 0 {
		// 	command = fmt.Sprintf("iwlist %s scan", wifiInterface[0])
		// }
	}
	stdout, _, err := runner.runCommand(TimeLimit, command)
	if err != nil {
		return []Wifi{}, err
	}
	wifiList, err = Parse(stdout, os)
	return
}

func (runner *runner) runCommand(tDuration time.Duration, commands string) (stdout, stderr string, err error) {
	command := strings.Fields(commands)
	cmd := runner.exec.Command(command[0])
	if len(command) > 0 {
		cmd = runner.exec.Command(command[0], command[1:]...)
	}
	var stdOut, stdErr bytes.Buffer
	cmd.SetStdout(&stdOut)
	cmd.SetStderr(&stdErr)
	err = cmd.Start()
	if err != nil {
		return
	}
	done := make(chan error, 1)
	go func() {
		done <- cmd.Wait()
	}()
	select {
	case <-time.After(tDuration):
		err = cmd.GetProcess().Kill()
	case err = <-done:
		if runtime.GOOS == "windows" {
			stdout = ConvertToString(stdOut.String(), "gbk", "utf8")
			stderr = ConvertToString(stdErr.String(), "gbk", "utf8")
		} else {
			stdout = stdOut.String()
			stderr = stdErr.String()
		}
	}
	return
}

func (runner *runner) getInterfacesLinux() (interfaces []string, err error) {
	stdout, _, err := runner.runCommand(TimeLimit, "ip address")
	if err != nil {
		return
	}
	return getInteracesFromString(stdout)
}

func getInteracesFromString(s string) (interfaces []string, err error) {
	for _, line := range strings.Split(s, "\n") {
		if !strings.Contains(line, "BROADCAST") {
			continue
		}
		cols := strings.Split(line, ":")
		if len(cols) < 3 {
			continue
		}
		_, errConvert := strconv.Atoi(cols[0])
		if errConvert != nil {
			continue
		}
		if strings.Contains(cols[1], "@") || strings.Contains(cols[1], "docker") {
			continue
		}
		interfaces = append(interfaces, strings.TrimSpace(cols[1]))
	}
	return
}
