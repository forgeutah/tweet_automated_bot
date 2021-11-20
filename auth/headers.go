package auth

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"math/big"
	"net/url"
	"os"
	"strconv"
	"time"
)

// OauthParams are needed for this request header
/*Authorization:
OAuth oauth_consumer_key="xvz1evFS4wEEPTGEFPHBog",
oauth_nonce="kYjzVBB8Y0ZFabxSWbWovY3uYSQ2pTgmZeNu2VS4cg",
oauth_signature="tnnArxj06cWHq44gCs1OSKk%2FjLY%3D",
oauth_signature_method="HMAC-SHA1",
oauth_timestamp="1318622958",
oauth_token="370773112-GmHxMAgYyLbNEtIKZeRNFsMKPR9EyMZeS9weJAEb",
oauth_version="1.0"
*/
type OauthParams struct {
	oauthConsumerKey     string
	oauthConsumerSecret  string
	oauthNonce           string
	oauthSignature       string
	oauthSignatureMethod string
	oauthTimestamp       string
	oauthToken           string
	oauthTokenSecret     string
	oauthVersion         string
}

// maybe change to twitter auth

func GetTwitterOauthHeader(tweetUrl string, tweetPayload string) (string, error) {
	// oauthNonce needs to be a 32bit base 64encode random alphanumeric string.
	randBytes, err := generateRandomString(32) // 32 = length of []byte
	if err != nil {
		return "", fmt.Errorf("failed to generate random string for nonce| %w", err)
	}
	oauthParams := OauthParams{
		oauthVersion:         "1.0",
		oauthConsumerKey:     os.Getenv("OAUTH_CONSUMER_KEY"),
		oauthConsumerSecret:  os.Getenv("OAUTH_CONSUMER_SECRET"),
		oauthNonce:           base64.StdEncoding.EncodeToString(randBytes),
		oauthSignatureMethod: "HMAC-SHA1",
		oauthTimestamp:       strconv.FormatInt(time.Now().Unix(), 10),
		oauthToken:           os.Getenv("OAUTH_ACCESS_TOKEN"),
		oauthTokenSecret:     os.Getenv("OAUTH_ACCESS_SECRET"),
	}

	paramString := makeParameterString(tweetPayload, oauthParams)
	oauthParams.oauthSignature = oauthParams.getTwitterSignature("POST", tweetUrl, paramString)
	if err != nil {
		return "", fmt.Errorf("failed to create signature: %w", err)
	}

	return makeOauthString(oauthParams), nil
}

func makeOauthString(oauthParams OauthParams) string {
	// format oauth header
	return fmt.Sprintf("OAuth oauth_consumer_key=%s, oauth_nonce=%s, oauth_signature=%s,oauth_signature_method=%s, oauth_timestamp=%s, oauth_token=%s, oauth_version=%s", oauthParams.oauthConsumerKey, oauthParams.oauthNonce, oauthParams.oauthSignature, oauthParams.oauthSignatureMethod, oauthParams.oauthTimestamp, oauthParams.oauthToken, oauthParams.oauthVersion)
}

func generateRandomString(n int) ([]byte, error) {
	const letters = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
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

// generate signature
func (oauth *OauthParams) getTwitterSignature(httpsMethod, url, paramstring string) string {
	signatureBase := makeBaseString(httpsMethod, url, paramstring)
	signingKey := oauth.generateSigningKey()
	mac := hmac.New(sha1.New, []byte(signingKey))
	mac.Write([]byte(signatureBase))
	signature := base64.StdEncoding.EncodeToString(mac.Sum(nil))
	return signature
}

func makeParameterString(requestPayload string, oauthParams OauthParams) string {
	params := url.Values{}
	params.Add("include_entities", "true")
	params.Add("oauth_consumer_key", oauthParams.oauthConsumerKey)
	params.Add("oauth_nonce", oauthParams.oauthNonce)
	params.Add("oauth_signature_method", oauthParams.oauthSignatureMethod)
	params.Add("oauth_timestamp", oauthParams.oauthTimestamp)
	params.Add("oauth_token", oauthParams.oauthToken)
	params.Add("oauth_version", oauthParams.oauthVersion)
	paramEncode := params.Encode()
	parameterStr := fmt.Sprintf("%s&status=%s", paramEncode, url.PathEscape(requestPayload))
	return parameterStr
}

func makeBaseString(httpsMethod, twitterUrl, paramstring string) string {
	baseString := fmt.Sprintf("%s&%s&%s", httpsMethod, url.QueryEscape(twitterUrl), url.QueryEscape(paramstring))
	return baseString
}

func (oauth *OauthParams) generateSigningKey() []byte {
	signingKey := fmt.Sprintf("%s&%s", oauth.oauthConsumerSecret, oauth.oauthTokenSecret)
	return []byte(signingKey)
}
