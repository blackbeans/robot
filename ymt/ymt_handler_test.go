package ymt

import (
	"testing"
	"time"

	"gopkg.in/redis.v3"

	log "github.com/blackbeans/log4go"
	"github.com/blackbeans/turbo/pipe"
)

func TestPublishHandler(t *testing.T) {

	log.LoadConfiguration("../log.xml")

	redisclient := redis.NewClient(&redis.Options{
		Addr: "localhost:6379"})

	line := pipe.NewDefaultPipeline()
	line.RegisteHandler("login", NewLoginHandler("login", "http://app.ymatou.com/api/Auth/LoginAuth"))
	line.RegisteHandler("activities", NewChannelHandler("activities", "http://app.ymatou.com/api/activity/GetCountryGroupList", redisclient))
	line.RegisteHandler("im", NewPublishHandler("im", "http://app.ymatou.com/api/Letter/AddMessage", "hi", redisclient))

	req := &LoginReq{}
	req.Username = ""
	req.Password = ""
	line.FireWork(req)
	time.Sleep(2 * time.Second)

}
