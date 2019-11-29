package shell

import (
	"encoding/json"
	"fmt"
	gp "github.com/c-bata/go-prompt"
	"os/exec"
	"time"
)

// naive cache which holds subscribers data for imsiFilterSuggestions
var devicesCache []gp.Suggest

var deviceIdFilterSuggestions = func(word, specifiedProfileName, specifiedCoverageType, providedAPIKey, providedAPIToken string) []gp.Suggest {
	c := make(chan []gp.Suggest, 1024)
	if len(devicesCache) == 0 {
		go getDevices(c, specifiedProfileName, specifiedCoverageType, providedAPIKey, providedAPIToken)
		select {
		case res := <-c:
			devicesCache = res
		case <-time.After(10 * time.Second):
			return []gp.Suggest{{
				Text:        "Downloading device information in background",
				Description: "Hit space to see latest",
			}}
		}
	}
	return filterFunc(devicesCache, word, filterTextOrDescriptionFuzzy)
}

var getDevices = func(c chan<- []gp.Suggest, specifiedProfileName, specifiedCoverageType, providedAPIKey, providedAPIToken string) {
	var r []gp.Suggest
	command := "soracom "
	if specifiedProfileName != "" {
		command = fmt.Sprintf("%s --profile %s ", command, specifiedProfileName)
	}
	if specifiedCoverageType != "" {
		command = fmt.Sprintf("%s --coverage-type %s ", command, specifiedCoverageType)
	}
	if providedAPIKey != "" {
		command = fmt.Sprintf("%s --api-key %s ", command, providedAPIKey)
	}
	if providedAPIToken != "" {
		command = fmt.Sprintf("%s --api-token %s ", command, providedAPIToken)
	}
	command = fmt.Sprintf("%s devices list --fetch-all", command)
	cmd := exec.Command("/bin/sh", "-c", command)
	//cmd := exec.Command("/bin/sh", "-c", "soracom subscribers list --fetch-all")
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		c <- []gp.Suggest{{
			Text:        "Error " + command,
			Description: err.Error(),
		}}
	}
	if err := cmd.Start(); err != nil {
		c <- []gp.Suggest{{
			Text:        "Error " + command,
			Description: err.Error(),
		}}
	}
	if err := json.NewDecoder(stdout).Decode(&devices); err != nil {
		c <- []gp.Suggest{{
			Text:        "Error " + command,
			Description: err.Error(),
		}}
	}
	if err := cmd.Wait(); err != nil {
		c <- []gp.Suggest{{
			Text:        "Error " + command,
			Description: err.Error(),
		}}
	}
	for _, device := range devices {
		online := "offline"
		if device.Online {
			online = "online"
		}
		r = append(r, gp.Suggest{
			Text: device.DeviceId,
			Description: fmt.Sprintf("%-14s | %-10s | %-7s | %-8s | %15s | %15s | %s",
				trunc(device.Endpoint, 14),
				device.Status,
				online,
				device.Imei,
				device.Imsi,
				device.IPAddress,
				device.Manufacturer,
			),
		})
	}
	c <- r
}
