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
	// line.RegisteHandler("shop_channel", NewChannelHandler("shop_channel", "http://v.lehe.com/shop/get_dimensions"))
	// line.RegisteHandler("shop_more", NewShopMoreHandler("shop_more", "http://v.lehe.com/shop/get_more", redisclient))
	// line.RegisteHandler("shop_follow", NewShopFollowHandler("shop_follow", "http://im.lehe.com/im/open_group_add", redisclient))

	req := &LoginReq{}
	req.Username = "18612372884"
	req.Password = "0301151313"
	line.FireWork(req)
	time.Sleep(2 * time.Second)

}
