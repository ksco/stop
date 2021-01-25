package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
)

type Text struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

type Block struct {
	Type string `json:"type"`
	Text *Text  `json:"text,omitempty"`
}

type Blocks struct {
	Blocks []Block `json:"blocks"`
}

func Push(f *Frame) error {
	token := os.Getenv("SLACK_TOKEN")
	u := url.URL{
		Scheme: "https",
		Host:   "hooks.slack.com",
		Path:   fmt.Sprintf("services/%s", token),
	}
	payload := Blocks{
		Blocks: []Block{{
			Type: "header",
			Text: &Text{
				Type: "plain_text",
				Text: f.Name,
			},
		}, {
			Type: "divider",
		}, {
			Type: "section",
			Text: &Text{
				Type: "mrkdwn",
				Text: "*Memory Usage*: " + progressBar(f.MemoryUsage()),
			},
		}},
	}
	data, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	println(string(data))
	req, err := http.NewRequest("POST", u.String(), bytes.NewBuffer(data))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer func() { _ = resp.Body.Close() }()

	body, _ := ioutil.ReadAll(resp.Body)
	println(string(body))
	return nil
}

func progressBar(percent float64) string {
	return fmt.Sprintf("[%-50s]%.2f%%", strings.Repeat("█", int(percent/2)), percent)
}
