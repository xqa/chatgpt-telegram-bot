package main

import (
	bot "github.com/xqa/chatgpt-telegram-bot/bot"
)

func main() {
	bot.InitConfig()
	bot.InitChatMode()
	bot.Start()
}
