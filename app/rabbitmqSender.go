// Copyright 2020 The Lunar.Music.Web AUTHORS. All rights reserved.
//
// Use of this source code is governed by a license that can be found in the LICENSE file.

package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/streadway/amqp"
)

var rabbitMQServerConnectionString = os.Getenv("RABBITMQ_SERVER_CONNECTION_STRING")
var rabbitMQChannelName = os.Getenv("RABBITMQ_CHANNEL_NAME")

func sendCommandToRaspberryPi(context *gin.Context) {
	r := context.Request

	methodQueryStrings, hasParam := r.URL.Query()["method"]
	if !hasParam {
		context.JSON(200, gin.H{
			"message": "Missing parameter \"method\" in the query string.",
		})

		return
	}

	descriptionQueryStrings, hasParam := r.URL.Query()["description"]
	if !hasParam {
		context.JSON(200, gin.H{
			"message": "Missing parameter \"description\" in the query string.",
		})

		return
	}

	conn, err := amqp.Dial(rabbitMQServerConnectionString)
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		rabbitMQChannelName, // name
		false,               // durable
		false,               // delete when unused
		false,               // exclusive
		false,               // no-wait
		nil,                 // arguments
	)
	failOnError(err, "Failed to declare a queue")

	body := "play-all"
	if methodQueryStrings[0] == "download" {
		body = methodQueryStrings[0] + ":" + descriptionQueryStrings[0]
	}

	err = ch.Publish(
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
		})
	log.Printf(" [x] Sent %s", body)
	failOnError(err, "Failed to publish a message")

	context.JSON(200, gin.H{
		"message": body,
	})
}
