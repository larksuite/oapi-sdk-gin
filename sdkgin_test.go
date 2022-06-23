package sdkgin

import (
	"context"
	"fmt"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/larksuite/oapi-sdk-go/card"
	"github.com/larksuite/oapi-sdk-go/core"
	"github.com/larksuite/oapi-sdk-go/event/dispatcher"
	"github.com/larksuite/oapi-sdk-go/service/contact/v3"
	"github.com/larksuite/oapi-sdk-go/service/im/v1"
)

func TestStartGin(t *testing.T) {

	handler := dispatcher.NewEventDispatcher("v", "1212121212").OnMessageReceiveV1(func(ctx context.Context, event *im.MessageReceiveEvent) error {
		fmt.Println(core.Prettify(event))
		return nil
	}).OnMessageMessageReadV1(func(ctx context.Context, event *im.MessageMessageReadEvent) error {
		fmt.Println(core.Prettify(event))
		return nil
	}).OnUserCreatedV3(func(ctx context.Context, event *contact.UserCreatedEvent) error {
		fmt.Println(core.Prettify(event))
		return nil
	})

	// 创建card处理器
	cardHandler := card.NewCardActionHandler("12", "", func(ctx context.Context, cardAction *card.CardAction) (interface{}, error) {
		fmt.Println(core.Prettify(cardAction))

		// 返回卡片消息
		//return getCard(),nil

		//custom resp
		//return getCustomResp(),nil

		// 无返回值
		return nil, nil
	})

	g := gin.Default()

	g.POST("/webhook/event", NewEventHandlerFunc(handler))
	g.POST("/webhook/card", NewCardActionHandlerFunc(cardHandler))

	err := g.Run(":9999")
	if err != nil {
		panic(err)
	}

}
