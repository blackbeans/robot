package higo

import (
	"encoding/json"
	"time"

	"gopkg.in/redis.v3"

	log "github.com/blackbeans/log4go"
	"github.com/blackbeans/turbo/pipe"
)

//shop follow  请求
type ShopFollowReq struct {
	HigoSession
	pipe.IForwardEvent
	ctx *RobotContext

	GroupId     string `uri:"group_id"`
	HigoGroupId string `uri:"higo_group_id"`
	HigoId      string `uri:"higo_id"`
}

type ShopFollowResp struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

//shop more  请求
type Shop struct {
	ShopId    string `json:"shop_id"`
	GroupId   string `json:"group_id"`
	GroupName string `json:"group_name"`
}

type ShopFollowHandler struct {
	pipe.BaseForwardHandler
	url         string
	redisClient *redis.Client
}

func NewShopFollowHandler(name, url string, redisClient *redis.Client) *ShopFollowHandler {

	handler := &ShopFollowHandler{}
	handler.url = url
	handler.BaseForwardHandler = pipe.NewBaseForwardHandler(name, handler)
	handler.redisClient = redisClient
	return handler
}

func (self *ShopFollowHandler) TypeAssert(event pipe.IEvent) bool {
	_, ok := self.cast(event)
	return ok
}

func (self *ShopFollowHandler) cast(event pipe.IEvent) (val *ShopFollowReq, ok bool) {
	val, ok = event.(*ShopFollowReq)
	return
}

func (self *ShopFollowHandler) Process(ctx *pipe.DefaultPipelineContext, event pipe.IEvent) error {

	ae, ok := self.cast(event)
	if !ok {
		return pipe.ERROR_INVALID_EVENT_TYPE
	}

	ae.HigoSession = *ae.ctx.session

	// //try open
	resp, err := HttpReq(ae.ctx.client, "POST", self.url, *ae)

	//if code eq 0 ,login success
	if nil == err && resp.Code == 0 {

		var shopResp ShopFollowResp
		err = json.Unmarshal(resp.Data, &shopResp)
		if nil != err {
			log.ErrorLog("robot_handler", "ShopFollowHandler|Follow|FAIL|%s|%s", err, resp.Data)
		} else {
			self.redisClient.ZAdd("_higo_group_followd", redis.Z{Score: float64(time.Now().Unix()), Member: ae.HigoGroupId})
			log.InfoLog("robot_handler", "ShopFollowHandler|Follow|SUCC|%d|%s", resp.Code, resp.Message)
		}

	} else {
		log.WarnLog("robot_handler", "ShopFollowHandler|Shop|HttpReq|FAIL|%s|%s", resp.Code, resp.Message)
	}

	return nil

}
