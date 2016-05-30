package ymt

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"reflect"
	"strings"

	log "github.com/blackbeans/log4go"

	"github.com/blackbeans/go-uuid"
)

type YmtCookie struct {
	DeviceId    string `uri:"DeviceId",json:"DeviceId"`
	CKId        string `uri:"CKId",json:"CKId"`
	ClientType  string `uri:"ClientType",json:"ClientType"`
	CookieId    string `uri:"Cookieid",json:"Cookieid"`
	WIFI        string `uri:"WIFI",json:"WIFI"`
	DeviceToken string `uri:"DeviceToken",json:"DeviceToken"`
	ClientId    string `uri:"ClientId",json:"ClientId"`
	IDFA        string `uri:"IDFA",json:"IDFA"`
	Guid        string `uri:"Guid",json:"Guid"`
	VersionInfo string `uri:"versionInfo",json:"versionInfo"`
	AppName     string `uri:"AppName",json:"AppName"`
	Yid         string `uri:"yid",json:"yid"`
	AccessToken string `uri:"AccessToken",json:"AccessToken"`
}

type BaseResp struct {
	Status  int             `json:"Status"`
	Message string          `json:"Msg"`
	Result  json.RawMessage `json:"Result"`
}

type YmtSession struct {
	YmtCookie
	UserId int64 `uri:"UserId",json:"UserId"`
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
			case reflect.Slice:
				data, _ := json.Marshal(at.Field(i).Interface())
				s.WriteString(string(data))
				s.WriteString("&")
			}

		} else if f.Type.Kind() == reflect.Struct {
			s.WriteString(WrapReq2Buff(at.Field(i).Interface()).String())
		}

	}
	return s
}

func WrapBuff2HttpRequest(method string, url string, buff *bytes.Buffer) *http.Request {
	var req *http.Request
	if strings.ToUpper(method) == "GET" {
		url += "?"
		url += buff.String()
		req, _ = http.NewRequest(method, url, nil)
	} else {
		req, _ = http.NewRequest(method, url, buff)
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	// req.Header.Add("User-Agent", "Ymt/5.0.0 (iPhone; iOS 9.3.2; Scale/3.00)")
	req.Header.Set("User-Agent", "01ios====9.3.2=e420a82d-9fb9-4f3a-a20a-bc1215c145e8=======")
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

func HttpReq(client *http.Client, method string, url string, req interface{}) (*BaseResp, error) {
	buff := WrapReq2Buff(req)
	log.DebugLog("robot_handler", "HttpReq|%s", buff.String())
	httpreq := WrapBuff2HttpRequest(method, url, buff)
	r, err := client.Do(httpreq)
	if nil != err {
		log.ErrorLog("robot_handler", "HttpReq|Try Open |FAIL|%s|%v", err, httpreq.PostForm)
		return nil, err
	}

	resp, err := UnmarshalResponse(r)
	if nil != err {
		log.ErrorLog("robot_handler", "HttpReq|Try Open|UnmarshalResponse |FAIL|%s|%v", err, httpreq.PostForm)
		return nil, err
	}
	return &resp, nil
}

func HttpReqAndDecode(client *http.Client, httpreq *http.Request) (*BaseResp, error) {
	r, err := client.Do(httpreq)
	if nil != err {
		log.ErrorLog("robot_handler", "HttpReq|Try Open |FAIL|%s|%v", err, httpreq.PostForm)
		return nil, err
	}

	resp, err := UnmarshalResponse(r)
	if nil != err {
		log.ErrorLog("robot_handler", "HttpReq|Try Open|UnmarshalResponse |FAIL|%s|%v", err, httpreq.PostForm)
		return nil, err
	}
	return &resp, nil
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
