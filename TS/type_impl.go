package TS

import (
	"fmt"
	"ion-go/Token"
	"reflect"
	"unsafe"
)

func (b *BaseType) isType()                   {}
func (b *BaseType) String() string            { return "BaseType" }
func (b *BaseType) Underlying() Type          { return b.underlying }
func (b *BaseType) Size() int                 { return 0 }
func (b *BaseType) IsInteger() bool           { return false }
func (b *BaseType) IsString() bool            { return false }
func (b *BaseType) IsStruct() bool            { return false }
func (b *BaseType) IsPointer() bool           { return false }
func (b *BaseType) IsArray() bool             { return false }
func (b *BaseType) IsFixedSizeArray() bool    { return false }
func (b *BaseType) IsInferredSizeArray() bool { return false }
func (b *BaseType) IsFloat() bool             { return false }

func (b *BaseType) IsFunction() bool { return false }
func AddAlias(t Type, strict bool, name string) Type {
	return &AliasType{
		BaseType:  BaseType{t},
		Strict:    strict,
		AliasName: name,
	}
}

func AddStaticArray(t Type, count int) Type {
	return &StaticArrayType{
		BaseType: BaseType{t},
		Count:    count,
	}
}

func AddPointer(t Type) Type {
	return &PointerType{
		BaseType: BaseType{t},
	}
}

func RemoveModifier(t Type) Type {
	return AllocateDeepCopyType(t.Underlying())
}

func AllocateDeepCopyType(t Type) Type {
	if t == nil {
		return nil
	}

	switch v := t.(type) {
	case *CharType:
		return &CharType{
			BaseType: BaseType{AllocateDeepCopyType(t.Underlying())},
			Signed:   v.Signed,
		}
	case *IntegerType:
		return &IntegerType{
			BaseType: BaseType{AllocateDeepCopyType(t.Underlying())},
			Signed:   v.Signed,
			Bytes:    v.Bytes,
		}
	case *FloatType:
		return &FloatType{
			BaseType: BaseType{AllocateDeepCopyType(t.Underlying())},
			Bytes:    v.Bytes,
		}
	case *StaticArrayType:
		return &StaticArrayType{
			BaseType: BaseType{AllocateDeepCopyType(t.Underlying())},
			Count:    v.Count,
		}
	case *StructType:
		return &StructType{
			BaseType:   BaseType{AllocateDeepCopyType(t.Underlying())},
			StructName: v.StructName,
			Members:    v.Members,
		}
	case *FunctionType:
		return &FunctionType{
			BaseType:   BaseType{AllocateDeepCopyType(t.Underlying())},
			ReturnType: v.ReturnType,
			Params:     v.Params,
		}
	case *AliasType:
		return &AliasType{
			BaseType:  BaseType{AllocateDeepCopyType(t.Underlying())},
			Strict:    v.Strict,
			AliasName: v.AliasName,
		}
	case *PointerType:
		return &PointerType{
			BaseType{AllocateDeepCopyType(t.Underlying())},
		}
	default:
		panic(fmt.Sprintf("Failed to allocate a new type %T", v))
	}

	return nil
}

func (b *BaseType) RemoveModifier() Type {
	return AllocateDeepCopyType(b.Underlying())
}

func (v *VoidType) isType() {}
func (v *VoidType) String() string {
	return "void"
}

func (b *BoolType) isType()        {}
func (b *BoolType) Size() int      { return 1 }
func (b *BoolType) String() string { return "bool" }

func (c *CharType) isType()   {}
func (c *CharType) Size() int { return 1 }
func (c *CharType) String() string {
	u := "u"
	if c.Signed {
		u = ""
	}

	return fmt.Sprintf("%schar", u)
}

func (i *IntegerType) isType()         {}
func (i *IntegerType) Size() int       { return i.Bytes }
func (i *IntegerType) IsInteger() bool { return true }
func (i *IntegerType) String() string {
	c := "u"
	if i.Signed {
		c = "s"
	}

	return fmt.Sprintf("%s%d", c, i.Bytes*8)
}

func (f *FloatType) isType()       {}
func (f *FloatType) Size() int     { return f.Bytes }
func (f *FloatType) IsFloat() bool { return true }
func (f *FloatType) String() string {
	return fmt.Sprintf("f%d", f.Bytes*8)
}

