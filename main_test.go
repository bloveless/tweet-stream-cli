package main

import (
	"testing"
)

func TestGetSigningString(t *testing.T) {
	a := auth{
		oauthConsumerKey:       "xvz1evFS4wEEPTGEFPHBog",
		oauthConsumerSecret:    "kAcSOqF21Fu85e7zjz7ZN2U4ZRhfV3WpwPAoE3Z7kBw",
		oauthAccessToken:       "370773112-GmHxMAgYyLbNEtIKZeRNFsMKPR9EyMZeS9weJAEb",
		oauthAccessTokenSecret: "LswwdoUaIvS8ltyTt5jkRh4J50vUPVVHtR2YPi5kE",
	}

	sr := signatureRequest{
		method:    "POST",
		uri:       "https://api.twitter.com/1.1/statuses/update.json?include_entities=true",
		nonce:     "kYjzVBB8Y0ZFabxSWbWovY3uYSQ2pTgmZeNu2VS4cg",
		timestamp: "1318622958",
		params: map[string]string{
			"status":           "Hello Ladies + Gentlemen, a signed OAuth request!",
			"include_entities": "true",
		},
	}

	s := a.signature(sr)
	expected := "hCtSmYh+iHYCEqBWrE7C7hYmtUk="
	if s != expected {
		t.Fatalf("%s != %s", s, expected)
	}
}

func TestGetSigningStringNoOAuthToken(t *testing.T) {
	a := auth{
		oauthConsumerKey:    "xvz1evFS4wEEPTGEFPHBog",
		oauthConsumerSecret: "kAcSOqF21Fu85e7zjz7ZN2U4ZRhfV3WpwPAoE3Z7kBw",
	}

	sr := signatureRequest{
		method:    "POST",
		uri:       "https://api.twitter.com/1.1/statuses/update.json?include_entities=true",
		nonce:     "kYjzVBB8Y0ZFabxSWbWovY3uYSQ2pTgmZeNu2VS4cg",
		timestamp: "1318622958",
		params: map[string]string{
			"status":           "Hello Ladies + Gentlemen, a signed OAuth request!",
			"include_entities": "true",
		},
	}

	s := a.signature(sr)
	expected := "SeAnFOsJg0uDVE8Coxfv5QdLNII="
	if s != expected {
		t.Fatalf("%s != %s", s, expected)
	}
}

func TestGetOauthAuthorizationHeader(t *testing.T) {
	a := auth{
		oauthConsumerKey:       "xvz1evFS4wEEPTGEFPHBog",
		oauthConsumerSecret:    "kAcSOqF21Fu85e7zjz7ZN2U4ZRhfV3WpwPAoE3Z7kBw",
		oauthAccessToken:       "370773112-GmHxMAgYyLbNEtIKZeRNFsMKPR9EyMZeS9weJAEb",
		oauthAccessTokenSecret: "LswwdoUaIvS8ltyTt5jkRh4J50vUPVVHtR2YPi5kE",
	}

	hp := headerParameters{
		oauthNonce:     "kYjzVBB8Y0ZFabxSWbWovY3uYSQ2pTgmZeNu2VS4cg",
		oauthSignature: "tnnArxj06cWHq44gCs1OSKk/jLY=",
		oauthTimestamp: "1318622958",
	}

	ah := a.getOauthAuthorizationHeader(hp)
	expected := "OAuth oauth_consumer_key=\"xvz1evFS4wEEPTGEFPHBog\", oauth_nonce=\"kYjzVBB8Y0ZFabxSWbWovY3uYSQ2pTgmZeNu2VS4cg\", oauth_signature=\"tnnArxj06cWHq44gCs1OSKk%2FjLY%3D\", oauth_signature_method=\"HMAC-SHA1\", oauth_timestamp=\"1318622958\", oauth_token=\"370773112-GmHxMAgYyLbNEtIKZeRNFsMKPR9EyMZeS9weJAEb\", oauth_version=\"1.0\""

	if ah != expected {
		t.Fatalf("%s != %s", ah, expected)
	}
}

func TestGetOauthAuthorizationHeaderNoOAuthToken(t *testing.T) {
	a := auth{
		oauthConsumerKey:    "xvz1evFS4wEEPTGEFPHBog",
		oauthConsumerSecret: "kAcSOqF21Fu85e7zjz7ZN2U4ZRhfV3WpwPAoE3Z7kBw",
	}

	hp := headerParameters{
		oauthNonce:     "kYjzVBB8Y0ZFabxSWbWovY3uYSQ2pTgmZeNu2VS4cg",
		oauthSignature: "tnnArxj06cWHq44gCs1OSKk/jLY=",
		oauthTimestamp: "1318622958",
	}

	ah := a.getOauthAuthorizationHeader(hp)
	expected := "OAuth oauth_consumer_key=\"xvz1evFS4wEEPTGEFPHBog\", oauth_nonce=\"kYjzVBB8Y0ZFabxSWbWovY3uYSQ2pTgmZeNu2VS4cg\", oauth_signature=\"tnnArxj06cWHq44gCs1OSKk%2FjLY%3D\", oauth_signature_method=\"HMAC-SHA1\", oauth_timestamp=\"1318622958\", oauth_version=\"1.0\""

	if ah != expected {
		t.Fatalf("%s != %s", ah, expected)
	}
}
