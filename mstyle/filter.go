package mapboxglstyle

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/jamesrr39/goutil/errorsx"
	"github.com/jamesrr39/ownmap-app/ownmap"
	"github.com/jamesrr39/ownmap-app/ownmapdal"
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

/*

    "filter": [
        "all",
        ["==", "$type", "Polygon"],
		["in", "class", "residential", "suburb", "neighbourhood"]
	]

	"filter": ["==", "$type", "Point"],

	"filter": ["all",["==","$type","Polygon"],["in","class","residential","suburb","neighbourhood"]]
*/

// https://openmaptiles.org/schema/#transportation
const (
	SourceLayerTransportation string = "transportation"
	SourceLayerWaterway       string = "waterway"
	SourceLayerPlace          string = "place"
)

type FilterType int

const (
	// FilterTypeLogical is a filter containing 2 (or more) sub-filters
	FilterTypeLogical FilterType = iota
	// FilterTypeComparative is a filter containing comparisions on e.g. object classes (it does not contain sub-filters)
	FilterTypeComparative
)

type Filter interface {
	Type() FilterType
	Operator() string
	GetAllPossibleTagKeys(sourceLayer string) []*ownmapdal.TagKeyWithType
	IsObjectShown(sourceLayer string, tags []*ownmap.OSMTag, objectFilterType ownmap.ObjectType) bool
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

func (f *LogicalFilter) GetAllPossibleTagKeys(sourceLayer string) []*ownmapdal.TagKeyWithType {
	var tagKeys []*ownmapdal.TagKeyWithType
	for _, subFilter := range f.subFilters {
		tagKeys = append(tagKeys, subFilter.GetAllPossibleTagKeys(sourceLayer)...)
	}
	return tagKeys
}

func (f *LogicalFilter) IsObjectShown(sourceLayer string, tags []*ownmap.OSMTag, objectFilterType ownmap.ObjectType) bool {
	switch f.operator {
	case FilterOperatorAll:
		for _, subfilter := range f.subFilters {
			shown := subfilter.IsObjectShown(sourceLayer, tags, objectFilterType)
			if !shown {
				return false
			}
		}
		return true
	default:
		panic("unhandled logical filter operator: " + f.operator)
	}
}

type ComparativeFilter struct {
	operator string
	category Category // class, $type, subclass, brunnel
}

func interfaceListToStringList(interfaceList []interface{}) ([]string, errorsx.Error) {
	var stringList []string
	for _, item := range interfaceList {
		itemStr, ok := item.(string)
		if !ok {
			return nil, errorsx.Errorf("skipping operand; hoped for string but got %T :: %v\n", item, item)
		}
		stringList = append(stringList, itemStr)
	}

	return stringList, nil
}

func (cf *ComparativeFilter) UnmarshalJSON(data []byte) error {
	var items []interface{}

	err := json.Unmarshal(data, &items)
	if err != nil {
		return errorsx.Wrap(err)
	}

	cf.operator = items[0].(string)

	switch items[1].(string) {
	case FilterCategoryClass:
		operands, err := interfaceListToStringList(items[2:])
		if err != nil {
			return errorsx.Wrap(err)
		}

		cf.category = ClassCategory{
			operands: operands,
		}
	case FilterCategorySubclass:
		operands, err := interfaceListToStringList(items[2:])
		if err != nil {
			return errorsx.Wrap(err)
		}

		cf.category = SubClassCategory{
			operands: operands,
		}
	case FilterCategoryType:
		operands, err := interfaceListToStringList(items[2:])
		if err != nil {
			return errorsx.Wrap(err)
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
	Matches(operandIdx int, sourceLayer string, objectTags []*ownmap.OSMTag, objectType ownmap.ObjectType) bool
	OperandCount() int
	GetAllPossibleTagKeys(sourceLayer string) []*ownmapdal.TagKeyWithType
}

type ClassCategory struct {
	operands []string
}

func (cc ClassCategory) Matches(operandIdx int, sourceLayer string, objectTags []*ownmap.OSMTag, objectType ownmap.ObjectType) bool {
	return isClassShown(cc.operands[operandIdx], sourceLayer, objectTags)
}

func (cc ClassCategory) GetAllPossibleTagKeys(sourceLayer string) []*ownmapdal.TagKeyWithType {
	return getOSMTagsWithTypesFromMapboxGLTypes(sourceLayer, cc.operands, mapMapboxGLClassToOSMTags)
}

func (cc ClassCategory) OperandCount() int {
	return len(cc.operands)
}

type SubClassCategory struct {
	operands []string
}

func (scc SubClassCategory) Matches(operandIdx int, sourceLayer string, objectTags []*ownmap.OSMTag, objectType ownmap.ObjectType) bool {
	return isSubclassShown(scc.operands[operandIdx], sourceLayer, objectTags)
}

func (scc SubClassCategory) GetAllPossibleTagKeys(sourceLayer string) []*ownmapdal.TagKeyWithType {
	return getOSMTagsWithTypesFromMapboxGLTypes(sourceLayer, scc.operands, mapMapboxGLSubclassToOSMTags)
}

func (scc SubClassCategory) OperandCount() int {
	return len(scc.operands)
}

type TypeCategory struct {
	operands []string
}

func (tc TypeCategory) Matches(operandIdx int, sourceLayer string, objectTags []*ownmap.OSMTag, objectType ownmap.ObjectType) bool {
	operand := tc.operands[operandIdx]
	switch operand {
	case ObjectTypePoint:
		return objectType == ownmap.ObjectTypeNode
	case ObjectTypePolygon, ObjectTypeLineString:
		return objectType == ownmap.ObjectTypeWay
	}
	panic("unhandled type category: " + operand)
}

func (TypeCategory) GetAllPossibleTagKeys(sourceLayer string) []*ownmapdal.TagKeyWithType {
	// TODO, do we need to get all tags relating to the type? Or are they picked up in other filters?
	return nil
}

func (tc TypeCategory) OperandCount() int {
	return len(tc.operands)
}

type UnhandledCategory struct{}

func (UnhandledCategory) Matches(operandIdx int, sourceLayer string, objectTags []*ownmap.OSMTag, objectType ownmap.ObjectType) bool {
	return false
}

func (UnhandledCategory) GetAllPossibleTagKeys(sourceLayer string) []*ownmapdal.TagKeyWithType {
	return nil
}
func (cc UnhandledCategory) OperandCount() int {
	return 0
}

func (f *ComparativeFilter) Type() FilterType {
	return FilterTypeComparative
}

func (f *ComparativeFilter) Operator() string {
	return f.operator
}

func (f *ComparativeFilter) IsObjectShown(sourceLayer string, tags []*ownmap.OSMTag, objectFilterType ownmap.ObjectType) bool {
	switch f.operator {
	case "==":
		return f.category.Matches(0, sourceLayer, tags, objectFilterType)
	case "!=":
		return !f.category.Matches(0, sourceLayer, tags, objectFilterType)
	case "in":
		for operandIdx := 0; operandIdx < f.category.OperandCount(); operandIdx++ {
			if f.category.Matches(operandIdx, sourceLayer, tags, objectFilterType) {
				return true
			}
		}
		return false
	case "!in":
		for operandIdx := 0; operandIdx < f.category.OperandCount(); operandIdx++ {
			if f.category.Matches(operandIdx, sourceLayer, tags, objectFilterType) {
				return false
			}
		}
		return true
	case "all":
		for operandIdx := 0; operandIdx < f.category.OperandCount(); operandIdx++ {
			if !f.category.Matches(operandIdx, sourceLayer, tags, objectFilterType) {
				return false
			}
		}
		return true
	default:
		log.Printf("WARN: unhandled operator: %q\n", f.operator)
		return false
	}
}

func (f *ComparativeFilter) GetAllPossibleTagKeys(sourceLayer string) []*ownmapdal.TagKeyWithType {
	return f.category.GetAllPossibleTagKeys(sourceLayer)
}

type FilterContainer struct {
	filter Filter
}

func (f *FilterContainer) UnmarshalJSON(data []byte) error {
	filter, err := unmarshalFilter(data)
	if err != nil {
		return errorsx.Wrap(err)
	}

	f.filter = filter
	return nil
}

func unmarshalFilter(data []byte) (Filter, error) {
	var fi []json.RawMessage
	err := json.Unmarshal(data, &fi)
	if err != nil {
		return nil, errorsx.Wrap(err)
	}

	var operator string
	err = json.Unmarshal(fi[0], &operator)
	if err != nil {
		return nil, errorsx.Wrap(err)
	}

	if len(fi) == 1 {
		// fill-in for possibly erronous filter of ["all"]
		return &LogicalFilter{
			operator: operator,
		}, nil
	}

	var categoryInterface interface{} // either string (comparative filter), or array (logical filter)
	err = json.Unmarshal(fi[1], &categoryInterface)
	if err != nil {
		return nil, errorsx.Wrap(err)
	}

	switch category := categoryInterface.(type) {
	case []interface{}:
		filter := &LogicalFilter{
			operator:   operator,
			subFilters: []Filter{}, // initialise an object for [] serialisation to JSON
		}
		for _, subFilterJSON := range fi[1:] {
			subFilter, err := unmarshalFilter(subFilterJSON)
			if err != nil {
				return nil, errorsx.Wrap(err)
			}
			filter.subFilters = append(filter.subFilters, subFilter)
		}
		return filter, nil
	case string:
		cf := new(ComparativeFilter)
		err = json.Unmarshal(data, cf)
		if err != nil {
			return nil, errorsx.Wrap(err)
		}
		return cf, nil
	default:
		panic(fmt.Sprintf("not implemented: %s :: %T :: %s", category, category, data))
	}
}

func (f FilterContainer) IsObjectShown(sourceLayer string, tags []*ownmap.OSMTag, objectType ownmap.ObjectType) bool {
	if f.filter == nil {
		return true
	}

	return f.filter.IsObjectShown(sourceLayer, tags, objectType)
}

func (f FilterContainer) GetTagKeysToFetch(sourceLayer string) []*ownmapdal.TagKeyWithType {

	if f.filter == nil {
		// if no filter, return no keys to fetch
		// this seems to happen on "building" layers
		return nil
	}
	return f.filter.GetAllPossibleTagKeys(sourceLayer)
}

func isSubclassShown(subclass, sourceLayer string, tags []*ownmap.OSMTag) bool {
	return isClassTypeShown(subclass, sourceLayer, tags, mapMapboxGLSubclassToOSMTags)
}

func isClassShown(className, sourceLayer string, tags []*ownmap.OSMTag) bool {
	return isClassTypeShown(className, sourceLayer, tags, mapMapboxGLClassToOSMTags)
}

// https://docs.mapbox.com/vector-tiles/reference/mapbox-streets-v8/
func isClassTypeShown(
	className,
	sourceLayer string,
	objectTags []*ownmap.OSMTag,
	mapperFunc func(className, sourceLayer string) []*ownmap.OSMTag,
) bool {
	lookingForOsmTags := mapperFunc(className, sourceLayer)
	for _, needleTag := range lookingForOsmTags {
		for _, objectTag := range objectTags {
			if objectTag.Key == needleTag.Key {
				if needleTag.Value == "*" || needleTag.Value == objectTag.Value {
					return true
				}
			}
		}
	}
	return false
}

type mapboxGLToOSMTagsMapperFuncType = func(className string, sourceLayer string) []*ownmap.OSMTag

func getOSMTagsWithTypesFromMapboxGLTypes(sourceLayer string, operands []string, mapperFunc mapboxGLToOSMTagsMapperFuncType) []*ownmapdal.TagKeyWithType {
	var tagsWithTypes []*ownmapdal.TagKeyWithType
	for _, operand := range operands {
		tags := mapperFunc(operand, sourceLayer)
		for _, tag := range tags {
			types := getPossibleObjectTypesForOSMTags(tag)
			for _, objectType := range types {
				tagsWithTypes = append(tagsWithTypes, &ownmapdal.TagKeyWithType{
					ObjectType: objectType,
					TagKey:     ownmap.TagKey(tag.Key),
				})
			}
		}
	}
	return tagsWithTypes
}
