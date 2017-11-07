// Copyright 2017 The go-libvirt Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//   http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package lvgen

// The libvirt API is divided into several categories. (Gallia est omnis divisa
// in partes tres.) The generator will output code for each category in a
// package underneath the go-libvirt directory.

import (
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"text/template"
	"unicode"
	"unicode/utf8"
)

var keywords = map[string]int{
	"hyper":    HYPER,
	"int":      INT,
	"short":    SHORT,
	"char":     CHAR,
	"bool":     BOOL,
	"case":     CASE,
	"const":    CONST,
	"default":  DEFAULT,
	"double":   DOUBLE,
	"enum":     ENUM,
	"float":    FLOAT,
	"opaque":   OPAQUE,
	"string":   STRING,
	"struct":   STRUCT,
	"switch":   SWITCH,
	"typedef":  TYPEDEF,
	"union":    UNION,
	"unsigned": UNSIGNED,
	"void":     VOID,
	"program":  PROGRAM,
	"version":  VERSION,
}

// ConstItem stores an const's symbol and value from the parser. This struct is
// also used for enums.
type ConstItem struct {
	Name string
	Val  string
}

// Generator holds all the information parsed out of the protocol file.
type Generator struct {
	// Enums holds the enum declarations. The type of enums is always int32.
	Enums []Decl
	// EnumVals holds the list of enum values found by the parser. In sunrpc as
	// in go, these are not separately namespaced.
	EnumVals []ConstItem
	// Consts holds all the const items found by the parser.
	Consts []ConstItem
	// Structs holds a list of all the structs found by the parser
	Structs []Structure
	// Typedefs hold all the type definitions from 'typedef ...' lines.
	Typedefs []Typedef
	// Unions hold all the discriminated unions
	Unions []Union
}

// Gen accumulates items as the parser runs, and is then used to produce the
// output.
var Gen Generator

// CurrentEnumVal is the auto-incrementing value assigned to enums that aren't
// explicitly given a value.
var CurrentEnumVal int64

// oneRuneTokens lists the runes the lexer will consider to be tokens when it
// finds them. These are returned to the parser using the integer value of their
// runes.
var oneRuneTokens = `{}[]<>(),=;:*`

var reservedIdentifiers = map[string]string{
	"type":   "lvtype",
	"string": "lvstring",
	"error":  "lverror",
}

// Decl records a declaration, like 'int x' or 'remote_nonnull_string str'
type Decl struct {
	Name, Type string
}

// Structure records the name and members of a struct definition.
type Structure struct {
	Name    string
	Members []Decl
}

// Typedef holds the name and underlying type for a typedef.
type Typedef struct {
	Name string
	Type string
}

// Union holds a "discriminated union", which consists of a discriminant, which
// tells you what kind of thing you're looking at, and a number of encodings.
type Union struct {
	Name             string
	DiscriminantType string
	Cases            []Case
}

// Case holds a single case of a discriminated union.
type Case struct {
	DiscriminantVal string
	Type            Decl
}

// CurrentStruct will point to a struct record if we're in a struct declaration.
// When the parser adds a declaration, it will be added to the open struct if
// there is one.
var CurrentStruct *Structure

// CurrentTypedef will point to a typedef record if we're parsing one. Typedefs
// can define a struct or union type, but the preferred for is struct xxx{...},
// so we may never see the typedef form in practice.
var CurrentTypedef *Typedef

// CurrentUnion holds the current discriminated union record.
var CurrentUnion *Union

// CurrentCase holds the current case record while the parser is in a union and
// a case statement.
var CurrentCase *Case

// Generate will output go bindings for libvirt. The lvPath parameter should be
// the path to the root of the libvirt source directory to use for the
// generation.
func Generate(proto io.Reader) error {
	lexer, err := NewLexer(proto)
	if err != nil {
		return err
	}
	go lexer.Run()
	parser := yyNewParser()
	yyErrorVerbose = true
	// Turn this on if you're debugging.
	// yyDebug = 3
	rv := parser.Parse(lexer)
	if rv != 0 {
		return fmt.Errorf("failed to parse libvirt protocol: %v", rv)
	}

	// Generate and write the output.
	constFile, err := os.Create("../constants/constants.gen.go")
	if err != nil {
		return err
	}
	defer constFile.Close()
	procFile, err := os.Create("../../libvirt.gen.go")
	if err != nil {
		return err
	}
	defer procFile.Close()

	err = genGo(constFile, procFile)

	return err
}

