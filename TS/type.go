package TS

const (
	TYPE_VOID = "void"
	TYPE_BOOL = "bool"
	TYPE_CHAR = "%schar"

	TYPE_INTEGER  = "%s%d"
	TYPE_FLOATING = "f%d"

	TYPE_STRING       = "string"
	TYPE_MATH_VECTOR2 = "vec2"
	TYPE_MATH_VECTOR3 = "vec3"
	TYPE_MATH_VECTOR4 = "vec4"
	TYPE_MATH_MAT4    = "mat4"

	TYPE_ALIAS         = "TYPE_ALIAS"
	TYPE_STATIC_ARRAY  = "[%d]"
	TYPE_DYNAMIC_ARRAY = "[..]"
	TYPE_POINTER       = "*"

	TYPE_STRUCT   = "STRUCT"
	TYPE_FUNCTION = "FUNCTION"
)

/*
struct StringType {
	u8* data;
	u64 length;
};

struct DynamicArrayType {
	void* data;
	u64 length;
};
*/

type Type interface {
	isType()
	Underlying() Type
	Size() int
	IsFunction() bool
	IsStruct() bool
	IsPointer() bool
	IsArray() bool
	RemoveModifier() Type
	RemoveStaticArray(capacity int) Type
	RemoveDynamicArray(capacity int) Type
	RemovePointer() Type
	String() string
}

type BaseType struct {
	underlying Type
}

type VoidType struct {
	BaseType
}

type BoolType struct {
	BaseType
}

type CharType struct {
	BaseType
	Signed bool
}

type IntegerType struct {
	BaseType
	Signed bool
	Bytes  int // 1: u8, 2: u16, 4: u32, 8: u64
}

type FloatType struct {
	BaseType
	Bytes int // 1: u8, 2: u16, 4: u32, 8: u64
}

type StringType struct {
	BaseType
}

type StructType struct {
	BaseType
	StructName  string
	MemberTypes []Type
}

type AliasType struct {
	BaseType
	Strict    bool
	AliasName string
}

type FunctionType struct {
	BaseType
	ReturnType Type
	ParamTypes []Type
}

type StaticArrayType struct {
	BaseType
	Capacity int // not byte capacity but the like size capacity int[4] is not the same as u64[2] even if the capacity in bytes is the same
}

type DynamicArrayType struct {
	BaseType
}

type PointerType struct {
	BaseType
}
