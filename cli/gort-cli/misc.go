package main

import (
	"fmt"
	"golang.org/x/term"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"syscall"
	"time"
)

var in_time_format string = `2006-01-02 15:04:05`
var out_time_format string = `Mon Jan 02 15:04:05 2006`

func is_string_in_slice(s string, sl []string) bool {
	for _, v := range sl {
		if v == s {
			return true
		}
	}

	return false
}

func suggest_help(action string, type_name string) {
	var msg string = "rt: For help, run 'rt help %s'.\n"

	if action != "" {
		fmt.Fprintf(os.Stderr, msg, action)
	}

	if type_name != "" {
		fmt.Fprintf(os.Stderr, msg, type_name)
	}
}

func read_passwd() string {
	var pass string
	var bytepw []byte
	var err error

	fmt.Printf("Password: ")
	bytepw, err = term.ReadPassword(int(syscall.Stdin))
	if err != nil {
		os.Exit(1)
	}
	pass = string(bytepw)
	return pass
}

func prettyshow(forms []form_t) {
	var form form_t
	var ok bool
	var str string

	// In the slice of forms, there will be only one element containing a map with
	// the key "Queue". I don't understand why this is true, and I'd like to find out
	// why this is. But all my experiments show this is true.

	// Dates are in local time zone
	for _, form = range forms {
		if _, ok = form.k["Queue"]; ok {
			var tmp_str string

			fmt.Printf("Date: %s\n", form.k["Created"][0])
			fmt.Printf("From: %s\n", form.k["Requestors"][0])
			if len(form.k["Cc"]) > 0 {
				if str = form.k["Cc"][0]; len(str) > 0 {
					fmt.Printf("Cc: %s\n", str)
				}
			}
			if len(form.k["AdminCc"]) > 0 {
				if str = form.k["AdminCc"][0]; len(str) > 0 {
					fmt.Printf("X-AdminCc: %s\n", str)
				}
			}
			fmt.Printf("X-Queue: %s\n", form.k["Queue"][0])
			if len(form.k["Subject"]) > 0 {
				tmp_str = form.k["Subject"][0]
			} else {
				tmp_str = ""
			}
			fmt.Printf("Subject: [rt #%s] %s\n\n", form.k["id"][0], tmp_str)
			break
		}
	}

	// Dates in these attributes are in GMT and will be converted into the local time zone.
	for _, form = range forms {
		var str string
		var id_str, queue_str string

		if len(form.k["id"]) > 0 {
			id_str = form.k["id"][0]
		}
		if len(form.k["Queue"]) > 0 {
			queue_str = form.k["Queue"][0]
			if (len(id_str) != 0) || (len(queue_str) != 0) {
				continue
			}
		}

		if len(form.k["Created"]) > 0 {
			if str = form.k["Created"][0]; str != "" {
				var localLoc *time.Location
				var localDateTime time.Time
				var parsed_time time.Time
				var err error

				parsed_time, err = time.ParseInLocation(in_time_format, str, time.UTC)
				if err != nil {
					log.Fatal(err)
				}
				localLoc, err = time.LoadLocation("Local")
				if err != nil {
					log.Fatal(`Failed to load location "Local"`)
				}
				localDateTime = parsed_time.In(localLoc)
				if len(form.k["Description"]) > 0 {
					if str = form.k["Description"][0]; len(str) != 0 {
						fmt.Printf("===> %s on %s\n", str, localDateTime.Format(out_time_format))
					}
				}
			}
		}

		if len(form.k["Content"]) > 0 {
			str = form.k["Content"][0]

		}
		if (len(str) != 0) && !strings.HasSuffix(str, "to have no content") {
			var tmp_str string

			if len(form.k["Type"]) > 0 {
				tmp_str = form.k["Type"][0]
			} else {
				tmp_str = ""
			}
			if tmp_str != "EmailRecord" {
				fmt.Println(str)
			}
		}

		if len(form.k["Attachments"]) > 0 {
			str = form.k["Attachments"][0]
			if len(str) != 0 {
				fmt.Println(str)
			}
		}
	}
}

