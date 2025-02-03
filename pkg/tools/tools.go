package tools

import (
	"errors"
	"fmt"
	"net/http"
	"time"
)

func CheckURL(url string) bool {
	resp, err := http.Get(url)
	if err != nil {
		return false
	}
	defer resp.Body.Close()
	return resp.StatusCode >= 200 && resp.StatusCode < 400
}

func WaitForURL(url string, attempts int, delay time.Duration) error {
	for i := 0; i < attempts; i++ {
		if CheckURL(url) {
			fmt.Println(url, "is responding")
			return nil
		}
		fmt.Println("Waiting for", url)
		time.Sleep(delay)
	}
	return errors.New("URL did not respond within the given attempts")
}