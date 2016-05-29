package robot

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/blackbeans/go-uuid"
	"io/ioutil"
	"net/http"
	"net/url"
	"reflect"
)

type HigoCookie struct {
	App         string `uri:"app"`
	AccessToken string `uri:"access_token"`
	ClientId    int    `uri:"client_id"`
	Cver        string `uri:"cver"`
	DeviceId    string `uri:"device_id"`
	GetConfig   string `uri:"getConfig"`
	Qudaoid     int    `uri:"qudaoid"`
	Source      string `uri:"source"`

	Backup      int     `uri:"backup"`
	DeviceToken string  `uri:"device_token"`
	IDFA        string  `uri:"idfa"`
	OpenUdid    string  `uri:"open_udid"`
	Ver         float32 `uri:"ver"`
	Via         string  `uri:"via"`
	UUID        string  `uri:"uuid"`
}

type BaseResp struct {
	Code    int             `json:"code"`
	Message string          `json:"message"`
	Data    json.RawMessage `json:"data"`
}

type HigoSession struct {
	HigoCookie
	MlsUserId     string
	AccountId     string
	AccountMobile string
}

func WrapReq2Buff(greq interface{}) *bytes.Buffer {

	s := bytes.NewBuffer(make([]byte, 0, 1024))

	at := reflect.ValueOf(greq)
	t := reflect.TypeOf(greq)

	count := at.NumField()
	for i := 0; i < count; i++ {
		f := t.Field(i)

		name := f.Tag.Get("uri")
		if len(name) > 0 {
			s.WriteString(name)
			s.WriteString("=")
			k := f.Type.Kind()
			switch k {
			case reflect.Int, reflect.Int64:
				s.WriteString(fmt.Sprintf("%d", at.Field(i).Int()))
				s.WriteString("&")
			case reflect.Float32, reflect.Float64:
				s.WriteString(fmt.Sprintf("%.1f", at.Field(i).Float()))
				s.WriteString("&")
			case reflect.String:
				fs := at.Field(i).Interface().(string)
				s.WriteString(url.QueryEscape(fs))
				s.WriteString("&")
			}

		} else if f.Type.Kind() == reflect.Struct {
			s.WriteString(WrapReq2Buff(at.Field(i).Interface()).String())
		}

	}
	return s
}

func WrapBuff2HttpRequest(url string, buff *bytes.Buffer) *http.Request {
	req, _ := http.NewRequest("POST", url, buff)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("User-Agent", "HIGO/5.0.0 (iPhone; iOS 9.3.2; Scale/3.00)")

	return req
}

func UnmarshalResponse(resp *http.Response) (BaseResp, error) {
	var baseResp BaseResp

	body, err := ioutil.ReadAll(resp.Body)
	if nil != err {
		return baseResp, err
	}
	defer resp.Body.Close()

	err = json.Unmarshal(body, &baseResp)
	if nil != err {
		return baseResp, err
	}

	return baseResp, nil
}

//生成messageId uuid
func UUID() string {
	id := uuid.NewRandom()
	if id == nil || len(id) != 16 {
		return ""
	}
	b := []byte(id)
	return fmt.Sprintf("%08x%04x%04x%04x%012x",
		b[:4], b[4:6], b[6:8], b[8:10], b[10:])
}