func prettylist(forms []form_t) {
	var form form_t
	var heading string = "Ticket Owner Queue    Age   Told Status Requestor Subject\n"
	var ok bool
	var line string
	var open, me []string

	heading += (strings.Repeat("-", 80) + "\n")

	for _, form = range forms {
		var id, owner, queue, subject, age, told, status, requestor string

		if _, ok = form.k["id"]; !ok {
			continue
		}

		if heading != "" {
			fmt.Printf("%s", heading)
		}
		heading = ""

		id = strings.TrimPrefix(form.k["id"][0], "ticket/")

		if form.k["Owner"][0] == "Nobody" {
			owner = ""
		} else {
			owner = form.k["Owner"][0]
		}

		if len(form.k["Queue"]) > 0 {
			queue = form.k["Queue"][0]
		}

		if len(form.k["Subject"]) > 0 {
			subject = form.k["Subject"][0]
		}

		if len(form.k["Created"]) > 0 {
			age = date_diff(form.k["Created"][0])
		}

		if len(form.k["Told"]) > 0 {
			if form.k["Told"][0] == "Not set" {
				told = ""
			} else {
				told = date_diff(form.k["Told"][0])
			}
		}

		if len(form.k["Status"]) > 0 {
			status = form.k["Status"][0]
		}

		if len(form.k["Requestors"]) > 0 {
			requestor = form.k["Requestors"][0]
		}

		line = fmt.Sprintf("%6s %5s %.5s %6s %6s %-6s %-.9s %-30s\n",
			id, owner, queue, age, told, status, requestor, subject)
		if form.k["Owner"][0] == "Nobody" {
			open = append(open, line)
		} else if form.k["Owner"][0] == config["user"] {
			me = append(me, line)
		} else {
			fmt.Println(line)
		}
	}

	if heading != "" {
		fmt.Printf("No matches found\n")
	}

	if len(me) > 0 {
		fmt.Printf("========== my %2d open tickets ==========\n", len(me))
		for _, line = range me {
			fmt.Printf(line)
		}
	}

	if len(open) > 0 {
		fmt.Printf("========== %2d unowned tickets ==========\n", len(open))
		for _, line = range open {
			fmt.Printf(line)
		}
	}
}

func is_object_spec(spec string, type_name string) string {
	var pattern string
	var re *regexp.Regexp
	var str string

	// Add a "/" to the end of the type_name, if one is given.
	if type_name != "" {
		// I don't do the typical conditional pattern compilation here because
		// the text of the pattern can change each time this function is called.
		pattern = `^(?:` + type_name + `/)?`
		re = regexp.MustCompile(pattern)
		spec = re.ReplaceAllString(spec, type_name+"/")
	}
	pattern = `^` + name + `/(?:` + idlist + `|` + labels + `)(?:/.*)?$`
	if spec_re == nil {
		spec_re = regexp.MustCompile(pattern)
	}

	str = spec_re.FindString(spec)

	return str
}

// Get arguments of the style var=val from the command line.
// This presumes cmd_line_ndx is already indexing the next argument
// on the command line.
func get_var_argument(data data_t) bool {
	var kv string
	var equal_ndx int

	// If there are no more arguments on the command line.
	if cmd_line_ndx >= argc {
		fmt.Fprintf(os.Stderr, "No variable argument specified with -S.\n")
		return false
	}

	kv = os.Args[cmd_line_ndx]
	// fmt.Printf("get_var_arguments looking at %s\n", cur_arg)

	equal_ndx = strings.Index(kv, "=")
	if equal_ndx != -1 {
		data[kv[:equal_ndx]] = append(data[kv[:equal_ndx]], kv[equal_ndx+1:])
		// fmt.Printf("var %s=%s\n", cur_arg[:equal_ndx], data[cur_arg[:equal_ndx]])
		return true
	} else {
		fmt.Fprintf(os.Stderr, "Invalid variable specification: %s\n", kv)
		return false
	}
}

// Checks that the current argument is a valid type name.
func get_type_argument() string {
	var type_name string

	// If there are still unhandled command line arguments.
	if cmd_line_ndx < argc {
		var s string
		var ok bool

		type_name = os.Args[cmd_line_ndx]
		// fmt.Printf("found type %s\n", type_name)

		// If the last rune in type_name is "s", then remove it.
		s = type_name[len(type_name)-1:]
		if s == "s" {
			type_name = strings.TrimSuffix(type_name, "s")
			// fmt.Printf("type without s =  %s\n", type_name)
		}

		ok = is_string_in_slice(type_name, types)
		if !ok {
			fmt.Fprintf(os.Stderr, "Invalid type '%s' specified\n", type_name)
			return ""
		} else {
			// fmt.Printf("type = %s\n", type_name)
			return type_name
		}
	} else {
		fmt.Fprintf(os.Stderr, "No type argument specified with -t.\n")
		return ""
	}
}

