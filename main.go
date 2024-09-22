package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/alecthomas/participle/v2"
  "github.com/alecthomas/participle/v2/lexer"
)

// Define the lexer rules using MustStateful
var lex = lexer.MustStateful(lexer.Rules{
    "Root": {
        {"Comment", `[#;!].*`, nil},
        {"Whitespace", `[ \t\r\n]+`, nil},
        {"Ident", `[A-Za-z0-9._-]+`, nil},
        {"Equal", `=`, nil},
        {"LParen", `\(`, nil},
        {"RParen", `\)`, nil},
        {"String", `[^()=#\s][^()=#]*`, nil},
    },
})
// TNSFile represents the entire tnsnames.ora file
type TNSFile struct {
	Entries []*Entry `{ @@ }`
}

// Entry represents a service entry
type Entry struct {
	Name        string       `@Ident "="`
	Description *Description `@@`
}

// Description represents the description of a service
type Description struct {
	LParen   string     `"(" "DESCRIPTION" "="`
	Elements []*Element `{ @@ }`
	RParen   string     `")"`
}

// Element can be an Address, AddressList, or ConnectData
type Element struct {
	Address     *Address     `( @@`
	AddressList *AddressList `| @@`
	ConnectData *ConnectData `| @@ )`
}

// Address represents an address block
type Address struct {
	LParen string      `"(" "ADDRESS" "="`
	Params []*KeyValue `{ @@ }`
	RParen string      `")"`
}

// AddressList represents an address list block
type AddressList struct {
	LParen   string     `"(" "ADDRESS_LIST" "="`
	Elements []*Address `{ @@ }`
	RParen   string     `")"`
}

// ConnectData represents the connect data block
type ConnectData struct {
	LParen string      `"(" "CONNECT_DATA" "="`
	Params []*KeyValue `{ @@ }`
	RParen string      `")"`
}

// KeyValue represents a key-value pair
type KeyValue struct {
	LParen string `"("`
	Key    string `@Ident`
	Equal  string `"="`
	Value  *Value `@@`
	RParen string `")"`
}

// Value represents a value which can be a nested KeyValue or a string
type Value struct {
	String   string    `( @Ident | @String`
	KeyValue *KeyValue `| @@ )`
}

// Create the parser
var parser = participle.MustBuild[TNSFile](
    participle.Lexer(lex),
    participle.Unquote(),
    participle.Elide("Whitespace", "Comment"),
)
func parseTNSFile(filename string) (*TNSFile, error) {
    content, err := os.ReadFile(filename)
    if err != nil {
        return nil, err
    }
    // Remove any carriage returns (for Windows compatibility)
    input := strings.ReplaceAll(string(content), "\r", "")
    // Parse the content
    tnsFile, err := parser.ParseString("", input)
    if err != nil {
        return nil, err
    }
    return tnsFile, nil
}

// Function to search entries by pattern
func searchEntries(tnsFile *TNSFile, pattern string) ([]*Entry, error) {
	var matches []*Entry
	for _, entry := range tnsFile.Entries {
		match, err := filepath.Match(pattern, entry.Name)
		if err != nil {
			return nil, err
		}
		if match {
			matches = append(matches, entry)
		}
	}
	return matches, nil
}

func printEntry(entry *Entry) {
	fmt.Printf("%s =\n", entry.Name)
	printDescription(entry.Description, "  ")
}

func printDescription(desc *Description, indent string) {
	fmt.Printf("%s(DESCRIPTION =\n", indent)
	for _, elem := range desc.Elements {
		if elem.Address != nil {
			printAddress(elem.Address, indent+"  ")
		} else if elem.AddressList != nil {
			printAddressList(elem.AddressList, indent+"  ")
		} else if elem.ConnectData != nil {
			printConnectData(elem.ConnectData, indent+"  ")
		}
	}
	fmt.Printf("%s)\n", indent)
}

func printAddress(addr *Address, indent string) {
	fmt.Printf("%s(ADDRESS =\n", indent)
	for _, param := range addr.Params {
		printKeyValue(param, indent+"  ")
	}
	fmt.Printf("%s)\n", indent)
}

func printAddressList(addrList *AddressList, indent string) {
	fmt.Printf("%s(ADDRESS_LIST =\n", indent)
	for _, addr := range addrList.Elements {
		printAddress(addr, indent+"  ")
	}
	fmt.Printf("%s)\n", indent)
}

func printConnectData(connectData *ConnectData, indent string) {
	fmt.Printf("%s(CONNECT_DATA =\n", indent)
	for _, param := range connectData.Params {
		printKeyValue(param, indent+"  ")
	}
	fmt.Printf("%s)\n", indent)
}

func printKeyValue(kv *KeyValue, indent string) {
	if kv.Value.String != "" {
		fmt.Printf("%s(%s = %s)\n", indent, kv.Key, kv.Value.String)
	} else if kv.Value.KeyValue != nil {
		fmt.Printf("%s(%s =\n", indent, kv.Key)
		printKeyValue(kv.Value.KeyValue, indent+"  ")
		fmt.Printf("%s)\n", indent)
	}
}

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: go run main.go <pattern> <tnsnames.ora>")
		return
	}
	pattern := os.Args[1]
	filename := os.Args[2]

	tnsFile, err := parseTNSFile(filename)
	if err != nil {
		fmt.Println("Error parsing tnsnames.ora:", err)
		return
	}

	matches, err := searchEntries(tnsFile, pattern)
	if err != nil {
		fmt.Println("Error searching entries:", err)
		return
	}

	if len(matches) == 0 {
		fmt.Println("No matching entries found.")
		return
	}

	for _, entry := range matches {
		printEntry(entry)
		fmt.Println()
	}
}
