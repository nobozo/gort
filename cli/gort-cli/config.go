package main

import (
	"bufio"
	"errors"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

var conf_re *regexp.Regexp

func config_from_env() {
	var conf_vars = [...]string{"EXTERNALAUTH", "AUTH", "DEBUG", "USER", "PASSWD", "SERVER", "QUERY", "ORDERBY"}
	var var_name, env_val string

	for _, conf_var := range conf_vars {
		var_name = "RT" + conf_var
		env_val = os.Getenv(var_name)
		conf_var = strings.ToLower(conf_var)
		if env_val != "" {
			config[conf_var] = env_val
		}
	}
}

func config_from_file(conf_file string) string {
	var file, cwd string
	var err error
	var loc int

	if len(conf_file) == 0 {
		return ""
	}

	if (len(conf_file) >= 1) && filepath.IsAbs(conf_file) {
		return parse_config_file(conf_file)
	} else {
		cwd, err = os.Getwd()
		if err != nil {
			log.Println(err)
		}
	}

	loc = strings.LastIndex(cwd, "/")
	for loc >= 0 {
		var err error
		file = cwd + "/" + conf_file
		if _, err = os.Stat(file); !errors.Is(err, fs.ErrNotExist) {
			return parse_config_file(file)
		}
		cwd = cwd[0:loc]
		loc = strings.LastIndex(cwd, "/")
	}

	return ""
}

// Makes a hash of the specified configuration file.
func parse_config_file(config_file string) string {
	var readFile *os.File
	var fileScanner *bufio.Scanner
	var err error

	readFile, err = os.Open(config_file)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fileScanner = bufio.NewScanner(readFile)

	fileScanner.Split(bufio.ScanLines)

	if conf_re == nil {
		conf_re = regexp.MustCompile(`^(externalauth|auth|user|passwd|server|query|orderby|queue)\s+(.*)\s?$`)
	}
	for fileScanner.Scan() {
		var line string
		var matches [][]string

		line = fileScanner.Text()

		// Get rid of comment lines or lines that are all whitespace.
		if (line[0:1] == "#") || (strings.TrimSpace(line) == "") {
			continue
		}

		matches = conf_re.FindAllStringSubmatch(line, -1)
		if len(matches) > 0 {
			config[matches[0][1]] = matches[0][2]
		} else {
			fmt.Fprintf(os.Stderr, "rt: %s:%s: unknown configuration directive.\n", config_file, line)
			os.Exit(1)
		}
	}

	readFile.Close()
	return ""
}
