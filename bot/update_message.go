package bot

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"github.com/MakeNowJust/heredoc"
	"strconv"
	"strings"
)

type UpdateMessage struct {
	Update_Id int
	Message struct {
		Text string
		From struct {
			Id int
			Is_Bot bool
		}
	}
}

type CoronaInfo struct {
	Country string
	Cases struct {
		New string
		Active int
		Critical int
		Recovered int
		Total int
	}
	Deaths struct {
		New string
		Total int
	}
	Tests struct {
		Total int
	}
	Day string
	Time string
}

var client = &http.Client{}

func (update *UpdateMessage) Process(bot *Bot) {
	req, err := http.NewRequest("GET", "https://covid-193.p.rapidapi.com/history?country=" + strings.Replace(update.Message.Text, " ", "%20", -1), nil)
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Add("X-RapidAPI-Key", bot.config.XRapidAPIKey)
	req.Header.Add("x-rapidapi-host", "covid-193.p.rapidapi.com")

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		bot.PublishMessage("Что-то пошло не так, возможно данной страны не существует", update.Message.From.Id)
		return
	}
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	var responseMap map[string]json.RawMessage
	if err := json.Unmarshal(respBody, &responseMap); err != nil {
		log.Fatal(err)
	}

	resultsCnt, err := strconv.Atoi(string(responseMap["results"]))
	if err != nil {
		log.Fatal(err)
	}

	if resultsCnt == 0 {
		bot.PublishMessage("Нет такой страны", update.Message.From.Id)
		return
	}

	var arrayResponse []json.RawMessage
	if err := json.Unmarshal(responseMap["response"], &arrayResponse); err != nil {
		log.Fatal(err)
	}

	info := new(CoronaInfo)
	if err := json.Unmarshal(arrayResponse[0], &info); err != nil {
		log.Fatal(err)
	}

	message := heredoc.Docf(`
		%s
		Всего зараженных: %v
		Новых: %s
		Выздоровело: %v
		Смертей: %v
		_____
		На дату: %s
	`, info.Country, info.Cases.Total, info.Cases.New, info.Cases.Recovered, info.Deaths.Total, info.Day)
	bot.PublishMessage(message, update.Message.From.Id)
}
