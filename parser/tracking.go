package parser

import (
	"github.com/sirupsen/logrus"
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
	//dbCli := a.DbCli
	tgCli := a.Telegram.Client

	updates := mqCli.Consume()

	for update := range updates {
		msg := mq.UnmarshalRabbitBody(update.Body)

		logrus.Infof("Получено сообщение с id %v из rabbit \n\n", msg.MessageID)

		if msg.Date <= *stopTime {
			// message expired so we remove the message from the parser queue
			update.Ack(false)
		}

		//TODO update

		tgCli.GetMessage(msg.ChatID, msg.MessageID)

		// -------------
		//	update.Nack(false, false)

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
