package style

import "encoding/json"

type numericStopType struct {
	ZoomLevel ZoomLevel
	Value     float64
}

func (n *numericStopType) UnmarshalJSON(data []byte) error {
	type internalType [2]interface{}

	var i internalType
	err := json.Unmarshal(data, &i)
	if err != nil {
		return err
	}

	n.ZoomLevel = ZoomLevel(i[0].(float64))
	n.Value = i[1].(float64)

	return nil
}
