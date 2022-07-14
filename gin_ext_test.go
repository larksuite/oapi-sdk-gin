package sdkginext

import (
	"context"
	"fmt"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/larksuite/oapi-sdk-go/v3/card"
	"github.com/larksuite/oapi-sdk-go/v3/core"
	"github.com/larksuite/oapi-sdk-go/v3/event/dispatcher"
	"github.com/larksuite/oapi-sdk-go/v3/service/contact/v3"
	"github.com/larksuite/oapi-sdk-go/v3/service/im/v1"
)

func TestStartGin(t *testing.T) {

	// 注册消息处理器
	handler := dispatcher.NewEventDispatcher("v", "").OnP2MessageReceiveV1(func(ctx context.Context, event *larkim.P2MessageReceiveV1) error {
		fmt.Println(larkcore.Prettify(event))
		fmt.Println(event.RequestId())
		return nil
	}).OnP2MessageReadV1(func(ctx context.Context, event *larkim.P2MessageReadV1) error {
		fmt.Println(larkcore.Prettify(event))
		fmt.Println(event.RequestId())
		return nil
	}).OnP2UserCreatedV3(func(ctx context.Context, event *larkcontact.P2UserCreatedV3) error {
		fmt.Println(larkcore.Prettify(event))
		fmt.Println(event.RequestId())
		return nil
	})

	// 创建card处理器
	cardHandler := larkcard.NewCardActionHandler("v", "", func(ctx context.Context, cardAction *larkcard.CardAction) (interface{}, error) {
		fmt.Println(larkcore.Prettify(cardAction))

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
