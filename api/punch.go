package api

import (
	"fmt"
	"net/http"
)

// curl 'https://blueprint.cyberlogitec.com.vn/api/checkInOut/insert' -X POST -H 'User-Agent: Mozilla/5.0 (X11; Linux x86_64; rv:126.0) Gecko/20100101 Firefox/126.0' -H 'Accept: application/json, text/plain, */*' -H 'Accept-Language: en-US,en;q=0.5' -H 'Accept-Encoding: gzip, deflate, br, zstd' -H 'Origin: https://blueprint.cyberlogitec.com.vn' -H 'Connection: keep-alive' -H 'Referer: https://blueprint.cyberlogitec.com.vn/UI_TAT_028' -H 'Cookie: OAuth_Token_Request_State=39e7d9fc-9094-4d59-9758-c8aea2bffad9; JSESSIONID=F6740F92D29DE500FEB5F1CC3C6DF199; _ga=GA1.3.432067738.1721062976; _gid=GA1.3.1711277899.1721062976; _gat_gtag_UA_180475649_1=1' -H 'Sec-Fetch-Dest: empty' -H 'Sec-Fetch-Mode: cors' -H 'Sec-Fetch-Site: same-origin' -H 'Priority: u=1' -H 'Content-Length: 0'
func Punch() error {
	url := baseURL + "/api/checkInOut/insert"
	method := "POST"
	req, _ := http.NewRequest(method, url, nil)

	req.Header.Set("Accept", "application/json, text/plain, */*")
	req.Header.Set("Content-Type", "application/json;charset=utf-8")

	res, err := authClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

    fmt.Println(res.Status)
    return nil
}
