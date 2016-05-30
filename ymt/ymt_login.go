package ymt

import (
	"encoding/json"
	"net/http"

	"github.com/blackbeans/go-uuid"
	log "github.com/blackbeans/log4go"
	"github.com/blackbeans/turbo/pipe"
)

type YmtAccountReq struct {
	YmtCookie
	AccountId string `uri:"Ymt_account_ids"`
}

type YmtAccountResp struct {
	Body []struct {
		AccountId string `json:"Ymt_account_id"`
		MlsUserId string `json:"mls_user_id"`
	} `json:"map"`
}

//login
type LoginReq struct {
	YmtCookie
	pipe.IForwardEvent
	Username string `uri:"Username"`
	Password string `uri:"Password"`
}

type LoginResp struct {
	UserId      int64  `json:"UserId"`
	AccessToken string `json:"AccessToken"`
}

type RobotContext struct {
	client  *http.Client
	session *YmtSession
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

	ae.DeviceId = uuid.NewRandom().String()
	ae.CKId = UUID()
	ae.CookieId = UUID()
	ae.WIFI = "0"
	ae.ClientType = "1"
	ae.ClientId = UUID()
	ae.IDFA = uuid.NewRandom().String()
	ae.Guid = uuid.NewRandom().String()
	ae.VersionInfo = "2.6.9"
	ae.AppName = "Buyer"
	ae.Yid = uuid.NewRandom().String()

	c := &http.Client{}
	context := &RobotContext{}
	context.client = c

	resp, err := HttpReq(c, "POST", self.url, *ae)

	//if code eq 0 ,login success
	if nil == err && resp.Status == 200 {

		var lresp LoginResp
		//get token
		err = json.Unmarshal(resp.Result, &lresp)
		if nil != err {
			log.ErrorLog("robot_handler", "LoginHandler| Login|FAIL|Unmarshal|FAIL|%s|%s", err, string(resp.Message))
			return err
		}

		ae.AccessToken = lresp.AccessToken

		//set session
		session := &YmtSession{}
		session.YmtCookie = ae.YmtCookie
		session.UserId = lresp.UserId
		session.AccessToken = lresp.AccessToken

		context.session = session
		context.client = c

		chanReq := &ChannelReq{}
		chanReq.YmtSession = *session
		chanReq.ctx = context
		log.InfoLog("robot_handler", "LoginHandler| Login|SUCC|%v", session)

		ctx.SendForward(chanReq)

		// log.InfoLog("robot_handler", "LoginHandler| Login|SUCC|%s|%v", session)
	} else {
		log.InfoLog("robot_handler", "LoginHandler| Login|FAIL|%s|%s", resp.Status, resp.Result)
	}

	return nil

}
