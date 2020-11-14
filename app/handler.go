// Copyright 2020 The Lunar.Music.Web AUTHORS. All rights reserved.
//
// Use of this source code is governed by a license that can be found in the LICENSE file.

package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

// init is invoked before main()
func init() {
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
}

func main() {
	customHandlerPort, exists := os.LookupEnv("FUNCTIONS_CUSTOMHANDLER_PORT")
	if !exists {
		customHandlerPort = "8080"
	}

	router := gin.Default()
	router.LoadHTMLGlob("templates/*.tmpl.html")

	router.GET("/api/HttpTrigger/login-url", getLoginURL)

	router.GET("/api/HttpTrigger/auth/logout", showLogoutPage)

	router.GET("/api/HttpTrigger/auth/callback", showLoginCallbackPage)

	router.GET("/api/HttpTrigger/", showMusicListPage)

	router.GET("/api/HttpTrigger/send-command-to-raspberrypi", sendCommandToRaspberryPi)

	router.Run(":" + customHandlerPort)
}