func (s *StructType) isType()        {}
func (s *StructType) IsString() bool { return s.StructName == "String" } // very hacky way to do this
func (s *StructType) IsStruct() bool { return true }
func (s *StructType) Size() int {
	if s.IsString() {
		return 16
	}

	totalSizeWithoutPadding := 0
	for _, m := range s.Members {
		totalSizeWithoutPadding += m.DeclType.Size()
	}

	// TODO(Jovanni): Alignment and padding
	return totalSizeWithoutPadding
}
func (s *StructType) String() string {
	return fmt.Sprintf(s.StructName)
}

func (a *AliasType) isType()   {}
func (a *AliasType) Size() int { return a.Underlying().Size() }
func (a *AliasType) String() string {
	return fmt.Sprintf(a.AliasName)
}

func (f *FunctionType) isType()          {}
func (f *FunctionType) IsFunction() bool { return true }
func (f *FunctionType) String() string {
	ret := "fn("
	for i, param := range f.Params {
		ret += param.DeclType.String()
		if i < len(f.Params)-1 {
			ret += ", "
		}
	}
	ret += fmt.Sprintf(") -> %s", f.ReturnType.String())

	return ret
}

func (arr *StaticArrayType) isType()                   {}
func (arr *StaticArrayType) IsArray() bool             { return true }
func (arr *StaticArrayType) IsFixedSizeArray() bool    { return arr.Count > 0 }
func (arr *StaticArrayType) IsInferredSizeArray() bool { return arr.Count == 0 }
func (arr *StaticArrayType) Size() int {
	return arr.Underlying().Size() * arr.Count
}
func (arr *StaticArrayType) String() string {
	return fmt.Sprintf("[%d]", arr.Count) + arr.Underlying().String()
}

func (p *PointerType) isType()         {}
func (p *PointerType) IsPointer() bool { return true }
func (p *PointerType) Size() int {
	var t *int = nil
	return int(unsafe.Sizeof(t))
}
func (p *PointerType) String() string {
	return "*" + p.Underlying().String()
}

func NewTypeVoid() Type {
	return &VoidType{
		BaseType{nil},
	}
}

func NewTypeBool() Type {
	return &BoolType{
		BaseType{nil},
	}
}

func NewTypeChar(signed bool) Type {
	return &CharType{
		BaseType{nil},
		signed,
	}
}

func NewTypeInteger(signed bool, bytes int) Type {
	return &IntegerType{
		BaseType{nil},
		signed,
		bytes,
	}
}

func NewTypeFloat(bytes int) Type {
	return &FloatType{
		BaseType{nil},
		bytes,
	}
}

func NewTypeString() Type {
	typeU8Star := AddPointer(NewTypeInteger(false, 1))
	typeU64 := NewTypeInteger(false, 8)
	members := []Member{
		{
			Token.CreateToken(Token.IDENTIFIER, "data", 0),
			typeU8Star,
		},
		{
			Token.CreateToken(Token.IDENTIFIER, "length", 0),
			typeU64,
		},
	}

	return NewTypeStruct("String", members)
}

func NewTypeStruct(name string, members []Member) Type {
	return &StructType{
		BaseType:   BaseType{nil},
		StructName: name,
		Members:    members,
	}
}

func NewTypeFunction(returnType Type, params []Parameter) Type {
	return &FunctionType{
		BaseType:   BaseType{nil},
		ReturnType: returnType,
		Params:     params,
	}
}

func getFirstStrictAlias(t Type) Type {
	c := t
	for c != nil {
		v, ok := c.(*AliasType)
		if ok && v.Strict {
			break
		}

		c = c.Underlying()
	}

	return c
}

func getFirstNonTypeAlias(t Type) Type {
	c := t
	for c != nil {
		_, ok := c.(*AliasType)
		if !ok {
			break
		}

		c = c.Underlying()
	}

	return c
}

func typeEqual(t1 Type, t2 Type) bool {
	return reflect.TypeOf(t1) == reflect.TypeOf(t2)
}

// TypeStructuralEquivalence(t1 Type, t2 Type) (bool, error)
// TypeStrictNameEquivalence(t1 Type, t2 Type) (bool, error)
// TypeStrictCompatible(t1 Type, t2 Type) (bool, error) // does unsigned checks, and byte check
// TypeLooseCompatible(t1 Type, t2 Type) (bool, error) // does unsigned checks

