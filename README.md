# JSON Comparison Tool (jcmp)

## Description

`jcmp` is a command-line utility designed to compare two JSON files, highlighting differences with precision and clarity. It supports various visual styles to represent differences, including additions, deletions, and type mismatches. The tool is particularly useful for developers and QA engineers who need to track changes in JSON data structures or verify JSON outputs.

## Features

- **Deep Comparison**: Recursively compares nested JSON objects and arrays.
- **Color-coded Output**: Differences are displayed with color coding for easy identification.
- **Blacklist Paths**: Ability to ignore specified JSON paths during the comparison.
- **Strict Mode**: Option to hide paths that have equal values in both JSON files.
- **Customizable Output**: Supports disabling colored output for better compatibility with logs or files.

## Installation

To use `jcmp`, clone the repository and build the binary using Go:

```
git clone https://github.com/phplego/jcmp.git
cd jcmp
go build
```

## Usage

After building the tool, you can run it directly from the command line:

```
./jcmp <file1.json> <file2.json> [options]
```

![image](https://github.com/user-attachments/assets/ab7663d3-3449-412c-a6d0-c48bb7ee0748)

Output:

![image](https://github.com/user-attachments/assets/21f15e69-fbdf-4718-83af-cfce61cf02db)


## Output Explanation

- `+ADD`: Path exists only in the second file.
- `-DEL`: Path exists only in the first file.
- `!TYP`: Type mismatch at the given path.
- `:EXS`: Path exists in both files (shown only when not in strict mode).
- `=EQL`: Path values are equal (shown only when not in strict mode).
- `!EQL`: Path values are different.
- `!BLK`: Path is blacklisted.


### Options

- `-p`: Path within the JSON files to specifically compare. If omitted, compares entire files.
- `-bl`: Comma-separated list of JSON paths to ignore in the comparison.
- `-s`: Enable strict mode, which omits paths that have equal values.
- `-nc`: Disable color output, useful for redirecting output to a text file or log.

### Example

```
./jcmp data1.json data2.json -s -bl "meta,config.version"
```

This command compares `data1.json` and `data2.json` in strict mode, ignoring changes in the `meta` and `config.version` paths.

## Contributions

Contributions are welcome. Please fork the repository, make your changes, and submit a pull request.

## License

MIT License.
