package main

import "github.com/0legovich/corona/bot"

func main() {
   config := bot.NewConfig()
   bot := bot.New(config)
   bot.Start()
}
