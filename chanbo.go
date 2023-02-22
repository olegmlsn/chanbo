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

func (c Chanbo) SendMessage(message string) error {
	urlTemplate := "https://api.telegram.org/bot%s/%s?chat_id=%s&%s=%s"
	url := fmt.Sprintf(urlTemplate, c.BotApiKey, "sendMessage", c.ChannelName, "text", message)

	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("chabot: SendMesage: http.Get error: %w", err)
	}
	if resp.Status != "200 OK" {
		return fmt.Errorf("Chabot: SendMessage: return code not 200: %w", err)
	}
	defer resp.Body.Close()

	var result Response
	body, err := ioutil.ReadAll(resp.Body)
	err = json.Unmarshal(body, &result)
	if err != nil {
		return fmt.Errorf("Chabot: SendMessage: json unmarshal error: %w", err)
	}

	fmt.Println(result)

	return nil
}

func (c Chanbo) SendPhoto(path string) error {
	urlTempl := "https://api.telegram.org/bot%s/%s?chat_id=%s"
	url := fmt.Sprintf(urlTempl, c.BotApiKey, "sendPhoto", c.ChannelName)

	file, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("chabot: SendPhoto: http.Get error: %w", err)
	}

	client := &http.Client{
		Timeout: time.Second * 10,
	}

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	fw, err := writer.CreateFormFile("photo", filepath.Base(path))
	if err != nil {
		return fmt.Errorf("chabot: CreateFormFile: error: %w", err)
	}

	_, err = io.Copy(fw, file)
	if err != nil {
		return fmt.Errorf("chabot: io.Copy: error: %w", err)
	}

	writer.Close()
	req, err := http.NewRequest("POST", url, bytes.NewReader(body.Bytes()))
	if err != nil {
		return fmt.Errorf("chabot: NewRequest: error: %w", err)
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("chabot: SendPhoto: client.Do error: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("chabot: NewRequest: error: %w", err)
	}

	defer resp.Body.Close()

	var result Response
	respBody, err := ioutil.ReadAll(resp.Body)
	fmt.Println(string(respBody))
	err = json.Unmarshal(respBody, &result)
	if err != nil {
		return fmt.Errorf("Chabot: SendPhoto: json unmarshal error: %w", err)
	}

	fmt.Println(result)

	return nil
}
