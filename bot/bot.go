package bot

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	_"github.com/davecgh/go-spew/spew"
)

type Bot struct {
	config *Config
}

func New(config *Config) *Bot {
	return &Bot{
		config: config,
	}
}

func (b *Bot) Start() {
	if b.config.Token == "" {
		log.Fatal(errors.New("Bot does not have token"))
	}
	staticUrl := "https://api.telegram.org/bot" + b.config.Token + "/getUpdates?timeout=" + strconv.Itoa(b.config.Timeout) + "&offset="
	offset := 0

	for {
		resp, err := http.Get(staticUrl + strconv.Itoa(offset + 1))
		if err != nil {
			log.Fatal(err)
		}
		defer resp.Body.Close()

		updatesJson, err := parseBody(&resp.Body)
		if err != nil {
			log.Fatal(err)
		}

		for index, updateJson := range updatesJson {
			update := new(UpdateMessage)
			err = json.Unmarshal([]byte(updateJson), update)
			if err != nil {
				log.Fatal(err.Error())
			}

			update.Process(b)
			if index == len(updatesJson) -1 {
				offset = update.Update_Id
			}
		}
	}
}

func (b *Bot) PublishMessage(message string, userId int) {
	url := "https://api.telegram.org/bot" + b.config.Token + "/sendMessage"
	bodyString := []byte(`{"text": "` + message + `", "chat_id": ` + strconv.Itoa(userId) + `"}`)
	body := bytes.NewReader(bodyString)

	resp, err := http.Post(url, "application/json", body)
	if err != nil {
		log.Fatal(err.Error())
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		respBody, _ := ioutil.ReadAll(resp.Body)
		fmt.Printf("%s", respBody)
	}
}

func parseBody(httpBody *io.ReadCloser) ([]json.RawMessage, error) {
	decoder := json.NewDecoder(*httpBody)

	var responseMap map[string]json.RawMessage
	if err := decoder.Decode(&responseMap); err == io.EOF {
		return nil, err
	} else if err != nil {
		return nil, err
	}

	var updates []json.RawMessage
	if err := json.Unmarshal(responseMap["result"], &updates); err != nil {
		return nil, err
	}

	return updates, nil
}
