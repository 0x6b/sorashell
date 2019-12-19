package sorashell

import (
	"encoding/json"
	"fmt"
	"github.com/c-bata/go-prompt"
	"strings"
	"time"
)

// naive cache which holds subscribers data for imsiFilterSuggestions
var inventoryDevicesCache []prompt.Suggest

func (s *SoracomCompleter) inventoryDeviceIdFilterSuggestions(word string) []prompt.Suggest {
	c := make(chan []prompt.Suggest, 1024)
	if len(inventoryDevicesCache) == 0 {
		go getInventoryDevices(c, s.worker)
		select {
		case res := <-c:
			inventoryDevicesCache = res
		case <-time.After(10 * time.Second):
			return []prompt.Suggest{{
				Text:        "Downloading device information in background",
				Description: "Hit space to see latest",
			}}
		}
	}
	return filterFunc(inventoryDevicesCache, word, filterTextOrDescriptionFuzzy)
}

var getInventoryDevices = func(c chan<- []prompt.Suggest, worker *SoracomWorker) {
	var r []prompt.Suggest

	result := worker.Execute("devices list --fetch-all")
	if err := json.NewDecoder(strings.NewReader(result)).Decode(&inventoryDevices); err != nil {
		c <- []prompt.Suggest{{
			Text:        "Error while running 'devices list --fetch-all'",
			Description: err.Error(),
		}}
	}
	for _, device := range inventoryDevices {
		online := "offline"
		if device.Online {
			online = "online"
		}
		r = append(r, prompt.Suggest{
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
