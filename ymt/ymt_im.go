package ymt

import (
	"bytes"
	"encoding/json"

	log "github.com/blackbeans/log4go"
	"github.com/blackbeans/turbo/pipe"
)

//im  请求
type PublishReq struct {
	YmtSession
	pipe.IForwardEvent
	ctx      *RobotContext
	ToUserId int64  `uri:"ToUserId",json:"ToUserId"`
	Message  string `uri:"Message",json:"Message"`
}

type PublishHandler struct {
	pipe.BaseForwardHandler
	url     string
	message string
}

func NewPublishHandler(name, url string, message string) *PublishHandler {

	handler := &PublishHandler{}
	handler.url = url
	handler.BaseForwardHandler = pipe.NewBaseForwardHandler(name, handler)
	handler.message = message
	return handler
}

func (self *PublishHandler) TypeAssert(event pipe.IEvent) bool {
	_, ok := self.cast(event)
	return ok
}

func (self *PublishHandler) cast(event pipe.IEvent) (val *PublishReq, ok bool) {
	val, ok = event.(*PublishReq)
	return
}

func (self *PublishHandler) Process(ctx *pipe.DefaultPipelineContext, event pipe.IEvent) error {

	ae, ok := self.cast(event)
	if !ok {
		return pipe.ERROR_INVALID_EVENT_TYPE
	}

	ae.Message = self.message
	data, _ := json.Marshal(*ae)

	buff := WrapReq2Buff(*ae)

	request := WrapBuff2HttpRequest("POST", self.url+"?"+buff.String(), bytes.NewBuffer(data))
	request.Header.Set("Content-Type", "application/json")

	resp, err := HttpReqAndDecode(ae.ctx.client, request)
	if nil == err && resp.Status == 200 {
		log.InfoLog("robot_handler", "PublishHandler|Publish Message|SUCC|%d|%s", ae.ToUserId, ae.Message)

	} else {
		log.WarnLog("robot_handler", "PublishHandler|Publish Message|FAIL|%s|%d|%s", resp.Message, ae.ToUserId, ae.Message)
	}

	return nil

}
