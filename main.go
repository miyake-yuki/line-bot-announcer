package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/line/line-bot-sdk-go/linebot"
	"github.com/line/line-bot-sdk-go/linebot/httphandler"
)

func main() {
	// HTTP Handlerの初期化
	handler, err := httphandler.New(
		os.Getenv("CHANNEL_SECRET"),
		os.Getenv("CHANNEL_TOKEN"),
	)
	if err != nil {
		log.Fatal(err)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("Defaulting to port %s", port)
	}

	// 実際にRequestを受け取った時に処理を行うHandle関数を定義し、handlerに登録
	handler.HandleEvents(func(events []*linebot.Event, r *http.Request) {
		bot, err := handler.NewClient()
		if err != nil {
			log.Print(err)
			return
		}

		for _, event := range events {
			if event.Type != linebot.EventTypeMessage {
				return
			}

			switch message := event.Message.(type) {
			case *linebot.TextMessage:
				if message.Text == "leave" {
					leaveGroup(bot, event.Source)
					return
				}

				replyText := message.Text
				if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(replyText)).Do(); err != nil {
					log.Print(err)
				}
			}
		}
	})

	// /callback にエンドポイントの定義
	http.Handle("/callback", handler)
	// HTTPサーバの起動
	log.Printf("Listening on port %s", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}

// グループ・トークルームから退出させる
func leaveGroup(bot *linebot.Client, eventSource *linebot.EventSource) {
	switch eventSource.Type {
	case linebot.EventSourceTypeGroup:
		if _, err := bot.LeaveGroup(eventSource.GroupID).Do(); err != nil {
			log.Print(err)
		}
	case linebot.EventSourceTypeRoom:
		if _, err := bot.LeaveRoom(eventSource.RoomID).Do(); err != nil {
			log.Print(err)
		}
	}
}

// データベースにグループIDを登録する
func registerGroupID() {}

// データベースからグループIDを削除する
func deleteGropID() {}
