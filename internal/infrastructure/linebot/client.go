package linebot

import (
	"bytes"

	"github.com/tian841224/stock-bot/config"
	"github.com/tian841224/stock-bot/internal/infrastructure/imgbb"
	"github.com/tian841224/stock-bot/pkg/logger"

	"github.com/line/line-bot-sdk-go/linebot"
	"go.uber.org/zap"
)

type LineBotClient struct {
	Client *linebot.Client
	logger logger.Logger
}

// NewBot 初始化 LINE Bot
func NewBot(cfg config.Config, log logger.Logger) (*LineBotClient, error) {
	client, err := linebot.New(cfg.CHANNEL_SECRET, cfg.CHANNEL_ACCESS_TOKEN)
	if err != nil {
		return nil, err
	}
	return &LineBotClient{Client: client, logger: log}, nil
}

// ReplyMessage 回覆文字訊息
func (b *LineBotClient) ReplyMessage(replyToken, text string) error {
	_, err := b.Client.ReplyMessage(replyToken, linebot.NewTextMessage(text)).Do()
	if err != nil {
		b.logger.Error("發送訊息失敗", zap.Error(err))
	}
	return err
}

// ReplyMessageWithButtons 回覆帶有按鈕的訊息
func (b *LineBotClient) ReplyMessageWithButtons(replyToken, text string, buttons []linebot.TemplateAction) error {
	if len(buttons) == 0 {
		return b.ReplyMessage(replyToken, text)
	}

	template := linebot.NewButtonsTemplate(
		"", "", text, buttons...,
	)

	_, err := b.Client.ReplyMessage(replyToken, linebot.NewTemplateMessage("按鈕", template)).Do()
	if err != nil {
		b.logger.Error("發送帶有按鈕的訊息失敗", zap.Error(err))
	}
	return err
}

// ReplyImage 回覆圖片訊息
func (b *LineBotClient) ReplyImage(replyToken, imageURL string) error {
	imageMessage := linebot.NewImageMessage(imageURL, imageURL)
	_, err := b.Client.ReplyMessage(replyToken, imageMessage).Do()
	if err != nil {
		b.logger.Error("發送圖片訊息失敗", zap.Error(err))
	}
	return err
}

// ReplyPhoto 上傳圖片並回覆（需要 ImgBB 客戶端）
func (b *LineBotClient) ReplyPhoto(replyToken string, data []byte, caption string, imgbbClient *imgbb.ImgBBClient) error {
	// 如果沒有 ImgBB 客戶端，只發送文字訊息
	if imgbbClient == nil {
		b.logger.Warn("ImgBB 客戶端未設定，只發送文字訊息")
		return b.ReplyMessage(replyToken, caption)
	}

	// 上傳圖片到 ImgBB
	options := &imgbb.UploadOptions{
		Name: "stock_chart",
	}

	reader := bytes.NewReader(data)
	resp, err := imgbbClient.UploadFromFile(reader, "chart.png", options)
	if err != nil {
		b.logger.Error("上傳圖片到 ImgBB 失敗", zap.Error(err))
		// 如果上傳失敗，只發送文字訊息
		return b.ReplyMessage(replyToken, caption)
	}

	// 發送圖片
	return b.ReplyImage(replyToken, resp.Data.URL)
}
