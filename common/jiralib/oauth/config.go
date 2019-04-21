package oauth

import (
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"github.com/dghubble/oauth1"
	"golang.org/x/net/context"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/user"
	"path/filepath"
	"strings"
)

var (
	jiraUrl = url.URL{Host: "https://staging.tools.adidas-group.com/jira/"}
)

/*
   $ openssl genrsa -out jira.pem 1024
   $ openssl rsa -in jira.pem -pubout -out jira.pub
*/
const jiraPrivateKey = `-----BEGIN RSA PRIVATE KEY-----
MIICXQIBAAKBgQDA5mbgFpVZUhCBuDgmNaF8HiMdeHi7hxqeArP2M/BPbZSIZ3D7
kgm5hQ7Fs8gXin8e27RRVTdxndVanI9k7skz7ETGsIEzPRMv6u2885Kmx5a1UwQ1
mt1mXshgPLl11UWykvsy34WwfxbKH0txjWhZeTDTjp6jA7VwxWSkznbDWQIDAQAB
AoGACogJsc5J1RiP4iUmm59t85LJpABBxys3Hs1S+ewYAJ4g79mF55Yvhbtn9Q89
q1taWVrxW0dlwYQ2c738bixDO8Q0qX4HzyeAyTebBUghadBWhBsQZF5cufRkDffh
K+QbNsHiCcZ15D6KCzvSfn7YXnzIhThml2iAbYMWohg7T8kCQQDem3mks9Ny6CnR
8HFH/4KYF9ceCxperAuPFr8Z1lPVIIgLz3EtXG5SPnwE7LuLGkROZk5rhRS6mp/z
IoTSa6anAkEA3dYcoT7D3Q+kDGFBttRlvdprC0vbXCsJKo9CsnmLOnMXj46EjtcS
EFkbvheRta+47jBnl+FVHJ0IYKEOdwGF/wJBALXjWcB/Ar33/vvAN/95Qg7eI/Iz
ZkeG0icHkfwdiQAzBZaI2FQVGztuPM2VVSQywS9CHr9xzN8wKpNyWA7K0S8CQEGL
lgH+rZiPmoUd53DB6R3jf2VjEHl3LcopcieRyhWHFBsSnRAnc+roqU3NYPwx445d
Nv6lUaSWsXb7n26CQLkCQQCmbk+36fKwLNW9FMu6j+GsUhnPsFsmZuYQXwfII1uB
vfIXN9woI0Qh//4F2nd/GogwM/OSRHyyu3pZPQyegmld
-----END RSA PRIVATE KEY-----`

const jiraConsumerKey = "OauthKey"

func getJIRAHTTPClient(ctx context.Context, config *oauth1.Config) *http.Client {
	cacheFile, err := jiraTokenCacheFile()
	if err != nil {
		log.Fatalf("Unable to get path to cached credential file. %v", err)
	}
	tok, err := jiraTokenFromFile(cacheFile)
	if err != nil {
		tok = getJIRATokenFromWeb(config)
		saveJIRAToken(cacheFile, tok)
	}
	return config.Client(ctx, tok)
}

func getJIRATokenFromWeb(config *oauth1.Config) *oauth1.Token {
	requestToken, requestSecret, err := config.RequestToken()
	if err != nil {
		log.Fatalf("Unable to get request token. %v", err)
	}
	authorizationURL, err := config.AuthorizationURL(requestToken)
	if err != nil {
		log.Fatalf("Unable to get authorization url. %v", err)
	}
	fmt.Printf("Go to the following link in your browser then type the "+
		"authorization code: \n%v\n", authorizationURL.String())

	var code string
	if _, err := fmt.Scan(&code); err != nil {
		log.Fatalf("Unable to read authorization code. %v", err)
	}

	accessToken, accessSecret, err := config.AccessToken(requestToken, requestSecret, code)
	if err != nil {
		log.Fatalf("Unable to get access token. %v", err)
	}
	return oauth1.NewToken(accessToken, accessSecret)
}

func jiraTokenCacheFile() (string, error) {
	_, err := user.Current()
	if err != nil {
		return "", err
	}
	tokenCacheDir := filepath.Join("./resources", ".credentials")
	os.MkdirAll(tokenCacheDir, 0700)
	return filepath.Join(tokenCacheDir,
		url.QueryEscape(jiraUrl.Host+".json")), err
}

func jiraTokenFromFile(file string) (*oauth1.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	t := &oauth1.Token{}
	err = json.NewDecoder(f).Decode(t)
	defer f.Close()
	return t, err
}

func saveJIRAToken(file string, token *oauth1.Token) {
	fmt.Printf("Saving credential file to: %s\n", file)
	f, err := os.OpenFile(file, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}

func GetJIRAClient() *http.Client {
	ctx := context.Background()
	keyDERBlock, _ := pem.Decode([]byte(jiraPrivateKey))
	if keyDERBlock == nil {
		log.Fatal("unable to decode key PEM block")
	}
	if !(keyDERBlock.Type == "PRIVATE KEY" || strings.HasSuffix(keyDERBlock.Type, " PRIVATE KEY")) {
		log.Fatalf("unexpected key DER block type: %s", keyDERBlock.Type)
	}

	privateKey, err := x509.ParsePKCS1PrivateKey(keyDERBlock.Bytes)
	if err != nil {
		log.Fatalf("unable to parse PKCS1 private key. %v", err)
	}
	config := oauth1.Config{
		ConsumerKey: jiraConsumerKey,
		CallbackURL: "oob", /* for command line usage */
		Endpoint: oauth1.Endpoint{
			RequestTokenURL: jiraUrl.Host + "plugins/servlet/oauth/request-token",
			AuthorizeURL:    jiraUrl.Host + "plugins/servlet/oauth/authorize",
			AccessTokenURL:  jiraUrl.Host + "plugins/servlet/oauth/access-token",
		},
		Signer: &oauth1.RSASigner{
			PrivateKey: privateKey,
		},
	}
	jiraClient := getJIRAHTTPClient(ctx, &config)
	if err != nil {
		log.Fatalf("unable to create new JIRA client. %v", err)
	}
	return jiraClient
}
