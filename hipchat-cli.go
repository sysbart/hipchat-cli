package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"os"
)

// Based on tbruyelle/hipchat-go

// notificationRequest represents a HipChat room notification request.
type notificationRequest struct {
	Color         string `json:"color,omitempty"`
	Message       string `json:"message,omitempty"`
	Notify        bool   `json:"notify,omitempty"`
	MessageFormat string `json:"message_format,omitempty"`
}

// sendMessage represents an Hipchat message sent
func sendMessage(r string, n interface{}, t string) (*http.Response, error) {

	url := fmt.Sprintf("https://api.hipchat.com/v2/room/%s/notification?auth_token=%s", r, t)

	buf := new(bytes.Buffer)
	if n != nil {
		err := json.NewEncoder(buf).Encode(n)
		if err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequest("POST", url, buf)
	req.Header.Set("Content-Type", "application/json")

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	if c := resp.StatusCode; c < 200 || c > 299 {
		return resp, fmt.Errorf("Server returns status %d", c)
	}

	return resp, err
}

func main() {
	token := flag.String("token", "", "Hipchat Token")
	roomID := flag.String("roomID", "", "RoomID of the Hipchat room")
	msg := flag.String("msg", "", "Message you want to send")
	color := flag.String("color", "yellow", "Color of the message")

	flag.Parse()

	if *token == "" || *roomID == "" || *msg == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}

	notificationRequest := notificationRequest{Color: *color, Message: *msg, Notify: true}
	fmt.Printf("%+v\n", notificationRequest)

	resp, err := sendMessage(*roomID, &notificationRequest, *token)

	if err != nil {
		fmt.Fprintf(os.Stderr, "The message was not sent properly : %q\n", err)
		fmt.Fprintf(os.Stderr, "Server returns %+v\n", resp)
		os.Exit(1)
	}

}
