package main

import (
	"fmt"
	"os"
	"os/user"
	"regexp"
)

type data_t map[string][]string

var config = map[string]string{
	"debug":   "false",
	"user":    "",
	"passwd":  "",
	"server":  "http://localhost:8081/",
	"query":   "Status!='resolved' and Status!='rejected'",
	"orderby": "id",
	"queue":   "Unknown_Queue",
	"auth":    "rt",
}

// The RT type names.
var types = []string{"group", "queue", "ticket", "user"}

// Number of arguments on the command line.
var argc int = len(os.Args)

// User's home directory
var home_dir string

// The action specified on the command line.
var action string

var version string = "0.02"

// Regular expressions used by various parsing routines.
var cf_name string = `[^,]+?`
var name string = `[\w.-]+`
var label string = `[^,\\/]+`
var labels string = `(?:` + label + `,)*` + label
var idlist string = `(?:(?:\d+-)?\d+,)*(?:\d+-)?\d+`
var field string = `(?i:[a-z][a-z0-9_-]*|C(?:ustom)?F(?:ield)?-` + cf_name + `|CF\.\{` + cf_name + `\})`
var fields string = `^(?:(?:` + field + `,)*` + field + `)$`

// The command line element that is currently being examined.
// Once a function is finished doing so, it is responsible for
// incrementing this variable. Its initial value points to the specified
// action, if any.
var cmd_line_ndx int = 1

// The REST endpoint on the web server.
var rest string

var prompt string = "rt> "

var spec_re *regexp.Regexp

func main() {
	var username string
	var usr *user.User
	var fn func()
	var ok bool

	// As a special case, if "--help" or "-h" are given on the command line
	// just show the help message, then exit.
	if argc == 2 {
		action = os.Args[cmd_line_ndx]
		if (action == "--help") || (action == "-h") {
			usage()
			os.Exit(1)
		}
	}

	// If it turns out that fn_shell doesn't need home_dir, then
	// this section can be moved lower.
	home_dir, _ = os.UserHomeDir()
	if home_dir == "" {
		home_dir = os.Getenv("HOME")
	}
	if home_dir == "" {
		home_dir = os.Getenv("LOGDIR")
	}
	if home_dir == "" {
		home_dir = os.Getenv("HOMEPATH")
	}

	usr, _ = user.Current()
	username = usr.Username
	if username == "" {
		username = os.Getenv("USER")
		if username == "" {
			username = os.Getenv("HOMEPATH")
		}
	}
	config["user"] = username

	// Values from .rtrc are processed last so they take
	// precidence over RTCONFIG.
	config_from_file(os.Getenv("RTCONFIG"))
	config_from_file(".rtrc")
	config_from_env()

	if _, ok = config["externalauth"]; ok {
		delete(config, "externalauth")
		config["auth"] = "basic"
	}

	session = newsession(home_dir + "/.rt_sessions")

	rest = config["server"] + "/REST/1.0"

	// As another special case, if no actions are given, assume
	// that the action is "shell".
	if argc == 1 || (os.Args[1] == "shell") {
		fn_shell()
		os.Exit(1)
	}

	// We've checked whether argc is 1 or 2. So, now we know where
	// the action starts.
	action = os.Args[cmd_line_ndx]

	// Look up the action in the actions map to get the function to execute.
	// fmt.Printf("Action = %s\n", action)
	fn, ok = actions[action]

	// If this action is valid, execute it. It's the responsibility
	// of each handler to parse the rest of the command line.
	if ok {
		cmd_line_ndx++
		fn()

		// If any of the session data has changed, write it to the session file.
		if session_save {
			session.save()
		}
	} else {
		fmt.Printf("gort: Unknown command '%s'\n", action)
		fmt.Printf("gort: For help, run 'gort help'\n")
	}
}
