package mbstyle

import (
	"github.com/flywave/flywave-server/formats/mapbox"
)

const (
	CoalesceExpression              = 0
	CompoundExpression              = 1
	LiteralExpression               = 2
	AtExpression                    = 3
	InterpolateExpression           = 4
	AssertionExpression             = 5
	LengthExpression                = 6
	StepExpression                  = 7
	LetExpression                   = 8
	VarExpression                   = 9
	CollatorExpression              = 10
	CoercionExpression              = 11
	MatchExpression                 = 12
	ErrorExpression                 = 13
	CaseExpression                  = 14
	AnyExpression                   = 15
	AllExpression                   = 16
	ComparisonExpression            = 17
	FormatExpressionExpression      = 18
	FormatSectionOverrideExpression = 19
	NumberFormatExpression          = 20
)

const (
	NullType      = "null"
	NumberType    = "number"
	BooleanType   = "boolean"
	StringType    = "string"
	ColorType     = "color"
	ObjectType    = "object"
	ErrorType     = "error"
	ValueType     = "value"
	CollatorType  = "collator"
	FormattedType = "formatted"
)

type EvaluationContext struct {
	Zoom               *float64
	Accumulated        interface{}
	Feature            *mapbox.GeometryTileFeature
	ColorRampParameter *float64
	FormattedSection   map[string]interface{}
}

type Expression interface {
	Evaluate(ctx *EvaluationContext) (interface{}, error)
	EachChild(fun *func(*Expression))
	GetKind() int
	GetType() string
	GetOperator() string
}
