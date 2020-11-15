// Copyright 2020 The Lunar.Music.Web AUTHORS. All rights reserved.
//
// Use of this source code is governed by a license that can be found in the LICENSE file.

package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/streadway/amqp"
)

type Command struct {
	Tasks []*Task `json:"tasks"`
}

type Task struct {
	Name    string   `json:"name"`
	Content []string `json:"content"`
}

var rabbitMQServerConnectionString = os.Getenv("RABBITMQ_SERVER_CONNECTION_STRING")
var rabbitMQChannelName = os.Getenv("RABBITMQ_CHANNEL_NAME")
var rabbitMQAllowedMicrosoftUserId = os.Getenv("RABBITMQ_ALLOWED_MICROSOFT_USER_ID")

func sendCommandToRaspberryPi(context *gin.Context) {
	requestBody, err := context.GetRawData()
	if err != nil {
		context.AbortWithStatusJSON(400, gin.H{
			"message": "The request body is invalid.",
		})

		return
	}

	currentUserId, err := getOneDriveOwnerUserId(context)
	if err != nil {
		context.AbortWithStatusJSON(500, gin.H{
			"message": fmt.Sprintf("Error in getting the user id: %v.", err.Error()),
		})

		return
	}

	if currentUserId != rabbitMQAllowedMicrosoftUserId {
		context.AbortWithStatusJSON(403, gin.H{
			"message": "The current logged in user is not allowed to send command to the Raspberry Pi.",
		})

		return
	}

	var command *Command
	json.Unmarshal(requestBody, &command)

	for _, task := range command.Tasks {
		if task.Name == "play-single" {
			itemId := task.Content[0]
			itemDownloadUrl, err := getOneDriveItemDownloadUrl(context, itemId)
			if err != nil {
				context.AbortWithStatusJSON(500, gin.H{
					"message": fmt.Sprintf("The Download URL for item %v is invalid. Command will not be sent.", itemId),
				})

				return
			}

			task.Content = append(task.Content, itemDownloadUrl)
		}
	}

	messageSent := sendCommandToRabbitMQServer(command)

	context.JSON(200, gin.H{
		"message": string(messageSent),
	})
}

func sendCommandToRabbitMQServer(command *Command) []byte {
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

	body, err := json.Marshal(command)
	failOnError(err, "Failed to convert the command JSON to a byte slice")

	err = ch.Publish(
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		})
	failOnError(err, "Failed to publish a message")

	return body
}
