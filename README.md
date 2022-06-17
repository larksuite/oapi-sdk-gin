# oapi-sdk-gin
an  oapi-sdk-go extension package that integrates the Gin Web framework


# 使用示例

```go
package main

import (
	"context"
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/larksuite/oapi-sdk-gin"
	"github.com/larksuite/oapi-sdk-go"
	"github.com/larksuite/oapi-sdk-go/card"
	"github.com/larksuite/oapi-sdk-go/core"
	"github.com/larksuite/oapi-sdk-go/dispatcher"
	"github.com/larksuite/oapi-sdk-go/service/contact/v3"
	"github.com/larksuite/oapi-sdk-go/service/im/v1"
)


func main() {

	// 创建消息事件处理器
	handler := dispatcher.NewEventReqDispatcher("v", "1212121212").OnMessageReceiveV1(func(ctx context.Context, event *im.MessageReceiveEvent) error {
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

	// 创建gin服务
	g := gin.Default()

	// 配置路由
	g.POST("/webhook/event", sdkgin.NewEventReqHandlerFunc(handler))
	g.POST("/webhook/card", sdkgin.NewCardActionHandlerFunc(cardHandler))

	// 启动服务
	err := g.Run(":9999")
	if err != nil {
		panic(err)
	}
}


```
