package db

import (
	_ "github.com/lib/pq"
)

// TODO переписать все тесты в соответствии с изменениями
//
//func TestMain(m *testing.M) {
//	var db *sql.DB
//	var err error
//	database := "postgres"
//
//	pool, err := dockertest.NewPool("")
//	if err != nil {
//		log.Fatalf("Could not connect to docker: %s", err)
//	}
//
//	resource, err := pool.Run("postgres", "12.5", []string{"POSTGRES_PASSWORD=secret", "POSTGRES_DB=" + database})
//	if err != nil {
//		log.Fatalf("Could not start resource: %s", err)
//	}
//
//	if err = pool.Retry(func() error {
//		var err error
//		testPort = resource.GetPort("5432/tcp")
//		db, err = sql.Open("postgres", fmt.Sprintf("postgres://postgres:secret@localhost:%s/%s?sslmode=disable", testPort, database))
//		if err != nil {
//			return err
//		}
//		return db.Ping()
//	}); err != nil {
//		log.Fatalf("Could not connect to docker: %s", err)
//	}
//
//	code := m.Run()
//
//	if err = pool.Purge(resource); err != nil {
//		log.Fatalf("Could not purge resource: %s", err)
//	}
//
//	os.Exit(code)
//
//	//	TODO удалять контейнер после использования
//}
//
////    ---------------------------------INSERT-----------------------------------------
//
//func TestPostgresClient_Insert(t *testing.T) {
//	pg, err := postgresTestConnection()
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	// In the first cycle, all values must be entered into the table without errors
//	for _, testMsg := range testInsert {
//		err = pg.InsertClient(testMsg)
//		if err != nil {
//			t.Errorf(
//				"The first sycle.\nUnexpected error while entering %v into the database.\nERROR expected %v\nERROR got %v\n",
//				color.MagentaString("%#v", testMsg),
//				color.GreenString("<nil>"), color.RedString("%#v", err.Error()),
//			)
//		}
//	}
//
//	// The values of the second cycle should not be entered into the table, since it already contains rows with this set of testMsg ID and chat ID. An error must be returned
//	for _, testMsg := range testInsert {
//		err = pg.InsertClient(testMsg)
//		if err == nil {
//			t.Errorf(
//				"The second cycle. No error received while entering %v into the database.\nERROR expected %v\nERROR got %v\n",
//				color.MagentaString("%#v", testMsg),
//				color.GreenString("pq: duplicate key value violates unique constraint 'tg_parser_pkey'"),
//				color.RedString("<nil>"),
//			)
//		}
//	}
//
//	count, err := postgresTestGetRowsCount(pg)
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	if count != len(testInsert) {
//		t.Errorf(
//			"Expected number of rows %v\ngot %v",
//			color.GreenString("%v", len(testInsert)), color.RedString("%v", count),
//		)
//	}
//}
//
////    --------------------------------GETMESSAGEBYID----------------------------------
//
//func TestPostgresClient_GetMessageById(t *testing.T) {
//	pg, err := postgresTestConnection()
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	for _, test := range testGetMessageByIds {
//		message, err := pg.GetByID(test.ChatID, test.MessageID)
//
//		stringErr := fmt.Sprintf(
//			"For MessageID = %v and ChatID = %v\nMESSAGE expected %v\nMESSAGE got %v\nERROR expected %v\nERROR got %v\n",
//			color.MagentaString("%v", test.MessageID), color.MagentaString("%v", test.ChatID),
//			color.GreenString("%#v", test.ExpectedMessage), color.RedString("%#v", message),
//			color.GreenString("%#v", test.ErrorText), color.RedString("%#v", err),
//		)
//
//		if err != nil && !test.HasError {
//			t.Error(stringErr)
//			continue
//		}
//
//		if err == nil && test.ExpectedMessage == nil {
//			t.Error(stringErr)
//			continue
//		}
//
//		if !cmp.Equal(message, test.ExpectedMessage) {
//			t.Error(stringErr)
//			continue
//		}
//	}
//}
//
////    ------------------------------------UPDATE--------------------------------------
//
//func TestPostgresClient_Update(t *testing.T) {
//	pg, err := postgresTestConnection()
//	if err != nil {
//		log.Fatal("postgresTestConnection error: ", err)
//	}
//
//	for index, test := range testUpdates {
//
//		// Removing previous values in database
//		err = postgresTestDeleteRow(pg, test.InitialMessage.ChatID, test.InitialMessage.MessageID)
//		if err != nil {
//			log.Fatalf("RANGE INDEX %v\npostgresTestDeleteRow error: %v \n",
//				color.MagentaString("%v", index), color.RedString("%v", err.Error()))
//		}
//
//		var gotMessage *Message
//		var expectedMessage Message
//		var updateValues UpdateRow
//
//		gotMessage = test.InitialMessage
//
//		expectedMessage = *test.InitialMessage
//		updateValues = *test.UpdateRow
//
//		// if updates are expected to be accepted, the expected message will change according to the updateValues
//		if test.ExpectedUpdateCount == 1 {
//			expectedMessage.ChatTitle = updateValues.NewChatTitle
//			expectedMessage.Content = updateValues.NewContent
//			expectedMessage.Views = updateValues.NewViews
//			expectedMessage.Forwards = updateValues.NewForwards
//			expectedMessage.Replies = updateValues.NewReplies
//			expectedMessage.Date = updateValues.NewDate
//		}
//
//		// MessageToUpdate is constant
//		err = pg.InsertClient(test.InitialMessage)
//		if err != nil {
//			log.Fatalf(
//				"RANGE INDEX %v\nInsert error: %v\n",
//				color.MagentaString("%v", index), color.RedString("%v", err.Error()))
//		}
//
//		updateCount, err := pg.Update(test.UpdateRow)
//		if (err != nil && test.ExpectedUpdateCount == 1) || (updateCount != test.ExpectedUpdateCount) {
//			// 1 row update expected. Instead, we got an error, or the number of updated rows is not equal to 1
//			t.Error(fmt.Sprintf(
//				"RANGE INDEX %v\n"+
//					"Init message: %v\nUpdate values: %v\n"+
//					"RESULT expected %v\nRESULT got %v\n"+
//					"ERROR expected %v\nERROR got %v\n"+
//					"COUNT expected %v\nCOUNT got %v\n",
//				color.MagentaString("%#v", index),
//				color.MagentaString("%#v", test.InitialMessage), color.MagentaString("%#v", test.UpdateRow),
//				color.GreenString("%#v", expectedMessage), color.RedString("%#v", gotMessage),
//				color.GreenString("%#v", test.ErrorText), color.RedString("%#v", err),
//				color.GreenString("%#v", test.ExpectedUpdateCount), color.RedString("%#v", updateCount),
//			))
//
//			continue
//		} else if err != nil && test.ExpectedUpdateCount == 0 {
//			// the error that we made a mistake received, the test passed
//			continue
//		}
//
//		// receive updated message from DB
//		gotMessage, err = pg.GetByID(test.InitialMessage.ChatID, test.InitialMessage.MessageID)
//		if err != nil {
//			t.Errorf("RANGE INDEX %v\nUnexpected behavior. Error \"%v\" returned from method %v",
//				color.MagentaString("%#v", index),
//				color.RedString("%v", err.Error()), color.MagentaString("GetByID"))
//			continue
//		}
//
//		// checking the received message for compliance with expectations
//		if !cmp.Equal(*gotMessage, expectedMessage) {
//			t.Error(fmt.Sprintf(
//				"RANGE INDEX %v\nInit message: %v\nUpdate values: %v\nRESULT expected %v\nRESULT got %v\nERROR expected %v\nERROR got %v\n",
//				color.MagentaString("%#v", index),
//				color.MagentaString("%#v", test.InitialMessage), color.MagentaString("%#v", test.UpdateRow),
//				color.GreenString("%#v", expectedMessage), color.RedString("%#v", gotMessage),
//				color.GreenString("%#v", test.ErrorText), color.RedString("%#v", err),
//			))
//		}
//	}
//}

//func TestPostgresClient_GetMessagesForATimePeriod(t *testing.T) {
//	pg, err := postgresTestConnection()
//	if err != nil {
//		log.Fatal("postgresTestConnection error: ", err)
//	}
//
//}
