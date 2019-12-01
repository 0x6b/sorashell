package shell

import (
	"encoding/json"
	"fmt"
	gp "github.com/c-bata/go-prompt"
	"strings"
	"time"
)

// naive cache which holds subscribers data for imsiFilterSuggestions
var devicesCache []gp.Suggest

var deviceIdFilterSuggestions = func(word string, worker *SoracomWorker) []gp.Suggest {
	c := make(chan []gp.Suggest, 1024)
	if len(devicesCache) == 0 {
		go getDevices(c, worker)
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

var getDevices = func(c chan<- []gp.Suggest, worker *SoracomWorker) {
	var r []gp.Suggest

	result := worker.Execute("devices list --fetch-all")
	if err := json.NewDecoder(strings.NewReader(result)).Decode(&devices); err != nil {
		c <- []gp.Suggest{{
			Text:        "Error while running 'devices list --fetch-all'",
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
			Description: fmt.Sprintf("%-14s | %-8s | %-8s | %-15s | %15s | %15s | %s",
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
