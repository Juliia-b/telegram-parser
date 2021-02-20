package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/jinzhu/now"
	"net/http"
	"sync"
	"telegram-parser/db"
	"time"
)

/*---------------------------------STRUCTURES----------------------------------------*/

type handler struct {
	dbCli      db.DB
	ws         *ws
	CookieName string
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

// getBestInPeriod returns the best posts for the specified period. Limit is 50. Less can be returned.
func (h *handler) getBestInPeriod(w http.ResponseWriter, r *http.Request) {
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

// handlerInit initializes structure handler.
func handlerInit(dbCli db.DB) *handler {
	var wsConn = make(map[*websocket.Conn]int)
	return &handler{
		dbCli: dbCli,
		ws: &ws{
			connections: wsConn,
			rwMutex:     &sync.RWMutex{},
		},
		CookieName: "u.v1",
	}
}

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

/*----------------------------------DEPRECATE----------------------------------------*/

// DEPRECATE
// sessionMiddleware
//func (h *server) sessionMiddleware(next http.Handler) http.Handler {
//	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//		cookieName := h.CookieName
//
//		_, err := getCookie(r, cookieName)
//		if err == nil {
//			// cookie presents
//			next.ServeHTTP(w, r)
//		}
//
//		cookie := generateCookie(cookieName)
//
//		//TODO обдумать для чего использовать таблицу client
//
//		//h.dbCli.InsertClient(db.Client{
//		//	ID:     0,
//		//	Cookie: cookie,
//		//})
//
//		http.SetCookie(w, cookie)
//
//		next.ServeHTTP(w, r)
//	})
//}

// getCookie gets the value of the set cookie by name.
//func getCookie(r *http.Request, cookieName string) (cookie string, err error) {
//	c, err := r.Cookie(cookieName)
//	if c != nil {
//		cookie = c.Value
//	}
//
//	return cookie, err
//}

//// generateCookie returns *http.Cookie with filled fields.
//func generateCookie(cookieName string) (cookie *http.Cookie) {
//	val := generateString()
//
//	return &http.Cookie{
//		Name:   cookieName,
//		Value:  val, // Some encoded value
//		Path:   "/", // Otherwise it defaults to the /login if you create this on /login (standard cookie behaviour)
//		MaxAge: 0,   // MaxAge=0 means no 'Max-Age' attribute specified.
//	}
//}

//// generateString generates a fixed length string from unix time.
//func generateString() string {
//	rand.Seed(time.Now().UnixNano())
//
//	//Only lowercase
//	var charSet = "abcdedfghijkluywxzmnopqrst"
//	var result string
//	var resultStringLen = 12
//
//	for i := 0; i < resultStringLen; i++ {
//		randomIndex := rand.Intn(len(charSet))
//		randomChar := charSet[randomIndex]
//		result += string(randomChar)
//	}
//
//	return result
//}
