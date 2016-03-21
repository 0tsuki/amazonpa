package amazonpa

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/xml"
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

type ItemLookupResponse struct {
	XMLName xml.Name `xml:"ItemLookupResponse"`
	Items   Items
}

type Items struct {
	Request Request
	Item    Item
}

type Request struct {
	IsValid           bool
	ItemLookupRequest ItemLookupRequest
}

type ItemLookupRequest struct {
	IdType         string
	ItemId         string
	ResponseGroups []string `xml:"ResponseGroup"`
	VariationPage  string
}

type Item struct {
	ASIN           string
	DetailPageURL  string
	ItemLinks      []ItemLink `xml:"ItemLinks>ItemLink"`
	ItemAttributes ItemAttributes
	OfferSummary   OfferSummary
	Offers         Offers
	SalesRank      int
	BrowseNodes    BrowseNodes
}

type BrowseNodes struct {
	BrowseNode []BrowseNode
}

type BrowseNode struct {
	BrowseNodeId   int
	Name           string
	IsCategoryRoot bool
	Ancestors      Ancestors
}

type Ancestors struct {
	BrowseNode *BrowseNode
}

type ItemLink struct {
	Description string
	URL         string
}

type ItemAttributes struct {
	Binding           string
	Brand             string
	Color             string
	EAN               string
	EANList           []string `xml:"EANList>EANListElement"`
	Feature           []string
	ItemDimensions    Dimensions
	Label             string
	ListPrice         Price
	Manufacturer      string
	Model             string
	MPN               string
	PackageDimensions Dimensions
	PartNumber        string
	ProductGroup      string
	ProductTypeName   string
	Publisher         string
	Studio            string
	Title             string
	UPC               string
	UPCList           []string `xml:"UPCList>UPCListElement"`
}

type Dimensions struct {
	Height int
	Length int
	Weight int
	Width  int
}

type OfferSummary struct {
	LowestNewPrice   Price
	LowestUsedPrice  Price
	TotalNew         int
	TotalUsed        int
	TotalCollectible int
	TotalRefurbished int
}

type Price struct {
	Amount         int
	CurrencyCode   string
	FormattedPrice string
}

type Offers struct {
	TotalOffers     int
	TotalOfferPages int
	MoreOffersUrl   string
	Offer           Offer
}

type Offer struct {
	OfferAttributes OfferAttributes
	OfferListing    OfferListing
}

type OfferAttributes struct {
	Condition string
}

type OfferListing struct {
	AvailabilityAttributes          AvailabilityAttributes
	OfferListingId                  string
	Price                           Price
	AmountSaved                     Price
	PercentageSaved                 int
	IsEligibleForSuperSaverShipping int
	IsEligibleForPrime              int
}

type AvailabilityAttributes struct {
	AvailabilityType string
	MinimumHours     int
	MaximumHours     int
}

const Path = "/onca/xml"

func (c *Client) ItemLookup(itemId string, idType string, responsedGroup []string) (*ItemLookupResponse, error) {
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

	var ires ItemLookupResponse
	resp, err := http.Get(reqUrl.String())
	log.Println(reqUrl.String())
	if err != nil {
		return &ires, err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return &ires, err
	}

	xml.Unmarshal(body, &ires)

	return &ires, nil
}
