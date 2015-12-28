package htmltest

import (
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/rtfb/go-html-transform/css/selector"
	"golang.org/x/net/html"
)

var (
	tclient *http.Client
	tserver *httptest.Server
)

func initClient() {
	jar, err := cookiejar.New(nil)
	if err != nil {
		panic(err)
	}
	tclient = &http.Client{
		Jar: jar,
	}
}

func initServer(router http.Handler) {
	tserver = httptest.NewServer(router)
}

func Init(router http.Handler) {
	initClient()
	initServer(router)
}

func Client() *http.Client {
	return tclient
}

func CssSelect(t *testing.T, node *html.Node, query string) []*html.Node {
	chain, err := selector.Selector(query)
	if err != nil {
		t.Fatalf("Error: query=%q, err=%s", query, err.Error())
	}
	return chain.Find(node)
}

func curlParam(url string, method func(string) (*http.Response, error)) string {
	if r, err := method(url); err == nil {
		b, err := ioutil.ReadAll(r.Body)
		r.Body.Close()
		if err == nil {
			return string(b)
		}
		println(err.Error())
	} else {
		println(err.Error())
	}
	return ""
}

func Curl(url string) string {
	return curlParam(url, tclientGet)
}

func CurlPost(url string) string {
	return curlParam(url, tclientPostForm)
}

func PostForm(path string, values *url.Values) (string, error) {
	resp, err := tclient.PostForm(localhostURL(path), *values)
	defer resp.Body.Close()
	if err != nil {
		return "", err
	}
	body, err := ioutil.ReadAll(resp.Body)
	return string(body), err
}

func localhostURL(u string) string {
	if u == "" {
		return tserver.URL
	} else if u[0] == '/' {
		return tserver.URL + u
	} else {
		return tserver.URL + "/" + u
	}
}

func tclientGet(rqURL string) (*http.Response, error) {
	return tclient.Get(localhostURL(rqURL))
}

func tclientPostForm(rqURL string) (*http.Response, error) {
	return tclient.PostForm(localhostURL(rqURL), url.Values{})
}
