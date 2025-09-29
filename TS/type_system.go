package TS

import (
	"ion-go/Token"
)

type TypeKind string

const (
	INVALID_TYPE TypeKind = "INVALID_TYPE"
	VOID                  = "void"
	INTEGER               = "int"
	FLOAT                 = "float"
	BOOL                  = "bool"
	STRING                = "string"
	ARRAY                 = "[]"
	STRUCT                = ""
	POINTER               = "*"
	FUNCTION              = "fn(...) -> "
)

type Parameter struct {
	Tok      Token.Token
	DeclType *Type
}

type Type struct {
	Kind       TypeKind
	Next       *Type // For Functions the return type is the last node in the next chain
	Parameters []Parameter
}

func NewType(kind TypeKind, next *Type, parameters []Parameter) *Type {
	if kind != FUNCTION && parameters != nil {
		panic("Attempted to give parameters to a non function type")
	}

	return &Type{
		Kind:       kind,
		Next:       next,
		Parameters: parameters,
	}
}

func (t *Type) IsPointer() bool {
	return t.Kind == POINTER
}

func (t *Type) IsArray() bool {
	return t.Kind == ARRAY
}

func (t *Type) IsStruct() bool {
	return t.Kind == STRUCT
}

func (t *Type) IsFunction() bool {
	return t.Kind == FUNCTION
}

func (t *Type) GetReturnType() *Type {
	if t.Kind != FUNCTION {
		panic("Return type is not a function")
	}

	return t.Next
}

func (t *Type) GetArrayDimensions() int {
	if t.Kind != ARRAY {
		panic("Expected ARRAY type")
	}

	dimensions := 0
	current := t
	for current.IsArray() {
		dimensions++
		current = current.Next
	}

	return dimensions
}

func (t *Type) SetBaseType(kind TypeKind) {
	current := t
	for current.Next != nil {
		current = current.Next
	}

	current.Next = NewType(kind, nil, nil)
}

func (t *Type) AddArrayModifier() *Type {
	current := NewType(ARRAY, t, nil)

	return current
}

func (t *Type) RemoveArrayModifier() *Type {
	if t.Kind != ARRAY {
		panic("Expected ARRAY type")
	}

	current := NewType(t.Kind, t.Next, t.Parameters)

	current.Kind = current.Next.Kind
	current.Next = current.Next.Next

	return current
}

func (t *Type) AddStructModifier() *Type {
	if t.Kind == STRUCT {
		panic("Type is already a struct")
	}

	current := NewType(STRUCT, t, nil)

	return current
}

func (t *Type) RemoveStructModifier() *Type {
	if t.Kind != STRUCT {
		panic("Expected ARRAY type")
	}

	current := NewType(t.Kind, t.Next, t.Parameters)

	current.Kind = current.Next.Kind
	current.Next = current.Next.Next

	return current
}

func (t *Type) String() string {
	ret := ""
	current := t
	for current != nil {
		ret += string(current.Kind)
		current = current.Next
	}

	return ret
}

// TypeCompare NOTE(Jovanni):
// Later on this might have like subtype and type group implications so it probably
// won't just be a bool it will be some type of int 0 is exact type 1 is super type, -1 is not equal
func TypeCompare(c1, c2 *Type) bool {
	if c1 == nil || c2 == nil {
		return false
	}

	for c1 != nil && c2 != nil {
		if c1.Kind != c2.Kind {
			return false
		} else if len(c1.Parameters) != len(c2.Parameters) {
			return false
		} else {
			for i := 0; i < len(c1.Parameters); i++ {
				param1 := c1.Parameters[i]
				param2 := c2.Parameters[i]
				if param1.Tok != param2.Tok {
					return false
				} else if !TypeCompare(param1.DeclType, param2.DeclType) {
					return false
				}
			}
		}

		c1 = c1.Next
		c2 = c2.Next
	}

	return true
}

type BinaryQuery struct {
	Op    string
	Left  TypeKind
	Right TypeKind
}

