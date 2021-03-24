## oauth Example

> 用于请求oauthapi和微服务api方式

```go
config := apiclient.Config{
    ClientID:     "YOUR_CLIENT_ID",
    ClientSecret: "YOUR_CLIENT_SECRET",
    AuthURL:      "https://provider.com/o/oauth2/auth",
    TokenURL:     "https://provider.com/o/oauth2/token",
    Scopes:       []string{"email", "avatar"},
    Mode:         3,
}

// create a client
client := apiclient.NewClient(http.DefaultClient, config)

// url to fetch the code
url := client.AuthCodeURL("state")
fmt.Printf("Visit the URL for the auth dialog: %v", url)

// Use the authorization code that is pushed to the redirect URL.
// Exchange will do the handshake to retrieve the initial access token.
var code string
if _, err := fmt.Scan(&code); err != nil {
    log.Fatal(err)
}

// get a token
token, err := client.Exchange(context.Background(), code)
if err != nil {
    panic(err)
}

var _ string = token.AccessToken  // OAuth2 token
var _ string = token.TokenType    // type of the token
var _ string = token.RefreshToken // token for a refresh
var _ time.Time = token.Expiry    // token expiration time
var _ bool = token.IsExpired()    // have token expired?
```

## MSA Example
```go
timeout := time.Duration(1 * time.Second)

headers := map[string]string{
	"Authorization": "bearer hRYr75hQJLwPIQJJBq56BRrtDssYNyoX",
	"content-type": "application/json",
},
clis := apiclient.NewApiClient(
	"http://localhost/api/",
	headers,
	timeout,
	false, 
)

res,_ := clis.GET("me/profile", nil)
b,_ := ioutil.ReadAll( res.Body)
```