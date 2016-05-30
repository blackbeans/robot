package main

import (
	"flag"
	"robot/core"
	"time"

	"github.com/blackbeans/go-uuid"
	log "github.com/blackbeans/log4go"
	"github.com/blackbeans/turbo/pipe"
	"gopkg.in/redis.v3"
)

func main() {
	redisHost := flag.String("redis", "localhost:6379", "-redis=localhost:6379")
	mobile := flag.String("mobile", "", "-mobile=1862222222")
	password := flag.String("password", "", "-password=1234")

	flag.Parse()
	log.LoadConfiguration("./log.xml")
	redisclient := redis.NewClient(&redis.Options{
		Addr: *redisHost})

	line := pipe.NewDefaultPipeline()
	line.RegisteHandler("login", core.NewLoginHandler("login", "http://v.lehe.com/account/login"))
	line.RegisteHandler("shop_channel", core.NewChannelHandler("shop_channel", "http://v.lehe.com/shop/get_dimensions"))
	line.RegisteHandler("shop_more", core.NewShopMoreHandler("shop_more", "http://v.lehe.com/shop/get_more", redisclient))
	line.RegisteHandler("shop_follow", core.NewShopFollowHandler("shop_follow", "http://im.lehe.com/im/open_group_add", redisclient))

	req := &core.LoginReq{}
	req.IDFA = "84FBA21D-C514-4D0E-82BE-1831912A0963"
	req.Mobile = *mobile
	req.OpenUdid = uuid.NewRandom().String()
	req.Password = *password
	line.FireWork(req)
	time.Sleep(2 * time.Second)
}
