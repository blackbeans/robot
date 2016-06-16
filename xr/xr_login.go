package xr

import (
	"encoding/json"
	"net/http"

	// "github.com/blackbeans/go-uuid"
	log "github.com/blackbeans/log4go"
	"github.com/blackbeans/turbo/pipe"
)

//login
type LoginReq struct {
	XRCookie
	pipe.IForwardEvent
	Mobile   string `uri:"mobile"`
	Password string `uri:"password"`
}

type LoginResp struct {
	Token string `json:"token"`
	User  struct {
		Id        string  `json:"id"`
		Latitude  float32 `json:"latitude"`
		Longitude float32 `json:"longtiude"`
		Sign      string  `json:"sign"`
	} `json:"user"`
}

type RobotContext struct {
	client  *http.Client
	session *XRSession
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

	c := &http.Client{}
	context := &RobotContext{}
	context.client = c

	resp, err := HttpReq(c, "POST", self.url, *ae)

	//if code eq 0 ,login success
	if nil == err && resp.Status == 100 {

		var lresp LoginResp
		//get token
		err = json.Unmarshal(resp.Result, &lresp)
		if nil != err {
			log.ErrorLog("robot_handler", "LoginHandler| Login|FAIL|Unmarshal|FAIL|%s|%s", err, string(resp.Message))
			return err
		}

		cookies := XRCookie{}
		cookies.Latitude = lresp.User.Latitude
		cookies.Longitude = lresp.User.Longitude
		cookies.Sign = lresp.User.Sign
		cookies.Token = lresp.Token
		cookies.UserId = lresp.User.Id
		cookies.Version = "2.0.3"

		// //set session
		session := &XRSession{}
		session.XRCookie = cookies

		context.session = session
		context.client = c

		hotReq := &XRHotReq{}
		hotReq.XRCookie = cookies
		hotReq.PageNo = 0

		log.InfoLog("robot_handler", "LoginHandler| Login|SUCC|%v", lresp)

		ctx.SendForward(hotReq)

		// log.InfoLog("robot_handler", "LoginHandler| Login|SUCC|%s|%v", session)
	} else {
		log.InfoLog("robot_handler", "LoginHandler| Login|SUCC|%s|%s", resp.Status, resp.Result)
	}

	return nil

}
