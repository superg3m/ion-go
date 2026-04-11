package TS

import (
	"ion-go/Token"
)

/*
// These are just structs now that I think about it:
	TYPE_STRING       = "String"
	TYPE_MATH_VECTOR2 = "Vector2"
	TYPE_MATH_VECTOR3 = "Vector3"
	TYPE_MATH_VECTOR4 = "Vector4"
	TYPE_MATH_MAT4    = "Matrix4"

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
	Align() int
	IsFunction() bool
	IsStruct() bool
	IsPointer() bool
	IsArray() bool
	IsFixedSizeArray() bool
	IsInferredSizeArray() bool
	IsString() bool
	IsInteger() bool
	IsFloat() bool
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

type Member struct {
	Tok      Token.Token
	DeclType Type
}

type Parameter struct {
	Tok      Token.Token
	DeclType Type
}

type StructType struct {
	BaseType
	StructName string
	Members    []Member
}

type AliasType struct {
	BaseType
	Strict    bool
	AliasName string
}

type FunctionType struct {
	BaseType
	ReturnType Type
	Params     []Parameter
}

type StaticArrayType struct {
	BaseType
	Count int
}

type PointerType struct {
	BaseType
}
