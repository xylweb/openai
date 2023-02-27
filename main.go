package main

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	gogpt "github.com/sashabaranov/go-gpt3"
	"github.com/xylweb/xdb"
)

var (
	c      *gogpt.Client
	ctx    context.Context
	APIKEY = ""
	db     *xdb.Xdbase[string]
)

func init() {
	c = gogpt.NewClient(APIKEY)
	ctx = context.Background()
	db = xdb.NewXdb[string]()
	db.SetParams(xdb.Config{DbPath: "./data/", DbName: "base", IsIndex: true, Pass: "adjflwfkfa134asfa"}).Open()
	apikey, ok := db.Get("apikey")
	if ok {
		APIKEY = apikey.(string)
	}
}
func Get(text, model string, token int) string {
	req := gogpt.CompletionRequest{
		Model:            model,
		Prompt:           text,
		MaxTokens:        token,
		Temperature:      0.9,
		TopP:             0,
		FrequencyPenalty: 0,
		PresencePenalty:  0,
	}
	resp, err := c.CreateCompletion(ctx, req)
	if err != nil {
		return err.Error()
	}
	if len(resp.Choices) > 0 {
		return strings.TrimSpace(resp.Choices[0].Text)
	} else {
		return ""
	}
}
func main() {
	models := []string{
		gogpt.GPT3TextDavinci003,
		gogpt.GPT3TextDavinci002,
		gogpt.GPT3TextCurie001,
		gogpt.GPT3TextBabbage001,
		gogpt.GPT3TextAda001,
		gogpt.GPT3TextDavinci001,
		gogpt.GPT3DavinciInstructBeta,
		gogpt.GPT3Davinci,
		gogpt.GPT3CurieInstructBeta,
		gogpt.GPT3Curie,
		gogpt.GPT3Ada,
		gogpt.GPT3Babbage,
	}
	model := gogpt.GPT3TextCurie001
	token := uint64(512)
	for {
		fmt.Print("I: ")
		var text string
		fmt.Scan(&text)
		if text != "" {
			ms, ok := db.Get("model")
			if ok {
				model = ms.(string)
			}
			ts, ok := db.Get("token")
			if ok {
				token = ts.(uint64)
			}
			if text == "exit" {
				os.Exit(0)
				continue
			}
			if text == "models" {
				fmt.Printf("He: now:%s %+v\n", model, models)
				continue
			}
			if text == "cmd" {
				fmt.Printf("He: cmd:%s\n", "gpt0... tokens512... apikey...")
				continue
			}
			if strings.HasPrefix(text, "gpt") {
				if text == "gpt" {
					fmt.Printf("He: now:%s\n", model)
					continue
				}
				numstr := strings.ReplaceAll(text, "gpt", "")
				num, _ := strconv.Atoi(numstr)
				if num < 12 {
					model = models[num]
					db.Add("model", model)
				}
				continue
			}
			if strings.HasPrefix(text, "tokens") {
				if text == "tokens" {
					fmt.Printf("He: now:%d\n", token)
					continue
				}
				tokenstr := strings.ReplaceAll(text, "tokens", "")
				tokeni, _ := strconv.Atoi(tokenstr)
				if tokeni > 0 {
					token = uint64(tokeni)
					db.Add("token", token)
				}
				continue
			}
			if strings.HasPrefix(text, "apikey") {
				if text == "apikey" {
					fmt.Printf("He: now:%s\n", APIKEY)
					continue
				}
				apikey := strings.ReplaceAll(text, "apikey", "")
				if apikey != "" {
					APIKEY = apikey
					db.Add("apikey", APIKEY)
				}
				continue
			}
			fmt.Printf("\nHe:\n%s\n\n", Get(text, model, int(token)))
		}
		text = ""
		time.Sleep(time.Second)
	}
}
