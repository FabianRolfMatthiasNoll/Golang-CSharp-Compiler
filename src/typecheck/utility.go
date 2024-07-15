package typecheck

import "fmt"

func (tc *TypeChecker) isTypeCompatible(a, b string) bool {
	if a == "char" && b == "int" {
		return true
	} else if a == "int" && b == "char" {
		return true
	} else if tc.isUserObject(a) && b == "null" {
		return true
	} else if a == "null" && tc.isUserObject(b) {
		return true
	} else if a == b {
		return true
	}
	return false
}

func (tc *TypeChecker) isUserObject(typ string) bool {
	_, ok := tc.classes[typ]
	return ok
}

// Helper function to find the upper bound of a list of types
func (tc *TypeChecker) upperBound(types []string) string {
	if len(types) == 0 {
		return "void"
	}
	if len(types) == 1 {
		return types[0]
	}

	// Example of a type hierarchy for determining upper bounds
	typeHierarchy := map[string]int{
		"int":     1,
		"float":   2,
		"string":  3,
		"object":  4,
		"null":    5,
		"boolean": 6,
		"char":    7,
	}

	maxRank := 0
	upperType := "void"

	for _, t := range types {
		if rank, exists := typeHierarchy[t]; exists && rank > maxRank {
			maxRank = rank
			upperType = t
		}
	}

	return upperType
}

func (tc *TypeChecker) errorf(line, column int, format string, args ...interface{}) {
	message := fmt.Sprintf(format, args...)
	panic(fmt.Sprintf("Error at line %d, column %d: %s", line, column, message))
}