// TypeCompare This is a structural + strict equivalence
func TypeCompare(t1 Type, t2 Type) (bool, error) {
	for {
		strictAlias1 := getFirstStrictAlias(t1)
		strictAlias2 := getFirstStrictAlias(t2)

		if (strictAlias1 == nil) || (strictAlias2 == nil) {
			break
		}

		if strictAlias1.(*AliasType).AliasName != strictAlias2.(*AliasType).AliasName {
			return false, fmt.Errorf("strict alias types are different")
		}

		t1 = strictAlias1.Underlying()
		t2 = strictAlias2.Underlying()
	}

	t1 = getFirstNonTypeAlias(t1)
	t2 = getFirstNonTypeAlias(t2)
	if !typeEqual(t1, t2) {
		return false, fmt.Errorf("")
	}

	switch v1 := t1.(type) {
	case *IntegerType:
		v2 := t2.(*IntegerType)
		if (v1.Bytes == v2.Bytes) && (v1.Signed == v2.Signed) {
			return true, nil
		}

		return false, fmt.Errorf("")
	case *FloatType:
		v2 := t2.(*FloatType)
		if v1.Bytes == v2.Bytes {
			return true, nil
		}

		return false, fmt.Errorf("")
	case *PointerType:
		return TypeCompare(t1.Underlying(), t2.Underlying())
	case *StaticArrayType:
		v2 := t2.(*StaticArrayType)
		ok, err := TypeCompare(t1.Underlying(), t2.Underlying())

		return ok && (v1.Count == v2.Count), err
	case *StructType:
		// This is a structural equivalence
		v2 := t2.(*StructType)
		if len(v1.Members) != len(v2.Members) {
			return false, fmt.Errorf("")
		}

		for i := 0; i < len(v1.Members); i++ {
			m1 := v1.Members[i]
			m2 := v2.Members[i]

			if ok, err := TypeCompare(m1.DeclType, m2.DeclType); !ok {
				return false, err
			}
		}
	case *FunctionType:
		v2 := t2.(*FunctionType)
		if ok, err := TypeCompare(v1.ReturnType, v2.ReturnType); !ok {
			return false, err
		}

		if len(v1.Params) != len(v2.Params) {
			return false, fmt.Errorf("")
		}

		for i := 0; i < len(v1.Params); i++ {
			p1 := v1.Params[i]
			p2 := v2.Params[i]

			if ok, err := TypeCompare(p1.DeclType, p2.DeclType); !ok {
				return false, err
			}
		}
	}

	return true, nil
}

// TypeSafeStructurallyStrictCompatible Structurally and strict equivalent and safe implicit casting
func TypeSafeStructurallyStrictCompatible(t1 Type, t2 Type) (bool, error) {
	for {
		strictAlias1 := getFirstStrictAlias(t1)
		strictAlias2 := getFirstStrictAlias(t2)

		if (strictAlias1 == nil) || (strictAlias2 == nil) {
			break
		}

		if strictAlias1.(*AliasType).AliasName != strictAlias2.(*AliasType).AliasName {
			return false, fmt.Errorf("strict alias types are different")
		}

		t1 = strictAlias1.Underlying()
		t2 = strictAlias2.Underlying()
	}

	t1 = getFirstNonTypeAlias(t1)
	t2 = getFirstNonTypeAlias(t2)
	if !typeEqual(t1, t2) {
		return false, fmt.Errorf("")
	}

	switch v1 := t1.(type) {
	case *IntegerType:
		v2 := t2.(*IntegerType)
		// narrowing cast
		if v1.Signed == v2.Signed {
			return true, nil
		}

		return false, fmt.Errorf("integer sign mismatch")
	case *FloatType:
		return true, nil
	case *PointerType:
		return TypeSafeStructurallyStrictCompatible(t1.Underlying(), t2.Underlying())
	case *StructType:
		return TypeCompare(t1, t2)
	case *StaticArrayType:
		v2 := t2.(*StaticArrayType)
		countCheck := (v1.IsInferredSizeArray() || v2.IsInferredSizeArray()) || (v1.Count == v2.Count)
		ok, err := TypeCompare(t1.Underlying(), t2.Underlying())

		return countCheck && ok, err
	default:
		panic(fmt.Sprintf("TypeCompatible: unknown type %T", t1))
	}

	return true, nil
}

