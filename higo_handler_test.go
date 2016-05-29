package robot

import (
	log "github.com/blackbeans/log4go"
	"github.com/blackbeans/turbo/pipe"
	"testing"
	"time"
)

func TestLoginHandler(t *testing.T) {

	log.LoadConfiguration("./log.xml")

	line := pipe.NewDefaultPipeline()
	line.RegisteHandler("login", NewLoginHandler("login", "http://v.lehe.com/account/login"))

	req := &LoginReq{}
	req.IDFA = "84FBA21D-C514-4D0E-82BE-1831912A0963"
	req.Mobile = "18612372884"
	req.OpenUdid = "e34fa1eebdea2c7f4cf51e1ea3839ae303519b6b"
	req.Password = "116295838"

	line.FireWork(req)

	time.Sleep(2 * time.Second)

}

func TestPublishHandler(t *testing.T) {

	log.LoadConfiguration("./log.xml")

	line := pipe.NewDefaultPipeline()
	line.RegisteHandler("login", NewLoginHandler("login", "http://v.lehe.com/account/login"))
	line.RegisteHandler("openIm", NewOpenImHandler("openIm", "http://v.lehe.com/group_chat/getGroupNotice"))
	line.RegisteHandler("publish", NewPublishHandler("publish", "http://im.lehe.com/im/publish"))

	req := &LoginReq{}
	req.IDFA = "84FBA21D-C514-4D0E-82BE-1831912A0963"
	req.Mobile = "18612372884"
	req.OpenUdid = "e34fa1eebdea2c7f4cf51e1ea3839ae303519b6b"
	req.Password = "116295838"
	line.FireWork(req)

	time.Sleep(2 * time.Second)

}
