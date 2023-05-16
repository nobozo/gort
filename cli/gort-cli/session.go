package main

import (
	"bufio"
	"fmt"
	"net/http"
	"os"
	"regexp"
	"strings"
)

// Map a user name to a cookie.
type user_map map[string]string

// Map a server name to a user map.
type server_map map[string]user_map

type session_t struct {
	session_file string     // The file that stores the session info.
	sm           server_map // Given a server name, find the users who connect to it.
}

var session *session_t

var session_server, session_user string

// True if any of the session state was changed, thus requiring that it be
// saved to the session file bu calling save().
var session_save bool = false

var (
	url_re,
	cookie_re *regexp.Regexp
)

// Adds a Cookie header to an outgoing HTTP request.
func (s *session_t) add_cookie_header(req *http.Request) {
	var cookie string

	// If there's a session cookie, add it to the request header.
	cookie = s.cookie()
	if cookie != "" {
		req.Header.Set("Cookie", cookie)
	}
}

// Returns the current session cookie.
func (s *session_t) cookie() string {
	var cookie string

	cookie = s.sm[session_server][session_user]
	if cookie != "" {
		return "RT_SID_" + cookie
	} else {
		return cookie
	}
}

// Deletes the current session cookie.
func (s *session_t) del() {
	s.sm[session_server][session_user] = ""
}

// Loads the session cache from the specific file.
func (s *session_t) load() {
	var um user_map
	var sm server_map
	var readFile *os.File
	var fileScanner *bufio.Scanner
	var err error

	um = make(user_map)
	sm = make(server_map)

	// If the session cache file already exists, open it. Otherwise, create it.
	readFile, err = os.OpenFile(s.session_file, os.O_CREATE|os.O_RDONLY, 0644)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fileScanner = bufio.NewScanner(readFile)
	fileScanner.Split(bufio.ScanLines)

	if url_re == nil {
		url_re = regexp.MustCompile(`^https?://[^ ]+ \w+ [^;,\s]+=[0-9A-Fa-f]+$`)
	}

	for fileScanner.Scan() {
		var line, matched string
		var matches []string
		var server, user, cookie string

		line = fileScanner.Text()
		if (line[0:1] == "#") || line[0:1] == "$" {
			continue
		}
		matched = url_re.FindString(line)

		// If the pattern wasn't matched, get the next line.
		if matched == "" {
			continue
		} else {
			matches = strings.Split(matched, " ")
		}

		if len(matches) != 3 {
			fmt.Printf("rt: %s: %s - bad input line\n", s.session_file, line)
			os.Exit(1)
		}
		server = matches[0]
		user = matches[1]
		cookie = matches[2]

		um[user] = cookie
		sm[server] = um
		s.sm = sm
	}
	readFile.Close()
}

// "f" is the file the session info is kept in.
func newsession(f string) *session_t {
	var s *session_t

	// Create the maps used by the session cache.
	s = new(session_t)
	s.sm = make(server_map)

	// Figure out the name of the sessions file if it wasn't supplied
	// as an argument.
	if f != "" {
		s.session_file = f
	} else {
		s.session_file = home_dir + "/.rt_sessions"
	}

	session_server = config["server"]
	session_user = config["user"]
	s.load()

	return s
}

// Writes the current session cache to the specified file.
// Only one line will exist in the file - is that a problem?
func (s *session_t) save() {
	var writeFile *os.File
	var err error

	writeFile, err = os.OpenFile(s.session_file, os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Fprintf(writeFile, "%s %s %s\n", session_server, session_user, s.sm[session_server][session_user])
	writeFile.Close()
}

// Updates session information.
func (s *session_t) update(cookie string) {
	var matched [][]string

	if cookie_re == nil {
		cookie_re = regexp.MustCompile(`^RT_SID_(.[^;,\s]+=[0-9A-Fa-f]+);`)
	}

	if cookie != "" {
		matched = cookie_re.FindAllStringSubmatch(cookie, -1)
		if len(matched) > 0 {
			if s.sm[session_server] == nil {
				s.sm[session_server] = make(user_map)
			}
			s.sm[session_server][session_user] = matched[0][1]
		}

		// Show that the cache has been modified so that it needs to be saved
		// in the session file.
		session_save = true
	}
}