func CanExplicitCast(caster Type, castee Type) (bool, error) {
	for {
		strictAlias1 := getFirstStrictAlias(caster)
		strictAlias2 := getFirstStrictAlias(castee)

		if (strictAlias1 == nil) || (strictAlias2 == nil) {
			break
		}

		if strictAlias1.(*AliasType).AliasName != strictAlias2.(*AliasType).AliasName {
			return false, fmt.Errorf("strict alias types are different")
		}

		caster = strictAlias1.Underlying()
		castee = strictAlias2.Underlying()
	}

	{
		_, isCasterPointer := caster.(*PointerType)
		if isCasterPointer {
			_, isCasterVoidPointer := caster.Underlying().(*PointerType)
			_, isCasteePointer := castee.(*PointerType)
			if isCasterVoidPointer && isCasteePointer {
				return true, nil
			}
		}
	}

	{
		v, isCasterStruct := caster.(*StructType)
		if isCasterStruct && v.IsString() {
			return true, nil
		}

		_, isCasterInteger := caster.(*IntegerType)
		_, isCasteeFloat := castee.(*FloatType)
		if isCasterInteger && isCasteeFloat {
			return true, nil
		}

		_, isCasterFloat := caster.(*FloatType)
		_, isCasteeInteger := castee.(*IntegerType)
		if isCasterFloat && isCasteeInteger {
			return true, nil
		}
	}

	/*
		{
			bool is_number_to_number_cast_allowed = (
			(type_is_signed(caster)   && type_is_unsigned(castee)) ||
				(type_is_unsigned(caster) && type_is_unsigned(castee)) ||
				(type_is_unsigned(caster) && type_is_signed(castee))   ||
				(type_is_signed(caster)   && type_is_signed(castee))
			);
			if (is_number_to_number_cast_allowed) return true;
		}
	*/

	return TypeCompare(caster, castee)
}

func GetBuiltin(s string) (Type, bool) {
	m := map[string]Type{
		"void": NewTypeVoid(),
		"bool": NewTypeBool(),

		"uchar": NewTypeChar(false),
		"char":  NewTypeChar(true),

		"u8":  NewTypeInteger(false, 1),
		"u16": NewTypeInteger(false, 2),
		"u32": NewTypeInteger(false, 4),
		"u64": NewTypeInteger(false, 8),

		"s8":  NewTypeInteger(true, 1),
		"s16": NewTypeInteger(true, 2),
		"s32": NewTypeInteger(true, 4),
		"s64": NewTypeInteger(true, 8),

		"f32": NewTypeFloat(4),
		"f64": NewTypeFloat(8),

		"string": NewTypeString(),
	}

	ret, ok := m[s]
	return ret, ok
}

func CanDereference(t Type) bool {
	p, isPointer := t.(*PointerType)
	if isPointer {
		_, isVoidPointer := p.Underlying().(*VoidType)
		return !isVoidPointer
	}

	return false
}

// CanImplicitCast Does NOT take structural equivalence into account
func CanImplicitCast(caster Type, castee Type) (bool, error) {
	for {
		strictAlias1 := getFirstStrictAlias(caster)
		strictAlias2 := getFirstStrictAlias(castee)

		if (strictAlias1 == nil) || (strictAlias2 == nil) {
			break
		}

		if strictAlias1.(*AliasType).AliasName != strictAlias2.(*AliasType).AliasName {
			return false, fmt.Errorf("can't perform implicitly to type `%s` from `%s`\n", caster.String(), castee.String())
		}

		caster = strictAlias1.Underlying()
		castee = strictAlias2.Underlying()
	}

	{
		_, isCasterPointer := caster.(*PointerType)
		if isCasterPointer {
			_, isCasterVoidPointer := caster.Underlying().(*PointerType).Underlying().(*VoidType)
			_, isCasteePointer := castee.(*PointerType)
			if isCasterVoidPointer && isCasteePointer {
				return true, nil
			}
		}
	}

	i1, isCasterInteger := caster.(*IntegerType)
	i2, isCasteeInteger := castee.(*IntegerType)
	if isCasterInteger && isCasteeInteger {
		return i1.Bytes >= i2.Bytes, fmt.Errorf("can't perform narrowing implicit cast to `%s` from `%s`\n", caster.String(), castee.String())
	}

	f1, isCasterFloat := caster.(*FloatType)
	f2, isCasteeFloat := castee.(*FloatType)
	if isCasterFloat && isCasteeFloat {
		return f1.Bytes >= f2.Bytes, fmt.Errorf("can't perform narrowing implicit cast to `%s` from `%s`\n", caster.String(), castee.String())
	}

	a1, isCasterArray := caster.(*StaticArrayType)
	a2, isCasteeArray := castee.(*StaticArrayType)
	if isCasterArray && isCasteeArray {
		if a1.IsInferredSizeArray() || a2.IsInferredSizeArray() {
			return true, nil
		}

		return false, fmt.Errorf("can't perform implicit cast to `%s` from `%s`\n", caster.String(), castee.String())
	}

	s1, isCasterStruct := caster.(*StructType)
	s2, isCasteeStruct := castee.(*StructType)
	if isCasterStruct && isCasteeStruct {
		if s1.StructName == s2.StructName {
			return true, nil
		}

		return false, fmt.Errorf("can't perform implicit cast to `%s` from `%s`\n", caster.String(), castee.String())
	}

	// if castee is a IntegerType or a FloatingType
	// and the caster is a Integer Type or Floating Type

	// this case is like if its a parameter
	// f(a: u64)
	// CanImplicitCast(param.DeclType, argument.DeclType, argument.Expr)

	return TypeSafeStructurallyStrictCompatible(caster, castee)
}

