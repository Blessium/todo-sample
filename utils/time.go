package utils

import (
    "time"
    "encoding/json"
    "errors"
)

type CustomTime struct {
    time.Time
}

const (
    customTimeLayout = "02-01-06 15:04"
)

func (c CustomTime) MarshalJSON() ([]byte, error) {
	formatted := c.Time.Format(customTimeLayout)
	return json.Marshal(formatted)
}

func (c *CustomTime) UnmarshalJSON(data []byte) error {
	var formatted string
	if err := json.Unmarshal(data, &formatted); err != nil {
		return err
	}

	parsed, err := time.Parse(customTimeLayout, formatted)
	if err != nil {
        return errors.New("date format required DD-MM-YY HH:mm")
	}

	c.Time = parsed
	return nil
}