// GetPromotedType This is strictly for binary operations
func GetPromotedType(op Token.Token, leftType, rightType *Type) TypeKind {
	var typeMap = map[BinaryQuery]TypeKind{
		{"+", INTEGER, INTEGER}: INTEGER,
		{"+", INTEGER, FLOAT}:   FLOAT,
		{"+", FLOAT, INTEGER}:   FLOAT,
		{"+", FLOAT, FLOAT}:     FLOAT,

		{"-", INTEGER, INTEGER}: INTEGER,
		{"-", INTEGER, FLOAT}:   FLOAT,
		{"-", FLOAT, INTEGER}:   FLOAT,
		{"-", FLOAT, FLOAT}:     FLOAT,

		{"*", INTEGER, INTEGER}: INTEGER,
		{"*", INTEGER, FLOAT}:   FLOAT,
		{"*", FLOAT, INTEGER}:   FLOAT,
		{"*", FLOAT, FLOAT}:     FLOAT,

		{"/", INTEGER, INTEGER}: INTEGER,
		{"/", INTEGER, FLOAT}:   FLOAT,
		{"/", FLOAT, INTEGER}:   FLOAT,
		{"/", FLOAT, FLOAT}:     FLOAT,

		{"%", INTEGER, INTEGER}: INTEGER,

		{"+", STRING, STRING}:  STRING,
		{"+", STRING, INTEGER}: STRING,
		{"+", STRING, FLOAT}:   STRING,
		{"+", INTEGER, STRING}: STRING,
		{"+", FLOAT, STRING}:   STRING,

		{"==", INTEGER, INTEGER}: BOOL,
		{"==", INTEGER, FLOAT}:   BOOL,
		{"==", FLOAT, INTEGER}:   BOOL,
		{"==", FLOAT, FLOAT}:     BOOL,

		{"!=", INTEGER, INTEGER}: BOOL,
		{"!=", INTEGER, FLOAT}:   BOOL,
		{"!=", FLOAT, INTEGER}:   BOOL,
		{"!=", FLOAT, FLOAT}:     BOOL,

		{"||", BOOL, BOOL}: BOOL,
		{"&&", BOOL, BOOL}: BOOL,

		{"<", INTEGER, INTEGER}: BOOL,
		{"<", INTEGER, FLOAT}:   BOOL,
		{"<", FLOAT, INTEGER}:   BOOL,
		{"<", FLOAT, FLOAT}:     BOOL,

		{"<=", INTEGER, INTEGER}: BOOL,
		{"<=", INTEGER, FLOAT}:   BOOL,
		{"<=", FLOAT, INTEGER}:   BOOL,
		{"<=", FLOAT, FLOAT}:     BOOL,

		{">", INTEGER, INTEGER}: BOOL,
		{">", INTEGER, FLOAT}:   BOOL,
		{">", FLOAT, INTEGER}:   BOOL,
		{">", FLOAT, FLOAT}:     BOOL,

		{">=", INTEGER, INTEGER}: BOOL,
		{">=", INTEGER, FLOAT}:   BOOL,
		{">=", FLOAT, INTEGER}:   BOOL,
		{">=", FLOAT, FLOAT}:     BOOL,
	}
	query := BinaryQuery{
		Op:    op.Lexeme,
		Left:  leftType.Kind,
		Right: rightType.Kind,
	}

	if ret, ok := typeMap[query]; ok {
		return ret
	}

	return INVALID_TYPE
}

type TypeCastQuery struct {
	TypeCast TypeKind
	ExprType TypeKind
}

func CanCastType(castType *Type, exprType *Type) bool {
	// Trivial case
	if TypeCompare(castType, exprType) {
		return true
	}

	var castMap = map[TypeCastQuery]bool{
		{INTEGER, FLOAT}:  true,
		{FLOAT, INTEGER}:  true,
		{STRING, INTEGER}: true,
		{STRING, FLOAT}:   true,
	}

	query := TypeCastQuery{
		TypeCast: castType.Kind,
		ExprType: exprType.Kind,
	}

	if _, ok := castMap[query]; ok {
		return true
	}

	return false
}
