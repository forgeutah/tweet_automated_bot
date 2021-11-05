package main

import (
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/big"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
)

// Tweet is twitter payload for tweets endpoint.
type Tweet struct {
	Text string `json:"text"`
}

func main() {

	// define http client
	client := http.Client{
		Timeout: 10 * time.Second,
	}

	newTweet := Tweet{
		Text: "We love #golang",
	}

	tweetPayload, err := json.Marshal(&newTweet)
	if err != nil {
		panic(err)
	}

	tweetUrl := "https://api.twitter.com/2/tweets"
	req, err := http.NewRequest("POST", tweetUrl, bytes.NewReader(tweetPayload))
	if err != nil {
		panic(err)
	}

	// example authorization
	/*Authorization:
	OAuth oauth_consumer_key="xvz1evFS4wEEPTGEFPHBog",
	oauth_nonce="kYjzVBB8Y0ZFabxSWbWovY3uYSQ2pTgmZeNu2VS4cg",
	oauth_signature="tnnArxj06cWHq44gCs1OSKk%2FjLY%3D",
	oauth_signature_method="HMAC-SHA1",
	oauth_timestamp="1318622958",
	oauth_token="370773112-GmHxMAgYyLbNEtIKZeRNFsMKPR9EyMZeS9weJAEb",
	oauth_version="1.0"
	*/

	var oauthConsumerKey, oauthNonce, oauthSignature, oauthSignatureMethod, oauthTimestamp, oauthToken, oauthVersion string

	// require for v2 tweets endpoint
	oauthVersion = "1.0"
	oauthConsumerKey = os.Getenv("OAUTH_CONSUMER_KEY")

	// oauthNonce needs to be a 32bit base 64encode random alphanumeric string.
	randBytes, err := generateRandomString(32) // 32 = length of []byte
	if err != nil {
		panic(err)
	}
	oauthNonce = base64.StdEncoding.EncodeToString(randBytes)

	oauthSignatureMethod = "HMAC-SHA1"

	oauthTimestamp = strconv.FormatInt(time.Now().Unix(), 10)

	oauthToken = os.Getenv("OAUTH_ACCESS_TOKEN")

	// generate signature
	oauthSignature, err = getTwitterSignature("POST", tweetUrl, tweetPayload, oauthConsumerKey, oauthNonce, oauthSignatureMethod, oauthTimestamp, oauthToken, oauthVersion)
	if err != nil {
		panic(err)
	}
	// format oauth header
	oauthString := fmt.Sprintf("OAuth oauth_consumer_key=%s, oauth_nonce=%s, oauth_signature=%s,oauth_signature_method=%s, oauth_timestamp=%s, oauth_token=%s, oauth_version=%s", oauthConsumerKey, oauthNonce, oauthSignature, oauthSignatureMethod, oauthTimestamp, oauthToken, oauthVersion)

	req.Header.Set("Authorization", oauthString)

	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(body))
}

func generateRandomString(n int) ([]byte, error) {
	const letters = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz-"
	ret := make([]byte, n)
	for i := 0; i < n; i++ {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(letters))))
		if err != nil {
			return ret, err
		}
		ret[i] = letters[num.Int64()]
	}

	return ret, nil
}

func getTwitterSignature(httpMethod, twitterURL string, requestPayload []byte, oauthParams ...string) (string, error) {
	signatureStr := url.Values{}
	signatureStr.Add(httpMethod, twitterURL)
	signatureStr.Add(string(requestPayload), strings.Join(oauthParams, ","))
	signature := signatureStr.Encode()

	return signature, nil
}
