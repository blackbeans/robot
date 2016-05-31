package higo

import (
	"encoding/json"
	"strconv"
	"time"

	log "github.com/blackbeans/log4go"
	"github.com/blackbeans/turbo/pipe"
	"gopkg.in/redis.v3"
)

//shop more  请求
type ShopMoreReq struct {
	HigoSession
	pipe.IForwardEvent
	ctx *RobotContext

	ID     string `uri:"id"`
	PageNo int    `uri:"p"`
	Size   int    `uri:"size"`
}

type ShopMoreResp struct {
	BaseResp
	HigoGroupIds []struct {
		HigoGroupId string `json:"group_id"`
	} `json:"list"`
	Total  string `json:"total"`
	PageNo string `json:"p"`
	Size   string `json:"size"`
}

//shop detail
type ShopDetailReq struct {
	HigoSession
	HigoGroupId string `uri:"group_id"`
}

//shop detail
type ShopDetail struct {
	ID           string `json"id"`          //群组的ID
	HigoGroupId  string `json:"group_id"`   //higo群组ID
	AccountId    string `json:"account_id"` //卖家账号ID
	ShopId       string `json:"shop_id"`
	GroupName    string `json:"group_name"`
	MlsAccountId int64  `json:"mls_account_id"`
}

type ShopMoreHandler struct {
	pipe.BaseForwardHandler
	url         string
	redisClient *redis.Client
}

func NewShopMoreHandler(name, url string, redisClient *redis.Client) *ShopMoreHandler {

	handler := &ShopMoreHandler{}
	handler.url = url
	handler.BaseForwardHandler = pipe.NewBaseForwardHandler(name, handler)
	handler.redisClient = redisClient
	return handler
}

func (self *ShopMoreHandler) TypeAssert(event pipe.IEvent) bool {
	_, ok := self.cast(event)
	return ok
}

func (self *ShopMoreHandler) cast(event pipe.IEvent) (val *ShopMoreReq, ok bool) {
	val, ok = event.(*ShopMoreReq)
	return
}

func (self *ShopMoreHandler) Process(ctx *pipe.DefaultPipelineContext, event pipe.IEvent) error {

	ae, ok := self.cast(event)
	if !ok {
		return pipe.ERROR_INVALID_EVENT_TYPE
	}

	ae.PageNo = 1
	// ae.ID = 66
	ae.Size = 20

	shops := make([]ShopDetail, 0, 100)
	ae.HigoSession = *ae.ctx.session
	hasMore := true
	for hasMore {
		//try open
		resp, err := HttpReq(ae.ctx.client, "POST", self.url, *ae)
		//if code eq 0 ,login success
		if nil == err && resp.Code == 0 {

			var shopResp ShopMoreResp
			err = json.Unmarshal(resp.Data, &shopResp)
			if nil != err {
				log.ErrorLog("robot_handler", "ShopMoreHandler|Open|FAIL|%s|%s", err, resp.Data)
			} else {
				for _, shop := range shopResp.HigoGroupIds {
					//check group has been followed
					result := self.redisClient.ZScore("_higo_group_followed", shop.HigoGroupId)
					if err = result.Err(); nil == err && result.Val() > 0 {
						//skipped
						log.WarnLog("robot_handler", "ShopMoreHandler|Followed|SKIPPED|%s", shop.HigoGroupId)
						continue
					}

					//get shop detail
					detailReq := &ShopDetailReq{}
					detailReq.HigoSession = ae.HigoSession
					detailReq.HigoGroupId = shop.HigoGroupId

					resp, err = HttpReq(ae.ctx.client, "GET", "http://v.lehe.com/shop/Get_group_detail", *detailReq)
					if nil == err && resp.Code == 0 {
						var shopDetail ShopDetail
						err = json.Unmarshal(resp.Data, &shopDetail)
						if nil != err {
							log.ErrorLog("robot_handler", "ShopMoreHandler|ShopDetail|FAIL|%s|%s", err, resp.Data)
						} else {

							//save higo group info
							self.redisClient.Set("higo:"+shop.HigoGroupId+":info", string(resp.Data), -1)

							shops = append(shops, shopDetail)
							log.InfoLog("robot_handler", "ShopMoreHandler|Shop|%v", shopDetail)
							showFollow := &ShopFollowReq{}
							showFollow.ctx = ae.ctx
							showFollow.HigoSession = ae.HigoSession
							showFollow.HigoGroupId = shopDetail.HigoGroupId
							showFollow.GroupId = shopDetail.ID
							showFollow.HigoId = ae.AccountId
							ctx.SendForward(showFollow)

							//100 ms follow shop
							time.Sleep(100 * time.Millisecond)
						}
					} else {
						log.WarnLog("robot_handler", "ShopMoreHandler|ShopDetail|FAIL|%s|HigoGroupId:%s", resp.Code, shop.HigoGroupId)
					}
				}
			}
			//hasMore
			pageNo, _ := strconv.Atoi(shopResp.PageNo)
			pageSize, _ := strconv.Atoi(shopResp.Size)
			total, _ := strconv.Atoi(shopResp.Total)
			hasMore = pageNo*pageSize < total
		} else {
			log.WarnLog("robot_handler", "ShopMoreHandler|Shop|FAIL|%s|%s", resp.Code, resp.Message)
		}

		ae.PageNo++

	}

	return nil

}
