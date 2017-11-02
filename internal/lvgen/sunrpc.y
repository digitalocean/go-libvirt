%{

package lvgen

import (
    //"fmt"
)

%}

// SymType
%union{
    val string
}

// XDR tokens:
%token BOOL CASE CONST DEFAULT DOUBLE ENUM FLOAT OPAQUE STRING STRUCT
%token SWITCH TYPEDEF UNION UNSIGNED VOID HYPER INT SHORT CHAR
%token IDENTIFIER CONSTANT ERROR
// RPCL additional tokens:
%token PROGRAM VERSION

%%

specification
    : definition_list
    ;

value
    : IDENTIFIER
    | CONSTANT
    ;

definition_list
    : definition ';'
    | definition ';' definition_list
    ;

definition
    : enum_definition
    | const_definition
    | typedef_definition
    | struct_definition
    | union_definition
    | program_definition
    ;

enum_definition
    : ENUM enum_ident '{' enum_value_list '}' { StartEnum() }
    ;

enum_value_list
    : enum_value
    | enum_value ',' enum_value_list
    ;

enum_value
    : enum_value_ident {
        err := AddEnumAutoVal($1.val)
        if err != nil {
            yylex.Error(err.Error())
            return 1
        }
    }
    | enum_value_ident '=' value {
        err := AddEnum($1.val, $3.val)
        if err != nil {
            yylex.Error(err.Error())
            return 1
        }
    }
    ;

enum_ident
    : IDENTIFIER
    ;

enum_value_ident
    : IDENTIFIER
    ;

// Ignore consts that are set to IDENTIFIERs - this isn't allowed by the spec,
// but occurs in the file because libvirt runs the pre-processor on the protocol
// file, and it handles replacing the identifier with it's #defined value.
const_definition
    : CONST const_ident '=' IDENTIFIER
    | CONST const_ident '=' CONSTANT {
        err := AddConst($2.val, $4.val)
        if err != nil {
            yylex.Error(err.Error())
            return 1
        }
    }
    ;

const_ident
    : IDENTIFIER
    ;

typedef_definition
    : TYPEDEF declaration
    ;

declaration
    : simple_declaration
    | fixed_array_declaration
    | variable_array_declaration
    | pointer_declaration
    ;

simple_declaration
    : type_specifier variable_ident
    ;

type_specifier
    : int_spec
    | UNSIGNED int_spec
    | FLOAT
    | DOUBLE
    | BOOL
    | STRING
    | OPAQUE
    | enum_definition
    | struct_definition
    | union_definition
    | IDENTIFIER
    ;

int_spec
    : HYPER
    | INT
    | SHORT
    | CHAR
    ;

variable_ident
    : IDENTIFIER
    ;

fixed_array_declaration
    : type_specifier variable_ident '[' value ']'
    ;

variable_array_declaration
    : type_specifier variable_ident '<' value '>'
    | type_specifier variable_ident '<' '>'
    ;

pointer_declaration
    : type_specifier '*' variable_ident
    ;

struct_definition
    : STRUCT struct_ident '{' declaration_list '}'
    ;

struct_ident
    : IDENTIFIER
    ;

declaration_list
    : declaration ';'
    | declaration ';' declaration_list
    ;

union_definition
    : UNION union_ident SWITCH '(' simple_declaration ')' '{' case_list '}'
    ;

union_ident
    : IDENTIFIER
    ;

case_list
    : case ';'
    | case ';' case_list
    ;

case
    : CASE value ':' declaration
    | DEFAULT ':' declaration
    ;

program_definition
    : PROGRAM program_ident '{' version_list '}' '=' value
    ;

program_ident
    : IDENTIFIER
    ;

version_list
    : version ';'
    | version ';' version_list
    ;

version
    : VERSION version_ident '{' procedure_list '}' '=' value ';'
    ;

version_ident
    : IDENTIFIER
    ;

procedure_list
    : procedure ';'
    | procedure ';' procedure_list
    ;

procedure
    : type_specifier procedure_ident '(' type_specifier ')' '=' value ';'
    ;

procedure_ident
    : IDENTIFIER
    ;

%%
