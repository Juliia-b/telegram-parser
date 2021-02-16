package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/jinzhu/now"
	"github.com/sirupsen/logrus"
	"net/http"
	"telegram-parser/db"
	"time"
)

//    --------------------------------------------------------------------------------
//                                  STRUCTURES
//    --------------------------------------------------------------------------------

type handler struct {
	dbCli db.DB
}

type timePeriods struct {
	Today              string
	Yesterday          string
	DayBeforeYesterday string
	ThisWeek           string
	LastWeek           string
	ThisMonth          string
	Whole              string // Denotes the entire period from 1970-01-01T00: 00: 00Z to the present
}

//    --------------------------------------------------------------------------------
//                                     METHODS
//    --------------------------------------------------------------------------------

//func (h *handler) GetWorstInPeriod(w http.ResponseWriter, r *http.Request) {
//	//vars := mux.Vars(r)
//
//}

func (h *handler) GetBestInPeriod(w http.ResponseWriter, r *http.Request) {
	var limit int = 50
	period := r.FormValue("period")

	logrus.Infof("Period : %v ; Limit : %v\n", period, limit)

	from, to, err := dateCalculation(period)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("\nPeriod is not valid.\n\n"))
		return
	}

	messages, err := h.dbCli.GetMessagesForATimePeriod(from, to, limit)
	if err != nil {
		w.WriteHeader(http.StatusNotImplemented)
		w.Write([]byte("\nSomething went wrong.\n\n"))
		return
	}

	payload, err := json.Marshal(messages)
	if err != nil {
		w.WriteHeader(http.StatusNotImplemented)
		w.Write([]byte(fmt.Sprintf("\nSomething went wrong with error %v.\n\n", err.Error())))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(payload)
}

//    --------------------------------------------------------------------------------
//                                     HELPERS
//    --------------------------------------------------------------------------------

func getTimePeriods() *timePeriods {
	return &timePeriods{
		Today:              "today",
		Yesterday:          "yesterday",
		DayBeforeYesterday: "daybeforeyesterday",
		ThisWeek:           "thisweek",
		LastWeek:           "lastweek",
		ThisMonth:          "thismonth",
		Whole:              "whole",
	}
}

// dateCalculation calculates the start and end unix times for the period.
func dateCalculation(period string) (from int64, to int64, err error) {
	p := getTimePeriods()
	t := time.Now()

	conf := &now.Config{
		WeekStartDay: time.Monday,
	}

	switch period {
	case p.Today:
		from = conf.With(t).BeginningOfDay().Unix()
		to = conf.With(t).EndOfDay().Unix()
	case p.Yesterday:
		t := t.AddDate(0, 0, -1)
		from = conf.With(t).BeginningOfDay().Unix()
		to = conf.With(t).EndOfDay().Unix()
	case p.DayBeforeYesterday:
		t := t.AddDate(0, 0, -2)
		from = conf.With(t).BeginningOfDay().Unix()
		to = conf.With(t).EndOfDay().Unix()
	case p.ThisWeek:
		from = conf.With(t).BeginningOfWeek().Unix()
		to = conf.With(t).EndOfWeek().Unix()
	case p.LastWeek:
		t := t.AddDate(0, 0, -7)
		from = conf.With(t).BeginningOfWeek().Unix()
		to = conf.With(t).EndOfWeek().Unix()
	case p.ThisMonth:
		from = conf.With(t).BeginningOfMonth().Unix()
		to = conf.With(t).EndOfMonth().Unix()
	case p.Whole:
		from = 0
		to = t.Unix()
	default:
		err = errors.New(fmt.Sprintf("unknown time period %v", period))
	}

	return from, to, err
}
