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

type taskResponse struct {
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

type updateTaskRequest struct {
	ReqID       string `json:"reqId"`
	PjtID       string `json:"pjtId"`
	SubPjtID    string `json:"subPjtId"`
	Type        string `json:"type"`
	CmtCtnt     string `json:"cmtCtnt"`
	PrvsReqCntn string `json:"prvsReqCntn"`
	CrntReqCntn string `json:"crntReqCntn"`
	ReqCtnt     string `json:"reqCtnt"`
	Action      string `json:"action"`
	PstTpCd     string `json:"pstTpCd"`
}

type getTaskResponse struct {
	LstDayOff []struct {
		ClassName string `json:"className"`
		VacID     string `json:"vacId"`
		VacDt     string `json:"vacDt"`
		VacTpCd   string `json:"vacTpCd"`
		Mode      int    `json:"mode"`
	} `json:"lstDayOff"`
	LstUsrWtc  []interface{} `json:"lstUsrWtc"`
	LstUsrRole []struct {
		ClassName string `json:"className"`
		UsrID     string `json:"usrId"`
		PjtID     string `json:"pjtId"`
		RoleID    string `json:"roleId"`
		UsrNm     string `json:"usrNm"`
		ImgURL    string `json:"imgUrl"`
		ID        string `json:"id"`
		RoleNm    string `json:"roleNm"`
		Mode      int    `json:"mode"`
	} `json:"lstUsrRole"`
	DetailReqVO struct {
		CreateUser    string        `json:"createUser"`
		CreateDate    string        `json:"createDate"`
		ClassName     string        `json:"className"`
		ReqID         string        `json:"reqId"`
		ReqCtnt       string        `json:"reqCtnt"`
		ReqTitNm      string        `json:"reqTitNm"`
		PjtID         string        `json:"pjtId"`
		PmUsrID       string        `json:"pmUsrId"`
		CustUsrID     string        `json:"custUsrId"`
		ReqPhsCd      string        `json:"reqPhsCd"`
		BizProcID     string        `json:"bizProcId"`
		SubPjtID      string        `json:"subPjtId"`
		CateID        string        `json:"cateId"`
		JbTpCd        string        `json:"jbTpCd"`
		PntNo         int           `json:"pntNo"`
		PlnDueDt      string        `json:"plnDueDt"`
		ImptTpCd      string        `json:"imptTpCd"`
		PrntReqID     string        `json:"prntReqId"`
		ReqStsCd      string        `json:"reqStsCd"`
		ItrtnID       string        `json:"itrtnId"`
		SeqNo         int           `json:"seqNo"`
		PlnStDt       string        `json:"plnStDt"`
		ConfFlg       string        `json:"confFlg"`
		ArrFileRegist []interface{} `json:"arrFileRegist"`
		PjtNm         string        `json:"pjtNm"`
		CateNm        string        `json:"cateNm"`
		PhsNm         string        `json:"phsNm"`
		PmUsrNm       string        `json:"pmUsrNm"`
		FileKnt       string        `json:"fileKnt"`
		BizProcNm     string        `json:"bizProcNm"`
		Path          string        `json:"path"`
		JbTpNm        string        `json:"jbTpNm"`
		ImptNm        string        `json:"imptNm"`
		SumPctNo      int           `json:"sumPctNo"`
		StrPlnDueDt   string        `json:"strPlnDueDt"`
		CustFlg       string        `json:"custFlg"`
		RqstrNm       string        `json:"rqstrNm"`
		RqstrDt       string        `json:"rqstrDt"`
		ItrtnNm       string        `json:"itrtnNm"`
		StatusDueDate int           `json:"statusDueDate"`
		PhsCd         string        `json:"phsCd"`
		UsrID         string        `json:"usrId"`
		ReqStsNm      string        `json:"reqStsNm"`
		TotalLike     int           `json:"totalLike"`
		Mode          int           `json:"mode"`
	} `json:"detailReqVO"`
	LstSkdUsr []struct {
		ClassName  string        `json:"className"`
		SkdID      string        `json:"skdId"`
		UsrID      string        `json:"usrId"`
		ProcPhsID  string        `json:"procPhsId"`
		EstmDueDt  string        `json:"estmDueDt"`
		ActDueDt   string        `json:"actDueDt,omitempty"`
		PlnDueDt   string        `json:"plnDueDt"`
		ProcFlg    string        `json:"procFlg"`
		PlnStDt    string        `json:"plnStDt"`
		ImgURL     string        `json:"imgUrl"`
		PrntPhsID  string        `json:"prntPhsId"`
		PrntPhsCd  string        `json:"prntPhsCd"`
		PhsCd      string        `json:"phsCd"`
		PctNo      int           `json:"pctNo"`
		UsrNm      string        `json:"usrNm"`
		PhsNm      string        `json:"phsNm"`
		EfrtPctNo  int           `json:"efrtPctNo"`
		AsgnList   []interface{} `json:"asgnList"`
		IconCd     string        `json:"iconCd"`
		Mode       int           `json:"mode"`
		WrkHrNo    float64       `json:"wrkHrNo,omitempty"`
		RsrcRoleCd string        `json:"rsrcRoleCd,omitempty"`
	} `json:"lstSkdUsr"`
	LstEv []struct {
		ClassName string `json:"className"`
		ComCd     string `json:"comCd"`
		PrntCd    string `json:"prntCd"`
		CdNm      string `json:"cdNm"`
		CdShrtNm  string `json:"cdShrtNm"`
		ImgURL    string `json:"imgUrl"`
		Mode      int    `json:"mode"`
	} `json:"lstEv"`
	CompanyVO struct {
		ClassName string  `json:"className"`
		EndTmNo   float64 `json:"endTmNo"`
		LunHrNo   float64 `json:"lunHrNo"`
		StLunTmNo float64 `json:"stLunTmNo"`
		StTmNo    float64 `json:"stTmNo"`
		WrkHrNo   float64 `json:"wrkHrNo"`
		WrkDyCd   string  `json:"wrkDyCd"`
		UtcTmNo   float64 `json:"utcTmNo"`
		Mode      int     `json:"mode"`
	} `json:"companyVO"`
	LstJbDetails []interface{} `json:"lstJbDetails"`
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

	var taskResponse taskResponse

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

func CreateTask(template string, title string, content string) {
	// TODO: implement this function
	fmt.Printf("Creating task with template: %s, title: %s, content: %s\n", template, title, content)
}

func GetTask(taskId string) (getTaskResponse, error) {
	url := baseURL + "/api/searchRequirementDetails?reqId=" + taskId

	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		return getTaskResponse{}, err
	}

	req.Header.Set("Accept", "application/json, text/plain, */*")

	res, err := client.Do(req)
	if err != nil {
		return getTaskResponse{}, err
	}

	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)

	if err != nil {
		return getTaskResponse{}, err
	}

	var taskResponse getTaskResponse
	err = json.Unmarshal(body, &taskResponse)
	if err != nil {
		return getTaskResponse{}, err
	}

	return taskResponse, nil
}

func UpdateTaskContent(taskID string, content string) error {
	currentTask, err := GetTask(taskID)
	if err != nil {
		return err
	}

	url := baseURL + "/api/update-content"
	payload := updateTaskRequest{
		ReqID:       taskID,
		Type:        "reqCtnt", // Update content
		Action:      "REQ_WTC_CNG",
		PstTpCd:     "PST_TP_CDACT",
		PjtID:       currentTask.DetailReqVO.PjtID,
		SubPjtID:    currentTask.DetailReqVO.SubPjtID,

        CrntReqCntn: currentTask.DetailReqVO.ReqCtnt,
		PrvsReqCntn: currentTask.DetailReqVO.ReqCtnt,

		CmtCtnt: "<div class='system-comment'> • Changed Content: </div>",
		ReqCtnt: content,
	}

	jsonPayload, err := json.Marshal(payload)

	if err != nil {
		return err
	}

	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return err
	}

	req.Header.Set("Accept", "application/json, text/plain, */*")
    req.Header.Set("Content-Type", "application/json;charset=utf-8")
	res, err := client.Do(req)

	if err != nil {
		return err
	}

	defer res.Body.Close()
	return nil
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