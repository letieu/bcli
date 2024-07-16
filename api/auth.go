package api

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"os"
	"path"
	"regexp"
	"strings"

	"go.nhat.io/cookiejar"
)

const homePage = "https://blueprint.cyberlogitec.com.vn/"

var authClient *http.Client
var cookiePath string
var jar *cookiejar.PersistentJar

type UserInfo struct {
	UsrNo      interface{} `json:"usrNo"`
	UsrID      string      `json:"usrId"`
	UsrNm      string      `json:"usrNm"`
	UsrEml     string      `json:"usrEml"`
	ComUsrSx   string      `json:"comUsrSx"`
	EmpeNo     interface{} `json:"empeNo"`
	CntCd      string      `json:"cntCd"`
	ImgURL     string      `json:"imgUrl"`
	OrzID      string      `json:"orzId"`
	CoCd       string      `json:"coCd"`
	CoNm       interface{} `json:"coNm"`
	OfcCd      interface{} `json:"ofcCd"`
	OrzNm      string      `json:"orzNm"`
}

func init() {
	home, _ := os.UserHomeDir()
    cookiePath = path.Join(home, ".bcli", "cookies.json")

	jar = cookiejar.NewPersistentJar(
		cookiejar.WithFilePath(cookiePath),
		cookiejar.WithAutoSync(true),
	)

	authClient = &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
		Jar: jar,
	}
}

func getLoginPostTarget() (string, error) {
	loginPageUrl, err := getLoginPageUrl()
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("GET", loginPageUrl, nil)
	if err != nil {
		return "", err
	}

	resp, err := authClient.Do(req)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	authUrl, err := getAuthUrl(string(body))
	if err != nil {
		return "", err
	}

	return authUrl, nil
}

// go to home page
// redirect to ssologin page
// get redirect url from ssologin page
func getLoginPageUrl() (string, error) {
	/// 1. ==== Go to home page, set JSSESSIONID cookie
	req, err := http.NewRequest("GET", homePage, nil)
	if err != nil {
		return "", err
	}

	resp, err := authClient.Do(req)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()
    if resp.StatusCode == 200 {
        return "", errors.New("Already logged in")
    }

	ssoUrl := resp.Header.Get("Location")

	// 2. ==== Go to sso link, get redirect URL. set OAuth_Token_Request_State cookie
	req, err = http.NewRequest("GET", ssoUrl, nil)
	if err != nil {
		return "", err
	}

	req.Header.Set("Upgrade-Insecure-Requests", "1")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,*/*;q=0.8")
	req.Header.Set("Accept-Language", "en-US,en;q=0.5")
	req.Header.Set("Accept-Encoding", "gzip, deflate, br, zstd")
	req.Header.Set("DNT", "1")
	req.Header.Set("Sec-GPC", "1")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Sec-Fetch-Dest", "document")
	req.Header.Set("Sec-Fetch-Mode", "navigate")
	req.Header.Set("Sec-Fetch-Site", "none")
	req.Header.Set("Sec-Fetch-User", "?1")
	req.Header.Set("Priority", "u=1")

	resp, err = authClient.Do(req)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	redirectUrl := resp.Header.Get("Location")
	return redirectUrl, nil
}

func getAuthUrl(loginPage string) (string, error) {
	re := regexp.MustCompile(`action="([^"]+)"`)

	// Find the first match
	match := re.FindStringSubmatch(loginPage)
	if len(match) < 2 {
		return "", errors.New("no action URL found")
	}

	// The first submatch (index 1) is the URL
	authUrl := match[1]

	// Remove amp from the URL
	authUrl = strings.ReplaceAll(authUrl, "&amp;", "&")
	return authUrl, nil
}

func authenticate(username, password string, postUrl string) (string, error) {
	formData := url.Values{}
	formData.Set("username", username)
	formData.Set("password", password)
	formData.Set("credentialId", "")

	req, err := http.NewRequest("POST", postUrl, strings.NewReader(formData.Encode()))

	if err != nil {
		return "", err
	}

	req.Header.Set("Upgrade-Insecure-Requests", "1")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,*/*;q=0.8")
	req.Header.Set("Accept-Language", "en-US,en;q=0.5")
	req.Header.Set("Accept-Encoding", "gzip, deflate, br, zstd")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := authClient.Do(req)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	redirectUrl := resp.Header.Get("Location")
	return redirectUrl, nil
}

func ssoLogin(ssoUrl string) error {
	req, _ := http.NewRequest("GET", ssoUrl, nil)
	req.Header.Set("Upgrade-Insecure-Requests", "1")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,*/*;q=0.8")
	req.Header.Set("Accept-Language", "en-US,en;q=0.5")
	req.Header.Set("Accept-Encoding", "gzip, deflate, br, zstd")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Sec-Fetch-Dest", "document")
	req.Header.Set("Sec-Fetch-Mode", "navigate")
	req.Header.Set("Sec-Fetch-Site", "none")
	req.Header.Set("Sec-Fetch-User", "?1")
	req.Header.Set("Priority", "u=1")

	resp, err := authClient.Do(req)

	if err != nil {
		return err
	}

	defer resp.Body.Close()
	return nil
}

func Login(username, password string) error {
	postUrl, err := getLoginPostTarget()
	if err != nil {
		return err
	}

	ssoUrl, err := authenticate(username, password, postUrl)
	if err != nil {
		return err
	}

	err = ssoLogin(ssoUrl)
	if err != nil {
		return err
	}

	return nil
}

func Logout() error {
    err := os.Remove(cookiePath)
    if err != nil {
        return errors.New("Already logged out")
    }

    return nil
}

func GetUserInfo() (UserInfo, error) {
    req, _ := http.NewRequest("GET", "https://blueprint.cyberlogitec.com.vn/api/getUserInfo", nil)

    req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64; rv:126.0) Gecko/20100101 Firefox/126.0")
    req.Header.Set("Accept", "application/json, text/plain, */*")

    resp, err := authClient.Do(req)

    if err != nil {
        return UserInfo{}, err
    }

    defer resp.Body.Close()

    var userInfo UserInfo

    json.NewDecoder(resp.Body).Decode(&userInfo)
    return userInfo, nil
}
