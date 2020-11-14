// Copyright 2020 The Lunar.Music.Web AUTHORS. All rights reserved.
//
// Use of this source code is governed by a license that can be found in the LICENSE file.

package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/microsoft"
)

func getLoginURL(context *gin.Context) {
	w := context.Writer

	// Create oauthState cookie
	oauthState := generateStateOauthCookie(w)
	loginURL := oauthConfig.AuthCodeURL(oauthState)

	context.JSON(http.StatusOK, gin.H{
		"url": loginURL,
	})
}

func showLoginCallbackPage(context *gin.Context) {
	w := context.Writer
	r := context.Request

	code := r.FormValue("code")

	response, err := http.PostForm(microsoft.AzureADEndpoint("common").TokenURL, url.Values{
		"client_id":     {AzureADClientID},
		"redirect_uri":  {AzureADCallbackURL},
		"client_secret": {AzureADClientSecret},
		"code":          {code},
		"grant_type":    {"authorization_code"}})

	if err != nil {
		fmt.Print(err)
	}

	defer response.Body.Close()
	var token *oauth2.Token

	// Try to decode the request body into the struct. If there is an error,
	// respond to the client with the error message and a 400 status code.
	err = json.NewDecoder(response.Body).Decode(&token)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	setTokenCookie(context, token)

	context.HTML(http.StatusOK, "login-callback.tmpl.html", nil)
}

func showLogoutPage(context *gin.Context) {
	removeTokenCookie(context)

	context.HTML(http.StatusOK, "logout.tmpl.html", nil)
}
