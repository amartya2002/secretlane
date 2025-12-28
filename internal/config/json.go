package config

import "encoding/json"

// UnmarshalJSON is a tiny helper so other packages don't need to import encoding/json directly
// when working with generic JSON strings from the database.
func UnmarshalJSON(data string, v any) error {
	return json.Unmarshal([]byte(data), v)
}

