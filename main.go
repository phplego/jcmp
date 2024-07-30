package main

import (
	"flag"
	"fmt"
	c "github.com/fatih/color"
	"github.com/tidwall/gjson"
	"os"
	"strings"
)

type Styles struct {
	Eql, Neq, Add, Del, Typ, Exs, Df1, Df2, Blk, Sep *c.Color
}

var style = Styles{ // Define styles
	Eql: c.New(c.FgHiCyan),               // Equal
	Neq: c.New(c.FgHiYellow),             // Not equal
	Add: c.New(c.FgHiGreen),              // Added
	Del: c.New(c.FgHiRed),                // Deleted
	Typ: c.New(c.FgHiMagenta),            // Type mismatch
	Exs: c.New(c.FgHiBlue),               // Exists
	Df1: c.New(c.FgHiBlack, c.BgWhite),   // Different 1
	Df2: c.New(c.FgYellow, c.BgWhite),    // Different 2
	Blk: c.New(c.FgHiWhite, c.BgHiBlack), // Blacklisted
	Sep: c.New(c.FgHiYellow),             // Separator
}

var compared = make(map[string]bool) // Paths that have been compared
var blacklist []string               // Paths to ignore
var strictMode = false               //	Strict mode (dont print equal keys)

func pp(path string) string { // Prettify path for output
	parts := strings.Split(path, `\.`)
	for i, p := range parts {
		parts[i] = strings.ReplaceAll(p, `.`, " . ")
	}
	path = strings.Join(parts, ".")
	return path
}

// Recursively compare JSON objects
func compareJSON(json1, json2 string, path string) {
	if compared[path] {
		return // prevent duplicates
	}
	compared[path] = true

	for _, b := range blacklist {
		if strings.Contains(path, b) {
			style.Blk.Println("!BLK " + pp(path))
			return
		}
	}

	var results1, results2 gjson.Result
	if path == "" {
		// If the path is empty, use the whole JSON
		results1 = gjson.Parse(json1)
		results2 = gjson.Parse(json2)
	} else {
		// Otherwise, get the specific path in JSON
		results1 = gjson.Get(json1, path)
		results2 = gjson.Get(json2, path)
	}

	if !results1.Exists() && !results2.Exists() {
		style.Del.Println("!EBT " + pp(path) + " (path does not exist in both JSONs)")
		return // Stop recursion as parent does not exist
	}

	// Check if either result does not exist at this path
	if !results1.Exists() {
		style.Add.Println("+ADD " + pp(path))
		return // Stop recursion as parent does not exist
	} else if !results2.Exists() {
		style.Del.Println("-DEL " + pp(path))
		return // Stop recursion as parent does not exist
	}

	if results1.Type != results2.Type {
		style.Typ.Println("!TYP " + pp(path) + " " + results1.Type.String() + " vs " + results2.Type.String())
		return // stop recursion as types are different
	}

	if results1.Exists() && results2.Exists() {
		if !strictMode {
			style.Exs.Println(":EXS " + pp(path))
		}
	}

	if results1.String() == results2.String() {
		if !strictMode {
			style.Eql.Println("=EQL " + pp(path))
		}
		return // stop recursion as values are the same
	} else {
		if !results1.IsObject() && !results2.IsObject() && !results1.IsArray() && !results2.IsArray() {
			style.Neq.Println("!EQL " + pp(path) + ":")
			const MAX_LEN = 1000
			style.Df1.Println(cut(results1.String(), MAX_LEN))
			style.Sep.Println(" --- vs ---")
			style.Df2.Println(cut(results2.String(), MAX_LEN))
		}
	}

	if results1.IsObject() { // if it is an object, iterate over each key
		results1.ForEach(func(key, value gjson.Result) bool {
			compareJSON(json1, json2, appendPath(path, key.String()))
			return true // continue iterating
		})
	}

	if results2.IsObject() { // don't worry about duplicates, they will be skipped (see `compared` map)
		results2.ForEach(func(key, value gjson.Result) bool {
			compareJSON(json1, json2, appendPath(path, key.String()))
			return true // continue iterating
		})
	}

	if results1.IsArray() { // If it is an array, iterate over each element
		results1.ForEach(func(key, value gjson.Result) bool {
			compareJSON(json1, json2, appendPath(path, key.String()))
			return true // continue iterating
		})
	}

	if results2.IsArray() {
		results2.ForEach(func(key, value gjson.Result) bool {
			compareJSON(json1, json2, appendPath(path, key.String()))
			return true // continue iterating
		})
	}
}

func cut(s string, maxLen int) string {
	if len(s) > maxLen {
		return fmt.Sprintf("%s...(%d bytes)", s[:maxLen], len(s))
	}
	return s
}

// Helper function to append keys to the base path
func appendPath(basePath, key string) string {
	if basePath == "" {
		return escapePath(key)
	}
	p := basePath + "." + escapePath(key)
	return p
}

func escapePath(path string) string {
	return strings.ReplaceAll(path, ".", "\\.")
}

// Load JSON data from a file
func loadJSON(filename string) (string, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func main() {

	flagSet := flag.NewFlagSet("flags", flag.ContinueOnError)
	optPath := flagSet.String("p", "", "Path to compare")
	optBlacklist := flagSet.String("bl", "", "Comma-separated list of paths to blacklist")
	optStrict := flagSet.Bool("s", false, "Strict mode (don't print equal keys)")
	optNoColor := flagSet.Bool("nc", false, "Disable color output")

	if len(os.Args) < 3 {
		fmt.Println("Usage: jcmp <file1.json> <file2.json> [...named options]\n\nNamed options:")
		flagSet.PrintDefaults()
		return
	}
	fmt.Println("[CMD]", strings.Join(os.Args, " ")) // print command line arguments

	flagSet.Parse(os.Args[3:])

	// apply named options
	c.NoColor = *optNoColor
	strictMode = *optStrict
	if *optBlacklist != "" {
		blacklist = strings.Split(*optBlacklist, ",")
	}

	json1, err := loadJSON(os.Args[1])
	if err != nil {
		fmt.Println("Error loading first JSON file:", err)
		return
	}

	json2, err := loadJSON(os.Args[2])
	if err != nil {
		fmt.Println("Error loading second JSON file:", err)
		return
	}

	compareJSON(json1, json2, *optPath)
}
