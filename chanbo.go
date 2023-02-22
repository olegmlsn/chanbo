package chanbo

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

type Chanbo struct {
	BotApiKey   string
	ChannelName string
}

func New() *Chanbo {
	return &Chanbo{
		BotApiKey:   os.Getenv("BOT_API_KEY"),
		ChannelName: os.Getenv("CHANNEL_NAME"),
	}
}

type Response struct {
	Ok     bool `json:"ok"`
	Result struct {
		MessageID  int `json:"message_id"`
		SenderChat struct {
			ID       int64  `json:"id"`
			Title    string `json:"title"`
			Username string `json:"username"`
			Type     string `json:"type"`
		} `json:"sender_chat"`
		Chat struct {
			ID       int64  `json:"id"`
			Title    string `json:"title"`
			Username string `json:"username"`
			Type     string `json:"type"`
		} `json:"chat"`
		Date int    `json:"date"`
		Text string `json:"text"`
	} `json:"result"`
}

func (c Chanbo) SendMessage(message string, parseMode string) error {
	url := fmt.Sprintf("https://api.telegram.org/bot%s/%s?chat_id=%s&text=%s",
		c.BotApiKey, "sendMessage", c.ChannelName, message)

	switch parseMode {
	case "html":
		url += "&parse_mode=HTML"
	case "markdown":
		url += "&parse_mode=MarkdownV2"
	}

	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("chanbo: sendMesage: http.Get error: %w", err)
	}
	if resp.Status != "200 OK" {
		return fmt.Errorf("chanbo: sendMessage: return code not 200: %w", err)
	}
	defer resp.Body.Close()

	var result Response
	body, err := ioutil.ReadAll(resp.Body)
	err = json.Unmarshal(body, &result)
	if err != nil {
		return fmt.Errorf("chanbo: sendMessage: json unmarshal error: %w", err)
	}

	return nil
}

func (c Chanbo) SendPhoto(photo []byte, path string) error {
	url := fmt.Sprintf("https://api.telegram.org/bot%s/%s?chat_id=%s",
		c.BotApiKey, "sendPhoto", c.ChannelName)

	client := &http.Client{
		Timeout: time.Second * 10,
	}

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	fw, err := writer.CreateFormFile("photo", filepath.Base(path))
	if err != nil {
		return fmt.Errorf("chanbo: CreateFormFile: error: %w", err)
	}

	if photo != nil {
		_, err = fw.Write(photo)
		if err != nil {
			return fmt.Errorf("chanbo: fw.Write: error: %w", err)
		}
	} else {
		file, err := os.Open(path)
		if err != nil {
			return fmt.Errorf("chanbo: SendPhoto: http.Get error: %w", err)
		}
		_, err = io.Copy(fw, file)
		if err != nil {
			return fmt.Errorf("chanbo: io.Copy: error: %w", err)
		}
	}

	writer.Close()
	req, err := http.NewRequest("POST", url, bytes.NewReader(body.Bytes()))
	if err != nil {
		return fmt.Errorf("chanbo: NewRequest: error: %w", err)
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("chanbo: SendPhoto: client.Do error: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("chanbo: NewRequest: error: %w", err)
	}

	defer resp.Body.Close()

	var result Response
	respBody, err := ioutil.ReadAll(resp.Body)
	err = json.Unmarshal(respBody, &result)
	if err != nil {
		return fmt.Errorf("chanbo: SendPhoto: json unmarshal error: %w", err)
	}

	return nil
}
