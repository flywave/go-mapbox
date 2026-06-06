package style

// Expression operator names.
const (
	// Types
	ExpArray       = "array"
	ExpBoolean     = "boolean"
	ExpCollator    = "collator"
	ExpFormat      = "format"
	ExpImage       = "image"
	ExpLiteral     = "literal"
	ExpNumber      = "number"
	ExpNumberFmt   = "number-format"
	ExpObject      = "object"
	ExpString      = "string"
	ExpToBool      = "to-boolean"
	ExpToColor     = "to-color"
	ExpToNumber    = "to-number"
	ExpToString    = "to-string"
	ExpTypeOf      = "typeof"
	ExpToHSLA      = "to-hsla"
	ExpToRGBA      = "to-rgba"

	// Feature data
	ExpAccumulated    = "accumulated"
	ExpFeatureState   = "feature-state"
	ExpGeometryType   = "geometry-type"
	ExpID             = "id"
	ExpLineProgress   = "line-progress"
	ExpProperties     = "properties"

	// Lookup
	ExpAt             = "at"
	ExpAtInterpolated = "at-interpolated"
	ExpConfig         = "config"
	ExpGet            = "get"
	ExpHas            = "has"
	ExpIn             = "in"
	ExpIndexOf        = "index-of"
	ExpLength         = "length"
	ExpMeasureLight   = "measure-light"
	ExpSlice          = "slice"
	ExpSplit          = "split"
	ExpWorldview      = "worldview"

	// Decision
	ExpNot    = "!"
	ExpNEq    = "!="
	ExpLT     = "<"
	ExpLTE    = "<="
	ExpEQ     = "=="
	ExpGT     = ">"
	ExpGTE    = ">="
	ExpAll    = "all"
	ExpAny    = "any"
	ExpCase   = "case"
	ExpCoalesce = "coalesce"
	ExpMatch  = "match"
	ExpWithin = "within"

	// Ramps, scales, curves
	ExpInterpolate    = "interpolate"
	ExpInterpolateHCL = "interpolate-hcl"
	ExpInterpolateLab = "interpolate-lab"
	ExpStep           = "step"

	// Variable binding
	ExpLet = "let"
	ExpVar = "var"

	// String
	ExpConcat           = "concat"
	ExpDowncase         = "downcase"
	ExpIsSupportedScript = "is-supported-script"
	ExpResolvedLocale   = "resolved-locale"
	ExpUpcase           = "upcase"

	// Color
	ExpHSL  = "hsl"
	ExpHSLA = "hsla"
	ExpRGB  = "rgb"
	ExpRGBA = "rgba"

	// Math
	ExpSub = "-"
	ExpMul = "*"
	ExpDiv = "/"
	ExpMod = "%"
	ExpPow = "^"
	ExpAdd = "+"
	ExpAbs  = "abs"
	ExpAcos = "acos"
	ExpAsin = "asin"
	ExpAtan = "atan"
	ExpCeil = "ceil"
	ExpCos  = "cos"
	ExpDist = "distance"
	ExpE    = "e"
	ExpFloor = "floor"
	ExpLn   = "ln"
	ExpLn2  = "ln2"
	ExpLog10 = "log10"
	ExpLog2 = "log2"
	ExpMax  = "max"
	ExpMin  = "min"
	ExpPI   = "pi"
	ExpRand = "random"
	ExpRound = "round"
	ExpSin  = "sin"
	ExpSqrt = "sqrt"
	ExpTan  = "tan"

	// Camera
	ExpDistFromCenter = "distance-from-center"
	ExpPitch          = "pitch"
	ExpZoom           = "zoom"

	// Heatmap
	ExpHeatmapDensity = "heatmap-density"

	// Interpolation types (sub-operators for interpolate)
	ExpLinear     = "linear"
	ExpExponential = "exponential"
	ExpCubicBezier = "cubic-bezier"
)

// knownExpressions contains all valid expression operators mapped to their
// expected argument patterns for validation.
var knownExpressions = map[string]struct{}{
	// Types
	ExpArray: {}, ExpBoolean: {}, ExpCollator: {}, ExpFormat: {},
	ExpImage: {}, ExpLiteral: {}, ExpNumber: {}, ExpNumberFmt: {},
	ExpObject: {}, ExpString: {}, ExpToBool: {}, ExpToColor: {},
	ExpToNumber: {}, ExpToString: {}, ExpTypeOf: {}, ExpToHSLA: {}, ExpToRGBA: {},

	// Feature data
	ExpAccumulated: {}, ExpFeatureState: {}, ExpGeometryType: {},
	ExpID: {}, ExpLineProgress: {}, ExpProperties: {},

	// Lookup
	ExpAt: {}, ExpAtInterpolated: {}, ExpConfig: {}, ExpGet: {},
	ExpHas: {}, ExpIn: {}, ExpIndexOf: {}, ExpLength: {},
	ExpMeasureLight: {}, ExpSlice: {}, ExpSplit: {}, ExpWorldview: {},

	// Decision
	ExpNot: {}, ExpNEq: {}, ExpLT: {}, ExpLTE: {}, ExpEQ: {},
	ExpGT: {}, ExpGTE: {}, ExpAll: {}, ExpAny: {}, ExpCase: {},
	ExpCoalesce: {}, ExpMatch: {}, ExpWithin: {},

	// Ramps, scales, curves
	ExpInterpolate: {}, ExpInterpolateHCL: {}, ExpInterpolateLab: {},
	ExpStep: {},

	// Variable binding
	ExpLet: {}, ExpVar: {},

	// String
	ExpConcat: {}, ExpDowncase: {}, ExpIsSupportedScript: {},
	ExpResolvedLocale: {}, ExpUpcase: {},

	// Color
	ExpHSL: {}, ExpHSLA: {}, ExpRGB: {}, ExpRGBA: {},

	// Math
	ExpSub: {}, ExpMul: {}, ExpDiv: {}, ExpMod: {}, ExpPow: {},
	ExpAdd: {}, ExpAbs: {}, ExpAcos: {}, ExpAsin: {}, ExpAtan: {},
	ExpCeil: {}, ExpCos: {}, ExpDist: {}, ExpE: {}, ExpFloor: {},
	ExpLn: {}, ExpLn2: {}, ExpLog10: {}, ExpLog2: {}, ExpMax: {},
	ExpMin: {}, ExpPI: {}, ExpRand: {}, ExpRound: {}, ExpSin: {},
	ExpSqrt: {}, ExpTan: {},

	// Camera
	ExpDistFromCenter: {}, ExpPitch: {}, ExpZoom: {},

	// Heatmap
	ExpHeatmapDensity: {},

	// Interpolation sub-types
	ExpLinear: {}, ExpExponential: {}, ExpCubicBezier: {},
}

// IsKnownOperator returns true if the given string is a valid expression operator.
func IsKnownOperator(op string) bool {
	_, ok := knownExpressions[op]
	return ok
}
