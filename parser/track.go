package parser

import (
	"github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
	"strings"
	"telegram-parser/db"
	"telegram-parser/mq"
	"time"
)

// StartTrackingStatistics launches several handlers for updates from telegram.
// @param handlersCount is the number of message handlers that will be run. The minimum value is 10.
func (a *App) StartTrackingStatistics(handlersCount int) {
	if handlersCount < 10 {
		handlersCount = 10
	}

	// updates the maximum value of the message date
	now := time.Now().Unix()
	maxMsgTimeUpdater(&now)

	queue := a.MqCli.Consume()

	for i := 0; i < handlersCount; i++ {
		go a.trackingStatistics(queue, &now)
	}

	logrus.Infof("Launched %v goroutines to track statistics.", handlersCount)
}

// trackingStatistics starts continuous tracking of message statistics.
func (a *App) trackingStatistics(rabbitQueue <-chan amqp.Delivery, stopTime *int64) {
	mqCli := a.MqCli
	dbCli := a.DbCli
	tgCli := a.Telegram.Client

	for update := range rabbitQueue {
		var previousMsg *db.Message // last known state message
		var updatedMsg *db.Message  // updated message information

		previousMsg, err := mq.UnmarshalRabbitBody(update.Body)
		if err != nil {
			logrus.Errorf("Failed to parse the JSON-encoded data in *db.Message with err '%v'.", err.Error())
			continue
		}

		if isMsgExpired(previousMsg, stopTime) {
			// message expired so we remove the message from the parser queue
			update.Ack(false)
			continue
		}

		// GET CHAT INFORMATION
		chat, err := tgCli.GetChat(previousMsg.ChatID)
		if err != nil {
			logrus.Errorf("Failed to receive the message from telegram with CHAT TITLE = '%v', LINK = '%v'  message id = '%v' and chat id = '%v' with error = `%#v`. Message passed to message queue.", previousMsg.ChatTitle, previousMsg.Link, previousMsg.MessageID, previousMsg.ChatID, err.Error())

			// try to join chat
			_, err := tgCli.JoinChat(previousMsg.ChatID)
			if err != nil {
				// Failed to join channel. Remove it from messages queue
				logrus.Errorf("Failed to join channel '%v' with error = '%v'.", previousMsg.ChatID, err.Error())
				update.Ack(false)
				continue
			}

			update.Ack(false)
			publishMessage(mqCli, previousMsg)
			continue
		}

		// GET MESSAGE INFORMATION
		tdlibMsg, err := tgCli.GetMessage(previousMsg.ChatID, previousMsg.MessageID)
		if err != nil {

			if strings.Contains(err.Error(), "msg: Not Found") {
				logrus.Errorf("Message with id = '%v' and chat id = '%v' not found. Message will be deleted from queue.", previousMsg.MessageID, previousMsg.ChatID)
				//	TODO возможно стоит удалить и из базы
				update.Ack(false)

				continue
			}

			if strings.Contains(err.Error(), "msg: CHANNEL_PRIVATE") {
				logrus.Errorf("Message with id = '%v' and chat id = '%v' has error `CHANNEL_PRIVATE`. Message will be deleted from queue.", previousMsg.MessageID, previousMsg.ChatID)
				//	TODO возможно стоит удалить и из базы
				update.Ack(false)

				continue
			}

			logrus.Errorf("Failed to receive the message from telegram with message id = '%v' and chat id = '%v' with error = `%#v`. Message passed to message queue.", previousMsg.MessageID, previousMsg.ChatID, err.Error())

			// failed to get message information. The message will be returned to the queue
			update.Ack(false)
			publishMessage(mqCli, previousMsg)
			continue
		}

		// checking the validity of the message content
		if tdlibMsg.Content.GetMessageContentEnum() != "messageText" {
			logrus.Warnf("The content of the message with messageID = %v and chat id = %v is no longer text. The message is removed from the parser queue.", previousMsg.MessageID, previousMsg.ChatID)

			update.Ack(false)
			continue
		}

		// check message group type
		if chat.Type.GetChatTypeEnum() != "chatTypeSupergroup" {
			logrus.Warnf("The type of the chat with id = %v is no longer 'chatTypeSupergroup'. The message is removed from the parser queue.", previousMsg.ChatID)

			update.Ack(false)
			continue
		}

		// get message link
		link, err := tgCli.GetMessageLink(chat.ID, tdlibMsg.ID, true, true)
		if err != nil {
			logrus.Errorf("Could not get the link to the message with messageID = %v and chatID = %v with error = `%v` . Message passed to message queue.", tdlibMsg.ID, chat.ID, err.Error())

			// Could not get the link to the message. Perhaps it no longer exists.
			// The message will be returned to the queue.
			update.Ack(false)
			publishMessage(mqCli, previousMsg)
			continue
		}

		// TODO найти способ отлавливать ошибки если сообщение уже удалено из tg

		// transform the message received from the telegram into a structure *db.Message
		updatedMsg = db.NewMessage(tdlibMsg, chat.Title, link)

		// compare the old message with the updated information
		updates, hasUpdates := compareMessages(previousMsg, updatedMsg)
		if !hasUpdates {
			//logrus.Infof("Message with id = %v and chat id = %v does not differ from the last known version of this message.", previousMsg.MessageID, previousMsg.ChatID)

			// There are no changes. The message will be returned to the queue.
			update.Ack(false)
			publishMessage(mqCli, previousMsg)
			continue
		}

		// The message has updated data. It is necessary to add updates to the database.
		updateCount, err := dbCli.UpdateMessage(updates)
		if err != nil {
			logrus.Errorf("Failed to update db row with chatID = %v and messageID = %v with error = `%v`. Message passed to message queue.", previousMsg.ChatID, previousMsg.MessageID, err.Error())

			update.Ack(false)
			publishMessage(mqCli, previousMsg)
			continue
		}

		if updateCount > 1 {
			logrus.Errorf("Unexpected behavior. UpdateCount = %v (need 1) for message with chatID = %v and messageID = %v", updateCount, previousMsg.ChatID, previousMsg.MessageID)
		}

		// the post does not exist in the database. Need to add
		if updateCount == 0 {
			//logrus.Warnf("Message with id = %v and chat id = %v existed in the queue but did not exist in the database. The message has been added to the database.", previousMsg.MessageID, previousMsg.ChatID)
			a.DbCli.InsertMessage(updatedMsg)
		}

		logrus.Infof("Message with id = '%v' and chat id = '%v' updated with tracking.", updatedMsg.MessageID, updatedMsg.ChatID)

		//	Update the message in the message queue
		//  1. remove the message from the parser queue
		update.Ack(false)
		//	2. add the updated message to the queue
		publishMessage(mqCli, updatedMsg)
	}
}

