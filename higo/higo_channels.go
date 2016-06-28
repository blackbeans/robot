package higo

import (
	"encoding/json"
	log "github.com/blackbeans/log4go"
	"github.com/blackbeans/turbo/pipe"
	"strconv"
)

//shop more  请求
type ChannelReq struct {
	HigoSession
	pipe.IForwardEvent
	ctx *RobotContext
}

type ChannelResp struct {
	BaseResp
	Channels []struct {
		ID          string `json:"id"`
		ChannelName string `json:"channel_name"`
	} `json:"list"`
	Total string `json:"total"`
}

type ChannelHandler struct {
	pipe.BaseForwardHandler
	url        string
	channelMod int
}

func NewChannelHandler(name, url string, channelMod int) *ChannelHandler {

	handler := &ChannelHandler{}
	handler.url = url
	handler.channelMod = channelMod

	handler.BaseForwardHandler = pipe.NewBaseForwardHandler(name, handler)
	return handler
}

func (self *ChannelHandler) TypeAssert(event pipe.IEvent) bool {
	_, ok := self.cast(event)
	return ok
}

func (self *ChannelHandler) cast(event pipe.IEvent) (val *ChannelReq, ok bool) {
	val, ok = event.(*ChannelReq)
	return
}

func (self *ChannelHandler) Process(ctx *pipe.DefaultPipelineContext, event pipe.IEvent) error {

	ae, ok := self.cast(event)
	if !ok {
		return pipe.ERROR_INVALID_EVENT_TYPE
	}

	ae.HigoSession = *ae.ctx.session

	//try open
	resp, err := HttpReq(ae.ctx.client, "GET", self.url, *ae)
	//if code eq 0 ,login success
	if nil == err && resp.Code == 0 {

		var channelResp ChannelResp
		err = json.Unmarshal(resp.Data, &channelResp)
		if nil != err {
			log.ErrorLog("robot_handler", "ChannelHandler|Open|FAIL|%s|%s", err, resp.Data)
		} else {
			for _, channel := range channelResp.Channels {
				//crawl channel shop
				// shopMore
				v, _ := strconv.Atoi(channel.ID[len(channel.ID)-1:])
				v = v % 2
				if v == self.channelMod {
					shopMore := &ShopMoreReq{}
					shopMore.ctx = ae.ctx
					shopMore.ID = channel.ID
					log.InfoLog("robot_handler", "ChannelHandler|Start Channel|%s|%s", channel.ID, channel.ChannelName)
					ctx.SendForward(shopMore)

				}
			}
		}

	}

	return nil

}
