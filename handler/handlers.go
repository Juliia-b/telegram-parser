package handler

import (
	"errors"
	"fmt"
	"github.com/jinzhu/now"
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
}

//    --------------------------------------------------------------------------------
//                                     METHODS
//    --------------------------------------------------------------------------------

func (h *handler) GetMinInPeriod(w http.ResponseWriter, r *http.Request) {
	//vars := mux.Vars(r)

}

func (h *handler) GetMaxInPeriod(w http.ResponseWriter, r *http.Request) {
	//vars := mux.Vars(r)
	//h.dbCli.GetMessagesForATimePeriod()

}

//func ArticlesCategoryHandler(w http.ResponseWriter, router *http.Request) {
//	vars := mux.Vars(router)
//	w.WriteHeader(http.StatusOK)
//	fmt.Fprintf(w, "Category: %v\n", vars["category"])
//}

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
	default:
		err = errors.New(fmt.Sprintf("unknown time period %v", period))
	}

	return from, to, err
}