type TypeKind int

const (
	INVALID TypeKind = iota

	INTEGER
	FLOAT
	BOOL
	STRING
)

type BinaryQuery struct {
	Op    string
	Left  TypeKind
	Right TypeKind
}

func getTypeKind(t Type) TypeKind {
	switch v := t.(type) {
	case *IntegerType:
		return INTEGER
	case *FloatType:
		return FLOAT
	case *BoolType:
		return BOOL
	case *StructType:
		if v.IsString() {
			return STRING
		}
	}

	return INVALID
}

// GetPromotedType This is strictly for binary operations
func GetPromotedType(op Token.Token, leftType Type, rightType Type) Type {
	typeS32 := NewTypeInteger(true, 4)
	typeF32 := NewTypeFloat(4)
	typeString := NewTypeString()
	typeBool := NewTypeBool()

	var typeMap = map[BinaryQuery]Type{
		{"+", INTEGER, FLOAT}:   typeF32,
		{"-", INTEGER, FLOAT}:   typeF32,
		{"*", INTEGER, FLOAT}:   typeF32,
		{"/", INTEGER, FLOAT}:   typeF32,
		{"%", INTEGER, INTEGER}: typeS32,

		{"<", INTEGER, FLOAT}:  typeBool,
		{"<=", INTEGER, FLOAT}: typeBool,
		{">", INTEGER, FLOAT}:  typeBool,
		{">=", INTEGER, FLOAT}: typeBool,

		{"+", STRING, STRING}:  typeString,
		{"+", STRING, INTEGER}: typeString,
		{"+", STRING, FLOAT}:   typeString,

		{"||", BOOL, BOOL}: typeBool,
		{"&&", BOOL, BOOL}: typeBool,
	}

	arithmeticOperators := []string{"+", "-", "*", "/"}
	for _, operator := range arithmeticOperators {
		typeMap[BinaryQuery{operator, INTEGER, INTEGER}] = typeS32
		typeMap[BinaryQuery{operator, FLOAT, FLOAT}] = typeF32
	}

	comparisonOperators := []string{"<", "<=", ">", ">="}
	for _, operator := range comparisonOperators {
		typeMap[BinaryQuery{operator, INTEGER, INTEGER}] = typeBool
		typeMap[BinaryQuery{operator, FLOAT, FLOAT}] = typeBool
	}

	equalityOperators := []string{"==", "!="}
	for _, operator := range equalityOperators {
		typeMap[BinaryQuery{operator, INTEGER, INTEGER}] = typeBool
		typeMap[BinaryQuery{operator, FLOAT, FLOAT}] = typeBool
		typeMap[BinaryQuery{operator, STRING, STRING}] = typeBool
		typeMap[BinaryQuery{operator, BOOL, BOOL}] = typeBool
	}

	leftTypeKind := getTypeKind(leftType)
	if leftTypeKind == INVALID {
		panic("invalid left type")
	}

	rightTypeKind := getTypeKind(rightType)
	if rightTypeKind == INVALID {
		panic("invalid right type")
	}

	// NOTE(Jovanni): This is because I don't want to have to specify like pointer + u8 and u8 + pointer
	q1 := BinaryQuery{op.Lexeme, leftTypeKind, rightTypeKind}
	q2 := BinaryQuery{op.Lexeme, rightTypeKind, leftTypeKind}

	if ret, ok := typeMap[q1]; ok {
		return ret
	}

	if ret, ok := typeMap[q2]; ok {
		return ret
	}

	return nil
}
