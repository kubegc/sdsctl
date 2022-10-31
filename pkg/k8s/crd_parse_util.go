package k8s

import (
	"encoding/json"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
)

func AddPowerStatus(bytes []byte, message, reason string) ([]byte, error) {
	statusMap := map[string]interface{}{
		"message": message,
		"reason":  reason,
	}
	return UpdateCRDSpec(bytes, "status.conditions.state.waiting", statusMap)
}

func UpdateCRDSpec(spec []byte, key string, value interface{}) ([]byte, error) {
	bytes, err := sjson.SetBytes(spec, key, value)
	if err != nil {
		return []byte{}, err
	}
	return bytes, nil
}

func GetCRDSpec(spec []byte, key string) (map[string]string, error) {
	parse := gjson.ParseBytes(spec)
	value := parse.Get(key)
	res := make(map[string]string)
	if err := json.Unmarshal([]byte(value.Raw), &res); err != nil {
		return nil, err
	}
	return res, nil
}
