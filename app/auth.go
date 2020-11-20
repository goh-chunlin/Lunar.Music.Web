// Copyright 2020 The Lunar.Music.Web AUTHORS. All rights reserved.
//
// Use of this source code is governed by a license that can be found in the LICENSE file.

package main

import (
	"crypto/rand"
	"encoding/base64"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/securecookie"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/microsoft"
)

const ACCESS_AND_REFRESH_TOKENS_COOKIE_NAME = "Tokens"

var AzureADCallbackURL = os.Getenv("AZURE_AD_CALLBACK_URL")
var AzureADClientID = os.Getenv("AZURE_AD_CLIENT_ID")
var AzureADClientSecret = os.Getenv("AZURE_AD_CLIENT_SECRET")

var myClient = &http.Client{Timeout: 10 * time.Second}

var oauthConfig = &oauth2.Config{
	RedirectURL:  AzureADCallbackURL,
	ClientID:     AzureADClientID,
	ClientSecret: AzureADClientSecret,
	Scopes:       []string{"files.read offline_access"},
	Endpoint:     microsoft.AzureADEndpoint("common"),
}

var s = securecookie.New([]byte(os.Getenv("SECURECOOKIE_HASH_KEY")), []byte(os.Getenv("SECURECOOKIE_BLOCK_KEY")))

// tokenJSON is the struct representing the HTTP response from OAuth2
// providers, which is Azure AD in this case, returning a token in JSON form.
type tokenJSON struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
	Resource     string `json:"resource"`
	Scope        string `json:"scope"`
}

func (token *tokenJSON) expiry() (t time.Time) {
	return time.Now().Add(time.Duration(token.ExpiresIn) * time.Second)
}

func generateStateOauthCookie(w http.ResponseWriter) string {
	expiration := time.Now().Add(365 * 24 * time.Hour)

	b := make([]byte, 16)
	rand.Read(b)
	state := base64.URLEncoding.EncodeToString(b)
	cookie := http.Cookie{Name: "oauthstate", Value: state, Expires: expiration}
	http.SetCookie(w, &cookie)

	return state
}

func setTokenCookie(context *gin.Context, tokenJSON *tokenJSON) {

	var token = &oauth2.Token{
		AccessToken:  tokenJSON.AccessToken,
		TokenType:    tokenJSON.TokenType,
		RefreshToken: tokenJSON.RefreshToken,
		Expiry:       tokenJSON.expiry(),
	}

	encoded, err := s.Encode(ACCESS_AND_REFRESH_TOKENS_COOKIE_NAME, token)
	if err == nil {
		cookie := &http.Cookie{
			Name:     ACCESS_AND_REFRESH_TOKENS_COOKIE_NAME,
			Value:    encoded,
			Path:     "/",
			Secure:   true,
			HttpOnly: true,
			SameSite: http.SameSiteStrictMode,
		}
		http.SetCookie(context.Writer, cookie)
	}
}

func removeTokenCookie(context *gin.Context) {
	cookie := &http.Cookie{
		Name:     ACCESS_AND_REFRESH_TOKENS_COOKIE_NAME,
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	}
	http.SetCookie(context.Writer, cookie)
}

func getAccessAndRefreshTokenFromCookie(context *gin.Context) *oauth2.Token {
	cookie, err := context.Request.Cookie(ACCESS_AND_REFRESH_TOKENS_COOKIE_NAME)
	if err == nil {
		var tokensCookieValue *oauth2.Token
		err := s.Decode(ACCESS_AND_REFRESH_TOKENS_COOKIE_NAME, cookie.Value, &tokensCookieValue)
		if err == nil {
			return tokensCookieValue
		}
	}

	return nil
}
