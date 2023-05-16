package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"mime"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type fake_res_t struct {
	frs_code         int
	frs_content      string
	frs_content_type string
	frs_message      string
}

var (
	server_re,
	text_header_re,
	cred_re *regexp.Regexp
)

func submit(uri string, content data_t) *fake_res_t {
	var data []string
	var how, server, head, text, resp_body, content_type, status, message string
	var res *http.Response
	var fake_res *fake_res_t
	var body *bytes.Buffer
	var writer *multipart.Writer
	var code, attachment_num int
	var success bool
	var cookie, pass string
	var req *http.Request
	var client *http.Client
	var err error
	var m [][]string
	var resBody []byte
	var warning string = "   Password will be sent to %s %s\n   Press CTRL-C now if you do not want to continue\n"

	for key, value := range content {

		if len(value) == 1 {
			data = append(data, key, value[0])
		} else {
			var str string

			for _, str = range value {
				data = append(data, key, str)
			}
		}
	}

	// Should we send authentication information to start a new session?
	if strings.HasPrefix(config["server"], "https") {
		how = "over SSL"
	} else {
		how = "unencrypted"
	}

	// Extract the server name from the server configuration variable.
	if server_re == nil {
		server_re = regexp.MustCompile(`^.*//([^/]+)`)
	}
	server = server_re.FindString(config["server"])
	// fmt.Printf("server = %s how = %s \n", server, how)

	// GSSAPI isn't supported, for now.
	if config["auth"] == "gssapi" {
		fmt.Fprintf(os.Stderr, "GSSAPI support not available.\n")
		os.Exit(1)
	}

	pass = config["passwd"]
	if session.cookie() == "" {
		fmt.Printf(warning, server, how)
		data = append(data, "user", config["user"])
		if pass == "" {
			pass = read_passwd()
		}
		data = append(data, "pass", pass)
	}

	body = &bytes.Buffer{}
	writer = multipart.NewWriter(body)
	writer.SetBoundary("xYzZY")

	attachment_num = 1
	for i := 0; i < len(data); i += 2 {
		if data[i] == "attachment" {
			var filename, mimetype, fieldname string
			var h textproto.MIMEHeader
			var part io.Writer
			var file *os.File
			var err error

			filename = data[i+1]
			fieldname = "attachment_" + strconv.Itoa(attachment_num)
			mimetype = mime.TypeByExtension(filepath.Ext(filename))
			h = make(textproto.MIMEHeader)
			h.Set("Content-Disposition", fmt.Sprintf(`form-data; name="%s"; filename="%s"`, fieldname, filename))
			h.Set("Content-Type", mimetype)
			part, _ = writer.CreatePart(h)

			attachment_num++
			file, err = os.Open(filename)
			if err != nil {
				panic(err)
			}
			_, err = io.Copy(part, file)
			if err != nil {
				panic(err)
			}
			file.Close()
			continue
		}
		writer.WriteField(data[i], data[i+1])
	}
	writer.Close()

	if len(data) > 0 {
		req, err = http.NewRequest(http.MethodPost, uri, bytes.NewReader(body.Bytes()))
	} else {
		req, err = http.NewRequest(http.MethodGet, uri, nil)
	}

	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("User-Agent", "RT/3.0b")
	if config["auth"] == "basic" {
		fmt.Printf(warning, server, how)
		if pass == "" {
			pass = read_passwd()
		}
		req.SetBasicAuth(config["user"], pass)
	}

	client = &http.Client{
		// Set timeout to not be at mercy of microservice to respond and stall the server
		Timeout: time.Second * 30,
	}

	// Add the session key to the header.
	session.add_cookie_header(req)

	res, err = client.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	resBody, err = io.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	content_type = res.Header.Get("Content-Type")
	cookie = res.Header.Get("Set-Cookie")
	res.Body.Close()

	// resBody is a []byte so convert it to a string.
	resp_body = string(resBody)

	// The content of the response we get from the RT server consists
	// of an HTTP-like status line followed by optional header lines,
	// a blank line, and arbitrary text.

	// head is the header of the fake response.
	// text is the body of the fake response.
	head, text, _ = strings.Cut(resp_body, "\n\n")

	// fmt.Printf("Status: %d\n", res.StatusCode)
	// fmt.Printf("head = %s\ntext= %s\n", head, text)
	status, _, _ = strings.Cut(head, "\n")
	// fmt.Printf("status = %s\n", status)

	// Get rid of any trailing newlines from text except the last one.
	if text != "" {
		if text_header_re == nil {
			text_header_re = regexp.MustCompile(`\n*$`)
		}
		text = text_header_re.ReplaceAllString(text, "\n")
	}

	if cred_re == nil {
		cred_re = regexp.MustCompile(`^RT/\d+(?:\S+) (\d+) ([\w\s]+)$`)
	}

	// "RT/3.0.1 401 Credentials required"
	m = cred_re.FindAllStringSubmatch(status, -1)
	if len(m) == 0 {
		var debug_int int

		fmt.Fprintf(os.Stderr, "rt: Malformed RT response from %s.\n", server)

		debug_int, _ = strconv.Atoi(config["debug"])
		if debug_int < 3 {
			fmt.Fprintf(os.Stderr, "(Rerun with RTDEBUG=3 for details.)\n")
		}
		os.Exit(-1)
	}

	// Put the fake res assignments here.
	code, _ = strconv.Atoi(m[0][1])
	message = m[0][2]

	fake_res = new(fake_res_t)
	(*fake_res).frs_code = code
	(*fake_res).frs_message = message
	(*fake_res).frs_content = text
	(*fake_res).frs_content_type = content_type

	success = (code >= 200) && (code < 300)
	if success || (code != 401) {
		session.update(cookie)
	}

	if !success {
		// We can deal with authentication failures ourselves. Either
		// we sent invalid credentials, or our session has expired.
		if code == 401 {
			// Check to see if a user was specified.
			for _, str := range data {
				if str == "user" {
					fmt.Fprintf(os.Stderr, "rt: Incorrect username or password.\n")
					os.Exit(-1)
				}
			}

			if session.cookie() != "" {
				// We'll retry the request with credentials, unless
				// we only wanted to logout in the first place.
				session.del()
				if uri != (rest + "/logout") {
					return submit(uri, content)
				}
			}
		}
		//  Conflicts should be dealt with by the handler and user.
		// For anything else, we just die.
	} else if code != 409 {
		return fake_res
	}

	return fake_res
}
