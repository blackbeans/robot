package xr

import (
	"encoding/json"
	log "github.com/blackbeans/log4go"
	"github.com/blackbeans/turbo/pipe"
	"gopkg.in/redis.v3"
)

type XRHotReq struct {
	pipe.IForwardEvent
	ctx *RobotContext
	XRCookie
	PageNo int `json:"page"`
}

type XRHotResp struct {
	Total int `json:"total"`
	List  []struct {
		UserId    string  `json:"user_id"`
		LiveId    string  `json:"live_id"`
		Avatar    string  `json:"avatar"`
		NickName  string  `json:"nickname"`
		Latitude  float32 `json:"latitude"`
		Longitude float32 `json:"longtiude"`
	} `json:"list"`
}

type HotListHandler struct {
	pipe.BaseForwardHandler
	url         string
	redisClient *redis.Client
}

func NewHotListHandler(name, url string, redisClient *redis.Client) *HotListHandler {

	handler := &HotListHandler{}
	handler.url = url
	handler.BaseForwardHandler = pipe.NewBaseForwardHandler(name, handler)
	handler.redisClient = redisClient
	return handler
}

func (self *HotListHandler) TypeAssert(event pipe.IEvent) bool {
	_, ok := self.cast(event)
	return ok
}

func (self *HotListHandler) cast(event pipe.IEvent) (val *XRHotReq, ok bool) {
	val, ok = event.(*XRHotReq)
	return
}

func (self *HotListHandler) Process(ctx *pipe.DefaultPipelineContext, event pipe.IEvent) error {

	ae, ok := self.cast(event)
	if !ok {
		return pipe.ERROR_INVALID_EVENT_TYPE
	}

	for {

		resp, err := HttpReq(ae.ctx.client, "POST", self.url, *ae)
		//if code eq 0 ,HotList success
		if nil == err && resp.Status == 100 {

			var lresp XRHotResp
			err = json.Unmarshal(resp.Result, &lresp)
			if nil != err {
				log.ErrorLog("robot_handler", "HotListHandler| HotList|FAIL|Unmarshal|FAIL|%s|%s", err, string(resp.Message))
				return err
			}

			log.InfoLog("robot_handler", "HotListHandler| HotList|SUCC|%v", lresp)
			// ctx.SendForward(hotReq)

			// log.InfoLog("robot_handler", "HotListHandler| HotList|SUCC|%s|%v", session)
		} else {
			log.InfoLog("robot_handler", "HotListHandler| HotList|SUCC|%s|%s", resp.Status, resp.Result)
		}

		break
	}
	return nil

}
