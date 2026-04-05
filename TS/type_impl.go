package TS

import (
	"fmt"
	"reflect"
	"unsafe"
)

func (b *BaseType) isType()          {}
func (b *BaseType) String() string   { return "BaseType" }
func (b *BaseType) Underlying() Type { return b.underlying }
func (b *BaseType) Size() int        { return 0 }
func (b *BaseType) IsStruct() bool   { return false }
func (b *BaseType) IsPointer() bool  { return false }
func (b *BaseType) IsArray() bool    { return false }
func (b *BaseType) AddAlias(strict bool, name string) Type {
	return &AliasType{
		BaseType{b},
		strict,
		name,
	}
}

func (b *BaseType) AddStaticArray(capacity int) Type {
	return &StaticArrayType{
		BaseType{b},
		capacity,
	}
}

func (b *BaseType) AddDynamicArray(capacity int) Type {
	return &DynamicArrayType{
		BaseType{b},
	}
}

func (b *BaseType) AddPointer() Type {
	return &PointerType{
		BaseType{b},
	}
}

func AllocateDeepCopyType(t Type) Type {
	if t == nil {
		return nil
	}

	switch v := t.(type) {
	case *CharType:
		return &CharType{
			BaseType{AllocateDeepCopyType(t.Underlying())},
			v.Signed,
		}
	case *IntegerType:
		return &IntegerType{
			BaseType{AllocateDeepCopyType(t.Underlying())},
			v.Signed,
			v.Bytes,
		}
	case *FloatType:
		return &StaticArrayType{
			BaseType{AllocateDeepCopyType(t.Underlying())},
			v.Bytes,
		}
	case *StaticArrayType:
		return &StaticArrayType{
			BaseType{AllocateDeepCopyType(t.Underlying())},
			v.Capacity,
		}
	case *DynamicArrayType:
		return &DynamicArrayType{
			BaseType{AllocateDeepCopyType(t.Underlying())},
		}
	case *StructType:
		return &StructType{
			BaseType{AllocateDeepCopyType(t.Underlying())},
			v.StructName,
			v.MemberTypes,
		}
	case *FunctionType:
		return &FunctionType{
			BaseType{AllocateDeepCopyType(t.Underlying())},
			v.ReturnType,
			v.ParamTypes,
		}
	case *AliasType:
		return &AliasType{
			BaseType{AllocateDeepCopyType(t.Underlying())},
			v.Strict,
			v.AliasName,
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

func (b *BaseType) RemoveStaticArray(capacity int) Type {
	return &StaticArrayType{
		BaseType{b},
		capacity,
	}
}

func (b *BaseType) RemoveDynamicArray(capacity int) Type {
	return &DynamicArrayType{
		BaseType{b},
	}
}

func (b *BaseType) RemovePointer() Type {
	return &PointerType{
		BaseType{b},
	}
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

func (i *IntegerType) isType()   {}
func (i *IntegerType) Size() int { return i.Bytes }
func (i *IntegerType) String() string {
	c := "u"
	if i.Signed {
		c = "s"
	}

	return fmt.Sprintf(TYPE_INTEGER, c, i.Bytes*8)
}

func (f *FloatType) isType()   {}
func (f *FloatType) Size() int { return f.Bytes }
func (f *FloatType) String() string {
	return fmt.Sprintf(TYPE_FLOATING, f.Bytes*8)
}

func (s *StringType) isType()   {}
func (s *StringType) Size() int { return 16 }
func (s *StringType) String() string {
	return TYPE_STRING
}

func (s *StructType) isType()         {}
func (p *PointerType) IsStruct() bool { return true }
func (s *StructType) Size() int {
	totalSizeWithoutPadding := 0
	for _, m := range s.MemberTypes {
		totalSizeWithoutPadding += m.Size()
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

func (f *FunctionType) isType()      {}
func (b *BaseType) IsFunction() bool { return true }
func (f *FunctionType) String() string {
	ret := "fn("
	for i, param := range f.ParamTypes {
		ret += param.String()
		if i < len(f.ParamTypes)-1 {
			ret += ", "
		}
	}
	ret += fmt.Sprintf(") -> %s", f.ReturnType.String())

	return ret
}

func (arr *StaticArrayType) isType()       {}
func (arr *StaticArrayType) IsArray() bool { return true }
func (arr *StaticArrayType) Size() int {
	return arr.Underlying().Size() * arr.Capacity
}
func (arr *StaticArrayType) String() string {
	return fmt.Sprintf(TYPE_STATIC_ARRAY, arr.Capacity) + arr.Underlying().String()
}

func (dyn *DynamicArrayType) isType()       {}
func (dyn *DynamicArrayType) IsArray() bool { return true }
func (dyn *DynamicArrayType) Size() int {
	return 16 // void* + u64
}
func (dyn *DynamicArrayType) String() string {
	return TYPE_DYNAMIC_ARRAY + dyn.Underlying().String()
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
	return &StringType{
		BaseType{nil},
	}
}

func NewTypeStruct(name string, memberTypes []Type) Type {
	return &StructType{
		BaseType{nil},
		name,
		memberTypes,
	}
}

func NewTypeFunction(returnType Type, paramTypes []Type) Type {
	return &FunctionType{
		BaseType{nil},
		returnType,
		paramTypes,
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
		if len(v1.MemberTypes) != len(v2.MemberTypes) {
			return false
		}

		for i := 0; i < len(v1.MemberTypes); i++ {
			m1 := v1.MemberTypes[i]
			m2 := v2.MemberTypes[i]

			if !TypeCompare(m1, m2) {
				return false
			}
		}
	case *FunctionType:
		v2 := t2.(*FunctionType)
		if TypeCompare(v1.ReturnType, v2.ReturnType) {
			return false
		}

		if len(v1.ParamTypes) != len(v2.ParamTypes) {
			return false
		}

		for i := 0; i < len(v1.ParamTypes); i++ {
			p1 := v1.ParamTypes[i]
			p2 := v2.ParamTypes[i]

			if !TypeCompare(p1, p2) {
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

/*
func CanImplicitCast() bool {
	return true
}
*/
