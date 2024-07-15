package typecheck

type SymbolInfo struct {
	Type        string
	IsGlobal    bool
	IsField     bool
	IsParameter bool
}

type TypeEnvironment struct {
	symbols map[string]SymbolInfo
	outer   *TypeEnvironment
}

func NewTypeEnv(outer *TypeEnvironment) *TypeEnvironment {
	return &TypeEnvironment{
		symbols: make(map[string]SymbolInfo),
		outer:   outer,
	}
}

func (env *TypeEnvironment) Lookup(name string) (SymbolInfo, bool) {
	info, ok := env.symbols[name]
	if !ok && env.outer != nil {
		return env.outer.Lookup(name)
	}
	return info, ok
}

func (env *TypeEnvironment) Define(name, typ string, isGlobal, isField, isParameter bool) {
	env.symbols[name] = SymbolInfo{Type: typ, IsGlobal: isGlobal, IsField: isField, IsParameter: isParameter}
}
