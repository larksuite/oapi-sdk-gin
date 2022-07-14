package sdkginext

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/larksuite/oapi-sdk-go/v3/card"
	"github.com/larksuite/oapi-sdk-go/v3/core"
	"github.com/larksuite/oapi-sdk-go/v3/event"
	"github.com/larksuite/oapi-sdk-go/v3/event/dispatcher"
)

func doProcess(writer http.ResponseWriter, req *http.Request, reqHandler larkevent.IReqHandler, options ...larkevent.OptionFunc) {
	// 转换http请求对象为标准请求对象
	ctx := context.Background()
	eventReq, err := translate(ctx, req)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		writer.Write([]byte(err.Error()))
		return
	}

	//处理请求
	eventResp := reqHandler.Handle(ctx, eventReq)

	// 回写结果
	err = write(ctx, writer, eventResp)
	if err != nil {
		reqHandler.Logger().Error(ctx, fmt.Sprintf("write resp result error:%s", err.Error()))
	}
}

func NewCardActionHandlerFunc(cardActionHandler *larkcard.CardActionHandler, options ...larkevent.OptionFunc) func(c *gin.Context) {

	// 构建模板类
	cardActionHandler.InitConfig(options...)
	return func(c *gin.Context) {
		doProcess(c.Writer, c.Request, cardActionHandler, options...)
	}
}

func NewEventHandlerFunc(eventDispatcher *dispatcher.EventDispatcher, options ...larkevent.OptionFunc) func(c *gin.Context) {
	eventDispatcher.InitConfig(options...)
	return func(c *gin.Context) {
		doProcess(c.Writer, c.Request, eventDispatcher, options...)
	}
}

func processError(ctx context.Context, logger larkcore.Logger, path string, err error) *larkevent.EventResp {
	header := map[string][]string{}
	header[larkevent.ContentTypeHeader] = []string{larkevent.DefaultContentType}
	eventResp := &larkevent.EventResp{
		Header:     header,
		Body:       []byte(fmt.Sprintf(larkevent.WebhookResponseFormat, err.Error())),
		StatusCode: http.StatusInternalServerError,
	}
	logger.Error(ctx, fmt.Sprintf("event handle err:%s, %v", path, err))
	return eventResp
}

func write(ctx context.Context, writer http.ResponseWriter, eventResp *larkevent.EventResp) error {
	writer.WriteHeader(eventResp.StatusCode)
	for k, vs := range eventResp.Header {
		for _, v := range vs {
			writer.Header().Add(k, v)
		}
	}

	if len(eventResp.Body) > 0 {
		_, err := writer.Write(eventResp.Body)
		return err
	}
	return nil
}
func translate(ctx context.Context, req *http.Request) (*larkevent.EventReq, error) {
	rawBody, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return nil, err
	}
	eventReq := &larkevent.EventReq{
		Header:     req.Header,
		Body:       rawBody,
		RequestURI: req.RequestURI,
	}

	return eventReq, nil
}
