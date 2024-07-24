package typecheck

import "github.com/FabianRolfMatthiasNoll/Golang-CSharp-Compiler/src/ast"

// TODO: Implement rest of check expr but with some sort of structure to control this monster of code
func (tc *TypeChecker) CheckExpr(expr ast.Expr) ast.TypedExpr {
	switch e := expr.(type) {
	case ast.IntLiteralExpr:
		return ast.TypedExpr{Type: "int", Expr: e}
	case ast.BoolLiteralExpr:
		return ast.TypedExpr{Type: "bool", Expr: e}
	case ast.StringExpr:
		return ast.TypedExpr{Type: "string", Expr: e}
	case ast.IdentifierExpr:
		return tc.CheckIdentifierExpr(e)
	case ast.NullLiteralExpr:
		return ast.TypedExpr{Type: "null", Expr: e}
	case ast.CharLiteralExpr:
		return ast.TypedExpr{Type: "char", Expr: e}
	case ast.BinaryExpr:
		return tc.CheckBinaryExpr(e)
	case ast.MethodCallExpr:
		return tc.CheckMethodCallExpr(e)
	case ast.AssignmentExpr:
		assigneeType := tc.CheckExpr(e.Assignee)
		valueType := tc.CheckExpr(e.Value)
		if assigneeType != valueType {
			tc.errorf(e.Line, e.Column, "type mismatch: %s and %s", assigneeType, valueType)
		}
		return assigneeType
	case ast.PreDecrementExpr:
		return tc.CheckUnaryExpr(e)
	case ast.PreIncrementExpr:
		return tc.CheckUnaryExpr(e)
	case ast.PostDecrementExpr:
		return tc.CheckUnaryExpr(e)
	case ast.PostIncrementExpr:
		return tc.CheckUnaryExpr(e)
	default:
		tc.errorf(expr.GetLine(), expr.GetColumn(), "unexpected expression")
	}
	return ast.TypedExpr{Type: "void"}
}

func (tc *TypeChecker) CheckBinaryExpr(expr ast.BinaryExpr) ast.TypedExpr {
	expr.Left = tc.CheckExpr(expr.Left)
	expr.Right = tc.CheckExpr(expr.Right)
	if !tc.isBinaryCompatible(expr.Left.(ast.TypedExpr).Type, expr.Right.(ast.TypedExpr).Type) {
		tc.errorf(expr.Line, expr.Column, "type mismatch during binary expression: %s and %s", expr.Left.(ast.TypedExpr).Type, expr.Right.(ast.TypedExpr).Type)
	}
	return ast.TypedExpr{Expr: expr, Type: "bool"}
}

func (tc *TypeChecker) CheckMethodCallExpr(expr ast.MethodCallExpr) ast.TypedExpr {
	// TODO: Implement
	return ast.TypedExpr{Type: "void"}
}

func (tc *TypeChecker) CheckIdentifierExpr(expr ast.IdentifierExpr) ast.TypedExpr {
	info, ok := tc.env.Lookup(expr.Name)
	if !ok {
		tc.errorf(expr.Line, expr.Column, "undefined variable: %s", expr.Name)
	}
	if info.IsField || info.IsGlobal {
		return ast.TypedExpr{Type: info.Type, Expr: ast.FieldVarExpr(expr)}
	} else {
		return ast.TypedExpr{Type: info.Type, Expr: ast.LocalVarExpr(expr)}
	}
}

func (tc *TypeChecker) CheckUnaryExpr(expr ast.Expr) ast.TypedExpr {
	// TODO: Implement
	// switch e := expr.(type) {
	// default:
	// 	tc.errorf(e.GetLine(), e.GetColumn(), "unexpected unary expression")
	// }
	return ast.TypedExpr{Type: "void"}
}

func (tc *TypeChecker) checkBoolCondition(condition ast.Expr) ast.TypedExpr {
	condition = tc.CheckExpr(condition)

	if condition.(ast.TypedExpr).Type != "bool" {
		tc.errorf(condition.GetLine(), condition.GetColumn(), "type mismatch: expected boolean, got %s", condition.(ast.TypedExpr).Type)
	}

	return condition.(ast.TypedExpr)
}
