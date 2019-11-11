package shell

import (
	"encoding/json"
	"fmt"
	"github.com/c-bata/go-prompt"
	"log"
	"os/exec"
	"strings"
	"time"
)

// naive cache which holds subscribers data for imsiFilterSuggestions
var cache []gp.Suggest

var imsiFilterSuggestions = func(word string) []gp.Suggest {
	c := make(chan []gp.Suggest, 1024)
	if len(cache) == 0 {
		go getSubscribers(c)
		select {
		case res := <-c:
			cache = res
		case <-time.After(10 * time.Second):
		}
	}
	return filterFunc(cache, word, filterTextOrDescriptionFuzzy)
}

// filter by text or description based on
// https://github.com/c-bata/go-prompt/blob/f350bee28f376e06a9877a516ac4eabe01804013/filter.go#L31 (MIT)
var filterTextOrDescriptionFuzzy = func(suggestions []prompt.Suggest, sub string, ignoreCase bool) []prompt.Suggest {
	if sub == "" {
		return suggestions
	}
	if ignoreCase {
		sub = strings.ToUpper(sub)
	}

	ret := make([]prompt.Suggest, 0, len(suggestions))
	for i := range suggestions {
		t := suggestions[i].Text
		d := suggestions[i].Description
		if ignoreCase {
			t = strings.ToUpper(t)
			d = strings.ToUpper(d)
		}
		if fuzzyMatch(t, sub) || fuzzyMatch(d, sub) {
			ret = append(ret, suggestions[i])
		}
	}
	return ret
}

func fuzzyMatch(s, sub string) bool {
	sChars := []rune(s)
	subChars := []rune(sub)
	sIdx := 0

	for _, c := range subChars {
		found := false
		for ; sIdx < len(sChars); sIdx++ {
			if sChars[sIdx] == c {
				found = true
				sIdx++
				break
			}
		}
		if !found {
			return false
		}
	}
	return true
}

var getSubscribers = func(c chan<- []prompt.Suggest) {
	cmd := exec.Command("/bin/sh", "-c", "soracom subscribers list")
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatal(err)
	}
	if err := cmd.Start(); err != nil {
		log.Fatal(err)
	}
	if err := json.NewDecoder(stdout).Decode(&subscribers); err != nil {
		log.Fatal(err)
	}
	if err := cmd.Wait(); err != nil {
		log.Fatal(err)
	}
	var r []prompt.Suggest
	for _, subscriber := range subscribers {
		r = append(r, prompt.Suggest{
			Text: subscriber.Imsi,
			Description: fmt.Sprintf("%s | %s | %t | %s | %s | %s",
				subscriber.Subscription,
				subscriber.Status,
				subscriber.SessionStatus.Online,
				subscriber.ModuleType,
				subscriber.Tags.Name,
				subscriber.SpeedClass,
			),
		})
	}
	c <- r
}
