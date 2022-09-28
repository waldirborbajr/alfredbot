package entity

// Create a struct that mimics the webhook response body
// https://core.telegram.org/bots/api#update
type WebhookReqBody struct {
	Message struct {
		Text string `json:"text"`
		Chat struct {
			ID int64 `json:"id"`
		} `json:"chat"`
	} `json:"message"`
}

// Store data from the Api
type Response []struct {
	Word      string `json:"word"`
	Phonetic  string `json:"phonetic,omitempty"`
	Phonetics []struct {
		Text      string `json:"text"`
		Audio     string `json:"audio"`
		SourceURL string `json:"sourceUrl,omitempty"`
		License   struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"license,omitempty"`
	} `json:"phonetics"`
	Meanings []struct {
		PartOfSpeech string `json:"partOfSpeech"`
		Definitions  []struct {
			Definition string        `json:"definition"`
			Synonyms   []string      `json:"synonyms"`
			Antonyms   []interface{} `json:"antonyms"`
			Example    string        `json:"example,omitempty"`
		} `json:"definitions"`
		Synonyms []string      `json:"synonyms"`
		Antonyms []interface{} `json:"antonyms"`
	} `json:"meanings"`
	License struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"license"`
	SourceUrls []string `json:"sourceUrls"`
}

// The below code deals with the process of sending a response message
// to the user

type SendMessageReqBody struct {
	ChatID int64  `json:"chat_id"`
	Text   string `json:"text"`
}
