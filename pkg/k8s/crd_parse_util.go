package k8s

import "github.com/tidwall/sjson"

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
