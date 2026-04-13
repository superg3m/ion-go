package TS

import (
	"fmt"
	"ion-go/Token"
	"math"
	"reflect"
)

func (b *BaseType) isType()                   {}
func (b *BaseType) String() string            { return "BaseType" }
func (b *BaseType) Underlying() Type          { return b.underlying }
func (b *BaseType) Size() int                 { return 0 }
func (b *BaseType) Align() int                { return 0 }
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
func (b *BoolType) Align() int     { return 1 }
func (b *BoolType) String() string { return "bool" }

func (c *CharType) isType()    {}
func (c *CharType) Size() int  { return 1 }
func (c *CharType) Align() int { return 1 }
func (c *CharType) String() string {
	u := "u"
	if c.Signed {
		u = ""
	}

	return fmt.Sprintf("%schar", u)
}

func (i *IntegerType) isType()         {}
func (i *IntegerType) Size() int       { return i.Bytes }
func (i *IntegerType) Align() int      { return i.Bytes }
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
func (f *FloatType) Align() int    { return f.Bytes }
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

	offset := 0
	for _, m := range s.Members {
		size := m.DeclType.Size()
		align := m.DeclType.Align()

		offset = alignUp(offset, align)
		offset += size
	}

	return alignUp(offset, s.Align())
}

// NOTE(Jovanni): ONLY WORKS if alignment is a power of two
func alignUp(offset, alignment int) int {
	return (offset + alignment - 1) & ^(alignment - 1)
}

func (s *StructType) Align() int {
	if s.IsString() {
		return 16
	}

	maxTypeAlignment := 0
	for _, m := range s.Members {
		maxTypeAlignment = int(math.Max(float64(maxTypeAlignment), float64(m.DeclType.Align())))
	}

	return int(math.Min(float64(maxTypeAlignment), 32))
	// return int(math.Min(math.Max(float64(maxTypeAlignment), float64(s.Size())), 32))
}

func (s *StructType) String() string {
	return fmt.Sprintf(s.StructName)
}

func (a *AliasType) isType()    {}
func (a *AliasType) Size() int  { return a.Underlying().Size() }
func (a *AliasType) Align() int { return a.Underlying().Align() }
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
func (arr *StaticArrayType) Size() int                 { return arr.Underlying().Size() * arr.Count }
func (arr *StaticArrayType) Align() int                { return arr.Underlying().Align() }
func (arr *StaticArrayType) String() string {
	return fmt.Sprintf("[%d]", arr.Count) + arr.Underlying().String()
}

func (p *PointerType) isType()         {}
func (p *PointerType) IsPointer() bool { return true }
func (p *PointerType) Size() int       { return 8 } // just assume 64-bit architecture
func (p *PointerType) Align() int      { return 8 } // just assume 64-bit architecture
func (p *PointerType) String() string  { return "*" + p.Underlying().String() }

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

// TypeStructuralEquivalence this will also act as TypeLooseNameEquivalence()
func TypeStructuralEquivalence(t1 Type, t2 Type) bool {
	t1 = getFirstNonTypeAlias(t1)
	t2 = getFirstNonTypeAlias(t2)

	if !typeEqual(t1, t2) {
		return false
	}

	switch v1 := t1.(type) {
	case *IntegerType:
		v2 := t2.(*IntegerType)
		return v1.Bytes == v2.Bytes
	case *FloatType:
		v2 := t2.(*FloatType)
		return v1.Bytes == v2.Bytes
	case *PointerType:
		return TypeStructuralEquivalence(t1.Underlying(), t2.Underlying())
	case *StaticArrayType:
		v2 := t2.(*StaticArrayType)
		structurallyEquivalent := TypeStructuralEquivalence(t1.Underlying(), t2.Underlying())

		return (v1.Count == v2.Count) && structurallyEquivalent
	case *StructType:
		// This is a structural equivalence
		v2 := t2.(*StructType)
		if len(v1.Members) != len(v2.Members) {
			return false
		}

		for i := 0; i < len(v1.Members); i++ {
			m1 := v1.Members[i]
			m2 := v2.Members[i]

			if !TypeStructuralEquivalence(m1.DeclType, m2.DeclType) {
				return false
			}
		}
	case *FunctionType:
		v2 := t2.(*FunctionType)
		if !TypeStructuralEquivalence(v1.ReturnType, v2.ReturnType) {
			return false
		}

		if len(v1.Params) != len(v2.Params) {
			return false
		}

		for i := 0; i < len(v1.Params); i++ {
			p1 := v1.Params[i]
			p2 := v2.Params[i]

			if TypeStructuralEquivalence(p1.DeclType, p2.DeclType) {
				return false
			}
		}
	}

	return true
}

