package cryptobot

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

const (
	testnetAPIURL      = "https://testnet-pay.crypt.bot/api/"
	mainnetAPIURL      = "https://pay.crypt.bot/api/"
	apiTokenHeaderName = "Crypto-Pay-API-Token"
)

// Client for making requests to CryptoBot API methods
type Client struct {
	apiToken    string
	testingMode bool

	httpClient *http.Client
}

type Options struct {
	// API Token of your CryptoBot app (token from CryptoTestnetBot can also be used)
	APIToken string

	// Default false. determines if client will request Testnet url instead of Mainnet url
	Testing bool

	// Optional. Default is 30 seconds
	ClientTimeout time.Duration

	// Optional. Set Proxy URL, Default is nil
	ProxyUrl *url.URL
}

func NewClient(options Options) *Client {
	c := &Client{
		testingMode: options.Testing,
		apiToken:    options.APIToken,
	}

	clientTimeout := time.Second * 30
	if options.ClientTimeout != 0 {
		clientTimeout = options.ClientTimeout
	}
	if options.ProxyUrl != nil {
		transport := &http.Transport{
			Proxy: http.ProxyURL(options.ProxyUrl),
		}
		c.httpClient = &http.Client{
			Timeout:   clientTimeout,
			Transport: transport,
		}
	} else {
		c.httpClient = &http.Client{
			Timeout: clientTimeout,
		}
	}

	return c
}

func (c *Client) getRequestUrl() string {
	if c.testingMode {
		return testnetAPIURL
	} else {
		return mainnetAPIURL
	}
}

func (c *Client) request(path string, queryModifierFunc func(q url.Values) url.Values) (io.ReadCloser, error) {
	req, err := http.NewRequest(http.MethodGet, c.getRequestUrl()+path, nil)
	if err != nil {
		return nil, fmt.Errorf("error while creating a request: %w", err)
	}
	if queryModifierFunc != nil {
		req.URL.RawQuery = queryModifierFunc(req.URL.Query()).Encode()
	}

	req.Header.Set(apiTokenHeaderName, c.apiToken)
	r, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error while performing a request: %w", err)
	}

	return r.Body, nil
}

func (c *Client) decodeResponse(responseBodyReader io.Reader, targetPointer any) error {
	responseBody, err := ioutil.ReadAll(responseBodyReader)
	if err != nil {
		return fmt.Errorf("error while decoding response: %w", err)
	}
	if err := json.Unmarshal(responseBody, targetPointer); err != nil {
		return fmt.Errorf("error while unmarshaling response to the target: %w", err)
	}

	return nil
}
