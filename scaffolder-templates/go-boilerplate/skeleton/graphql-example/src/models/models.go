package models

import (
	"errors"
	"io"
	"strconv"
	"time"

	"github.com/99designs/gqlgen/graphql"
)

// UserDetail is a model to save user information
type UserDetail struct {
	ID        int       `json:"id" gorm:"primary_key;"`
	Name      string    `json:"name"`
	Email     string    `json:"email" gorm:"unique_index;"`
	Role      string    `json:"role" gorm:"not null;"`
	CreatedAt time.Time `json:"createdAt"  gorm:"not null;"`
	Password  string    `json:"-" gorm:"not null;"`
	Token     string    `json:"token" gorm:"-"`
	Address   Address   `json:"address" gorm:"-"`
}

// Address is a model used to save user address details
type Address struct {
	UserID  int    `json:"-" gorm:"not null;"`
	Country string `json:"country" gorm:"not null;"`
	State   string `json:"state" gorm:"not null;"`
	Zip     int    `json:"zip" gorm:"not null;"`
}

// MarshalTimestamp is a function that convert golang timestamp type to graphql string
func MarshalTimestamp(t time.Time) graphql.Marshaler {
	timestamp := t.Unix() * 1000

	return graphql.WriterFunc(func(w io.Writer) {
		io.WriteString(w, strconv.FormatInt(timestamp, 10))
	})
}

// UnmarshalTimestamp is a function that convert graphql string to golang time stamp ty[e]
func UnmarshalTimestamp(v interface{}) (time.Time, error) {
	if tmpStr, ok := v.(int); ok {
		return time.Unix(int64(tmpStr), 0), nil
	}
	return time.Time{}, errors.New("time should be a unix timestamp")
}
