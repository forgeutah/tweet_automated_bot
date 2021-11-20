package auth

import (
	"reflect"
	"testing"
)

// test to verify that the paramerstring matches the twitter api docs
// https://developer.twitter.com/en/docs/authentication/oauth-1-0a/creating-a-signature
func Test_makeParameterString(t *testing.T) {
	type args struct {
		requestPayload string
		oauthParams    OauthParams
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "make string",
			args: args{
				// httpMethod:     "POST",
				// twitterURL:     "https://api.twitter.com/1.1/statuses/update.json",
				requestPayload: "Hello Ladies and Gentlemen, a signed OAuth request!",
				oauthParams: OauthParams{
					oauthConsumerKey:     "xvz1evFS4wEEPTGEFPHBog",
					oauthNonce:           "kYjzVBB8Y0ZFabxSWbWovY3uYSQ2pTgmZeNu2VS4cg",
					oauthSignatureMethod: "HMAC-SHA1",
					oauthTimestamp:       "1318622958",
					oauthToken:           "370773112-GmHxMAgYyLbNEtIKZeRNFsMKPR9EyMZeS9weJAEb",
					oauthVersion:         "1.0",
				},
			},

			want: "include_entities=true&oauth_consumer_key=xvz1evFS4wEEPTGEFPHBog&oauth_nonce=kYjzVBB8Y0ZFabxSWbWovY3uYSQ2pTgmZeNu2VS4cg&oauth_signature_method=HMAC-SHA1&oauth_timestamp=1318622958&oauth_token=370773112-GmHxMAgYyLbNEtIKZeRNFsMKPR9EyMZeS9weJAEb&oauth_version=1.0&status=Hello%20Ladies%20and%20Gentlemen%2C%20a%20signed%20OAuth%20request%21",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := makeParameterString(tt.args.requestPayload, tt.args.oauthParams); got != tt.want {
				t.Errorf("generateSignatureBaseString() = \n%v\n, want \n%v\n", got, tt.want)
			}
		})
	}
}

func Test_makeBaseString(t *testing.T) {
	type args struct {
		httpsMethod string
		url         string
		paramstring string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "success: match twitter",
			args: args{
				httpsMethod: "POST",
				url:         "https://api.twitter.com/1.1/statuses/update.json",
				paramstring: "include_entities=true&oauth_consumer_key=xvz1evFS4wEEPTGEFPHBog&oauth_nonce=kYjzVBB8Y0ZFabxSWbWovY3uYSQ2pTgmZeNu2VS4cg&oauth_signature_method=HMAC-SHA1&oauth_timestamp=1318622958&oauth_token=370773112-GmHxMAgYyLbNEtIKZeRNFsMKPR9EyMZeS9weJAEb&oauth_version=1.0&status=Hello%20Ladies%20and%20Gentlemen%2C%20a%20signed%20OAuth%20request%21",
			},
			want: "POST&https%3A%2F%2Fapi.twitter.com%2F1.1%2Fstatuses%2Fupdate.json&include_entities%3Dtrue%26oauth_consumer_key%3Dxvz1evFS4wEEPTGEFPHBog%26oauth_nonce%3DkYjzVBB8Y0ZFabxSWbWovY3uYSQ2pTgmZeNu2VS4cg%26oauth_signature_method%3DHMAC-SHA1%26oauth_timestamp%3D1318622958%26oauth_token%3D370773112-GmHxMAgYyLbNEtIKZeRNFsMKPR9EyMZeS9weJAEb%26oauth_version%3D1.0%26status%3DHello%2520Ladies%2520and%2520Gentlemen%252C%2520a%2520signed%2520OAuth%2520request%2521",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := makeBaseString(tt.args.httpsMethod, tt.args.url, tt.args.paramstring); got != tt.want {
				t.Errorf("makeBaseString() = \n%v\n, want \n%v\n", got, tt.want)
			}
		})
	}
}

func TestOauthParams_generateSigningKey(t *testing.T) {
	tests := []struct {
		name  string
		oauth *OauthParams
		want  []byte
	}{
		{
			name: "success match twitter",
			oauth: &OauthParams{
				oauthConsumerSecret: "kAcSOqF21Fu85e7zjz7ZN2U4ZRhfV3WpwPAoE3Z7kBw",
				oauthTokenSecret:    "LswwdoUaIvS8ltyTt5jkRh4J50vUPVVHtR2YPi5kE",
			},
			want: []byte("kAcSOqF21Fu85e7zjz7ZN2U4ZRhfV3WpwPAoE3Z7kBw&LswwdoUaIvS8ltyTt5jkRh4J50vUPVVHtR2YPi5kE"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.oauth.generateSigningKey(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("OauthParams.generateSigningKey() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestOauthParams_getTwitterSignature(t *testing.T) {
	type args struct {
		httpsMethod string
		url         string
		paramstring string
	}
	tests := []struct {
		name  string
		oauth *OauthParams
		args  args
		want  string
	}{
		{
			name: "success: match twitter",
			oauth: &OauthParams{
				oauthConsumerSecret: "kAcSOqF21Fu85e7zjz7ZN2U4ZRhfV3WpwPAoE3Z7kBw",
				oauthTokenSecret:    "LswwdoUaIvS8ltyTt5jkRh4J50vUPVVHtR2YPi5kE",
			},
			args: args{
				httpsMethod: "POST",
				url:         "https://api.twitter.com/1.1/statuses/update.json",
				paramstring: "include_entities=true&oauth_consumer_key=xvz1evFS4wEEPTGEFPHBog&oauth_nonce=kYjzVBB8Y0ZFabxSWbWovY3uYSQ2pTgmZeNu2VS4cg&oauth_signature_method=HMAC-SHA1&oauth_timestamp=1318622958&oauth_token=370773112-GmHxMAgYyLbNEtIKZeRNFsMKPR9EyMZeS9weJAEb&oauth_version=1.0&status=Hello%20Ladies%20and%20Gentlemen%2C%20a%20signed%20OAuth%20request%21",
			},
			want: "hCtSmYh+iHYCEqBWrE7C7hYmtUk=",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.oauth.getTwitterSignature(tt.args.httpsMethod, tt.args.url, tt.args.paramstring); got != tt.want {
				t.Errorf("OauthParams.getTwitterSignature() = %v, want %v", got, tt.want)
			}
		})
	}
}
