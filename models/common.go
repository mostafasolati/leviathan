package models

import (
	"fmt"
	"strings"
	"time"
)

// App is the data-type for the client apps.
type App int

// List of client apps.
const (
	WebApp App = iota
	AndroidApp
	MobileWebApp
)

var appNames = map[App]string{
	WebApp:       "web",
	AndroidApp:   "android",
	MobileWebApp: "mobile-web",
}

// Name returns the client app's name.
func (app App) Name() string {
	if name, ok := appNames[app]; ok {
		return name
	}
	return "unknown"
}

// Error is Error response in http requests
type Error struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
}

// Pagination is an object contain pagination data
// The actual data stores in Data property other properties
// contain information about the pagination
type Pagination struct {
	Total       int         `json:"total"`
	PerPage     int         `json:"per_page"`
	CurrentPage int         `json:"current_page"`
	LastPage    int         `json:"last_page"`
	From        int         `json:"from"`
	To          int         `json:"to"`
	Data        interface{} `json:"data"`
}

// Base contains base properties common between all data models
type Base struct {
	ID        int        `json:"id"`
	CreatedAt *time.Time `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`
}

type Timestamp time.Time

func (t *Timestamp) UnmarshalJSON(b []byte) error {

	src := string(b)
	src = strings.Trim(src, `"`)
	layout := "2006-01-02"
	if len(src) > 10 {
		layout = "2006-01-02 15:04"
	}

	loc, _ := time.LoadLocation("Asia/Tehran")
	ts, err := time.ParseInLocation(layout, src, loc)
	*t = Timestamp(ts)

	return err
}

func (t *Timestamp) MarshalJSON() ([]byte, error) {
	d := time.Time(*t)
	return []byte(fmt.Sprintf(`"%d-%02d-%02d"`, d.Year(), d.Month(), d.Day())), nil
}
