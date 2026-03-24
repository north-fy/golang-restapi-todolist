package parse

import (
	"encoding/json"
)

func convertFromJSON(data interface{}) ([]byte, error) {
	switch v := data.(type) {
	case string:
		res, err := json.Marshal(v)
		if err != nil {
			return nil, err
		}

		return res, nil
	}

	return nil, nil
}
