package bot

import (
	"log"
	"strings"
	"time"

	tele "gopkg.in/telebot.v3"
	"gopkg.in/telebot.v3/middleware"
)

var (
	HELP_MESSAGE = `指令列表:
	/retry - 为先前的查询重新生成响应
	/new - 开启新一轮会话
	/mode - 选择会话模式
	/balance - 显示余额信息
	/help - 显示帮助信息
	`
)

func Start() {
	pref := tele.Settings{
		Token:  viperConfig.GetString("telegram_token"),
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
	}

	tgbot, err := tele.NewBot(pref)
	if err != nil {
		log.Fatal(err)
		return
	}
	tgbot.SetCommands([]tele.Command{
		{
			Text:        "retry",
			Description: "为先前的查询重新生成响应",
		},
		{
			Text:        "new",
			Description: "开启新一轮聊天",
		},
		{
			Text:        "mode",
			Description: "选择聊天模式",
		},
		{
			Text:        "balance",
			Description: "显示余额信息",
		},
		{
			Text:        "help",
			Description: "显示帮助信息",
		},
	})
	log.Printf("Authorized on account %s", tgbot.Me.Username)

	var allowedIds []int64
	for _, i := range viperConfig.GetIntSlice("allowed_ids") {
		allowedIds = append(allowedIds, int64(i))
	}
	tgbot.Use(middleware.Whitelist(allowedIds...))

	chatGPT := NewChatGPT()

	tgbot.Handle("/start", func(ctx tele.Context) error {
		reply_text := "你好！我是使用 GPT-3.5 OpenAI API 实现的 <b>ChatGPT</b> 机器人 🤖\n\n"
		reply_text += HELP_MESSAGE
		reply_text += "\n现在……问我任何事！"
		return ctx.Send(reply_text, tele.ModeHTML)
	})
	tgbot.Handle("/help", func(ctx tele.Context) error {
		return ctx.Send(HELP_MESSAGE)
	})
	tgbot.Handle("/new", chatGPT.OnNew)
	tgbot.Handle("/mode", chatGPT.OnMode)
	tgbot.Handle("/retry", chatGPT.OnRetry)
	tgbot.Handle("/balance", chatGPT.OnBalance)
	tgbot.Handle(tele.OnText, chatGPT.OnMessage)
	tgbot.Handle(tele.OnCallback, func(ctx tele.Context) error {
		if ctx.Get("type") == "mode" {
			return chatGPT.OnModeCallback(ctx)
		}
		// ...其他OnCallback
		return ctx.Send("你点击了一个按钮")
	}, OnCallback)
	tgbot.Start()
}

func OnCallback(next tele.HandlerFunc) tele.HandlerFunc {
	return func(c tele.Context) error {
		defer c.Respond()
		data := strings.TrimSpace(c.Callback().Data)
		arr := strings.Split(data, "|")
		c.Set("type", arr[0])
		c.Set("arg", arr[1])
		return next(c)
	}
}
