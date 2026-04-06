package TS

import (
	"fmt"
	"ion-go/Token"
	"reflect"
	"unsafe"
)

func (b *BaseType) isType()          {}
func (b *BaseType) String() string   { return "BaseType" }
func (b *BaseType) Underlying() Type { return b.underlying }
func (b *BaseType) Size() int        { return 0 }
func (b *BaseType) IsInteger() bool  { return false }
func (b *BaseType) IsString() bool   { return false }
func (b *BaseType) IsStruct() bool   { return false }
func (b *BaseType) IsPointer() bool  { return false }
func (b *BaseType) IsArray() bool    { return false }
func (b *BaseType) IsFloat() bool    { return false }

func (b *BaseType) IsFunction() bool { return false }
func (b *BaseType) AddAlias(strict bool, name string) Type {
	return &AliasType{
		BaseType:  BaseType{b},
		Strict:    strict,
		AliasName: name,
	}
}

func (b *BaseType) AddStaticArray(count int) Type {
	return &StaticArrayType{
		BaseType: BaseType{b},
		Count:    count,
	}
}

func (b *BaseType) AddPointer() Type {
	return &PointerType{
		BaseType: BaseType{b},
	}
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
	return AllocateDeepCopyType(b)
}

func (v *VoidType) isType() {}
func (v *VoidType) String() string {
	return string(TYPE_VOID)
}

func (b *BoolType) isType()        {}
func (b *BoolType) Size() int      { return 1 }
func (b *BoolType) String() string { return TYPE_BOOL }

func (c *CharType) isType()   {}
func (c *CharType) Size() int { return 1 }
func (c *CharType) String() string {
	u := "u"
	if c.Signed {
		u = ""
	}

	return fmt.Sprintf(TYPE_CHAR, u)
}

func (i *IntegerType) isType()         {}
func (i *IntegerType) Size() int       { return i.Bytes }
func (i *IntegerType) IsInteger() bool { return true }
func (i *IntegerType) String() string {
	c := "u"
	if i.Signed {
		c = "s"
	}

	return fmt.Sprintf(TYPE_INTEGER, c, i.Bytes*8)
}

