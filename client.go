package tissueabsuer

import (
	"math/rand"
	"net/http"
	"net/url"
)

type Client struct {
	Agents []string     // user agents
	Proxy  string       // http proxy
	client *http.Client // http client
}

// get new Client
func NewClient(proxy string) *Client {
	var c Client
	c.Proxy = proxy
	c.Agents = []string{
		// safari
		"Mozilla/5.0(Macintosh;U;IntelMacOSX10_6_8;en-us)AppleWebKit/534.50(KHTML,likeGecko)Version/5.1Safari/534.50",

		// firefox
		"Mozilla/5.0(Macintosh;IntelMacOSX10.6;rv:2.0.1)Gecko/20100101Firefox/4.0.1",

		// opera
		"Opera/9.80(Macintosh;IntelMacOSX10.6.8;U;en)Presto/2.8.131Version/11.11",

		// maxthon
		"Mozilla/4.0(compatible;MSIE7.0;WindowsNT5.1;Maxthon2.0)",

		// the world
		"Mozilla/4.0(compatible;MSIE7.0;WindowsNT5.1;TheWorld)",

		// iphone safari
		"Mozilla/5.0(iPhone;U;CPUiPhoneOS4_3_3likeMacOSX;en-us)AppleWebKit/533.17.9(KHTML,likeGecko)Version/5.0.2Mobile/8J2Safari/6533.18.5",

		// android
		"Mozilla/5.0(Linux;U;Android2.3.7;en-us;NexusOneBuild/FRF91)AppleWebKit/533.1(KHTML,likeGecko)Version/4.0MobileSafari/533.1",
	}

	c.client = &http.Client{}
	return &c
}

// use proxy
func (c *Client) UseProxy() {
	proxyFunc := func(r *http.Request) (*url.URL, error) {
		return url.Parse(c.Proxy)
	}

	transport := &http.Transport{
		Proxy: proxyFunc,
	}

	// set transport
	c.client.Transport = transport
}

// do http request
func (c *Client) DoRequest(r *http.Request, referer string) (*http.Response, error) {
	// use random user-agent
	randIndex := rand.Intn(len(c.Agents))
	r.Header.Add("User-Agent", c.Agents[randIndex])

	if referer != "" {
		r.Header.Add("Referer", referer)
	}

	return c.client.Do(r)
}