func TypeStrictNameEquivalence(t1 Type, t2 Type) bool {
	for {
		strictAlias1 := getFirstStrictAlias(t1)
		strictAlias2 := getFirstStrictAlias(t2)

		if (strictAlias1 == nil) || (strictAlias2 == nil) {
			break
		}

		if strictAlias1.(*AliasType).AliasName != strictAlias2.(*AliasType).AliasName {
			return false
		}

		t1 = strictAlias1.Underlying()
		t2 = strictAlias2.Underlying()
	}

	t1 = getFirstNonTypeAlias(t1)
	t2 = getFirstNonTypeAlias(t2)

	return t1.String() == t2.String()
}

// does unsigned vs signed checks, and byte check
func TypeStrictCompare(t1 Type, t2 Type) (bool, error) {
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
		if v1.Signed != v2.Signed {
			return false, fmt.Errorf("integer sign mismatch: '%s', '%s'", v1.String(), v2.String())
		}

		return v1.Bytes == v2.Bytes, fmt.Errorf("integer size mismatch: '%s', '%s'", v1.String(), v2.String())
	case *FloatType:
		v2 := t2.(*FloatType)
		return v1.Bytes == v2.Bytes, fmt.Errorf("float size mismatch: '%s', '%s'", v1.String(), v2.String())
	case *StaticArrayType:
		v2 := t2.(*StaticArrayType)
		ok, err := TypeStrictCompare(t1.Underlying(), t2.Underlying())

		return v1.Count == v2.Count && ok, err
	case *PointerType:
		return TypeStrictCompare(t1.Underlying(), t2.Underlying())
	case *StructType:
		return TypeStructuralEquivalence(t1, t2), fmt.Errorf("'%s' and '%s' and not structurally equivalent", t1.String(), t2.String())
	}

	if !TypeStructuralEquivalence(t1.Underlying(), t2.Underlying()) {
		return false, fmt.Errorf("'%s' and '%s' and not structurally equivalent", t1.String(), t2.String())
	}

	if !TypeStrictNameEquivalence(t1.Underlying(), t2.Underlying()) {
		return false, fmt.Errorf("'%s' and '%s' and not strictly name equivalent", t1.String(), t2.String())
	}

	return true, nil
}

// TypeLooseCompare does unsigned vs signed checks
func TypeLooseCompare(t1 Type, t2 Type) (bool, error) {
	t1 = getFirstNonTypeAlias(t1)
	t2 = getFirstNonTypeAlias(t2)
	if !typeEqual(t1, t2) {
		return false, fmt.Errorf("")
	}

	switch v1 := t1.(type) {
	case *IntegerType:
		v2 := t2.(*IntegerType)
		if v1.Signed != v2.Signed {
			return false, fmt.Errorf("integer sign mismatch: '%s', '%s'", v1.String(), v2.String())
		}

		return v1.Bytes == v2.Bytes, fmt.Errorf("integer signed mismatch: '%s', '%s'", v1.String(), v2.String())
	case *StaticArrayType:
		v2 := t2.(*StaticArrayType)
		countCheck := (v1.IsInferredSizeArray() || v2.IsInferredSizeArray()) || (v1.Count == v2.Count)
		ok, err := TypeLooseCompare(t1.Underlying(), t2.Underlying())

		return countCheck && ok, err
	}

	if !TypeStructuralEquivalence(t1.Underlying(), t2.Underlying()) {
		return false, fmt.Errorf("'%s' and '%s' and not structurally equivalent", t1.String(), t2.String())
	}

	return true, nil
}