func (f *FloatType) isType()       {}
func (f *FloatType) Size() int     { return f.Bytes }
func (f *FloatType) IsFloat() bool { return true }
func (f *FloatType) String() string {
	return fmt.Sprintf(TYPE_FLOATING, f.Bytes*8)
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

func (arr *StaticArrayType) isType()       {}
func (arr *StaticArrayType) IsArray() bool { return true }
func (arr *StaticArrayType) Size() int {
	return arr.Underlying().Size() * arr.Count
}
func (arr *StaticArrayType) String() string {
	return fmt.Sprintf(TYPE_STATIC_ARRAY, arr.Count) + arr.Underlying().String()
}

func (p *PointerType) isType()         {}
func (p *PointerType) IsPointer() bool { return true }
func (p *PointerType) Size() int {
	var t *int = nil
	return int(unsafe.Sizeof(t))
}
func (p *PointerType) String() string {
	return TYPE_POINTER + p.Underlying().String()
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
	typeU8Star := NewTypeInteger(false, 1).AddPointer()
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

func TypeEqual(t1 Type, t2 Type) bool {
	return reflect.TypeOf(t1) == reflect.TypeOf(t2)
}

func TypeCompare(t1 Type, t2 Type) bool {
	for {
		st := getFirstStrictAlias(t1)
		st2 := getFirstStrictAlias(t2)

		if (st == nil) || (st2 == nil) {
			break
		}

		if st.(*AliasType).AliasName != st2.(*AliasType).AliasName {
			return false
		}

		t1 = st.Underlying()
		t2 = st2.Underlying()
	}

	t1 = getFirstNonTypeAlias(t1)
	t2 = getFirstNonTypeAlias(t2)
	if !TypeEqual(t1, t2) {
		return false
	}

	switch v1 := t1.(type) {
	case *IntegerType:
		v2 := t2.(*IntegerType)
		return (v1.Bytes == v2.Bytes) && (v1.Signed == v2.Signed)
	case *FloatType:
		v2 := t2.(*FloatType)
		return v1.Bytes == v2.Bytes
	case *PointerType:
		return TypeCompare(t1.Underlying(), t2.Underlying())
	case *StructType:
		v2 := t2.(*StructType)
		if len(v1.Members) != len(v2.Members) {
			return false
		}

		for i := 0; i < len(v1.Members); i++ {
			m1 := v1.Members[i]
			m2 := v2.Members[i]

			if !TypeCompare(m1.DeclType, m2.DeclType) {
				return false
			}
		}
	case *FunctionType:
		v2 := t2.(*FunctionType)
		if TypeCompare(v1.ReturnType, v2.ReturnType) {
			return false
		}

		if len(v1.Params) != len(v2.Params) {
			return false
		}

		for i := 0; i < len(v1.Params); i++ {
			p1 := v1.Params[i]
			p2 := v2.Params[i]

			if !TypeCompare(p1.DeclType, p2.DeclType) {
				return false
			}
		}
	}

	return true
}

func CanExplicitCast(caster Type, castee Type) bool {
	for {
		st := getFirstStrictAlias(caster)
		st2 := getFirstStrictAlias(castee)

		if (st == nil) || (st2 == nil) {
			break
		}

		if st.(*AliasType).AliasName != st2.(*AliasType).AliasName {
			return false
		}

		caster = st.Underlying()
		castee = st2.Underlying()
	}

	{
		_, isCasterPointer := caster.(*PointerType)
		if isCasterPointer {
			_, isCasterVoidPointer := caster.Underlying().(*PointerType)
			_, isCasteePointer := castee.(*PointerType)
			if isCasterVoidPointer && isCasteePointer {
				return true
			}
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

/*
func CanImplicitCast() bool {
	return true
}
*/

/*
struct BinaryQuery {
TypeKind left;
const char* op;
TypeKind right;

bool operator==(BinaryQuery other) const {
return (left == other.left) && str_equal(this->op, other.op) && (right == other.right);
}
};

BinaryQuery binary_query_create(TypeKind left, const char* op, TypeKind right) {
BinaryQuery ret = {};
ret.left = left;
ret.op = op;
ret.right = right;

return ret;
}

#define X_ARITHMITIC_OPERATORS \
X("+")                     \
X("-")                     \
X("*")                     \
X("/")                     \
X("<")                     \
X("<=")                    \
X(">")                     \
X(">=")                    \

bool type_can_perform_binary_op(TypeKind t1, const char* op, TypeKind t2) {
LOCAL_PERSIST Hashmap<BinaryQuery, bool> binary_query_map = hashmap_create<BinaryQuery, bool>(allocator_general(), {
#define X(ENUM, STR, SIZE) {binary_query_create(ENUM, "==", ENUM), true}, {binary_query_create(ENUM, "!=", ENUM), true},
X_TYPE_BUILTIN
#undef X

#define X(OP) {binary_query_create(TYPE_INTEGER, OP, TYPE_INTEGER), true}, {binary_query_create(TYPE_INTEGER, OP, TYPE_FLOATING), true}, {binary_query_create(TYPE_FLOATING, OP, TYPE_FLOATING), true}, {binary_query_create(TYPE_INTEGER, OP, TYPE_FLOATING), true},
X_ARITHMITIC_OPERATORS
#undef X

{binary_query_create(TYPE_POINTER, "+", TYPE_INTEGER), true},
{binary_query_create(TYPE_POINTER, "-", TYPE_INTEGER), true},
});

// NOTE(Jovanni): This is because I don't want to have to specify like pointer + u8 and u8 + pointer
BinaryQuery q1 = binary_query_create(t1, op, t2);
BinaryQuery q2 = binary_query_create(t2, op, t1);
return hashmap_has(&binary_query_map, q1) || hashmap_has(&binary_query_map, q2);
}

*/
