package higo

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
	line.RegisteHandler("login", NewLoginHandler("login", "http://v.lehe.com/account/login"))
	line.RegisteHandler("shop_channel", NewChannelHandler("shop_channel", "http://v.lehe.com/shop/get_dimensions"))
	line.RegisteHandler("shop_more", NewShopMoreHandler("shop_more", "http://v.lehe.com/shop/get_more", redisclient))
	line.RegisteHandler("shop_follow", NewShopFollowHandler("shop_follow", "http://im.lehe.com/im/open_group_add", redisclient))

	req := &LoginReq{}
	req.IDFA = "84FBA21D-C514-4D0E-82BE-1831912A0963"
	req.Mobile = "18612372884"
	req.OpenUdid = "e34fa1eebdea2c7f4cf51e1ea3839ae303519b6b"
	req.Password = "116295838"
	line.FireWork(req)
	time.Sleep(2 * time.Second)

}
