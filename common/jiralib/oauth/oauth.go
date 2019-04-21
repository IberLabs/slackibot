package oauth

import (
	"pes_manager/auth"
	"pes_manager/resources"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

type Oauth struct {
	header auth.Header
	url    string
}

func (oauth Oauth) PreAuth() {
	url := oauth.url

	payload := strings.NewReader("{\"grant_type\":\"client_credentials\",\"client_id\": \"OauthKey\",\"client_secret\": \"vQ6L7y\"}")

	req, _ := http.NewRequest("POST", url, payload)

	req.Header.Add("content-type", "application/json")

	res, _ := http.DefaultClient.Do(req)
	body, _ := ioutil.ReadAll(res.Body)
	fmt.Println("Response Body:", string(body))

	defer res.Body.Close()
	setToken(res.Header["Set-Cookie"])
}

func setToken(header []string) {
	for i := 0; i < len(header); i++ {
		if strings.Contains(header[i], "atlassian.xsrf.token") {
			c := strings.Split(header[i], ";")
			fmt.Println("Taking token from Header", header[i])
			var oauth bytes.Buffer
			oauth.WriteString(string(c[0]))

			fmt.Println("Saving Oauth", oauth.String())
			oauthS := resources.Token{Token: oauth.String()}
			oauthJ, err := json.Marshal(oauthS)
			if err != nil {
				panic(err)
			}
			err = ioutil.WriteFile("./resources/oauth.json", oauthJ, 0644)
			if err != nil {
				panic(err)
			}
		}
	}
}

func (oauth Oauth) GetHeader() auth.Header {
	return oauth.header
}

func New() Oauth {
	var o Oauth
	o.header = auth.Header{}
	o.url = getUrl()
	return o
}

func getUrl() string {
	prop := resources.GetProperties("./resources/properties.json")
	return prop.BaseUrl
}
