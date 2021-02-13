package parser

import (
	"github.com/sirupsen/logrus"
	"telegram-parser/db"
	"telegram-parser/mq"
	"time"
)

// RunStatisticsTracking launches several handlers for updates from telegram.
// @param handlersCount is the number of message handlers that will be run. The minimum value is 10.
func (a *App) RunStatisticsTracking(handlersCount int) {
	if handlersCount < 10 {
		handlersCount = 10
	}

	// updates the maximum value of the message date
	now := time.Now().Unix()
	maxMsgTime(&now)

	for i := 0; i < handlersCount; i++ {
		go a.trackingStatistics(&now)
	}
}

// trackingStatistics starts continuous tracking of message statistics.
func (a *App) trackingStatistics(stopTime *int64) {
	mqCli := a.MqCli
	dbCli := a.DbCli
	tgCli := a.Telegram.Client

	for update := range mqCli.Consume() {
		mqMsg := mq.UnmarshalRabbitBody(update.Body)

		//logrus.Infof("Получено сообщение с id %v из rabbit \n\n", mqMsg.MessageID)

		// check if the message has expired
		if mqMsg.Date <= *stopTime {
			// message expired so we remove the message from the parser queue
			update.Ack(false)
			continue
		}

		// ask the Telegram for up-to-date information about the message
		tgMsg, err := tgCli.GetMessage(mqMsg.ChatID, mqMsg.MessageID)
		if err != nil {
			update.Nack(false, false)
			logrus.Errorf("Failed to receive message with chatID = %v and messageID = %v with error = `%v`. Message passed to message queue.", mqMsg.ChatID, mqMsg.MessageID, err.Error())
			continue
		}

		// checking the validity of the message content
		if tgMsg.Content.GetMessageContentEnum() != "messageText" {
			update.Ack(false)
			logrus.Warnf("The content of the message with messageID = %v is no longer text. The message is removed from the parser queue.", tgMsg.ID)
			continue
		}

		// request information about the chat from the telegram
		chat, err := tgCli.GetChat(mqMsg.ChatID)
		if err != nil {
			update.Nack(false, false)
			logrus.Errorf("Failed to receive chat with chatID = %v with error = `%v`. Message passed to message queue.", mqMsg.ChatID, err.Error())
			continue
		}

		// transform the message received from the telegram into a structure *db.Message
		updatedMsg := db.NewMessage(tgMsg, chat.Title)

		// compare the old message with the updated information
		updates, hasUpdates := compareMessages(mqMsg, updatedMsg)
		if !hasUpdates {
			// no updates available
			update.Nack(false, false)
			continue
		}

		// update table fields
		updateCount, err := dbCli.Update(updates)
		if err != nil {
			logrus.Errorf("Failed to update db row with chatID = %v and messageID = %v with error = `%v`. Message passed to message queue.", mqMsg.ChatID, mqMsg.MessageID, err.Error())
			update.Nack(false, false)
			continue
		}

		if updateCount != 1 {
			logrus.Warnf("Unexpected behavior. UpdateCount = %v (need 1) for message with chatID = %v and messageID = %v", updateCount, mqMsg.ChatID, mqMsg.MessageID)
		}

		//	Update the message in the message queue
		//  1. remove the message from the parser queue
		update.Ack(false)
		//  2. get a satisfactory structure
		msg := convertUpdateRawToMessage(updates)
		//	3. add the updated message to the queue
		err = mqCli.Publish(msg)
		if err != nil {
			logrus.Fatalf("Failed to send message with chatID = %v and messageID = %v to message queue with error = `%v`.", msg.ChatID, msg.MessageID, err.Error())
		}
	}
}

//    --------------------------------------------------------------------------------
//                                     HELPERS
//    --------------------------------------------------------------------------------

// maxMsgTime updates the maximum message creation time once a day (maximum message age is a month).
func maxMsgTime(t *int64) {
	*t = time.Now().AddDate(0, -1, 0).Unix()

	ticker := time.NewTicker(24 * time.Hour)

	go func() {
		for {
			<-ticker.C
			*t = time.Now().AddDate(0, -1, 0).Unix()
		}
	}()
}

// compareMessages compares messages by chat title and message statistics.
func compareMessages(msgFromMq, msgFromTg *db.Message) (updates *db.UpdateRow, hasUpdates bool) {
	var (
		mq = msgFromMq
		tg = msgFromTg
	)

	updates = &db.UpdateRow{}

	if tg.ChatTitle != mq.ChatTitle {
		updates.NewChatTitle = tg.ChatTitle
		hasUpdates = true
	}

	if tg.Date > mq.Date {
		updates.NewDate = tg.Date
		hasUpdates = true
	}

	if tg.Views > mq.Views {
		updates.NewViews = tg.Views
		hasUpdates = true
	}

	if tg.Replies > mq.Replies {
		updates.NewReplies = tg.Replies
		hasUpdates = true
	}

	if tg.Forwards > mq.Forwards {
		updates.NewForwards = tg.Forwards
		hasUpdates = true
	}

	if tg.Content != mq.Content {
		updates.NewContent = tg.Content
		hasUpdates = true
	}

	return updates, hasUpdates
}

// convertUpdateRawToMessage converts structure db.UpdateRow to structure T db.Message.
func convertUpdateRawToMessage(update *db.UpdateRow) *db.Message {
	var msg db.Message

	msg.MessageID = update.MessageID
	msg.ChatID = update.ChatID
	msg.Date = update.NewDate
	msg.ChatTitle = update.NewChatTitle
	msg.Content = update.NewContent
	msg.Views = update.NewViews
	msg.Forwards = update.NewForwards
	msg.Replies = update.NewReplies

	return &msg
}
