package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/jinzhu/now"
	"net/http"
	"telegram-parser/db"
	"time"
)

/*---------------------------------STRUCTURES----------------------------------------*/

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

/*-----------------------------------METHODS-----------------------------------------*/

// GetBestInPeriod returns the best posts for the specified period. Limit is 50. Less can be returned.
func (h *handler) GetBestInPeriod(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	var limit = 50
	var period = r.FormValue("period")

	from, to, err := dateCalculation(period)
	if err != nil {
		// Period is not valid
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	messages, err := h.dbCli.GetMessageWithPeriod(from, to, limit)
	if err != nil {
		w.WriteHeader(http.StatusNotImplemented)
		return
	}

	if len(messages) == 0 {
		v := map[string]int64{"from": from, "to": to}

		payload, err := json.Marshal(v)
		if err != nil {
			w.WriteHeader(http.StatusNotImplemented)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(payload)
		return
	}

	payload, err := json.Marshal(messages)
	if err != nil {
		w.WriteHeader(http.StatusNotImplemented)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(payload)
}

/*-----------------------------------HELPERS-----------------------------------------*/

// getTimePeriods returns a structure with constant names of periods.
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
