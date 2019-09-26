package storable

import (
	"database/sql/driver"
	"encoding/hex"
	"encoding/json"
	"time"
)

// ByteArray is a custom type that maps to a the database `bytea` fields
type ByteArray string

func (val *ByteArray) Scan(value interface{}) error {
	encoded := hex.EncodeToString(value.([]byte))
	*val = ByteArray(encoded)

	return nil
}

func (val ByteArray) Value() (driver.Value, error) {
	return hex.DecodeString(string(val))
}

func (val ByteArray) String() string {
	return string(val)
}

// JSONStringArray binds a slice of strings to a `jsonb` database field
type JSONStringArray []string

func (j *JSONStringArray) Scan(value interface{}) error {
	err := json.Unmarshal(value.([]byte), j)

	return err
}

func (j JSONStringArray) Value() (driver.Value, error) {
	return json.Marshal(j)
}

// JSONObject binds a map[string]interface{} to a `jsonb` database field
type JSONObject map[string]interface{}

func (obj *JSONObject) Scan(value interface{}) error {
	err := json.Unmarshal(value.([]byte), obj)

	return err
}

func (obj JSONObject) Value() (driver.Value, error) {
	return json.Marshal(obj)
}

// DatetimeToJSONUnix binds a time.Time to a `timestamp` database field
// when marshaled to JSON, outputs a unix timestamp
type DatetimeToJSONUnix time.Time

func (t DatetimeToJSONUnix) MarshalJSON() ([]byte, error) {
	return json.Marshal(time.Time(t).UTC().Unix())
}

func (t *DatetimeToJSONUnix) UnmarshalJSON(data []byte) error {
	var unix int64

	err := json.Unmarshal(data, &unix)
	if err != nil {
		return err
	}

	*t = DatetimeToJSONUnix(time.Unix(unix, 0).UTC())

	return nil
}

func (t *DatetimeToJSONUnix) Scan(value interface{}) error {
	*t = DatetimeToJSONUnix(value.(time.Time))

	return nil
}

func (t DatetimeToJSONUnix) Value() (driver.Value, error) {
	return time.Time(t), nil
}

func (t DatetimeToJSONUnix) String() string {
	return time.Time(t).UTC().String()
}
