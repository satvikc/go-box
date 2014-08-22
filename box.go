package box

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/golang/oauth2"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
)

// Box Client
type Box struct {
	APIURL    string
	UPLOADURL string
	config    *oauth2.Config
	token     *oauth2.Token
}

// NewBox gets the new Box object with appropriate APIURL.
func NewBox() *Box {
	box := &Box{
		APIURL:    "https://api.box.com/2.0",
		UPLOADURL: "https://upload.box.com/api/2.0",
	}
	return box
}

// SetAppInfo adds oauth2 app info
func (box *Box) SetAppInfo(clientid, clientsecret string) error {
	var err error
	box.config, err = oauth2.NewConfig(
		&oauth2.Options{
			ClientID:     clientid,
			ClientSecret: clientsecret,
		},
		"https://app.box.com/api/oauth2/authorize",
		"https://app.box.com/api/oauth2/token")
	return err
}

// SetAccessToken sets access token to avoid calling Auth method.
func (box *Box) SetAccessToken(accesstoken string) {
	box.token = &oauth2.Token{AccessToken: accesstoken}
}

// AccessToken returns the OAuth access token.
func (box *Box) AccessToken() string {
	return box.token.AccessToken
}

// Get the http client for further api accesses.
func (box *Box) client() *http.Client {
	var t *oauth2.Transport
	t = box.config.NewTransport()
	t.SetToken(box.token)
	return &http.Client{Transport: t}
}

// Auth displays the URL to authorize this application to connect to your account.
func (box *Box) Auth() error {
	var code string
	var t *oauth2.Transport
	var err error
	fmt.Printf("Please visit:\n%s\nEnter the code: ",
		box.config.AuthCodeURL(""))
	fmt.Scanln(&code)
	if t, err = box.config.NewTransportWithCode(code); err != nil {
		return err
	}
	box.token = t.Token()
	box.token.TokenType = "Bearer"
	return nil
}

// doRequest performs the request (GET or POST) using authorized http
// client. You can also pass params to encode them in the request url
// or body to place in the request body.
func (box *Box) doRequest(method, path string, params *url.Values, reqBody string) ([]byte, error) {
	var body []byte
	var rawurl string
	var response *http.Response
	var request *http.Request
	var err error
	var reqBodyReader io.Reader

	// If paramerters are nil then dont add `?` to the url
	if params == nil {
		rawurl = fmt.Sprintf("%s/%s", box.APIURL, urlEncode(path))
	} else {
		rawurl = fmt.Sprintf("%s/%s?%s", box.APIURL, urlEncode(path), params.Encode())
	}

	// If reqBody is empty then dont create new reader
	if reqBody != "" {
		reqBodyReader = bytes.NewReader([]byte(reqBody))
	}

	if request, err = http.NewRequest(method, rawurl, reqBodyReader); err != nil {
		return nil, err
	}
	if response, err = box.client().Do(request); err != nil {
		return nil, err
	}
	defer response.Body.Close()
	if body, err = getResponse(response); err != nil {
		return nil, err
	}
	return body, nil
}

func getResponse(r *http.Response) ([]byte, error) {
	var b []byte
	var err error
	if b, err = ioutil.ReadAll(r.Body); err != nil {
		return nil, err
	}
	switch r.StatusCode {
	case http.StatusOK, http.StatusCreated, http.StatusAccepted:
		return b, nil
	}
	return nil, errors.New(fmt.Sprintf("unexpected HTTP status code %d", r.StatusCode))
}

// urlEncode encodes s for url
func urlEncode(s string) string {
	encoded := url.QueryEscape(s)
	return encoded
}
