package main

import (
	"bufio"
	"github.com/bluele/slack"
	"github.com/utahta/go-linenotify"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
)

type Config struct {
	SlackInfo []SlackInfo
	LineInfo  []LineInfo
}

type SlackInfo struct {
	SlackToken   string
	SlackName    string
	SlackIconUrl string
	SlackChannel string
}

type LineInfo struct {
	LineToken string
}

func slackNotify(message string) {
	for _, row := range config.SlackInfo {
		hook := slack.NewWebHook(row.SlackToken)
		err := hook.PostMessage(&slack.WebHookPostPayload{
			Attachments: []*slack.Attachment{
				{Text: message, Color: "danger"},
			},
			Channel:  row.SlackChannel,
			Username: row.SlackName,
			IconUrl:  row.SlackIconUrl,
		})
		check(err)
	}
}

func lineNotify(message string) {
	for _, row := range config.LineInfo {
		c := linenotify.New()
		c.Notify(row.LineToken, message, "", "", nil)
	}
}

func notify(message string) {
	slackNotify(message)
	lineNotify(message)
}

func check(err error) {
	if err != nil {
		log.Fatalf("Fatal: %v", err)
	}
}

var config Config

func main() {
	buf, err := ioutil.ReadFile("config.yml")
	check(err)
	err = yaml.Unmarshal(buf, &config)
	check(err)

	stdin := bufio.NewScanner(os.Stdin)
	for stdin.Scan() {
		text := stdin.Text()
		switch text {
		case "open":
			notify("課金通知botの監視を始めました")
		case "close":
			notify("課金通知botが止まりました、管理者をしばいてください")
		}
	}
}
