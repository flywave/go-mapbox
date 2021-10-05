package mapboxglstyle

import (
	"encoding/json"

	"github.com/jamesrr39/goutil/errorsx"
	"github.com/jamesrr39/ownmap-app/ownmap"
)

type numericStopType struct {
	ZoomLevel ownmap.ZoomLevel
	Value     float64
}

func (n *numericStopType) UnmarshalJSON(data []byte) error {
	type internalType [2]interface{}

	var i internalType
	err := json.Unmarshal(data, &i)
	if err != nil {
		return errorsx.Wrap(err)
	}

	n.ZoomLevel = ownmap.ZoomLevel(i[0].(float64))
	n.Value = i[1].(float64)

	return nil
}
