package model

import (
	"encoding/json"
)

// Cover File or null
type Cover struct {
	Valid bool
	File
}

// MarshalJSON parse value or nil
func (v Cover) MarshalJSON() ([]byte, error) {
	if v.Valid {
		return json.Marshal(v.File)
	}

	return json.Marshal(nil)
}

// UnmarshalJSON returns value or nil
func (v *Cover) UnmarshalJSON(data []byte) error {
	// Unmarshalling into a pointer will let us detect null
	var x *File
	if err := json.Unmarshal(data, &x); err != nil {
		return err
	}

	if x != nil {
		v.Valid = true
		v.File = *x
	} else {
		v.Valid = false
	}

	return nil
}
