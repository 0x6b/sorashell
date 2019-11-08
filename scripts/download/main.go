package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"time"
)

func main() {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}

	target := filepath.Join(dir, "assets")
	fmt.Printf("downloading assets to %s\n", target)

	ch := make(chan string)
	urls := []string{
		"https://raw.githubusercontent.com/soracom/soracom-cli/master/generators/assets/soracom-api.en.yaml",
		"https://raw.githubusercontent.com/soracom/soracom-cli/master/generators/assets/soracom-api.ja.yaml",
		"https://raw.githubusercontent.com/soracom/soracom-cli/master/generators/assets/cli/en.yaml",
		"https://raw.githubusercontent.com/soracom/soracom-cli/master/generators/assets/cli/ja.yaml",
	}

	for _, url := range urls {
		go fetch(url, target, ch)
	}

	for range urls {
		fmt.Println(<-ch)
	}
}

func fetch(url, target string, ch chan<- string) {
	start := time.Now()
	client := &http.Client{}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		ch <- fmt.Sprint(err)
		return
	}

	res, err := client.Do(req)
	defer res.Body.Close()
	if err != nil {
		ch <- fmt.Sprint(err)
		return
	}

	filename := path.Join(target, path.Base(req.URL.Path))

	out, err := os.Create(filename)
	defer out.Close()
	if err != nil {
		ch <- fmt.Sprint(err)
		return
	}

	n, err := io.Copy(out, res.Body)
	if err != nil {
		ch <- fmt.Sprint(err)
		return
	}

	ch <- fmt.Sprintf("%.2fs %5dKB %s", time.Since(start).Seconds(), n/1024, filename)
}
