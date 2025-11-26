# ion-go
The ion programming language

1. FrontEnd: lexing + parsing + AST
2. Analysis: semantic analysis, typechecking
3. Interpreter: AST Treewalk

## Language Features
Ion supports:
- Structs
- Slices and multi-dimensional slices
- Functions with typed parameters and return values
- Type inference (:=)
- Struct literals and slice literals
- Indexing and nested indexing
- Casting
- Control flow (if, for, continue, return)
- defer blocks with LIFO execution
- Recursion
- Built-ins: len()
- Basic string concatenation and printing

### Examples
```go
// struct.ion
struct Foo {
    bar: float
}

struct Person {
    age: int,
    foo: []Foo,
    name: string
}

fn do_something(x: Person) -> void {
    println(x);
}

fn main() -> void {
    var test := []Foo.[Foo.{1.5}, Foo.{2.2}];
    var x: Person = Person.{23, []Foo.[Foo.{5.2}, Foo.{10.0}], "John"};

    println(test);
    test[0].bar = 1.2;
    test[1].bar = 404.2;
    println(test);

    println(test[0].bar);
    println(x.foo[0].bar);
    x.foo[0].bar = 6.2;

    test[cast(int)test[0].bar].bar = 53.2;

    do_something(x);
    println(test);
}

/* OUTPUT:
[{bar: float = 1.5}, {bar: float = 2.2}]
[{bar: float = 1.2}, {bar: float = 404.2}]
1.2
5.2
{
    age: int = 23,
    foo: []Foo = [{bar: float = 6.2}, {bar: float = 10}],
    name: string = John
}
[{bar: float = 1.2}, {bar: float = 53.2}]
 */
```

```go
// fib.ion
fn fib(x: int) -> int {
    if (x == 0 || x == 1) {
        return x;
    }
    return fib(x - 1) + fib(x - 2);
}

fn main() -> void {
    defer {
        defer print(fib(8) + "\n");
        print("Hello world\n");
    }

    defer print("THIS IS ACTUALLY BEFORE LAST DEFER!\n");

    if (true || fib(100) > 4) {
        print("First!\n");
        return;
    }

    print("NOT HAPPENING\n");
}
/* OUTPUT:
First!
THIS IS ACTUALLY BEFORE LAST DEFER!
Hello world
21
*/
```

BNF Key
--------------------------------
<foo>       Nonterminal (rule defined elsewhere)
"text"      Terminal (literal keyword or symbol)
::=         Defines a rule
|           Choice (OR)
()          Grouping

?           Optional (0 or 1 occurrence)
*           Repetition (0 or more occurrences)
+           Repetition (1 or more occurrences)
--------------------------------------------------

### HIGH-LEVEL STRUCTURE
<program> ::= (<function_decl> | <struct_decl> | <variable_decl>)*
<scope> ::= "{" (<node>)* "}"

<node> ::= (<statement> | <decleration> | <expression>)

### DECLARATIONS
<variable_decl> ::= "var" <identifier> ":" ((<type>)? ("=" <expression>)) | ((<type>) ("=" <expression>)?) ";"
// var test: int;
// var test := 5;
// var test: int = 5;

<function_decl> ::= "fn" <identifier> "(" <param_list>? ")" "->" <return_type> <scope>
<param_list> ::= <parameter> ("," <parameter>)*
<parameter> ::= <identifier> ":" <type>
<return_type> ::= <type> | "(" <type_list> ")" | "void"
<type_list> ::= <type> ("," <type>)*
/*
func get_value(a: int, b: int) -> void {}
*/

<struct_decl> ::= "struct" <identifier> "{" (<struct_member>)* "}"
<struct_member> ::= <identifier> ":" <type> ";"

### TYPES
<primitive_type> ::= "int" | "float" | "bool" | "string"

### STATEMENTS
<statement> ::= <assignment> |<return> | <if_else> | <while> |
                <continue> | <break>


<assignment> ::= <lhs> "=" <expression> ";"
<lhs> ::= <identifier> | <member_access> | <array_access>
// test = 4

<member_access> ::= <member_access> "." <member_access> | <identifier>
<array_access> ::= <identifier> ("[" <expression> "]")+

<return_stmt> ::= "return" <return_value>? ";"
<return_value> ::= <expression> | "(" <expression> ("," <expression>)* ")"

<if_stmt> ::= "if" "(" <expression> ")" <scope> ("else" <scope>)?
<while_stmt> ::= "while" "(" <expression> ")" <statement>

### EXPRESSIONS (Operator Precedence)
// └── Logical (||, &&)
//      └── Comparison (==, !=, <, >, etc.)
//          └── Additive (+, -, &, |, ^) (BinaryOp)
//              └── Multiplicative (*, /, %, <<, >>) (BinaryOp)
//                  └── Unary (+, -, !, ~, &, *)
//                      └── Primary (literals, identifiers, etc.)
<expression> ::= <logical>
<logical> ::= <comparison> (("||" | "&&") <comparison>)*
<comparison> ::= <additive> (("==" | "!=" | "<" | "<=" | ">" | ">=") <additive>)*
<additive> ::= <multiplicative> (("+" | "-" | "|" | "^") <multiplicative>)*
<multiplicative> ::= <unary> (("*" | "/" | "%" | "<<" | ">>") <unary>)*
<unary> ::= ("+" | "-" | "!" | "~" | "&" | "*") <unary> | <primary>
<primary> ::= <literal> | <identifier> | "(" <expression> ")" | <function_call> | <member_access> | <array_access>

<function_call> ::= <identifier> "(" <expression_list>? ")"
<expression_list> ::= <expression> ("," <expression>)*

### MOST GRANULAR COMPONENTS
<literal> ::= <integer_literal> | <float_literal> | <string_literal> | <bool_literal>
<integer_literal> ::= e.g (-1, 0, 1, 2, 3, ...)
<float_literal> ::= e.g (-1.01, 0.00, 1.01, 2.02, 3.03, ...)
<string_literal> ::= e.g ("Hello", "World")
<bool_literal> ::= "true" | "false"
<identifier> ::= e.g(name, test, foo, bar)