package higo

import (
	"encoding/json"
	"net/http"

	log "github.com/blackbeans/log4go"
	"github.com/blackbeans/turbo/pipe"
)

type HigoAccountReq struct {
	HigoCookie
	AccountId string `uri:"higo_account_ids"`
}

type HigoAccountResp struct {
	Body []struct {
		AccountId string `json:"higo_account_id"`
		MlsUserId string `json:"mls_user_id"`
	} `json:"map"`
}

//login
type LoginReq struct {
	HigoCookie
	pipe.IForwardEvent
	DeviceVersion string `uri:"device_version"`
	Mobile        string `uri:"mobile"`
	Password      string `uri:"password"`
}

type LoginResp struct {
	AccountInfo struct {
		AccountId     string `json:"account_id"`
		AccountMobile string `json:"account_mobile"`
		NickName      string `json:"nick_name"`
	} `json:"account_info"`
	Token string `json:"token"`
}

type RobotContext struct {
	client  *http.Client
	session *HigoSession
}

type LoginHandler struct {
	pipe.BaseForwardHandler
	url string
}

func NewLoginHandler(name, url string) *LoginHandler {

	handler := &LoginHandler{}
	handler.url = url
	handler.BaseForwardHandler = pipe.NewBaseForwardHandler(name, handler)
	return handler
}

func (self *LoginHandler) TypeAssert(event pipe.IEvent) bool {
	_, ok := self.cast(event)
	return ok
}

func (self *LoginHandler) cast(event pipe.IEvent) (val *LoginReq, ok bool) {
	val, ok = event.(*LoginReq)
	return
}

func (self *LoginHandler) Process(ctx *pipe.DefaultPipelineContext, event pipe.IEvent) error {

	ae, ok := self.cast(event)
	if !ok {
		return pipe.ERROR_INVALID_EVENT_TYPE
	}

	ae.UUID = UUID()
	ae.App = "higo"
	ae.Qudaoid = 10000
	ae.Backup = 2
	ae.Source = "mob"
	ae.ClientId = 1
	ae.Cver = "5.0.0"
	ae.DeviceVersion = "9.3.2"
	// ae.IDFA = "84FBA21D-C514-4D0E-82BE-1831912A0963"
	// ae.OpenUdid = "e34fa1eebdea2c7f4cf51e1ea3839ae303519b6b"
	ae.DeviceId = "oudid_" + ae.OpenUdid
	ae.Ver = 0.8
	ae.Via = "iphone"
	ae.Source = "mob"

	c := &http.Client{}
	context := &RobotContext{}
	context.client = c

	resp, err := HttpReq(c, "POST", self.url, *ae)

	//if code eq 0 ,login success
	if nil == err && resp.Code == 0 {

		var lresp LoginResp
		//get token
		err = json.Unmarshal(resp.Data, &lresp)
		if nil != err {
			log.ErrorLog("robot_handler", "LoginHandler| Login|SUCC|Unmarshal|FAIL|%s|%s", err, string(resp.Data))
			return err
		}

		token := lresp.Token

		ae.HigoCookie.AccessToken = token
		//get mls_user_id

		hareq := &HigoAccountReq{}
		hareq.HigoCookie = ae.HigoCookie
		hareq.AccountId = lresp.AccountInfo.AccountId

		//try mls_user_id
		resp, err := HttpReq(c, "POST", "http://v.lehe.com/account/GetHigoAccountId2MlsUserIdMap", *hareq)

		if nil == err && resp.Code == 0 {

			var haresp HigoAccountResp
			//get token
			err = json.Unmarshal(resp.Data, &haresp)
			if nil != err {
				log.ErrorLog("robot_handler", "LoginHandler| HigoAccountReq|SUCC|Unmarshal|FAIL|%s|%s", err, string(resp.Data))
				return err
			}

			log.InfoLog("robot_handler", "LoginHandler|Login|SUCC|%s", ae.Mobile)

			session := &HigoSession{}
			session.HigoCookie = ae.HigoCookie
			session.AccountId = lresp.AccountInfo.AccountId
			session.AccountMobile = lresp.AccountInfo.AccountMobile
			session.MlsUserId = haresp.Body[0].MlsUserId
			context.session = session

			//crawl channel
			channelReq := &ChannelReq{}
			channelReq.HigoSession = *session
			channelReq.ctx = context
			ctx.SendForward(channelReq)

		} else {
			log.WarnLog("robot_handler", "LoginHandler|Login|FAIL|HigoAccountReq|%s", resp.Message)
		}

		// log.InfoLog("robot_handler", "LoginHandler| Login|SUCC|%s|%v", session)
	} else {
		log.InfoLog("robot_handler", "LoginHandler| Login|FAIL|%s|%s", resp.Code, resp.Data)
	}

	//next send follow shopper

	return nil

}
