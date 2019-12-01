package shell

import (
	"encoding/json"
	"fmt"
	gp "github.com/c-bata/go-prompt"
	"strings"
	"time"
)

// naive cache which holds subscribers data for imsiFilterSuggestions
var subscribersCache []gp.Suggest

var imsiFilterSuggestions = func(word string, worker *SoracomWorker) []gp.Suggest {
	c := make(chan []gp.Suggest, 1024)
	if len(subscribersCache) == 0 {
		go getSubscribers(c, worker)
		select {
		case res := <-c:
			subscribersCache = res
		case <-time.After(10 * time.Second):
			return []gp.Suggest{{
				Text:        "Downloading IMSI in background",
				Description: "Hit space to see latest",
			}}
		}
	}
	return filterFunc(subscribersCache, word, filterTextOrDescriptionFuzzy)
}

var getSubscribers = func(c chan<- []gp.Suggest, worker *SoracomWorker) {
	var r []gp.Suggest

	result := worker.Execute("subscribers list --fetch-all")
	if err := json.NewDecoder(strings.NewReader(result)).Decode(&subscribers); err != nil {
		c <- []gp.Suggest{{
			Text:        "Error while running 'subscribers list --fetch-all'",
			Description: err.Error(),
		}}
	}
	for _, subscriber := range subscribers {
		online := "offline"
		if subscriber.SessionStatus.Online {
			online = "online"
		}
		r = append(r, gp.Suggest{
			Text: subscriber.Imsi,
			Description: fmt.Sprintf("%-12s | %-10s | %-7s | %-8s | %-11s | %s",
				trunc(subscriber.Subscription, 14),
				subscriber.Status,
				online,
				subscriber.ModuleType,
				trunc(subscriber.SpeedClass, 11),
				subscriber.Tags.Name,
			),
		})
	}
	c <- r
}
