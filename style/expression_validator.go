package style

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"
)

// ValidationMode controls how strictly expressions are validated.
type ValidationMode int

const (
	// ValidationBasic checks operator names and basic structure only.
	ValidationBasic ValidationMode = iota
	// ValidationStrict checks operator names, argument counts, and type constraints.
	ValidationStrict
)

// ValidateExpression checks whether an expression tree is well-formed.
// It recursively validates all sub-expressions.
func ValidateExpression(e *Expression, mode ValidationMode) error {
	if e == nil {
		return errors.New("expression is nil")
	}
	if e.IsLiteral {
		return nil
	}
	if e.Operator == "" {
		return errors.New("expression operator is empty")
	}
	if !IsKnownOperator(e.Operator) {
		return errors.Errorf("unknown expression operator: %q", e.Operator)
	}

	if mode >= ValidationStrict {
		if err := validateArgs(e); err != nil {
			return errors.Wrapf(err, "expression %q", e.Operator)
		}
	}

	for i, arg := range e.Args {
		if arg == nil {
			continue
		}
		if err := ValidateExpression(arg, mode); err != nil {
			return errors.Wrapf(err, "arg %d of %q", i, e.Operator)
		}
	}
	return nil
}

func validateArgs(e *Expression) error {
	switch e.Operator {
	// Types
	case ExpArray:
		// ["array", value] or ["array", type, value] or ["array", type, N, value]
		if len(e.Args) < 1 {
			return errors.New("requires at least 1 argument")
		}

	case ExpBoolean:
		// ["boolean", value, fallback?, ...]
		if len(e.Args) < 1 {
			return errors.New("requires at least 1 argument")
		}

	case ExpNumber:
		if len(e.Args) < 1 {
			return errors.New("requires at least 1 argument")
		}

	case ExpString:
		if len(e.Args) < 1 {
			return errors.New("requires at least 1 argument")
		}

	case ExpObject:
		if len(e.Args) < 1 {
			return errors.New("requires at least 1 argument")
		}

	case ExpToString, ExpToBool, ExpTypeOf:
		if len(e.Args) != 1 {
			return errors.Errorf("requires 1 argument, got %d", len(e.Args))
		}

	case ExpToNumber:
		if len(e.Args) < 1 {
			return errors.New("requires at least 1 argument")
		}

	case ExpToColor:
		if len(e.Args) < 1 {
			return errors.New("requires at least 1 argument")
		}

	case ExpLiteral:
		if len(e.Args) != 1 {
			return errors.Errorf("requires 1 argument, got %d", len(e.Args))
		}

	case ExpImage:
		if len(e.Args) < 1 {
			return errors.New("requires at least 1 argument (image name)")
		}

	case ExpCollator:
		if len(e.Args) > 1 {
			return errors.Errorf("requires 0 or 1 argument (options object), got %d", len(e.Args))
		}

	case ExpNumberFmt:
		if len(e.Args) < 1 || len(e.Args) > 2 {
			return errors.Errorf("requires 1-2 arguments, got %d", len(e.Args))
		}

	// Feature data
	case ExpAccumulated:
		if len(e.Args) != 0 {
			return errors.Errorf("requires 0 arguments, got %d", len(e.Args))
		}

	case ExpFeatureState:
		if len(e.Args) != 1 {
			return errors.Errorf("requires 1 argument (property name), got %d", len(e.Args))
		}

	case ExpGeometryType:
		if len(e.Args) != 0 {
			return errors.Errorf("requires 0 arguments, got %d", len(e.Args))
		}

	case ExpID:
		if len(e.Args) != 0 {
			return errors.Errorf("requires 0 arguments, got %d", len(e.Args))
		}

	case ExpLineProgress:
		if len(e.Args) != 0 {
			return errors.Errorf("requires 0 arguments, got %d", len(e.Args))
		}

	case ExpProperties:
		if len(e.Args) != 0 {
			return errors.Errorf("requires 0 arguments, got %d", len(e.Args))
		}

	case ExpHeatmapDensity:
		if len(e.Args) != 0 {
			return errors.Errorf("requires 0 arguments, got %d", len(e.Args))
		}

	// Lookup
	case ExpGet:
		if len(e.Args) < 1 || len(e.Args) > 2 {
			return errors.Errorf("requires 1-2 arguments, got %d", len(e.Args))
		}

	case ExpHas:
		if len(e.Args) < 1 || len(e.Args) > 2 {
			return errors.Errorf("requires 1-2 arguments, got %d", len(e.Args))
		}

	case ExpAt:
		if len(e.Args) != 2 {
			return errors.Errorf("requires 2 arguments (index, array), got %d", len(e.Args))
		}

	case ExpAtInterpolated:
		if len(e.Args) != 2 {
			return errors.Errorf("requires 2 arguments (index, array), got %d", len(e.Args))
		}

	case ExpLength:
		if len(e.Args) != 1 {
			return errors.Errorf("requires 1 argument, got %d", len(e.Args))
		}

	case ExpIn:
		if len(e.Args) != 2 {
			return errors.Errorf("requires 2 arguments (keyword, input), got %d", len(e.Args))
		}

	case ExpIndexOf:
		if len(e.Args) < 2 || len(e.Args) > 3 {
			return errors.Errorf("requires 2-3 arguments, got %d", len(e.Args))
		}

	case ExpSlice:
		if len(e.Args) < 2 || len(e.Args) > 3 {
			return errors.Errorf("requires 2-3 arguments, got %d", len(e.Args))
		}

	case ExpConfig:
		if len(e.Args) < 1 || len(e.Args) > 2 {
			return errors.Errorf("requires 1-2 arguments, got %d", len(e.Args))
		}

	case ExpMeasureLight:
		if len(e.Args) != 1 {
			return errors.Errorf("requires 1 argument, got %d", len(e.Args))
		}

	case ExpSplit:
		if len(e.Args) != 2 {
			return errors.Errorf("requires 2 arguments (string, delimiter), got %d", len(e.Args))
		}

	case ExpWorldview:
		if len(e.Args) != 0 {
			return errors.Errorf("requires 0 arguments, got %d", len(e.Args))
		}

	// Decision
	case ExpNot:
		if len(e.Args) != 1 {
			return errors.Errorf("requires 1 argument, got %d", len(e.Args))
		}

	case ExpNEq, ExpEQ:
		if len(e.Args) < 2 || len(e.Args) > 3 {
			return errors.Errorf("requires 2-3 arguments (value, value, collator?), got %d", len(e.Args))
		}

	case ExpLT, ExpLTE, ExpGT, ExpGTE:
		if len(e.Args) < 2 || len(e.Args) > 3 {
			return errors.Errorf("requires 2-3 arguments (value, value, collator?), got %d", len(e.Args))
		}

	case ExpAll:
		if len(e.Args) < 2 {
			return errors.Errorf("requires at least 2 arguments, got %d", len(e.Args))
		}

	case ExpAny:
		if len(e.Args) < 2 {
			return errors.Errorf("requires at least 2 arguments, got %d", len(e.Args))
		}

	case ExpCase:
		if len(e.Args) < 3 {
			return errors.Errorf("requires at least 3 arguments (cond, output, ..., fallback), got %d", len(e.Args))
		}
		if len(e.Args)%2 == 0 {
			return errors.Errorf("expected odd number of arguments (condition/output pairs + fallback), got %d", len(e.Args))
		}

	case ExpCoalesce:
		if len(e.Args) < 1 {
			return errors.New("requires at least 1 argument")
		}

	case ExpMatch:
		if len(e.Args) < 3 {
			return errors.Errorf("requires at least 3 arguments (input, label, output, ..., fallback), got %d", len(e.Args))
		}

	case ExpWithin:
		if len(e.Args) != 1 {
			return errors.Errorf("requires 1 argument (GeoJSON object), got %d", len(e.Args))
		}

	// Ramps
	case ExpInterpolate, ExpInterpolateHCL, ExpInterpolateLab:
		if len(e.Args) < 4 {
			return errors.Errorf("requires at least 4 arguments (interpolation, input, stop, output, ...), got %d", len(e.Args))
		}
		interp := e.Args[0]
		if interp != nil && !interp.IsLiteral {
			op := interp.Operator
			if op != ExpLinear && op != ExpExponential && op != ExpCubicBezier {
				return errors.Errorf("first argument must be a interpolation type (linear/exponential/cubic-bezier), got %q", op)
			}
		}
		if (len(e.Args)-2)%2 != 0 {
			return errors.Errorf("expected stop/output pairs after input, got %d arguments after input", len(e.Args)-2)
		}

	case ExpStep:
		if len(e.Args) < 3 {
			return errors.Errorf("requires at least 3 arguments (input, output, stop, output, ...), got %d", len(e.Args))
		}

	// Variable binding
	case ExpLet:
		if len(e.Args) < 3 {
			return errors.Errorf("requires at least 3 arguments (name, value, ..., result), got %d", len(e.Args))
		}
		if (len(e.Args)-1)%2 != 0 {
			return errors.Errorf("expected result expression after name/value pairs, got %d arguments", len(e.Args))
		}

	case ExpVar:
		if len(e.Args) != 1 {
			return errors.Errorf("requires 1 argument (variable name), got %d", len(e.Args))
		}

	// String
	case ExpConcat:
		if len(e.Args) < 1 {
			return errors.New("requires at least 1 argument")
		}

	case ExpDowncase, ExpUpcase:
		if len(e.Args) != 1 {
			return errors.Errorf("requires 1 argument, got %d", len(e.Args))
		}

	case ExpIsSupportedScript:
		if len(e.Args) != 1 {
			return errors.Errorf("requires 1 argument, got %d", len(e.Args))
		}

	case ExpResolvedLocale:
		if len(e.Args) != 1 {
			return errors.Errorf("requires 1 argument (collator), got %d", len(e.Args))
		}

	// Color
	case ExpRGB:
		if len(e.Args) != 3 {
			return errors.Errorf("requires 3 arguments (r, g, b), got %d", len(e.Args))
		}

	case ExpRGBA:
		if len(e.Args) != 4 {
			return errors.Errorf("requires 4 arguments (r, g, b, a), got %d", len(e.Args))
		}

	case ExpHSL:
		if len(e.Args) != 3 {
			return errors.Errorf("requires 3 arguments (h, s, l), got %d", len(e.Args))
		}

	case ExpHSLA:
		if len(e.Args) != 4 {
			return errors.Errorf("requires 4 arguments (h, s, l, a), got %d", len(e.Args))
		}

	// Math
	case ExpAdd, ExpMul, ExpSub, ExpDiv:
		if len(e.Args) < 2 {
			return errors.Errorf("requires at least 2 arguments, got %d", len(e.Args))
		}

	case ExpMod, ExpPow:
		if len(e.Args) != 2 {
			return errors.Errorf("requires 2 arguments, got %d", len(e.Args))
		}

	case ExpAbs, ExpAcos, ExpAsin, ExpAtan, ExpCeil, ExpCos,
		ExpFloor, ExpLn, ExpLog10, ExpLog2, ExpRound, ExpSin,
		ExpSqrt, ExpTan:
		if len(e.Args) != 1 {
			return errors.Errorf("requires 1 argument, got %d", len(e.Args))
		}

	case ExpDist:
		if len(e.Args) != 2 {
			return errors.Errorf("requires 2 arguments, got %d", len(e.Args))
		}

	case ExpMax, ExpMin:
		if len(e.Args) < 2 {
			return errors.Errorf("requires at least 2 arguments, got %d", len(e.Args))
		}

	case ExpRand:
		if len(e.Args) < 2 || len(e.Args) > 3 {
			return errors.Errorf("requires 2-3 arguments, got %d", len(e.Args))
		}

	case ExpE, ExpPI:
		if len(e.Args) != 0 {
			return errors.Errorf("requires 0 arguments, got %d", len(e.Args))
		}

	// Camera
	case ExpDistFromCenter:
		if len(e.Args) != 0 {
			return errors.Errorf("requires 0 arguments, got %d", len(e.Args))
		}

	case ExpPitch:
		if len(e.Args) != 0 {
			return errors.Errorf("requires 0 arguments, got %d", len(e.Args))
		}

	case ExpZoom:
		if len(e.Args) != 0 {
			return errors.Errorf("requires 0 arguments, got %d", len(e.Args))
		}

	// Interpolation sub-types (should not appear as top-level operators)
	case ExpLinear:
		if len(e.Args) != 0 && !(len(e.Args) == 0) {
			return errors.Errorf("unexpected arguments for linear interpolation")
		}

	case ExpExponential:
		if len(e.Args) != 1 {
			return errors.Errorf("exponential interpolation requires 1 argument (base), got %d", len(e.Args))
		}

	case ExpCubicBezier:
		if len(e.Args) != 4 {
			return errors.Errorf("cubic-bezier interpolation requires 4 arguments (x1, y1, x2, y2), got %d", len(e.Args))
		}

	case ExpToHSLA, ExpToRGBA:
		if len(e.Args) != 1 {
			return errors.Errorf("requires 1 argument, got %d", len(e.Args))
		}
	}

	return nil
}

