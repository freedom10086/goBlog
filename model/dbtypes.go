package model

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/lib/pq"
	"time"
	"database/sql/driver"
)

type NullString sql.NullString

// Scan implements the Scanner interface.
func (ns *NullString) Scan(value interface{}) error {
	ns.String, ns.Valid = value.(string)
	return nil
}

// Value implements the driver Valuer interface.
func (ns NullString) Value() (driver.Value, error) {
	if !ns.Valid {
		return nil, nil
	}
	return ns.String, nil
}

// MarshalJSON for NullString
func (ns *NullString) MarshalJSON() ([]byte, error) {
	if ns.Valid {
		return json.Marshal(ns.String)
	} else {
		return json.Marshal(nil)
	}
}

// UnmarshalJSON for NullString
func (ns *NullString) UnmarshalJSON(b []byte) error {
	err := json.Unmarshal(b, &ns.String)
	ns.Valid = err == nil
	return err
}

type NullTime pq.NullTime

// Scan implements the Scanner interface.
func (nt *NullTime) Scan(value interface{}) error {
	nt.Time, nt.Valid = value.(time.Time)
	return nil
}

// Value implements the driver Valuer interface.
func (nt NullTime) Value() (driver.Value, error) {
	if !nt.Valid {
		return nil, nil
	}
	return nt.Time, nil
}

func (nt *NullTime) MarshalJSON() ([]byte, error) {

	if !nt.Valid {
		return json.Marshal(nil)
		//return []byte("null"), nil
	}
	val := fmt.Sprintf("\"%s\"", nt.Time.Format(time.RFC3339))
	return []byte(val), nil
}

// UnmarshalJSON for NullTime
func (nt *NullTime) UnmarshalJSON(b []byte) error {
	/*
	s := string(b)
	x, err := time.Parse(time.RFC3339, s)
	if err != nil {
		nt.Valid = false
		return err
	}
	nt.Time = x
	nt.Valid = true
	return nil
	*/

	err := json.Unmarshal(b, &nt.Time)
	nt.Valid = err == nil
	return err
}
