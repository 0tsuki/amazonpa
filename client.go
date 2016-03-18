package amazonpa

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type Client struct {
	AccessKeyId  string
	SecretKey    string
	AssociateTag string
	Host         string
}

const Path = "/onca/xml"

func (c *Client) ItemLookup(itemId string, idType string, responsedGroup []string) (string, error) {
	v := url.Values{}
	v.Set("Service", "AWSECommerceService")
	v.Set("AWSAccessKeyId", c.AccessKeyId)
	v.Set("AssociateTag", c.AssociateTag)
	v.Set("Timestamp", time.Now().UTC().Format("2006-01-02T15:04:05Z"))
	v.Set("Operation", "ItemLookup")
	v.Set("ResponseGroup", strings.Join(responsedGroup, ","))
	v.Set("ItemId", itemId)
	v.Set("IdType", idType)
	strToSign := fmt.Sprintf("GET\n%s\n%s\n%s", c.Host, Path, v.Encode())
	h := hmac.New(sha256.New, []byte(c.SecretKey))
	h.Write([]byte(strToSign))

	signature := base64.StdEncoding.EncodeToString(h.Sum(nil))
	v.Set("Signature", signature)

	reqUrl := url.URL{
		Scheme:   "http",
		Host:     c.Host,
		Path:     Path,
		RawQuery: v.Encode(),
	}

	resp, err := http.Get(reqUrl.String())
	log.Println(reqUrl.String())
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}
