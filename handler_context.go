package tgbotapp

import (
	"errors"
	"log/slog"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/nexoratech2025/go-telegram-bot-app/session"
)

type HandlerContext struct {
	*BotContext
	Logger *slog.Logger
}

type HandlerParams struct {
	ChatID         int64
	TelegramUserID int64
	Session        session.Sessioner
}

func NewHandlerContext(ctx *BotContext, handlerName string) *HandlerContext {
	logger := ctx.Logger().With(
		slog.String("type", "handler"),
		slog.String("name", handlerName),
	)

	return &HandlerContext{
		BotContext: ctx,
		Logger:     logger,
	}
}

// Parse mode constants for Telegram messages
const (
	ParseModeHTML       = "HTML"
	ParseModeMarkdown   = "Markdown"
	ParseModeMarkdownV2 = "MarkdownV2"
)

func (h *HandlerContext) GetParams() HandlerParams {
	return HandlerParams{
		ChatID:         h.Update.FromChat().ChatConfig().ChatID,
		TelegramUserID: h.Update.SentFrom().ID,
		Session:        h.Session,
	}
}

func (h *HandlerContext) SendMessage(text string, parseMode ...string) error {
	msg := tgbotapi.NewMessage(h.Update.FromChat().ChatConfig().ChatID, text)

	if len(parseMode) > 0 && parseMode[0] != "" {
		mode := parseMode[0]
		validModes := []string{ParseModeHTML, ParseModeMarkdown, ParseModeMarkdownV2}
		isValid := false
		for _, validMode := range validModes {
			if mode == validMode {
				isValid = true
				break
			}
		}

		if !isValid {
			h.Logger.ErrorContext(h.Ctx, "Invalid parse mode", "parseMode", mode)
			h.SendError("Invalid parse mode provided")
			return errors.New("invalid parse mode provided")
		}

		msg.ParseMode = mode
	}

	_, err := h.BotAPI.Send(msg)
	if err != nil {
		h.HandleSendMessageError(err)
		return err
	}

	return nil
}

func (h *HandlerContext) SendMessageWithKeyboard(text string, keyboard interface{}) error {
	msg := tgbotapi.NewMessage(h.Update.FromChat().ChatConfig().ChatID, text)
	msg.ReplyMarkup = keyboard
	_, err := h.BotAPI.Send(msg)
	if err != nil {
		h.HandleSendMessageError(err)
		return err
	}

	return nil
}

func (h *HandlerContext) SendMessageWithInlineKeyboard(text string, keyboard tgbotapi.InlineKeyboardMarkup) error {
	msg := tgbotapi.NewMessage(h.Update.FromChat().ChatConfig().ChatID, text)
	msg.ReplyMarkup = keyboard
	_, err := h.BotAPI.Send(msg)
	if err != nil {
		h.HandleSendMessageError(err)
		return err
	}

	return nil
}

func (h *HandlerContext) SendError(message string) {
	h.Logger.ErrorContext(h.Ctx, message)
	h.SendMessage("Something went wrong. " + message)
}

func (h *HandlerContext) LogError(message string, err error) {
	h.Logger.ErrorContext(h.Ctx, message, "error", err)
}

func (h *HandlerContext) SetState(state string) {
	h.Session.SetState(state)
}

func (h *HandlerContext) GetSessionData(key string) (interface{}, bool) {
	return h.Session.Get(key)
}

func (h *HandlerContext) SetSessionData(key string, value interface{}) {
	h.Session.Set(key, value)
}

func (h *HandlerContext) ClearSessionData() {
	h.Session.ClearData()
}

func (h *HandlerContext) HasDocument() bool {
	return hasDocument(h.Update.Message)
}

func (h *HandlerContext) GetDocumentType() string {
	return getDocumentType(h.Update.Message)
}

func (h *HandlerContext) GetDocument() *tgbotapi.Document {
	return h.Update.Message.Document
}

func (h *HandlerContext) GetPhoto() []tgbotapi.PhotoSize {
	return h.Update.Message.Photo
}

func (h *HandlerContext) GetVideo() *tgbotapi.Video {
	return h.Update.Message.Video
}

func (h *HandlerContext) GetAudio() *tgbotapi.Audio {
	return h.Update.Message.Audio
}

func (h *HandlerContext) GetVoice() *tgbotapi.Voice {
	return h.Update.Message.Voice
}

func (h *HandlerContext) GetVideoNote() *tgbotapi.VideoNote {
	return h.Update.Message.VideoNote
}

func (h *HandlerContext) GetSticker() *tgbotapi.Sticker {
	return h.Update.Message.Sticker
}

func (h *HandlerContext) GetBestPhoto() *tgbotapi.PhotoSize {
	photos := h.Update.Message.Photo
	if len(photos) > 0 {
		return &photos[len(photos)-1]
	}
	return nil
}