func CanExplicitCast(caster Type, castee Type) (bool, error) {
	caster = getFirstNonTypeAlias(caster)
	castee = getFirstNonTypeAlias(castee)

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

	a1, isCasterArray := caster.(*StaticArrayType)
	a2, isCasteeArray := castee.(*StaticArrayType)
	if isCasterArray && isCasteeArray {
		if a1.IsInferredSizeArray() || a2.IsInferredSizeArray() {
			return true, nil
		}

		return a1.Count == a2.Count, fmt.Errorf("can't perform explicit cast to `%s` from `%s`", caster.String(), castee.String())
	}

	v, isCasterStruct := caster.(*StructType)
	if isCasterStruct && v.IsString() {
		if castee.IsStruct() || (!castee.IsStruct() && !castee.IsArray()) {
			return true, nil
		}
	}

	s1, isCasterStruct := caster.(*StructType)
	s2, isCasteeStruct := castee.(*StructType)
	if isCasterStruct && isCasteeStruct {
		if s1.StructName == s2.StructName {
			return true, nil
		}

		return TypeStructuralEquivalence(caster, castee), fmt.Errorf("can't perform explicit cast to `%s` from `%s` because they are not structurally equivalent", caster.String(), castee.String())
	}

	return true, nil
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
	caster = getFirstNonTypeAlias(caster)
	castee = getFirstNonTypeAlias(castee)

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

	v, isCasterStruct := caster.(*StructType)
	if isCasterStruct && v.IsString() {
		return true, nil
	}

	i1, isCasterInteger := caster.(*IntegerType)
	f1, isCasterFloat := caster.(*FloatType)

	i2, isCasteeInteger := castee.(*IntegerType)
	f2, isCasteeFloat := castee.(*FloatType)

	if isCasterInteger && isCasteeFloat {
		return i1.Bytes >= f2.Bytes, fmt.Errorf("can't perform narrowing implicit cast because of sign mismatch '%s', '%s'", i1.String(), f2.String())
	}

	if isCasterFloat && isCasteeInteger {
		return f1.Bytes >= i2.Bytes, fmt.Errorf("can't perform narrowing implicit cast because of sign mismatch '%s', '%s'", f1.String(), i2.String())
	}

	if isCasterInteger && isCasteeInteger {
		if i1.Signed != i2.Signed {
			return false, fmt.Errorf("can't perform narrowing implicit cast because of sign mismatch '%s', '%s'", i1.String(), i2.String())
		}

		return i1.Bytes >= i2.Bytes, fmt.Errorf("can't perform narrowing implicit cast because of sign mismatch '%s', '%s'", i1.String(), i2.String())
	}

	if isCasterFloat && isCasteeFloat {
		return f1.Bytes >= f2.Bytes, fmt.Errorf("can't perform narrowing implicit cast because of sign mismatch '%s', '%s'", f1.String(), f2.String())
	}

	a1, isCasterArray := caster.(*StaticArrayType)
	a2, isCasteeArray := castee.(*StaticArrayType)
	if isCasterArray && isCasteeArray {
		if a1.IsInferredSizeArray() || a2.IsInferredSizeArray() {
			return true, nil
		}

		return a1.Count == a2.Count, fmt.Errorf("can't perform implicit cast to `%s` from `%s`", caster.String(), castee.String())
	}

	return TypeStrictNameEquivalence(caster, castee), fmt.Errorf("")
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
