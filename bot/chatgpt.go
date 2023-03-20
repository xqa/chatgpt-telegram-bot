package bot

import (
	"context"
	"fmt"
	"log"
	"sort"
	"strings"
	"time"

	"github.com/avast/retry-go"
	openai "github.com/sashabaranov/go-openai"
	tele "gopkg.in/telebot.v3"
)

type ChatGPT struct {
	client *openai.Client
}

func NewChatGPT() *ChatGPT {
	return &ChatGPT{
		client: openai.NewClient(viperConfig.GetString("openai_api_key")),
	}
}

func (c *ChatGPT) Completion(req openai.ChatCompletionRequest) (openai.ChatCompletionMessage, error) {
	var message openai.ChatCompletionMessage
	err := retry.Do(
		func() error {
			resp, err := c.client.CreateChatCompletion(context.Background(), req)
			if err != nil {
				return err
			}
			message = resp.Choices[0].Message
			return nil
		},
		retry.Delay(time.Second),
		retry.Attempts(3),
		retry.DelayType(retry.FixedDelay),
		retry.OnRetry(func(n uint, err error) {
			log.Printf("retry %d times: %v", n, err)
		}),
	)
	if err != nil {
		return message, err
	}
	return message, nil
}

func (c *ChatGPT) reply(ctx tele.Context, user *User) error {
	if !strings.HasPrefix(ctx.Message().Text, "/retry") {
		user.AddUserMessage(ctx.Message().Text)
	}
	request := openai.ChatCompletionRequest{
		Model:    openai.GPT3Dot5Turbo0301,
		Messages: user.Messages,
	}
	completedMessage, err := c.Completion(request)
	if err != nil {
		return err
	}
	user.AddAssistantMessage(completedMessage.Content)
	return ctx.Send(completedMessage.Content, tele.ModeMarkdown)
}

func (c *ChatGPT) OnNew(ctx tele.Context) error {
	user := GetUser(ctx.Message().Sender.ID)
	user.ResetMessage()
	return ctx.Send(modeConfig.GetString(user.ChatMode+".welcome_message"), tele.ModeHTML)
}

func (c *ChatGPT) OnMode(ctx tele.Context) error {
	menu := &tele.ReplyMarkup{ResizeKeyboard: true}
	var rows []tele.Row
	keys := make([]string, 0, len(modeConfig.AllSettings()))
	for k := range modeConfig.AllSettings() {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		rows = append(rows, menu.Row(menu.Data(modeConfig.GetString(k+".name"), "mode", k)))
	}
	menu.Inline(rows...)
	return ctx.Send("选择聊天模式：", menu)
}
func (c *ChatGPT) OnModeCallback(ctx tele.Context) error {
	user := GetUser(ctx.Callback().Sender.ID)
	user.ResetMessage()
	user.SetMode(ctx.Get("arg").(string))
	return ctx.Send(modeConfig.GetString(user.ChatMode+".welcome_message"), tele.ModeHTML)
}

func (c *ChatGPT) OnRetry(ctx tele.Context) error {
	user := GetUser(ctx.Message().Sender.ID)
	if err := user.RetryAnswer(); err != nil {
		return ctx.Send(err.Error())
	}
	return c.reply(ctx, user)
}

func (c *ChatGPT) OnMessage(ctx tele.Context) error {
	user := GetUser(ctx.Message().Sender.ID)
	ctx.Notify(tele.Typing)
	return c.reply(ctx, user)
}

func (c *ChatGPT) OnBalance(ctx tele.Context) error {
	data, err := GetBalance()
	if err != nil {
		return ctx.Send(err.Error())
	}
	t1 := time.Unix(int64(data.Grants.Data[0].EffectiveAt), 0)
	t2 := time.Unix(int64(data.Grants.Data[0].ExpiresAt), 0)
	msg := fmt.Sprintf("💵 已用: 💲%v\n💵 剩余: 💲%v\n⏳ 有效时间: 从 %v 到 %v\n", fmt.Sprintf("%.2f", data.TotalUsed), fmt.Sprintf("%.2f", data.TotalAvailable), t1.Format("2006-01-02 15:04:05"), t2.Format("2006-01-02 15:04:05"))
	return ctx.Send(msg)
}