func genGo(constFile, procFile io.Writer) error {
	t, err := template.ParseFiles("constants.tmpl")
	if err != nil {
		return err
	}
	if err = t.Execute(constFile, Gen); err != nil {
		return err
	}

	t, err = template.ParseFiles("procedures.tmpl")
	if err != nil {
		return err
	}
	if err := t.Execute(procFile, Gen); err != nil {
		return err
	}
	// Now generate the wrappers for libvirt's various public API functions.
	// for _, c := range Gen.Enums {
	// This appears to be the name of a libvirt procedure, so sort it into
	// the right list based on the next part of its name.
	// segs := camelcase.Split(c.Name)
	// if len(segs) < 3 || segs[0] != "Proc" {
	// 	continue
	// }
	//category := segs[1]

	//fmt.Println(segs)
	// }

	return nil
}

// constNameTransform changes an upcased, snake-style name like
// REMOTE_PROTOCOL_VERSION to a comfortable Go name like ProtocolVersion. It
// also tries to upcase abbreviations so a name like DOMAIN_GET_XML becomes
// DomainGetXML, not DomainGetXml.
func constNameTransform(name string) string {
	nn := fromSnakeToCamel(strings.TrimPrefix(name, "REMOTE_"), true)
	nn = fixAbbrevs(nn)
	return nn
}

func identifierTransform(name string) string {
	nn := strings.TrimPrefix(name, "remote_")
	nn = fromSnakeToCamel(nn, false)
	nn = fixAbbrevs(nn)
	nn = checkIdentifier(nn)
	return nn
}

func typeTransform(name string) string {
	nn := strings.TrimLeft(name, "*")
	diff := len(name) - len(nn)
	nn = identifierTransform(nn)
	return name[0:diff] + nn
}

// fromSnakeToCamel transmutes a snake-cased string to a camel-cased one. All
// runes that follow an underscore are up-cased, and the underscores themselves
// are omitted.
//
// ex: "PROC_DOMAIN_GET_METADATA" -> "ProcDomainGetMetadata"
func fromSnakeToCamel(s string, public bool) string {
	buf := make([]rune, 0, len(s))
	// Start rune may be either upper or lower case.
	hump := public

	for _, r := range s {
		if r == '_' {
			hump = true
		} else {
			var transform func(rune) rune
			if hump == true {
				transform = unicode.ToUpper
			} else {
				transform = unicode.ToLower
			}
			buf = append(buf, transform(r))
			hump = false
		}
	}

	return string(buf)
}

// abbrevs is a list of abbreviations which should be all upper-case in a name.
// (This is really just to keep the go linters happy and to produce names that
// are intuitive to a go developer.)
var abbrevs = []string{"Xml", "Io", "Uuid", "Cpu", "Id", "Ip"}

// fixAbbrevs up-cases all instances of anything in the 'abbrevs' array. This
// would be a simple matter, but we don't want to upcase an abbreviation if it's
// actually part of a larger word, so it's not so simple.
func fixAbbrevs(s string) string {
	for _, a := range abbrevs {
		for loc := 0; ; {
			loc = strings.Index(s[loc:], a)
			if loc == -1 {
				break
			}
			r := 'A'
			if len(a) < len(s[loc:]) {
				r, _ = utf8.DecodeRune([]byte(s[loc+len(a):]))
			}
			if unicode.IsLower(r) == false {
				s = s[:loc] + strings.Replace(s[loc:], a, strings.ToUpper(a), 1)
			}
			loc++
		}
	}
	return s
}

//---------------------------------------------------------------------------
// Routines called by the parser's actions.
//---------------------------------------------------------------------------

