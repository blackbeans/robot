package ymt

import (
	"encoding/json"

	"math"

	log "github.com/blackbeans/log4go"
	"github.com/blackbeans/turbo/pipe"
)

//shop more  请求
type ChannelReq struct {
	YmtSession
	pipe.IForwardEvent
	ctx *RobotContext
}

type ChannelResp struct {
	CountryId   int    `json:"CountryId"`
	CountryName string `json:"CountryName"`
}

type ChannelActivityIdsReq struct {
	YmtSession
	CountryId  int    `uri:"CountryId"`
	FilterType int    `uri:"FilterType"`
	OnlyFollow string `uri:"OnlyFollow"`
}

type ChannelActivityIdsResp struct {
	ActivityCount int `json:"ActivityCount"`
	Activities    []struct {
		ActivityId int64 `json:"ActivityId"`
	} `json:"Activities"`
}
type ChannelHandler struct {
	pipe.BaseForwardHandler
	url string
}

func NewChannelHandler(name, url string) *ChannelHandler {

	handler := &ChannelHandler{}
	handler.url = url
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

	ae.YmtSession = *ae.ctx.session

	//try open
	resp, err := HttpReq(ae.ctx.client, "GET", self.url, *ae)
	//if code eq 0 ,login success
	if nil == err && resp.Status == 200 {

		var channelResp []ChannelResp
		err = json.Unmarshal(resp.Result, &channelResp)
		if nil != err {
			log.ErrorLog("robot_handler", "ChannelHandler|Open|FAIL|%s|%s", err, resp.Result)
		} else {
			for _, channel := range channelResp {
				// crawl channel shop
				// shopMore
				activityIds := &ChannelActivityIdsReq{}
				activityIds.YmtSession = ae.YmtSession
				activityIds.CountryId = channel.CountryId

				//try open
				resp, err := HttpReq(ae.ctx.client, "GET", "http://app.ymatou.com/api/Activity/ListInProgressActivityIds", *activityIds)
				if nil == err && resp.Status == 200 {
					var idsResp ChannelActivityIdsResp
					err = json.Unmarshal(resp.Result, &idsResp)
					if nil != err {
						log.ErrorLog("robot_handler", "ChannelHandler|Open|FAIL|%s|%s", err, resp.Result)
					} else {
						pno := 1
						pageSize := 10
						maxPno := idsResp.ActivityCount/pageSize + 1
						for pno < maxPno {
							endIdx := int(math.Min(float64(pno*pageSize), float64(maxPno)))
							tmp := idsResp.Activities[(pno-1)*pageSize : endIdx]
							idarr := make([]int64, 0, 10)
							for _, id := range tmp {
								idarr = append(idarr, id.ActivityId)
							}

							//get activityInfo
							log.InfoLog("robot_handler", "ChannelHandler|Start Channel|%v", idarr)

							pno++

						}
					}
				}

				// ctx.SendForward(shopMore)
			}
		}

	} else {
		log.ErrorLog("robot_handler", "ChannelHandler|Open|FAIL|%s", resp.Status, resp.Message)
	}

	return nil

}
