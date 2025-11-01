package main

// csv2json --input /etc/passwd --delimit : -names username,password,uid,gid,gecos,home_dir,shell

import (
	"encoding/csv"
	"flag"
	"fmt"
	ep "github.com/PeterHickman/expand_path"
	"github.com/PeterHickman/toolbox"
	"os"
	"strconv"
	"strings"
)

var delimiter = ','
var header = false
var names = false
var headers []string
var filename string

func dropdead(message string) {
	fmt.Println(message)
	os.Exit(1)
}

func maybe_comma(index int) {
	if index == 0 || (index == 1 && header) {
		fmt.Println()
	} else {
		fmt.Println(",")
	}
}

func is_int(value string) bool {
	_, err := strconv.ParseInt(value, 10, 64)
	return err == nil
}

func is_float(value string) bool {
	_, err := strconv.ParseFloat(value, 64)
	return err == nil
}

func has_quotes(value string) bool {
	i := strings.Index(value, "\"")
	return i != -1
}

func embed_quotes(value string) string {
	var new_value []byte

	new_value = append(new_value, '"')

	for i := 0; i < len(value); i++ {
		if value[i] == '"' {
			new_value = append(new_value, '\\')
		}
		new_value = append(new_value, value[i])
	}

	new_value = append(new_value, '"')

	return string(new_value)
}

func formatted_value(value string) string {
	lower_value := strings.ToLower(value)

	if lower_value == "true" || lower_value == "false" {
		return lower_value
	} else if lower_value == "nil" || lower_value == "null" {
		return "null"
	} else if is_int(lower_value) || is_float(lower_value) {
		return lower_value
	} else if has_quotes(value) {
		return embed_quotes(value)
	} else {
		return fmt.Sprintf("\"%s\"", value)
	}
}

func init() {
	h := flag.Bool("header", false, "The first row of the CSV file is the column names")
	d := flag.String("delimit", ",", "The character that delimit the columns")
	n := flag.String("names", "", "A comma separated list of names to use for the columns")
	f := flag.String("input", "", "The csv to read in")

	flag.Parse()

	if *f == "" {
		dropdead("Need to provide a file to read with --input")
	}

	var err error
	filename, err = ep.ExpandPath(*f)

	if err != nil {
		dropdead(fmt.Sprintf("Unable to read %s, %s\n", *f, err))
	}

	if !toolbox.FileExists(filename) {
		dropdead(fmt.Sprintf("Unable to read %s\n", filename))
	}

	header = *h

	if *d == "\\t" {
		*d = "\t"
	}

	if len(*d) > 1 {
		dropdead(fmt.Sprintf("Column delimiter should be a single character [%s]\n", *d))
	}

	delimiter = []rune(*d)[0]

	if *n != "" {
		if header {
			dropdead("Use of --names overrides --header")
		}
		headers = strings.Split(*n, ",")
		header = true
		names = true
	}
}

func main() {
	file, err := os.Open(filename)
	if err != nil {
		dropdead(fmt.Sprintf("Error while reading the file. Do all rows have the same number of columns? %s", err))
	}

	reader := csv.NewReader(file)

	reader.Comma = delimiter
	reader.Comment = '#'

	records, err := reader.ReadAll()
	if err != nil {
		dropdead(fmt.Sprintf("Error reading records: %s\n", err))
	}

	file.Close()

	fmt.Print("[")
	for i, eachrecord := range records {
		if i == 0 && header {
			if !names {
				headers = eachrecord
			}
			continue
		}

		if header && len(eachrecord) > len(headers) {
			dropdead(fmt.Sprintf("Data row has more columns (%d) than the headers (%d)\n", len(eachrecord), len(headers)))
		}

		maybe_comma(i)
		fmt.Print("  {")
		for k, v := range eachrecord {
			if k > 0 {
				fmt.Println(",")
			} else {
				fmt.Println()
			}
			if header {
				fmt.Printf("    \"%s\": %s", headers[k], formatted_value(v))
			} else {
				fmt.Printf("    \"column%d\": %s", k+1, formatted_value(v))
			}
		}
		fmt.Print("\n  }")
	}
	fmt.Println("\n]")
}