func (h *HandlerContext) GetCurrentState() string {
	return h.Session.CurrentState()
}

func (h *HandlerContext) DeleteSessionData(key string) {
	h.Session.Delete(key)
}

func (h *HandlerContext) GetAllSessionKeys() []string {
	return h.Session.GetAllKeys()
}

func (h *HandlerContext) GetText() string {
	if h.Update.Message != nil {
		return h.Update.Message.Text
	}
	return ""
}

func (h *HandlerContext) GetCommand() string {
	if h.Update.Message != nil && h.Update.Message.IsCommand() {
		return h.Update.Message.Command()
	}
	return ""
}

func (h *HandlerContext) GetCommandArguments() string {
	if h.Update.Message != nil && h.Update.Message.IsCommand() {
		return h.Update.Message.CommandArguments()
	}
	return ""
}

func (h *HandlerContext) GetCallbackData() string {
	if h.Update.CallbackQuery != nil {
		return h.Update.CallbackQuery.Data
	}
	return ""
}

func (h *HandlerContext) GetCallbackQuery() *tgbotapi.CallbackQuery {
	return h.Update.CallbackQuery
}

func (h *HandlerContext) AnswerCallbackQuery(text string) {
	if h.Update.CallbackQuery != nil {
		callback := tgbotapi.NewCallback(h.Update.CallbackQuery.ID, text)
		h.BotAPI.Send(callback)
	}
}

func (h *HandlerContext) AnswerCallbackQueryWithAlert(text string) {
	if h.Update.CallbackQuery != nil {
		callback := tgbotapi.NewCallbackWithAlert(h.Update.CallbackQuery.ID, text)
		h.BotAPI.Send(callback)
	}
}

func (h *HandlerContext) SendPhoto(photo tgbotapi.FileBytes, caption string) error {
	msg := tgbotapi.NewPhoto(h.Update.FromChat().ChatConfig().ChatID, photo)
	if caption != "" {
		msg.Caption = caption
	}
	_, err := h.BotAPI.Send(msg)
	if err != nil {
		h.HandleSendMessageError(err)
		return err
	}

	return nil
}

func (h *HandlerContext) SendDocument(document tgbotapi.FileBytes, caption string) error {
	msg := tgbotapi.NewDocument(h.Update.FromChat().ChatConfig().ChatID, document)
	if caption != "" {
		msg.Caption = caption
	}
	_, err := h.BotAPI.Send(msg)
	if err != nil {
		h.HandleSendMessageError(err)
		return err
	}

	return nil
}

func (h *HandlerContext) SendVideo(video tgbotapi.FileBytes, caption string) error {
	msg := tgbotapi.NewVideo(h.Update.FromChat().ChatConfig().ChatID, video)
	if caption != "" {
		msg.Caption = caption
	}
	_, err := h.BotAPI.Send(msg)
	if err != nil {
		h.HandleSendMessageError(err)
		return err
	}

	return nil
}

func (h *HandlerContext) SendAudio(audio tgbotapi.FileBytes, caption string) error {
	msg := tgbotapi.NewAudio(h.Update.FromChat().ChatConfig().ChatID, audio)
	if caption != "" {
		msg.Caption = caption
	}
	_, err := h.BotAPI.Send(msg)
	if err != nil {
		h.HandleSendMessageError(err)
		return err
	}

	return nil
}

func (h *HandlerContext) SendVoice(voice tgbotapi.FileBytes) error {
	msg := tgbotapi.NewVoice(h.Update.FromChat().ChatConfig().ChatID, voice)
	_, err := h.BotAPI.Send(msg)
	if err != nil {
		h.HandleSendMessageError(err)
		return err
	}

	return nil
}

func (h *HandlerContext) SendSticker(sticker tgbotapi.FileBytes) error {
	msg := tgbotapi.NewSticker(h.Update.FromChat().ChatConfig().ChatID, sticker)
	_, err := h.BotAPI.Send(msg)
	if err != nil {
		h.HandleSendMessageError(err)
		return err
	}

	return nil
}

func (h *HandlerContext) SendLocation(latitude, longitude float64) error {
	msg := tgbotapi.NewLocation(h.Update.FromChat().ChatConfig().ChatID, latitude, longitude)
	_, err := h.BotAPI.Send(msg)
	if err != nil {
		h.HandleSendMessageError(err)
		return err
	}

	return nil
}

func (h *HandlerContext) SendVenue(latitude, longitude float64, title, address string) error {
	msg := tgbotapi.NewVenue(h.Update.FromChat().ChatConfig().ChatID, title, address, latitude, longitude)
	_, err := h.BotAPI.Send(msg)
	if err != nil {
		h.HandleSendMessageError(err)
		return err
	}

	return nil
}

