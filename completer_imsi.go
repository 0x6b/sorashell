package shell

import (
	"encoding/json"
	"fmt"
	gp "github.com/c-bata/go-prompt"
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
var filterTextOrDescriptionFuzzy = func(suggestions []gp.Suggest, sub string, ignoreCase bool) []gp.Suggest {
	if sub == "" {
		return suggestions
	}
	if ignoreCase {
		sub = strings.ToUpper(sub)
	}

	ret := make([]gp.Suggest, 0, len(suggestions))
	for i := range suggestions {
		c := suggestions[i].Text + " " + suggestions[i].Description
		if ignoreCase {
			c = strings.ToUpper(c)
		}
		if fuzzyMatch(c, sub) {
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

var getSubscribers = func(c chan<- []gp.Suggest) {
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
	var r []gp.Suggest
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

func trunc(s string, n int) string {
	r := s
	if len(s) > n {
		if n > 3 {
			n -= 3
		}
		r = s[0:n] + "..."
	}
	return r
}
