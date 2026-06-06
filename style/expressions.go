package style

import (
	"encoding/json"
)

// Expression represents a Mapbox GL style expression.
// It can be a literal value (string, number, boolean, null, array, object)
// or a compound expression array: ["operator", arg1, arg2, ...].
type Expression struct {
	IsLiteral bool
	Value     interface{}
	Operator  string
	Args      []*Expression
}

func (e *Expression) UnmarshalJSON(data []byte) error {
	var raw interface{}
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	return e.decode(raw)
}

func (e *Expression) MarshalJSON() ([]byte, error) {
	if e.IsLiteral {
		return json.Marshal(e.Value)
	}
	arr := make([]interface{}, 0, len(e.Args)+1)
	arr = append(arr, e.Operator)
	for _, arg := range e.Args {
		if arg == nil {
			arr = append(arr, nil)
		} else if arg.IsLiteral {
			arr = append(arr, arg.Value)
		} else {
			sub, err := arg.MarshalJSON()
			if err != nil {
				return nil, err
			}
			var v interface{}
			if err := json.Unmarshal(sub, &v); err != nil {
				return nil, err
			}
			arr = append(arr, v)
		}
	}
	return json.Marshal(arr)
}

func (e *Expression) decode(raw interface{}) error {
	switch v := raw.(type) {
	case []interface{}:
		if len(v) == 0 {
			e.IsLiteral = true
			e.Value = v
			return nil
		}
		op, ok := v[0].(string)
		if !ok {
			e.IsLiteral = true
			e.Value = v
			return nil
		}
		e.IsLiteral = false
		e.Operator = op
		e.Args = make([]*Expression, len(v)-1)
		for i, arg := range v[1:] {
			sub := &Expression{}
			if err := sub.decode(arg); err != nil {
				return err
			}
			e.Args[i] = sub
		}
		return nil
	default:
		e.IsLiteral = true
		e.Value = v
		return nil
	}
}
