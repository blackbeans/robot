package xr

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
	line.RegisteHandler("login", NewLoginHandler("login", "https://api.xianroukeji.com/service/login"))
	line.RegisteHandler("hot", NewHotListHandler("hot", "	https://api.xianroukeji.com/service/top_1_6_list", redisclient))
	// line.RegisteHandler("im", NewPublishHandler("im", "http://app.ymatou.com/api/Letter/AddMessage", "hi", redisclient))
	//jpeg_fuzzy.jpg
	req := &LoginReq{}
	req.Mobile = "18612372884"
	req.Password = "854121"
	line.FireWork(req)
	time.Sleep(2 * time.Second)

}