// StartEnum is called when the parser has found a valid enum.
func StartEnum(name string) {
	// Enums are always signed 32-bit integers.
	name = identifierTransform(name)
	Gen.Enums = append(Gen.Enums, Decl{name, "int32"})
	// Set the automatic value var to -1; it will be incremented before being
	// assigned to an enum value.
	CurrentEnumVal = -1
}

// AddEnumVal will add a new enum value to the list.
func AddEnumVal(name, val string) error {
	ev, err := parseNumber(val)
	if err != nil {
		return fmt.Errorf("invalid enum value %v = %v", name, val)
	}
	return addEnumVal(name, ev)
}

// AddEnumAutoVal adds an enum to the list, using the automatically-incremented
// value. This is called when the parser finds an enum definition without an
// explicit value.
func AddEnumAutoVal(name string) error {
	CurrentEnumVal++
	return addEnumVal(name, CurrentEnumVal)
}

func addEnumVal(name string, val int64) error {
	name = constNameTransform(name)
	Gen.EnumVals = append(Gen.EnumVals, ConstItem{name, fmt.Sprintf("%d", val)})
	CurrentEnumVal = val
	return nil
}

// AddConst adds a new constant to the parser's list.
func AddConst(name, val string) error {
	_, err := parseNumber(val)
	if err != nil {
		return fmt.Errorf("invalid const value %v = %v", name, val)
	}
	name = constNameTransform(name)
	Gen.Consts = append(Gen.Consts, ConstItem{name, val})
	return nil
}

// parseNumber makes sure that a parsed numerical value can be parsed to a 64-
// bit integer.
func parseNumber(val string) (int64, error) {
	base := 10
	if strings.HasPrefix(val, "0x") {
		base = 16
		val = val[2:]
	}
	n, err := strconv.ParseInt(val, base, 64)
	return n, err
}

// StartStruct is called from the parser when a struct definition is found, but
// before the member declarations are processed.
func StartStruct(name string) {
	name = identifierTransform(name)
	CurrentStruct = &Structure{Name: name}
}

// AddStruct is called when the parser has finished parsing a struct. It adds
// the now-complete struct definition to the generator's list.
func AddStruct() {
	Gen.Structs = append(Gen.Structs, *CurrentStruct)
	CurrentStruct = nil
}

func StartTypedef() {
	CurrentTypedef = &Typedef{}
}

// TODO: remove before flight
func Beacon(name string) {
	fmt.Println(name)
}

// StartUnion is called by the parser when it finds a union declaraion.
func StartUnion(name string) {
	name = identifierTransform(name)
	CurrentUnion = &Union{Name: name}
}

// AddUnion is called by the parser when it has finished processing a union
// type. It adds the union to the generator's list and clears the CurrentUnion
// pointer.
func AddUnion() {
	Gen.Unions = append(Gen.Unions, *CurrentUnion)
	CurrentUnion = nil
}

func StartCase(dvalue string) {
	CurrentCase = &Case{DiscriminantVal: dvalue}
}

func AddCase() {
	CurrentUnion.Cases = append(CurrentUnion.Cases, *CurrentCase)
	CurrentCase = nil
}

// AddDeclaration is called by the parser when it find a declaration (int x).
// The declaration will be added to any open container (such as a struct, if the
// parser is working through a struct definition.)
func AddDeclaration(identifier, itype string) {
	// TODO: panic if not in a struct/union/typedef?
	// If the name is a reserved word, transform it so it isn't.
	identifier = identifierTransform(identifier)
	itype = typeTransform(itype)
	decl := &Decl{Name: identifier, Type: itype}
	if CurrentStruct != nil {
		CurrentStruct.Members = append(CurrentStruct.Members, *decl)
	} else if CurrentTypedef != nil {
		CurrentTypedef.Name = identifier
		CurrentTypedef.Type = itype
		Gen.Typedefs = append(Gen.Typedefs, *CurrentTypedef)
		CurrentTypedef = nil
	} else if CurrentCase != nil {
		CurrentCase.Type = *decl
	} else if CurrentUnion != nil {
		CurrentUnion.DiscriminantType = itype
	}
}

func checkIdentifier(i string) string {
	nn, reserved := reservedIdentifiers[i]
	if reserved {
		return nn
	}
	return i
}
