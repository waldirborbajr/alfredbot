package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"localhost/alfred/entity"

	"github.com/sirupsen/logrus"
)

// State Machine Control
var (
	stateMachine string
)

func init() {
	stateMachine = "START"
}

func setWebhook() {
	// curl -F "url=https://8e3d-2804-d55-433d-5600-10b-8c9-7d95-bc9b.ngrok.io" https://api.telegram.org/bot5343272189:AAF5_yv9adxzqsNrYCqAY5jakgb4GqZFGBc/setWebhook
}

func deleteWhebhook() {
	// curl https://api.telegram.org/bot5343272189:AAF5_yv9adxzqsNrYCqAY5jakgb4GqZFGBc/deleteWebhook
}

// Create a new instance of the logger. You can have any number of instances.
var log = logrus.New()

func sendMessage(chatID int64, text string) error {
	// Create the request body struct
	reqBody := &entity.SendMessageReqBody{
		ChatID: chatID,
		Text:   text,
	}
	// Create the JSON body from the struct
	reqBytes, err := json.Marshal(reqBody)
	if err != nil {
		return err
	}
	// Send a post request with your token
	res, err := http.Post("https://api.telegram.org/bot5343272189:AAF5_yv9adxzqsNrYCqAY5jakgb4GqZFGBc/sendMessage", "application/json", bytes.NewBuffer(reqBytes))
	if err != nil {
		return err
	}
	if res.StatusCode != http.StatusOK {
		return errors.New("unexpected status" + res.Status)
	}

	return nil
}

// This handler is called everytime telegram sends us a webhook event
func Handler(_ http.ResponseWriter, req *http.Request) {
	// First, decode the JSON response body
	body := &entity.WebhookReqBody{}
	if err := json.NewDecoder(req.Body).Decode(body); err != nil {
		fmt.Println("could not decode request body", err)
		return
	}

	chatID := body.Message.Chat.ID
	text := strings.ToLower(body.Message.Text)

	log.Info("body handle: \n")
	log.Info(body)
	log.Info("TEXT: " + text)

	switch stateMachine {
	case "START":
		log.Info("state = START")
		stateMachine = "ABOUT"
	case "ABOUT":
		log.Info("state = ABOUT")
		stateMachine = "END"
	case "END":
		log.Info("state = END")
		stateMachine = "START"
	default:
		sendMessage(chatID, "Type /stat to begin.")
	}

	// Verify if command start with "/"
	if !strings.Contains(text[0:1], "/") {
		sendMessage(chatID, "Command must start with /")
	} else {
		if strings.Contains(text, "/about") {
			err := about(chatID)
			if err != nil {
				fmt.Println("error in sending reply:", err)
			}
		} else if strings.Contains(text, "/start") {
			err := start(chatID)
			if err != nil {
				fmt.Println("error in sending reply:", err)
			}
		} else if text != "" {
			wordsPlay(chatID, text)
		}
	}

	// log a confirmation message if the message is sent successfully
	fmt.Println("reply sent")
}

func start(chatID int64) error {
	err := sendMessage(chatID, "Welcome to the Alfred BOT")
	if err != nil {
		log.Fatal(err)
	}
	return nil
}

func about(chatID int64) error {
	err := sendMessage(chatID, "Alfre BOT ; v0.1.0")
	if err != nil {
		log.Fatal(err)
	}
	return nil
}

func getWordData(url string) (entity.Response, error) {
	var result entity.Response

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		fmt.Printf("client: could not create request: %s\n", err)
		os.Exit(1)
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Printf("client: error making http request: %s\n", err)
		os.Exit(1)
	}

	fmt.Printf("client: got response!\n")
	fmt.Printf("client: status code: %d\n", res.StatusCode)

	body, err := ioutil.ReadAll(res.Body) // response body is []byte
	if err != nil {
		return result, err
	}

	// read json data into a Result struct
	err = json.Unmarshal(body, &result)
	if err != nil {
		return result, err
	}
	return result, nil
}

func formatResponse(response entity.Response) string {
	var strResponse string
	var elemStr string

	for _, element := range response[0].Meanings {
		elemStr = element.PartOfSpeech + " : " + element.Definitions[0].Definition + "\n"
		strResponse += elemStr

	}
	log.Info("RESPONSE: ")
	log.Info(strResponse)
	return strResponse
}

func wordsPlay(chatID int64, stringInput string) error {
	requestURL := "https://api.dictionaryapi.dev/api/v2/entries/en/" + stringInput
	response, err := getWordData(requestURL)
	if err != nil {
		sendMessage(chatID, "Meaning not found for the word: "+stringInput)
		return err
	}

	// Return the 1st definition that comes for each part of speech

	err = sendMessage(chatID, formatResponse(response))
	return err
}

func logInit() error {
	file, err := os.OpenFile("log_file.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o666)
	if err == nil {
		log.Out = file
	} else {
		log.Info("Failed to log to file, using default stderr")
	}
	return nil
}

func main() {
	// init the log file
	err := logInit()
	if err != nil {
		os.Exit(1)
	}

	http.ListenAndServe(":9090", http.HandlerFunc(Handler))
}
