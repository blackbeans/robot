package robot

import (
	log "github.com/blackbeans/log4go"
	"github.com/blackbeans/turbo/pipe"
	"time"
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

	//try publish
	buff := WrapReq2Buff(*ae)
	// buff.Reset()
	// buff.WriteString("access_token=32ed8bb66f484350dda9890338976c8c&app=higo&backup=2&channel_id=2&client_id=1&cver=5.0.0&device_id=oudid_e34fa1eebdea2c7f4cf51e1ea3839ae303519b6b&device_token=&from=1200267081&group_id=15615&idfa=84FBA21D-C514-4D0E-82BE-1831912A0963&open_udid=e34fa1eebdea2c7f4cf51e1ea3839ae303519b6b&qudaoid=10000&sign=e308746c317bbfb14d3c22e3b20cabf5&source=mob&text=%EF%BC%9F&timestamp=1464443273384&to=&type=text&uuid=3698df89bb4de75de605cac3207e750d&ver=0.8&via=iphone&")

	req := WrapBuff2HttpRequest(self.url, buff)
	log.DebugLog("robot_handler", "PublishHandler|Publish|%s", buff.String())

	r, err := ae.ctx.client.Do(req)
	if nil != err {
		log.ErrorLog("robot_handler", "PublishHandler|Try Publish |FAIL|%s|%v", err, req.Form)
		return err
	}

	resp, err := UnmarshalResponse(r)
	if nil != err {
		log.ErrorLog("robot_handler", "PublishHandler|Try Login|UnmarshalResponse |FAIL|%s|%v", err, req.Form)
		return err
	}

	//if code eq 0 ,login success
	if resp.Code == 0 {
		log.InfoLog("robot_handler", "PublishHandler| Publish|SUCC|%s", resp.Message)
	} else {
		log.InfoLog("robot_handler", "PublishHandler| Publish|FAIL|%s|%s", resp.Code, resp)
	}

	//next send follow shopper

	return nil

}
