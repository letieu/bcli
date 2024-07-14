package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"path"
	"strings"
)

var baseURL = "https://blueprint.cyberlogitec.com.vn"
var client *http.Client

func init() {
	jar, err := loadCookies()
	if err != nil {
		panic(err)
	}

	client = &http.Client{
		Jar: jar,
	}
}

// TaskResponse represents the structure of the entire JSON response
type TaskResponse struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
	Data    struct {
		LstTaskVO  []Task `json:"lstTaskVO"`  // Open
		LstWrkShft []Task `json:"lstWrkShft"` // In progress
		LstTskDn   []Task `json:"lstTskDn"`   // Done
		TotalTsk   int    `json:"totalTsk"`
	} `json:"data"`
}

type Task struct {
	ClassName     string `json:"className"`
	UsrID         string `json:"usrId"`
	PjtNm         string `json:"pjtNm"`
	ReqID         string `json:"reqId"`
	TaskNm        string `json:"taskNm"`
	PhsNm         string `json:"phsNm"`
	PlnDueDt      string `json:"plnDueDt"`
	EstmDueDt     string `json:"estmDueDt"`
	JbTpNm        string `json:"jbTpNm"`
	AsgneeNm      string `json:"asgneeNm"`
	OrdrByNo      int    `json:"ordByNo"`
	SeqNo         int    `json:"seqNo"`
	ImptTpCd      string `json:"imptTpCd"`
	WradFlg       string `json:"wradFlg"`
	StatusDueDate int    `json:"statusDueDate"`
	ImgUrl        string `json:"imgUrl"`
	ColrVal       string `json:"colrVal"`
	EstmFinishDt  string `json:"estmFinishDt"`
	ReqStsCd      string `json:"reqStsCd"`
	IsUpdJobType  bool   `json:"isUpdJobType"`
	IsUpdContent  bool   `json:"isUpdContent"`
	IsUpdateTit   bool   `json:"isUpdateTit"`
	IsUpdCate     bool   `json:"isUpdCate"`
	IsUpdPriority bool   `json:"isUpdPriority"`
	Mode          int    `json:"mode"`
}

type Tasks struct {
	Open []Task
	InP  []Task
	Done []Task
}

func ListTasks() (Tasks, error) {
	url := baseURL + "/api/home/search-tasks"
	method := "POST"

	payload := []byte(`{"pjtId":"","duraMon":0,"multiSearch":"","reqTitNm":"","picId":"","seqNo":""}`)

	req, err := http.NewRequest(method, url, bytes.NewBuffer(payload))
	if err != nil {
		return Tasks{}, err
	}

	req.Header.Set("Accept", "application/json, text/plain, */*")
	req.Header.Set("Content-Type", "application/json;charset=utf-8")

	res, err := client.Do(req)
	if err != nil {
		return Tasks{}, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)

	if err != nil {
		fmt.Println(err)
		return Tasks{}, err
	}

	var taskResponse TaskResponse

	err = json.Unmarshal(body, &taskResponse)

	if err != nil {
		return Tasks{}, err
	}

	return Tasks{
		Open: taskResponse.Data.LstTaskVO,
		InP:  taskResponse.Data.LstWrkShft,
		Done: taskResponse.Data.LstTskDn,
	}, nil
}

func loadCookies() (http.CookieJar, error) {
	home, _ := os.UserHomeDir()
	cookiePath := path.Join(home, ".bcli/cookie")

	file, err := os.Open(cookiePath)
	if err != nil {
		return nil, err
	}

	defer file.Close()

	buf := new(bytes.Buffer)
	buf.ReadFrom(file)

	str := buf.String()
	parts := strings.Split(str, ";")

	jar, err := cookiejar.New(nil)
	if err != nil {
		return nil, err
	}

	url, _ := url.Parse(baseURL)

	cookies := []*http.Cookie{}
	for _, c := range parts {
		cookie := &http.Cookie{}
		cookieStr := strings.TrimSpace(c)
		cookieParts := strings.Split(cookieStr, "=")
		cookie.Name = cookieParts[0]
		cookie.Value = cookieParts[1]

		cookies = append(cookies, cookie)
	}

	jar.SetCookies(url, cookies)

	return jar, nil
}
