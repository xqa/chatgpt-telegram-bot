package bot

import (
	"errors"

	openai "github.com/sashabaranov/go-openai"
)

type User struct {
	ChatMode      string
	SystemMessage openai.ChatCompletionMessage
	Messages      []openai.ChatCompletionMessage
}

var users = make(map[int64]*User)

func GetUser(userID int64) *User {
	if _, ok := users[userID]; !ok {
		mode := viperConfig.GetString("chat_mode")
		users[userID] = &User{
			ChatMode: mode,
			SystemMessage: openai.ChatCompletionMessage{
				Role:    openai.ChatMessageRoleSystem,
				Content: modeConfig.GetString(mode + ".prompt_start"),
			},
			Messages: []openai.ChatCompletionMessage{},
		}
	}
	return users[userID]
}

func (u *User) SetMode(mode string) {
	u.ChatMode = mode
	u.SystemMessage.Content = modeConfig.GetString(mode + ".prompt_start")
}
func (u *User) RetryAnswer() error {
	if len(u.Messages) == 0 {
		return errors.New("no message to retry")
	}
	u.Messages = u.Messages[:len(u.Messages)-1]
	return nil
}
func (u *User) ResetMessage() {
	u.Messages = []openai.ChatCompletionMessage{}
}
func (u *User) AddUserMessage(msg string) {
	u.Messages = append(u.Messages, openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleUser,
		Content: msg,
	})
}
func (u *User) AddAssistantMessage(msg string) {
	u.Messages = append(u.Messages, openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleAssistant,
		Content: msg,
	})
}
