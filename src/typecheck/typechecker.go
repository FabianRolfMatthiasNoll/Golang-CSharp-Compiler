package typecheck

import (
	"github.com/FabianRolfMatthiasNoll/Golang-CSharp-Compiler/src/ast"
)

type TypeChecker struct {
	env     *TypeEnvironment
	classes map[string]ast.ClassDeclStmt
}

func NewTypeChecker() *TypeChecker {
	return &TypeChecker{env: NewTypeEnv(nil)}
}

func (tc *TypeChecker) CheckProgram(prog *ast.Program) ast.Program {
	tc.classes = make(map[string]ast.ClassDeclStmt)
	for i := range prog.Classes {
		tc.classes[prog.Classes[i].Name] = prog.Classes[i]
	}

	for i := range prog.Classes {
		tc.CheckClassDeclStmt(&prog.Classes[i])
	}

	return *prog
}






