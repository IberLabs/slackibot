package cookie

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

type Cookie struct {
	header auth.Header
	login  resources.Login
	Cookie string
}

type message struct {
	UserName string `json:"username"`
	Password string `json:"password"`
}

func (cookie Cookie) PreAuth() {
	m := message{cookie.login.UserName, cookie.login.Password}
	body, err := json.Marshal(m)
	if err != nil {
		panic(err)
	}
	jsonStr := []byte(body)
	req, err := http.NewRequest("POST", cookie.login.Url, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	fmt.Println("Requesting cookie to:", cookie.login.Url)
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	fmt.Println("Response Status:", resp.Status)
	fmt.Println("Response Headers:", resp.Header)
	body, _ = ioutil.ReadAll(resp.Body)
	fmt.Println("Response Body:", string(body))

	setCookie(resp.Header["Set-Cookie"])

}

func setCookie(header []string) {
	for i := 0; i < 4; i++ {
		if strings.Contains(header[i], "crowd.token_key") {
			c := header[i][16]
			if c != '"' {
				fmt.Println("Taking cookie from Header", header[i])
				var cookieB bytes.Buffer
				cookieB.WriteString("crowd.token_key=")
				for j := 16; c != ';'; {
					cookieB.WriteString(string(c))
					j++
					c = header[i][j]
				}
				fmt.Println("Saving cookie", cookieB.String())
				cookieS := Cookie{Cookie: cookieB.String()}
				cookieJ, err := json.Marshal(cookieS)
				if err != nil {
					panic(err)
				}
				err = ioutil.WriteFile("./resources/cookie.json", cookieJ, 0644)
				if err != nil {
					panic(err)
				}
			}
		}
	}
}

func (cookie Cookie) GetHeader() auth.Header {
	return cookie.header
}

func New() Cookie {
	var c Cookie
	login := getLogin()
	c.header = auth.Header{}
	c.login = login
	return c
}

func getLogin() resources.Login {
	prop := resources.GetLogin()
	login := resources.Login{UserName: prop.UserName, Password: prop.Password}
	return login
}
