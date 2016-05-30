package higo

import (
	"time"

	log "github.com/blackbeans/log4go"
	"github.com/blackbeans/turbo/pipe"
)

//group  请求
type PublishReq struct {
	HigoSession
	pipe.IForwardEvent
	ctx         *RobotContext
	ChannelId   int    `uri:"channel_id"`
	DeviceToken string `uri:"device_token"`
	From        string `uri:"from"`
	To          string `uri:"to"`
	GroupId     string `uri:"group_id"`
	Sign        string `uri:"sign"`
	Text        string `uri:"text"`
	TimeStamp   int64  `uri:"timestamp"`
	Type        string `uri:"type"`
}

type PublishHandler struct {
	pipe.BaseForwardHandler
	url string
}

func NewPublishHandler(name, url string) *PublishHandler {

	handler := &PublishHandler{}
	handler.url = url
	handler.BaseForwardHandler = pipe.NewBaseForwardHandler(name, handler)
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

	ae.ChannelId = 2
	ae.From = ae.ctx.session.MlsUserId
	ae.Sign = UUID()
	ae.TimeStamp = time.Now().UnixNano() / 1000 / 1000
	ae.HigoSession = *ae.ctx.session

	resp, err := HttpReq(ae.ctx.client, "POST", self.url, *ae)

	//if code eq 0 ,login success
	if nil != err && resp.Code == 0 {
		log.InfoLog("robot_handler", "PublishHandler| Publish|SUCC|%s", resp.Message)
	} else {
		log.InfoLog("robot_handler", "PublishHandler| Publish|FAIL|%s|%s", resp.Code, resp)
	}

	//next send follow shopper

	return nil

}