func (h *HandlerContext) SendContact(phoneNumber, firstName string) error {
	msg := tgbotapi.NewContact(h.Update.FromChat().ChatConfig().ChatID, phoneNumber, firstName)
	_, err := h.BotAPI.Send(msg)
	if err != nil {
		h.HandleSendMessageError(err)
		return err
	}

	return nil
}

func (h *HandlerContext) SendPoll(question string, options []string) error {
	msg := tgbotapi.NewPoll(h.Update.FromChat().ChatConfig().ChatID, question, options...)
	_, err := h.BotAPI.Send(msg)
	if err != nil {
		h.HandleSendMessageError(err)
		return err
	}

	return nil
}

func (h *HandlerContext) DeleteMessage(messageID int) error {
	deleteMsg := tgbotapi.NewDeleteMessage(h.Update.FromChat().ChatConfig().ChatID, messageID)
	_, err := h.BotAPI.Send(deleteMsg)
	if err != nil {
		h.HandleSendMessageError(err)
		return err
	}

	return nil
}

func (h *HandlerContext) EditMessageText(text string, messageID int) error {
	editMsg := tgbotapi.NewEditMessageText(h.Update.FromChat().ChatConfig().ChatID, messageID, text)
	_, err := h.BotAPI.Send(editMsg)
	if err != nil {
		h.HandleSendMessageError(err)
		return err
	}

	return nil
}

func (h *HandlerContext) EditMessageReplyMarkup(replyMarkup tgbotapi.InlineKeyboardMarkup, messageID int) error {
	editMsg := tgbotapi.NewEditMessageReplyMarkup(h.Update.FromChat().ChatConfig().ChatID, messageID, replyMarkup)
	_, err := h.BotAPI.Send(editMsg)
	if err != nil {
		h.HandleSendMessageError(err)
		return err
	}

	return nil
}

func (h *HandlerContext) GetChatID() int64 {
	return h.Update.FromChat().ChatConfig().ChatID
}

func (h *HandlerContext) GetUserID() int64 {
	return h.Update.SentFrom().ID
}

func (h *HandlerContext) GetUsername() string {
	if h.Update.SentFrom() != nil {
		return h.Update.SentFrom().UserName
	}
	return ""
}

func (h *HandlerContext) GetFirstName() string {
	if h.Update.SentFrom() != nil {
		return h.Update.SentFrom().FirstName
	}
	return ""
}

func (h *HandlerContext) GetLastName() string {
	if h.Update.SentFrom() != nil {
		return h.Update.SentFrom().LastName
	}
	return ""
}

func (h *HandlerContext) GetFullName() string {
	if h.Update.SentFrom() != nil {
		user := h.Update.SentFrom()
		if user.LastName != "" {
			return user.FirstName + " " + user.LastName
		}
		return user.FirstName
	}
	return ""
}

func (h *HandlerContext) IsPrivateChat() bool {
	return h.Update.FromChat().Type == "private"
}

func (h *HandlerContext) IsGroupChat() bool {
	return h.Update.FromChat().Type == "group"
}

func (h *HandlerContext) IsSupergroupChat() bool {
	return h.Update.FromChat().Type == "supergroup"
}

func (h *HandlerContext) IsChannel() bool {
	return h.Update.FromChat().Type == "channel"
}

func (h *HandlerContext) HasMessage() bool {
	return h.Update.Message != nil
}

func (h *HandlerContext) HasCallbackQuery() bool {
	return h.Update.CallbackQuery != nil
}

func (h *HandlerContext) HasEditedMessage() bool {
	return h.Update.EditedMessage != nil
}

func (h *HandlerContext) HasChannelPost() bool {
	return h.Update.ChannelPost != nil
}

func (h *HandlerContext) HasEditedChannelPost() bool {
	return h.Update.EditedChannelPost != nil
}

func (h *HandlerContext) HasInlineQuery() bool {
	return h.Update.InlineQuery != nil
}

func (h *HandlerContext) HasChosenInlineResult() bool {
	return h.Update.ChosenInlineResult != nil
}

func (h *HandlerContext) HasShippingQuery() bool {
	return h.Update.ShippingQuery != nil
}

func (h *HandlerContext) HasPreCheckoutQuery() bool {
	return h.Update.PreCheckoutQuery != nil
}

func (h *HandlerContext) HasPoll() bool {
	return h.Update.Poll != nil
}

func (h *HandlerContext) HasPollAnswer() bool {
	return h.Update.PollAnswer != nil
}

func (h *HandlerContext) HasMyChatMember() bool {
	return h.Update.MyChatMember != nil
}

func (h *HandlerContext) HasChatMember() bool {
	return h.Update.ChatMember != nil
}

func (h *HandlerContext) HasChatJoinRequest() bool {
	return h.Update.ChatJoinRequest != nil
}

func (h *HandlerContext) HandleSendMessageError(err error) {
	h.Logger.ErrorContext(h.Ctx, "Failed to send message", "error", err)
	h.SendError("Failed to send message. Please try again later.")
}
