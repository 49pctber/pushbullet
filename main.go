package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"os"
)

const apiURL string = "https://api.pushbullet.com/v2/pushes"

var accessToken *string

type PushbulletMessage struct {
	Type  string `json:"type"`
	Title string `json:"title"`
	Body  string `json:"body"`
}

func (pbm PushbulletMessage) String() string {
	return fmt.Sprintf("%s (%s): %s", pbm.Title, pbm.Type, pbm.Body)
}

func sendPushbulletNotification(msg PushbulletMessage) error {

	jsonData, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %v", err)
	}

	req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")

	if accessToken == nil {
		return fmt.Errorf("accessToken was not initialized")
	} else if len(*accessToken) == 0 {
		return fmt.Errorf("accessToken is empty")
	}

	req.Header.Set("Access-Token", *accessToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		return nil
	}

	switch resp.StatusCode {
	case http.StatusTooManyRequests:
		fmt.Println("Have you already used your push quota for the month?")
	case http.StatusUnauthorized:
		fmt.Println("Is your API key valid?")
	}

	return fmt.Errorf("failed to send notification, status code: %d", resp.StatusCode)
}

func main() {

	// get user input
	title := flag.String("title", "", "the title of the notification")
	body := flag.String("body", "", "the body of the notification")
	accessToken = flag.String("token", "", "your Pushbullet access token. If blank, the environment variable PUSHBULLET_TOKEN will be used")
	flag.Parse()

	// check if token was provided on command line
	if len(*accessToken) == 0 {
		*accessToken = os.Getenv("PUSHBULLET_TOKEN")
		// check if token is available as environment variable
		if len(*accessToken) == 0 {
			fmt.Println("environment variable PUSHBULLET_TOKEN is not defined")
			os.Exit(1)
		}
	}

	// check for title
	if len(*title) == 0 {
		fmt.Println("define a title using the title flag")
		os.Exit(1)
	}

	// check for body of notification
	if len(*body) == 0 {
		fmt.Println("define the notification body using the body flag")
		os.Exit(1)
	}

	// send notification
	msg := PushbulletMessage{
		Type:  "note",
		Title: *title,
		Body:  *body,
	}
	err := sendPushbulletNotification(msg)
	if err != nil {
		fmt.Printf("(error) %v\n", err)
		os.Exit(1)
	} else {
		fmt.Printf("(success) %v\n", msg)
	}
}
