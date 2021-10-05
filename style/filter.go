package style

import (
	"encoding/json"
	"fmt"

	"github.com/pkg/errors"
)

const (
	FilterOperatorEquals   = "=="
	FilterOperatorNotEqual = "!="
	FilterOperatorAny      = "any"
	FilterOperatorAll      = "all"
	FilterOperatorIn       = "in"
	FilterOperatorNotIn    = "!in"
)

const (
	FilterCategoryType   = "$type"
	FilterTypePoint      = "Point"
	FilterTypeLineString = "LineString"
	FilterTypePolygon    = "Polygon"

	FilterCategoryClass        = "class"
	FilterCategorySubclass     = "subclass"
	FilterCategoryIntermittent = "intermittent"
	FilterCategoryBrunnel      = "brunnel"
	FilterCategoryAdminLevel   = "admin_level"
	FilterCategoryClaimedBy    = "claimed_by"
)

const (
	SourceLayerTransportation string = "transportation"
	SourceLayerWaterway       string = "waterway"
	SourceLayerPlace          string = "place"
)

type FilterType int

const (
	FilterTypeLogical FilterType = iota
	FilterTypeComparative
)

type Filter interface {
	Type() FilterType
	Operator() string
}

type LogicalFilter struct {
	operator   string
	subFilters []Filter
}

func (f *LogicalFilter) Type() FilterType {
	return FilterTypeLogical
}

func (f *LogicalFilter) Operator() string {
	return f.operator
}

type ComparativeFilter struct {
	operator string
	category Category
}

func interfaceListToStringList(interfaceList []interface{}) ([]string, error) {
	var stringList []string
	for _, item := range interfaceList {
		itemStr, ok := item.(string)
		if !ok {
			return nil, errors.Errorf("skipping operand; hoped for string but got %T :: %v\n", item, item)
		}
		stringList = append(stringList, itemStr)
	}

	return stringList, nil
}

func (cf *ComparativeFilter) UnmarshalJSON(data []byte) error {
	var items []interface{}

	err := json.Unmarshal(data, &items)
	if err != nil {
		return err
	}

	cf.operator = items[0].(string)

	switch items[1].(string) {
	case FilterCategoryClass:
		operands, err := interfaceListToStringList(items[2:])
		if err != nil {
			return err
		}

		cf.category = ClassCategory{
			operands: operands,
		}
	case FilterCategorySubclass:
		operands, err := interfaceListToStringList(items[2:])
		if err != nil {
			return err
		}

		cf.category = SubClassCategory{
			operands: operands,
		}
	case FilterCategoryType:
		operands, err := interfaceListToStringList(items[2:])
		if err != nil {
			return err
		}

		cf.category = TypeCategory{
			operands: operands,
		}
	default:
		cf.category = UnhandledCategory{}
	}

	return nil
}

type Category interface {
	OperandCount() int
}

type ClassCategory struct {
	operands []string
}

func (cc ClassCategory) OperandCount() int {
	return len(cc.operands)
}

type SubClassCategory struct {
	operands []string
}

func (scc SubClassCategory) OperandCount() int {
	return len(scc.operands)
}

type TypeCategory struct {
	operands []string
}

func (tc TypeCategory) OperandCount() int {
	return len(tc.operands)
}

type UnhandledCategory struct{}

func (cc UnhandledCategory) OperandCount() int {
	return 0
}

func (f *ComparativeFilter) Type() FilterType {
	return FilterTypeComparative
}

func (f *ComparativeFilter) Operator() string {
	return f.operator
}

type FilterContainer struct {
	filter Filter
}

func (f *FilterContainer) UnmarshalJSON(data []byte) error {
	filter, err := unmarshalFilter(data)
	if err != nil {
		return err
	}

	f.filter = filter
	return nil
}

func unmarshalFilter(data []byte) (Filter, error) {
	var fi []json.RawMessage
	err := json.Unmarshal(data, &fi)
	if err != nil {
		return nil, err
	}

	var operator string
	err = json.Unmarshal(fi[0], &operator)
	if err != nil {
		return nil, err
	}

	if len(fi) == 1 {
		return &LogicalFilter{
			operator: operator,
		}, nil
	}

	var categoryInterface interface{}
	err = json.Unmarshal(fi[1], &categoryInterface)
	if err != nil {
		return nil, err
	}

	switch category := categoryInterface.(type) {
	case []interface{}:
		filter := &LogicalFilter{
			operator:   operator,
			subFilters: []Filter{},
		}
		for _, subFilterJSON := range fi[1:] {
			subFilter, err := unmarshalFilter(subFilterJSON)
			if err != nil {
				return nil, err
			}
			filter.subFilters = append(filter.subFilters, subFilter)
		}
		return filter, nil
	case string:
		cf := new(ComparativeFilter)
		err = json.Unmarshal(data, cf)
		if err != nil {
			return nil, err
		}
		return cf, nil
	default:
		panic(fmt.Sprintf("not implemented: %s :: %T :: %s", category, category, data))
	}
}