// IsCameraExpression returns true if the expression or any sub-expression
// uses the zoom operator at the top level (for camera expressions).
func IsCameraExpression(e *Expression) bool {
	if e == nil || e.IsLiteral {
		return false
	}
	if e.Operator == ExpZoom {
		return true
	}
	for _, arg := range e.Args {
		if IsCameraExpression(arg) {
			return true
		}
	}
	return false
}

// IsDataExpression returns true if the expression or any sub-expression
// uses a feature data operator (get, has, id, geometry-type, properties, feature-state).
func IsDataExpression(e *Expression) bool {
	if e == nil || e.IsLiteral {
		return false
	}
	switch e.Operator {
	case ExpGet, ExpHas, ExpID, ExpGeometryType, ExpProperties, ExpFeatureState:
		return true
	}
	for _, arg := range e.Args {
		if IsDataExpression(arg) {
			return true
		}
	}
	return false
}

// FormatExpression returns a human-readable string representation of an expression.
func FormatExpression(e *Expression) string {
	if e == nil {
		return "null"
	}
	if e.IsLiteral {
		return fmt.Sprintf("%v", e.Value)
	}
	var args []string
	for _, arg := range e.Args {
		args = append(args, FormatExpression(arg))
	}
	return fmt.Sprintf("[%s %s]", e.Operator, strings.Join(args, " "))
}