func vi(text string) string {
	var editor string
	var file *os.File
	var cmd *exec.Cmd
	var b []byte
	var err error

	editor = os.Getenv("EDITOR")
	if editor == "" {
		editor = os.Getenv("VISUAL")
		if editor == "" {
			editor = "vi"
		}
	}

	file, err = os.CreateTemp("", "sample")

	if err != nil {
		log.Fatal(err)
	}

	if _, err = file.WriteString(text); err != nil {
		log.Fatal(err)
	}

	cmd = exec.Command(editor, file.Name())

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()

	if err != nil {
		log.Printf("Error while editing. Error: %v\n", err)
	}

	b, err = os.ReadFile(file.Name())
	if err != nil {
		fmt.Print(err)
	}

	file.Close()
	os.Remove(file.Name())
	return (string(b))

}

// Create an editor session with the default form. After the user finishes editing it
// check to make sure it's still a valid form. If not, go again.
func vi_form_while(text string, cb func(string, form_t) (int, string)) string {
	var err bool = false
	var forms []form_t
	var form form_t
	var msg string
	var status int
	var ntext string = `# Required: id, Queue

id: ticket/new
Queue: General
Requestor: root
Subject: 
Cc:
AdminCc:
Owner: 
Status: new
Priority: 
InitialPriority: 
FinalPriority: 
TimeEstimated: 0
Starts: 2023-04-07 18:53:11
Due: 
Attachment: 
Text: 
`

top:
	// Get the text of the new form.
	ntext = vi(text)

	// If there was an error and the new form is the same as the old form,
	// then return an empty string.
	if err && (ntext == text) {
		return ""
	}

	// Parse the new form.
	text = ntext
	forms = form_parse(text)
	form = forms[0]
	err = false
	if form.e != "" {
		err = true
		form.c = "# Syntax error."
		goto next
	} else if form.o == nil {
		return ""
	}

	status, msg = cb(text, form)
	if status == 0 {
		err = true
		form.c = "# " + msg
	}

next:
	text = form_compose(forms)
	if err {
		goto top
	}

	return text
}

// Return a string showing the difference in seconds between the time passed as an argument
// and the current time.
func date_diff(old string) string {
	var new_time time.Time
	var old_time time.Time
	var err error
	var diff, howmuch time.Duration
	var diff_secs int

	// I had originally created a single map of names to seconds. The trouble with that
	// was that I always had to sort the map keys when going through them since they're
	// in an undefined order. Using arrays doesn't have this problem, and actually
	// simplifies the code.
	var names = [...]string{"min", "hr", "day", "wk", "mth", "yr"}
	var seconds = [...]int{60, 60 * 60, 60 * 60 * 24, 60 * 60 * 24 * 7, 60 * 60 * 24 * 30, 60 * 60 * 24 * 365}
	var i, secs int
	var what string

	new_time = time.Now()
	// fmt.Printf("New time is %s\n", new_time.Format(out_time_format))
	old_time, err = time.ParseInLocation(out_time_format, old, time.Local)
	if err != nil {
		fmt.Println("Could not parse time:", err)
	}
	// fmt.Printf("Old time is %s\n", old_time.Format(out_time_format))

	diff = new_time.Sub(old_time)
	howmuch = diff
	diff = diff.Round(time.Second)
	diff_secs = int(diff.Seconds())

	for i, secs = range seconds {
		// fmt.Printf("looking at %s %d %d\n", names[i], seconds[i], diff_secs)
		if diff_secs < secs {
			break
		}
		what = names[i]
		howmuch = time.Duration(diff_secs / secs)
	}
	return (fmt.Sprintf("%d %s", howmuch, what))
}

func expand_list(list string) []string {
	var elts []string
	_ = elts

	elts = append(elts, list)
	return elts
}