/*-----------------------------------HELPERS-----------------------------------------*/

// isMsgExpired checks if the message has expired.
func isMsgExpired(message *db.Message, stopTime *int64) (expired bool) {
	if message.Date <= *stopTime {
		return true
	}

	return false
}

// publishMessage
func publishMessage(mqCli *mq.Rabbit, msg *db.Message) {
	err := mqCli.Publish(msg)
	if err != nil {
		logrus.Errorf("Failed to publish message with error '%v'.", err.Error())
	}
}

// maxMsgTimeUpdater updates the maximum message creation time once a day (maximum message age is a month).
func maxMsgTimeUpdater(t *int64) {
	var tickerPeriod = 24 * time.Hour

	*t = time.Now().AddDate(0, -1, 0).Unix()

	ticker := time.NewTicker(tickerPeriod)

	go func() {
		for {
			<-ticker.C
			*t = time.Now().AddDate(0, -1, 0).Unix()
		}
	}()
}

// compareMessages compares messages by chat title, link and message statistics.
func compareMessages(msgFromMq, msgFromTg *db.Message) (updates *db.UpdateRow, hasUpdates bool) {
	updates = &db.UpdateRow{}

	if msgFromTg.ChatTitle != msgFromMq.ChatTitle {
		updates.NewChatTitle = msgFromTg.ChatTitle
		hasUpdates = true
	}

	if msgFromTg.Date > msgFromMq.Date {
		updates.NewDate = msgFromTg.Date
		hasUpdates = true
	}

	if msgFromTg.Views > msgFromMq.Views {
		updates.NewViews = msgFromTg.Views
		hasUpdates = true
	}

	if msgFromTg.Replies > msgFromMq.Replies {
		updates.NewReplies = msgFromTg.Replies
		hasUpdates = true
	}

	if msgFromTg.Forwards > msgFromMq.Forwards {
		updates.NewForwards = msgFromTg.Forwards
		hasUpdates = true
	}

	if msgFromTg.Content != msgFromMq.Content {
		updates.NewContent = msgFromTg.Content
		hasUpdates = true
	}

	if msgFromTg.Link != msgFromMq.Link {
		updates.NewLink = msgFromTg.Link
		hasUpdates = true
	}

	return updates, hasUpdates
}
