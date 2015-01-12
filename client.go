package goamazonpaapi

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

const (
	defaultBaseURL = "http://ecs.amazonaws.com/onca/xml"
	defaultService = "AWSECommerceService"
	defaultVersion = "2013-08-01"
)

type Client struct {
	BaseURL *url.URL
	Service string
	Version string
	Params  *url.Values
	App     *App
}

type App struct {
	AccessKey       string
	SecretAccessKey string
	AssociateTag    string
}

func NewClient(accessKey, secretAccessKey, associateTag string) *Client {
	baseURL, _ := url.Parse(defaultBaseURL)

	return &Client{
		BaseURL: baseURL,
		Service: defaultService,
		Version: defaultVersion,
		App: &App{
			AccessKey:       accessKey,
			SecretAccessKey: secretAccessKey,
			AssociateTag:    associateTag,
		},
	}
}

func (c *Client) buildParams(operation string, params map[string]string) error {
	values := &url.Values{}
	values.Set("Service", c.Service)
	values.Set("Version", c.Version)
	values.Set("AssociateTag", c.App.AssociateTag)
	values.Set("Operation", operation)
	values.Set("Timestamp", time.Now().UTC().Format(time.RFC3339))
	values.Set("AWSAccessKeyId", c.App.AccessKey)
	for key, val := range params {
		values.Set(key, val)
	}

	c.Params = values

	return c.buildSignature()
}

func (c *Client) buildSignature() error {
	hasher := hmac.New(sha256.New, []byte(c.App.SecretAccessKey))

	strToSign := fmt.Sprintf("GET\n%s\n%s\n%s",
		c.BaseURL.Host,
		c.BaseURL.Path,
		c.Params.Encode(),
	)
	if _, err := hasher.Write([]byte(strToSign)); err != nil {
		return err
	}

	signature := base64.StdEncoding.EncodeToString(hasher.Sum(nil))
	c.Params.Set("Signature", signature)

	return nil
}

func (c *Client) do() ([]byte, error) {
	u := c.BaseURL
	u.RawQuery = c.Params.Encode()

	res, err := http.Get(u.String())
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	return ioutil.ReadAll(res.Body)
}

func (c *Client) ItemLookup(itemId string) ([]byte, error) {
	params := map[string]string{
		"ItemId": itemId,
	}
	if err := c.buildParams("ItemLookup", params); err != nil {
		return nil, err
	}

	return c.do()
}

func (c *Client) ItemSearchByKeyword(searchIndex, keyword, responseGroup string) ([]byte, error) {
	params := map[string]string{
		"SearchIndex":   searchIndex,
		"Keywords":      keyword,
		"ResponseGroup": responseGroup,
	}
	if err := c.buildParams("ItemSearch", params); err != nil {
		return nil, err
	}

	return c.do()
}
