package robot

import (
	log "github.com/blackbeans/log4go"
	"github.com/blackbeans/turbo/pipe"
)

//openIm  请求
type OpenImReq struct {
	HigoSession
	pipe.IForwardEvent
	ctx     *RobotContext
	GroupId string `uri:"group_id"`
	// Count  int `uri:"count"`
	// NextId int `uri:"next_id"`
}

type OpenImHandler struct {
	pipe.BaseForwardHandler
	url string
}

func NewOpenImHandler(name, url string) *OpenImHandler {

	handler := &OpenImHandler{}
	handler.url = url
	handler.BaseForwardHandler = pipe.NewBaseForwardHandler(name, handler)
	return handler
}

func (self *OpenImHandler) TypeAssert(event pipe.IEvent) bool {
	_, ok := self.cast(event)
	return ok
}

func (self *OpenImHandler) cast(event pipe.IEvent) (val *OpenImReq, ok bool) {
	val, ok = event.(*OpenImReq)
	return
}

func (self *OpenImHandler) Process(ctx *pipe.DefaultPipelineContext, event pipe.IEvent) error {

	ae, ok := self.cast(event)
	if !ok {
		return pipe.ERROR_INVALID_EVENT_TYPE
	}

	// ae.NextId = 0
	// ae.Count = 10
	ae.HigoSession = *ae.ctx.session
	ae.GroupId = "d2ce37777da4f013c"
	//try open
	buff := WrapReq2Buff(*ae)
	// log.DebugLog("robot_handler", "OpenImHandler|Open|%s", buff.String())

	req := WrapBuff2HttpRequest(self.url, buff)

	r, err := ae.ctx.client.Do(req)
	if nil != err {
		log.ErrorLog("robot_handler", "OpenImHandler|Try Open |FAIL|%s|%v", err, req.PostForm)
		return err
	}

	resp, err := UnmarshalResponse(r)
	if nil != err {
		log.ErrorLog("robot_handler", "OpenImHandler|Try Open|UnmarshalResponse |FAIL|%s|%v", err, req.PostForm)
		return err
	}

	//if code eq 0 ,login success
	if resp.Code == 0 {

		log.InfoLog("robot_handler", "OpenImHandler|Open|SUCC|%s", resp.Message)
		//test send message
		publishReq := &PublishReq{}
		publishReq.ctx = ae.ctx
		publishReq.Type = "text"
		publishReq.GroupId = "15615"
		publishReq.Text = "hello"

		ctx.SendForward(publishReq)

	} else {
		log.InfoLog("robot_handler", "OpenImHandler|Open|FAIL|%s|%s", resp.Code, resp.Message)
	}

	//next send follow shopper

	return nil

}
