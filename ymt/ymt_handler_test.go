package ymt

import (
	"testing"
	"time"

	log "github.com/blackbeans/log4go"
	"github.com/blackbeans/turbo/pipe"
)

func TestPublishHandler(t *testing.T) {

	log.LoadConfiguration("../log.xml")

	// redisclient := redis.NewClient(&redis.Options{
	// 	Addr: "localhost:6379"})

	line := pipe.NewDefaultPipeline()
	line.RegisteHandler("login", NewLoginHandler("login", "http://app.ymatou.com/api/Auth/LoginAuth"))
	line.RegisteHandler("activities", NewChannelHandler("activities", "http://app.ymatou.com/api/activity/GetCountryGroupList"))
	line.RegisteHandler("im", NewPublishHandler("im", "http://app.ymatou.com/api/Letter/AddMessage", "hi"))

	req := &LoginReq{}
	req.Username = "18612372884"
	req.Password = "0301151313"
	line.FireWork(req)
	time.Sleep(2 * time.Second)

}
