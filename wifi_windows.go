package network

import (
	"bytes"
	"encoding/hex"
	"errors"
	"io"
	"os"
	"os/exec"
	"strings"
	"text/template"
)

// Template .
type Template struct {
	ProfileName     string // replace by SSIDString
	SsidStringToHex string // replace by hex SSIDString
	IsSSIDBroadcast bool   // replace by SSIDBroadcast
	SecAuth         string // replace by securityType, if securityType=None : securityType=open
	OpenPasscode    string // replace by "networkKey"
	WifiKey         string // replace by wifiKey
	Encryption      string // replace by encryption
	CaToTrust       string // replace by sha1 fingerprint
}

const wifiProfileTemplate = `<WLANProfile xmlns="http://www.microsoft.com/networking/WLAN/profile/v1">
  <name>{{.ProfileName}}</name>
  <SSIDConfig>
    <SSID>
      <hex>{{.SsidStringToHex}}</hex>
      <name>{{.ProfileName}}</name>
    </SSID>
    <nonBroadcast>{{.IsSSIDBroadcast}}</nonBroadcast>
  </SSIDConfig>
  <connectionType>ESS</connectionType>
  <connectionMode>manual</connectionMode>
  <autoSwitch>false</autoSwitch>
  <MSM>
    <security>
      <authEncryption>
        <authentication>{{.SecAuth}}</authentication>
        <encryption>{{.Encryption}}</encryption>
        <useOneX>false</useOneX>
        <FIPSMode xmlns="http://www.microsoft.com/networking/WLAN/profile/v2">false</FIPSMode>
      </authEncryption>
      <sharedKey>
        <keyType>{{.OpenPasscode}}</keyType>
        <protected>false</protected>
        <keyMaterial>{{.WifiKey}}</keyMaterial>
      </sharedKey>
    </security>
  </MSM>
</WLANProfile>`

func (runner *runner) SetWifiProfile(ssid string, securityType string, wifiKey string, ssidBroadcast bool) (msg string, err error) {
	var (
		templateToFile              string
		elementsToReplaceInTemplate Template
	)

	// Get SSID information
	ssidString := ssid
	ssidStringToHex := hex.EncodeToString([]byte(ssidString))

	if securityType == "None" {
		securityType = "open"
	}

	tempPath := os.Getenv("tmp")
	if tempPath == "" {
		return "", errors.New("Temp path can't found")
	}
	profileFile := tempPath + "\\template-out.xml"

	addWLANProfileCommand := exec.Command("netsh", "wlan", "add", "profile", "filename="+profileFile, "user=all")

	switch securityType {
	case "WEP":
		elementsToReplaceInTemplate = Template{
			ProfileName:     ssidString,
			SsidStringToHex: ssidStringToHex,
			IsSSIDBroadcast: ssidBroadcast,
			SecAuth:         "open",
			OpenPasscode:    "passPhrase",
			WifiKey:         wifiKey,
			Encryption:      "WEP",
		}
	case "WPA":
		elementsToReplaceInTemplate = Template{
			ProfileName:     ssidString,
			SsidStringToHex: ssidStringToHex,
			IsSSIDBroadcast: ssidBroadcast,
			SecAuth:         "WPA2PSK",
			OpenPasscode:    "passPhrase",
			WifiKey:         wifiKey,
			Encryption:      "AES",
		}
	default:
		elementsToReplaceInTemplate = Template{
			ProfileName:     ssidString,
			SsidStringToHex: ssidStringToHex,
			IsSSIDBroadcast: ssidBroadcast,
			SecAuth:         "open",
			OpenPasscode:    "passPhrase",
			WifiKey:         wifiKey,
			Encryption:      "none",
		}
	}
	templateToFile, err = executeTemplate("wireless Open template", wifiProfileTemplate, elementsToReplaceInTemplate)
	if err != nil {
		return "", err
	}

	err = createProfileFile(templateToFile)
	if err != nil {
		return "", err
	}

	return addProfileToMachine(profileFile, addWLANProfileCommand)
}

// Create, parse and execute templates
func executeTemplate(nameTemplate, constTemplate string, templateToApply Template) (string, error) {
	newTemplate := template.New(nameTemplate)
	// parses template
	newTemplate, err := newTemplate.Parse(constTemplate)
	if err != nil {
		os.Remove("profile.xml")
		return "", err
	}
	// executes the template into the open file
	var templateBuffer bytes.Buffer
	err = newTemplate.Execute(&templateBuffer, templateToApply)
	if err != nil {
		os.Remove("profile.xml")
		return templateBuffer.String(), err
	}
	// handles error
	if err != nil {
		os.Remove("profile.xml")
		return templateBuffer.String(), err
	}
	return templateBuffer.String(), nil
}

// Create and write profile file into templateToFile folder
func createProfileFile(templateToFile string) error {
	tempPath := os.Getenv("tmp")
	// create and open file
	profileFilePath := tempPath + "\\" + "template-out.xml"
	profileFile, err := os.Create(profileFilePath)
	if err != nil {
		os.Remove("profile.xml")
		return err
	}
	// close file
	defer profileFile.Close()
	// write the template into the new file
	_, err = io.Copy(profileFile, strings.NewReader(templateToFile))
	if err != nil {
		os.Remove("profile.xml")
		os.Remove(profileFilePath)
		return err
	}

	os.Remove("profile.xml")
	return nil
}

// Add wired and wireless profiles to Windows
func addProfileToMachine(profileFile string, cmd *exec.Cmd) (string, error) {
	out, err := cmd.CombinedOutput()
	output := ConvertToString(string(out[:]), "gbk", "utf8")

	if err != nil {
		os.Remove(profileFile)
		return "", err
	}

	os.Remove(profileFile)
	return string(output[:]), nil
}
