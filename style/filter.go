package style

import (
	"encoding/json"
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

// FilterContainer wraps an expression used as a layer filter.
// A filter is simply an expression that evaluates to a boolean.
type FilterContainer struct {
	Expr *Expression
}

func (f *FilterContainer) UnmarshalJSON(data []byte) error {
	expr := &Expression{}
	if err := expr.UnmarshalJSON(data); err != nil {
		return err
	}
	f.Expr = expr
	return nil
}

func (f *FilterContainer) MarshalJSON() ([]byte, error) {
	if f.Expr == nil {
		return json.Marshal(nil)
	}
	return f.Expr.MarshalJSON()
}
