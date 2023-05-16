package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"regexp"
	"strconv"
	"strings"
)

// The names of the actions and the handlers to execute for them.
// fn_shell can't be in here because fn_shell itself refers to this.
var actions = map[string]func(){
	"cat":        fn_show,
	"comment":    fn_comment,
	"correspond": fn_comment,
	"create":     fn_edit,
	"del":        fn_setcommand,
	"delete":     fn_setcommand,
	"ed":         fn_edit,
	"edit":       fn_edit,
	"exit":       fn_quit,
	"give":       fn_setcommand,
	"grant":      fn_grant,
	"help":       fn_help,
	"link":       fn_link,
	"list":       fn_list,
	"ln":         fn_link,
	"ls":         fn_list,
	"logout":     fn_logout,
	"man":        fn_help,
	"merge":      fn_merge,
	"new":        fn_edit,
	"quit":       fn_quit,
	"res":        fn_setcommand,
	"resolve":    fn_setcommand,
	"revoke":     fn_grant,
	"search":     fn_list,
	"show":       fn_show,
	"steal":      fn_take,
	"subject":    fn_setcommand,
	"take":       fn_take,
	"untake":     fn_take,
	"ver":        fn_version,
	"version":    fn_version,
}

var (
	add_del_set_re,
	comment_re,
	content_re,
	content_type_re,
	fields_re,
	form_re,
	id_re,
	line_re,
	links_re,
	plus_re,
	queue_re,
	quote_re,
	ticket_re *regexp.Regexp
)

func fn_comment() {
	var cur_arg string
	var content_type, status, msg, wtime, id, text string
	var edit, matched, bad bool
	var files, bcc, cc []string
	var form form_t
	var forms []form_t
	var data = data_t{}
	var fake_res *fake_res_t

	if comment_re == nil {
		comment_re = regexp.MustCompile(`^(?:ticket/)?(` + idlist + `)$`)
	}

	for ; cmd_line_ndx < argc; cmd_line_ndx++ {
		cur_arg = os.Args[cmd_line_ndx]

		if cur_arg == "-e" {
			edit = true
		} else if matched, _ = regexp.MatchString(`^-(?:[abcmws]|ct)$`, cur_arg); matched {
			if cmd_line_ndx >= argc {
				fmt.Fprintf(os.Stderr, "No argument specified with %s\n", cur_arg)
				bad = true
				break
			}

			switch cur_arg {
			case "-a":
				cmd_line_ndx++
				files = append(files, os.Args[cmd_line_ndx])
			case "-ct":
				cmd_line_ndx++
				content_type = os.Args[cmd_line_ndx]
			case "-s":
				cmd_line_ndx++
				status = os.Args[cmd_line_ndx]
			case "-b":
				bcc = strings.Split(os.Args[cmd_line_ndx], ",")
			case "-c":
				cc = strings.Split(os.Args[cmd_line_ndx], ",")
			case "-m":
				cmd_line_ndx++
				msg = os.Args[cmd_line_ndx]
			case "-w":
				cmd_line_ndx++
				wtime = os.Args[cmd_line_ndx]
			}
		} else if id == "" {
			var m [][]string

			m = comment_re.FindAllStringSubmatch(cur_arg, -1)
			if len(m) > 0 {
				id = m[0][1]
			}
		} else {
			var datum string

			if cur_arg[0:1] == "-" {
				datum = "option"
			} else {
				datum = "argument"
			}
			fmt.Fprintf(os.Stderr, "rt edit: Unrecognized %s '%s'\n", datum, cur_arg)
			bad = true
			break
		}
	}

	if id == "" {
		fmt.Fprintf(os.Stderr, "No object specified.\n")
		bad = true
	}

	if bad {
		suggest_help(action, "ticket")
		return
	}

	form.o = append(form.o, "Ticket", "Action", "Cc", "Bcc", "Attachment", "TimeWorked", "Content-Type", "Text")
	form.k = make(key_t)
	form.k["Ticket"] = append(form.k["Ticket"], id)
	form.k["Action"] = append(form.k["Action"], action)
	form.k["Cc"] = append(form.k["Cc"], cc...)
	form.k["Bcc"] = append(form.k["Bcc"], bcc...)
	form.k["Attachment"] = append(form.k["Attachment"], files...)
	if wtime != "" {
		form.k["TimeWorked"] = append(form.k["TimeWorked"], wtime)
	} else {
		form.k["TimeWorked"] = append(form.k["TimeWorked"], "")
	}
	if content_type != "" {
		form.k["Content-Type"] = append(form.k["Content-Type"], content_type)
	} else {
		form.k["Content-Type"] = append(form.k["Content-Type"], "text/plain")
	}
	if msg != "" {
		form.k["Text"] = append(form.k["Text"], msg)
	} else {
		form.k["Text"] = append(form.k["Text"], "")
	}
	form.k["Status"] = append(form.k["Status"], status)

	if status != "" {
		form.o = append(form.o, "Status")
	}

	forms = append(forms, form)
	text = form_compose(forms)

	if edit || (msg == "") {
		var tmp string

		tmp = vi_form_while(text, func(s string, f form_t) (int, string) {
			fmt.Printf("in cb with string %s\n", s)
			return 1, "done"
		})

		if tmp != "" {
			return
		}
		text = tmp
	}

	for _, file := range files {
		data["attachment"] = append(data["attachment"], file)
	}
	data["content"] = append(data["content"], text)
	fake_res = submit(rest+"/ticket/"+id+"/comment", data)
	fmt.Printf("%s", (*fake_res).frs_content)
}

func fn_edit() {
	type act_t map[string]string

	var data = data_t{}
	var objects, new_objects []string
	var edit, input, output, slurped, bad, cl, matched, synerr bool
	var content_type, text, type_name, cur_arg, object string
	var set, add, del act_t
	// var files []string
	var fake_res *fake_res_t

	set = make(act_t)
	add = make(act_t)
	del = make(act_t)

	// Go through the rest of the command line after 'rt [edit|create|new]'.
	for ; cmd_line_ndx < argc; cmd_line_ndx++ {
		cur_arg = os.Args[cmd_line_ndx]

		// If cur_arg is a '#' followed by digits.
		if matched, _ = regexp.MatchString(`^#\d+`, cur_arg); matched {
			// Get rid of leading '#'.
			cur_arg = cur_arg[1:]
		}

		switch cur_arg {
		case "-e":
			edit = true
		case "-i":
			input = true
		case "-o":
			output = true
		case "-ct":
			cmd_line_ndx++
			if cmd_line_ndx < argc {
				content_type = os.Args[cmd_line_ndx]
			}
		case "-t":
			cmd_line_ndx++
			type_name = get_type_argument()
			if type_name == "" {
				bad = true
				break
			}
		case "-S":
			cmd_line_ndx++
			if !get_var_argument(data) {
				bad = true
				break
			}
		case "-":
			if !(slurped || input) {
				var scanner *bufio.Scanner

				// This exits when the first non-type name is entered. You don't
				// have to type EOF like in Perl.
				scanner = bufio.NewScanner(os.Stdin)
				for scanner.Scan() {
					object = scanner.Text()

					if is_object_spec(object, type_name) != "" {
						objects = append(objects, object)
					} else {
						fmt.Fprintf(os.Stderr, "Invalid object on STDIN: %s\n", object)
						bad = true
						break
					}
				}
				slurped = true
			}

		case "add", "del", "set":
			var vars bool
			var m [][]string

			if add_del_set_re == nil {
				add_del_set_re = regexp.MustCompile(`(?s)^(` + field + `)=(.*)$`)
			}

			for cmd_line_ndx++; cmd_line_ndx < argc; cmd_line_ndx++ {
				var tmp_arg string

				tmp_arg = os.Args[cmd_line_ndx]
				if m = add_del_set_re.FindAllStringSubmatch(tmp_arg, -1); len(m) > 0 {
					var key, val string

					key = m[0][1]
					val = m[0][2]
					key = strings.ToLower(key)
					switch cur_arg {
					case "del":
						del[key] = val
					case "add":
						add[key] = val
					case "set":
						set[key] = val
					}
					// fmt.Printf("key = %s val= %s\n", key, val)
					vars = true
				} else {
					cmd_line_ndx--
					break
				}
			}
			if vars == false {
				fmt.Fprintf(os.Stderr, "rt: edit: No variables to %s.\n", cur_arg)
				bad = true
				break
			}
			cl = vars
		default:
			var tmp_arg string

			if matched, _ = regexp.MatchString(`^\d+$`, cur_arg); matched {
				tmp_arg = "ticket/" + cur_arg

				tmp_arg = is_object_spec(tmp_arg, type_name)
				if tmp_arg != "" {
					objects = append(objects, tmp_arg)
				}
			} else if tmp_arg = is_object_spec(cur_arg, type_name); tmp_arg != "" {
				objects = append(objects, tmp_arg)
			} else {
				var datum string

				if cur_arg[0:1] == "-" {
					datum = "option"
				} else {
					datum = "argument"
				}
				fmt.Fprintf(os.Stderr, "rt edit: Unrecognized %s '%s'\n", datum, cur_arg)
				bad = true
				break
			}
		}
	}

	if (action == "ed") || (action == "edit") {
		if len(objects) == 0 {
			fmt.Fprintf(os.Stderr, "rt: edit: No objects specified.\n")
			bad = true
		}
	} else {
		if len(objects) > 0 {
			fmt.Fprintf(os.Stderr, "You shouldn't specify objects as arguments to %s.\n", action)
			bad = true
		}

		if type_name == "" {
			fmt.Fprintf(os.Stderr, "What type of object do you want to create?\n")
			bad = true
		}

		if type_name != "" {
			objects = append(objects, type_name+"/new")
		}
	}

	if bad {
		suggest_help(action, type_name)
		return
	}

	// We need a form to make changes to. We usually ask the server for
	// one, but we can avoid that if we are fed one on STDIN, or if the
	// user doesn't want to edit the form by hand, and the command line
	// specifies only simple variable assignments.  We *should* get a
	// form if we're creating a new ticket, so that the default values
	// get filled in properly.
	for _, object = range objects {
		if strings.HasSuffix(object, "/new") {
			new_objects = append(new_objects, object)
		}
	}

	if input {
		var stdin []byte
		var err error

		stdin, err = io.ReadAll(os.Stdin)
		if err != nil {
			panic(err)
		}

		text = string(stdin)
	} else if edit || (len(add) > 0) || (len(del) > 0) || !cl || (len(new_objects) > 0) {
		data["format"] = append(data["format"], "l")
		for _, str := range objects {
			data["id"] = append(data["id"], str)
		}
		fake_res = submit(rest+`/show`, data)
		text = (*fake_res).frs_content

		// Clear data in preparation for next call to submit().
		data = make(data_t)
	}

	// If any changes were specified on the command line, apply them.
	if cl {
		if text != "" {
			var forms []form_t
			var form form_t
			var where_form int

			//  We're updating forms from the server.
			forms = form_parse(text)

			for where_form, form = range forms {
				var key, val string
				var where_key int

				if (form.e != "") || (form.o == nil) {
					continue
				}

				// Make changes to existing fields. First start by going through each key in
				// the form to see if it will be added, deleted, or set. In all cases,
				// remove the key from the add, del, or set maps so that what's left
				// after we've done this are the key/values that aren't already in the
				// form. Very clever!
				for where_key, key = range form.o {
					var lc_key string
					var ok bool

					// Check to see if the key is in the add, del, and set maps. If so,
					// modify the current form.
					lc_key = strings.ToLower(key)
					if val, ok = add[lc_key]; ok {
						delete(add, lc_key)
						form.k[key] = append(form.k[key], val)
						if strings.ContainsAny(val, ",\n") {
							fmt.Fprintf(os.Stderr, "vsplit not implemented\n")
						}
						continue
					}

					if val, ok = del[lc_key]; ok {
						delete(del, lc_key)

						// Delete the key from the map in the form.
						delete(form.k, key)

						// Delete the key from the list of keys in the form.
						form.o = append(form.o[0:where_key-1], form.o[where_key:]...)
						continue
					}

					if val, ok = set[lc_key]; ok {
						delete(set, lc_key)
						form.k[key] = append(form.k[key], val)
						continue
					}
				}

				// Then update the others.
				// Variables in set are single-valued.
				for key, val = range set {
					form.o = append(form.o, key)
					form.k[key] = nil
					form.k[key] = append(form.k[key], val)
				}

				// Variables in add are multi-valued.
				for key, val = range add {
					form.o = append(form.o, key)
					form.k[key] = append(form.k[key], val)
				}
			}

			forms[where_form] = form

			text = form_compose(forms)
		} else {
			// We're rolling our own set of forms.
			var forms []form_t

			if form_re == nil {
				form_re = regexp.MustCompile(`^(` + name + `)/(` + idlist + `|` + labels + `)(?:(/.*))?$`)
			}

			for _, object = range objects {
				var type_name, ids, args, obj string
				var m [][]string

				m = form_re.FindAllStringSubmatch(object, -1)
				if len(m) > 0 {
					type_name = m[0][1]
					ids = m[0][2]
					args = m[0][3]
				}
				for _, obj = range expand_list(ids) {
					var key, value string
					var form form_t

					set["id"] = type_name + "/" + obj + args
					form.k = make(key_t)
					for key, value = range set {
						form.o = append(form.o, key)
						form.k[key] = append(form.k[key], value)
					}
					forms = append(forms, form)
				}
			}
			text = form_compose(forms)
		}
	}

	if output {
		fmt.Printf("%s\n", text)
		return
	}

	synerr = false

edit:
	if content_type_re == nil {
		content_type_re = regexp.MustCompile(`(?m)^Content-Type:`)
	}

	if (type_name == "ticket") && !content_type_re.MatchString(text) {
		if content_type != "text/plain" {
			text += "Content-Type: " + content_type
		}
	}

	if edit || (!input && (!cl)) {
		var newtext string

		newtext = vi_form_while(text, func(s string, f form_t) (int, string) {
			fmt.Printf("in cb with string %s\n", s)
			return 1, "done"
		})

		if newtext == "" {
			return
		}

		if synerr && (newtext == text) {
			text = ""
		} else {
			text = newtext
		}
	}

	/*
	  delete @data{ grep /^attachment_\d+$/, keys %data };
	    my $i = 1;
	    foreach my $file (@files) {
	        $data{"attachment_$i"} = bless([ $file ], "Attachment");
	        $i++;
	    }
	*/

	if text != "" {
		data["content"] = append(data["content"], text)
		fake_res = submit(rest+`/edit`, data)
		if (*fake_res).frs_code == 409 {
			if edit || (!input && !cl) {
				var content string

				content = (*fake_res).frs_content + "\n"
				_ = content
				synerr = true
				goto edit
			}
		}

		fmt.Printf((*fake_res).frs_content)
	}
}

func fn_grant() {
	fmt.Println("fn_grant")
	fmt.Fprintf(os.Stderr, "rt: grant: %s is unimplemented.\n", os.Args[1])
}

func fn_link() {
	var type_name string
	var bad bool
	var data = data_t{}
	var where int
	var cur_arg string
	var fake_res *fake_res_t
	var from, rel, to string
	var del string = "0"

	var ltypes = map[string]string{
		"dependson":    "DependsOn",
		"hasmember":    "HasMember",
		"dependedonby": "DependedOnBy",
		"refersto":     "RefersTo",
		"referredtoby": "ReferredToBy",
		"memberof":     "MemberOf",
	}

	// Go through the rest of the command line after 'rt link'.
	for ; cmd_line_ndx < argc; cmd_line_ndx++ {
		cur_arg = os.Args[cmd_line_ndx]

		if cur_arg[0:1] == "-" {
			if cur_arg == "-d" {
				del = "1"
			} else if cur_arg == "-t" {
				cmd_line_ndx++
				type_name = get_type_argument()
				if type_name == "" {
					bad = true
					break
				}
			} else {
				fmt.Fprintf(os.Stderr, "Unrecognized option: '%s'\n", cur_arg)
				bad = true
				break
			}
		} else {
			break
		}
	}

	if type_name == "" {
		type_name = "ticket"
	}

	where = argc - cmd_line_ndx

	if where == 3 {

		from = os.Args[cmd_line_ndx]
		cmd_line_ndx++

		rel = os.Args[cmd_line_ndx]
		cmd_line_ndx++
		to = os.Args[cmd_line_ndx]

		if (type_name == "ticket") && (ltypes[strings.ToLower(rel)] == "") {
			fmt.Fprintf(os.Stderr, "Invalid  link '%s' for type %s specified.", rel, type_name)
			bad = true
		}

		data["id"] = append(data["id"], from)
		data["rel"] = append(data["rel"], rel)
		data["to"] = append(data["to"], to)
		data["del"] = append(data["del"], del)
	} else {
		var evil string

		if where > 2 {
			evil = "many"
		} else {
			evil = "few"
		}
		fmt.Fprintf(os.Stderr, "Too %s arguments specified.\n", evil)
		bad = true
	}

	if bad {
		suggest_help("link", type_name)
		return
	}
	fake_res = submit(rest+"/"+type_name+"/link", data)
	fmt.Printf("%s", (*fake_res).frs_content)
}

func fn_list() {
	var matched, bad, rawprint bool
	var type_name string
	var cur_arg string
	var data = data_t{}
	var orderby string
	var queue string
	var reverse_sort bool
	var query string
	var forms []form_t
	var fake_res *fake_res_t

	orderby = config["orderby"]
	if orderby != "" {
		data["orderby"] = append(data["orderby"], orderby)
	}

	// Go through the rest of the command line after 'rt list'.
	for ; cmd_line_ndx < argc; cmd_line_ndx++ {
		cur_arg = os.Args[cmd_line_ndx]

		// Get type specification.
		if cur_arg == "-t" {
			cmd_line_ndx++
			type_name = get_type_argument()
			if type_name == "" {
				bad = true
				break
			}

		} else if cur_arg == "-S" {
			cmd_line_ndx++
			if !get_var_argument(data) {
				bad = true
				break
			}

		} else if cur_arg == "-o" {
			if cmd_line_ndx < (argc - 1) {
				cmd_line_ndx++
				cur_arg = os.Args[cmd_line_ndx]
			}

			// Don't append to what's already there.
			data["orderby"][0] = cur_arg

			// Set the formatting option.
		} else if (cur_arg == "-i") || (cur_arg == "-s") || (cur_arg == "-l") {
			// Pass over the initial "-".
			cur_arg = cur_arg[1:]
			data["format"] = append(data["format"], cur_arg)
			rawprint = true

		} else if cur_arg == "-q" {
			queue = cur_arg

		} else if cur_arg == "-r" {
			reverse_sort = true

		} else if cur_arg == "-f" {
			cmd_line_ndx++
			cur_arg = os.Args[cmd_line_ndx]

			if fields_re == nil {
				fields_re = regexp.MustCompile(fields)
			}
			if matched = fields_re.MatchString(cur_arg); matched {
				data["fields"] = append(data["fields"], cur_arg)
				// option f requires short raw listing format
				data["format"] = append(data["format"], "s")
				rawprint = true
				// fmt.Printf("matched=%v\n", matched)
			} else {
				fmt.Fprintf(os.Stderr, "rt: show: No valid field list in '-f %s'.\n", cur_arg)
				bad = true
				break
			}
		} else if (query == "") && (cur_arg[0:1] != "-") {
			query = cur_arg

		} else {
			var arg string

			if matched, _ = regexp.MatchString(`^-`, cur_arg); matched {
				arg = "option"
			} else {
				arg = "argument"
			}

			fmt.Fprintf(os.Stderr, "rt: list Unrecognized %s '%s'.\n", arg, cur_arg)
			bad = true
			break
		}
	}

	if !rawprint && (len(data["format"]) == 0) {
		data["format"] = append(data["format"], "l")
		data["fields"] = append(data["fields"], "subject,status,queue,created,told,owner,requestors")
	}

	if reverse_sort && strings.HasPrefix(data["format"][0], "-") {
		data["orderby"][0] = "+" + data["orderby"][0][1:]
	} else if reverse_sort {
		var m [][]string

		if plus_re == nil {
			plus_re = regexp.MustCompile(`^\+?(.*)`)
		}
		m = conf_re.FindAllStringSubmatch(data["orderby"][0], -1)
		if len(m) > 0 {
			data["orderby"] = append(data["orderby"], "-"+m[0][1])
		}
	}

	if type_name == "" {
		type_name = "ticket"
	}

	if query == "" {
		if type_name == "ticket" {
			query = config["query"]
		}
	}

	if type_name != "ticket" {
		rawprint = true
	}

	if query == "" {
		var item string

		if type_name != "" {
			item = "query string"
		} else {
			item = "object type"
		}
		fmt.Fprintf(os.Stderr, "No %s specified.\n", item)
		bad = true
	}

	// Get rid of leading hash.
	if query[:1] == "#" {
		query = query[1:]
	}

	if type_name == "ticket" {
		var m [][]string

		if matched, _ = regexp.MatchString(`^\d+$`, query); matched {
			query = "id=" + query
		} else {
			// a string only, take it as an owner or requestor (quoting done later)
			if matched, _ = regexp.MatchString(`^[\w\-]+$`, query); matched {
				query = "(Owner=" + query + "or Requestor like " + query + ") and " + config["query"]
			}

			// always add a query for a specific queue or (comma separated) queues
			if queue != "" {
				strings.Replace(queue, ",", " or Queue=", -1)
			}

			if queue_re == nil {
				queue_re = regexp.MustCompile(`Queue\s*=`)
			}

			if id_re == nil {
				id_re = regexp.MustCompile(`id\s*=`)
			}

			if (queue != "") && (query != "") && queue_re.MatchString(query) && id_re.MatchString(query) {
				query += " and (Queue=" + queue
			}
		}
		// correctly quote strings in a query
		if quote_re == nil {
			quote_re = regexp.MustCompile(`(.+)(=|like\s)\s*([^'\d\s]\S*)\b`)
		}

		m = quote_re.FindAllStringSubmatch(query, -1)
		if len(m) > 0 {
			query = m[0][1] + m[0][2] + `'` + m[0][3] + `'`
		}

		if bad {
			suggest_help(action, type_name)
			return
		}

		if !rawprint {
			fmt.Printf("Query:%s\n", query)
		}
	}

	data["query"] = append(data["query"], query)
	fake_res = submit(rest+`/search/`+type_name, data)

	if rawprint {
		fmt.Printf("%s", (*fake_res).frs_content)
	} else {
		forms = form_parse((*fake_res).frs_content)
		prettylist(forms)
	}
}

func fn_logout() {
	if session.cookie() != "" {
		_ = submit(rest+"/logout", nil)
	}
}

func fn_merge() {
	var bad bool
	var cur_arg string
	var id []string
	var matched bool
	var len_id int
	var fake_res *fake_res_t

	for ; cmd_line_ndx < argc; cmd_line_ndx++ {
		cur_arg = os.Args[cmd_line_ndx]

		// Get rid of leading hash.
		if cur_arg[:1] == "#" {
			cur_arg = cur_arg[1:]
		}
		if matched, _ = regexp.MatchString(`^\d+`, cur_arg); matched {
			id = append(id, cur_arg)
		} else {
			fmt.Fprintf(os.Stderr, "Unrecognized argument: %s.\n", cur_arg)
			bad = true
			break
		}
	}
	len_id = len(id)
	if len_id != 2 {
		var evil string

		if len_id > 2 {
			evil = "many"
		} else {
			evil = "few"
		}
		fmt.Fprintf(os.Stderr, "Too %s arguments specified.\n", evil)
		bad = true
	}

	if bad {
		suggest_help("merge", "ticket")
		return
	}

	fake_res = submit(rest+"/ticket/"+id[0]+"/merge/"+id[1], nil)
	fmt.Printf("%s", (*fake_res).frs_content)
}

func fn_quit() {
	fmt.Println("fn_quit")
	fn_logout()
	os.Exit(1)
}

func fn_setcommand() {
	fmt.Println("fn_setcommand")
}

func fn_shell() {
	var line string
	var err error
	var prompt string = "rt"
	var fn func()

	for reader := bufio.NewReader(os.Stdin); ; {
		var ok bool

		fmt.Printf("%s> ", prompt)
		line, err = reader.ReadString('\n')
		if err != nil {
			break
		}
		line = strings.TrimSuffix(line, "\n")
		os.Args = strings.Split(line, " ")
		argc = len(os.Args)
		os.Args = append(os.Args, os.Args[argc-1])
		for i := argc - 2; i >= 0; i-- {
			os.Args[i+1] = os.Args[i]
		}
		os.Args[0] = prompt
		action = os.Args[1]

		// Look up the action in the actions map to get the function to execute.
		// fmt.Printf("Action = %s\n", action)
		fn, ok = actions[action]

		// If this action is valid, execute it. It's the responsibility
		// of each handler to parse the rest of the command line.
		if ok {
			argc = len(os.Args)
			cmd_line_ndx = 2
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
}

func fn_show() {
	var matched, slurped, bad, rawprint bool
	var type_name string
	var cur_arg string
	var histspec string
	var objects []string
	var data = data_t{}
	var spec string
	var fake_res *fake_res_t

	// Go through the rest of the command line after 'rt show'.
	for ; cmd_line_ndx < argc; cmd_line_ndx++ {
		cur_arg = os.Args[cmd_line_ndx]

		// If cur_arg is a '#' followed by digits.
		if matched, _ = regexp.MatchString(`^#\d+`, cur_arg); matched {
			// Get rid of leading '#'.
			cur_arg = cur_arg[1:]
		}

		// Check for type specification.
		if cur_arg == "-t" {
			cmd_line_ndx++
			type_name = get_type_argument()
			if type_name == "" {
				bad = true
				break
			}

		} else if cur_arg == "-S" {
			cmd_line_ndx++
			if !get_var_argument(data) {
				bad = true
				break
			}

			// Set the formatting option.
		} else if (cur_arg == "-i") || (cur_arg == "-s") || (cur_arg == "-l") {
			// Ignore the initial "-".
			data["format"] = append(data["format"], cur_arg[1:])
			rawprint = true

		} else if (cur_arg == "-") && (slurped == false) {
			var scanner *bufio.Scanner

			scanner = bufio.NewScanner(os.Stdin)
			for scanner.Scan() {
				var object string

				object = scanner.Text()

				if is_object_spec(object, type_name) != "" {
					objects = append(objects, object)
				} else {
					fmt.Fprintf(os.Stderr, "Invalid object on STDIN: %s\n", object)
					bad = true
					break
				}
			}

			// We're done with reading the input. Is there a need to
			// close the scanner?
			fmt.Println(objects)
			slurped = true

		} else if cur_arg == "-f" {
			cmd_line_ndx++
			cur_arg = os.Args[cmd_line_ndx]

			if fields_re == nil {
				fields_re = regexp.MustCompile(fields)
			}
			if matched = fields_re.MatchString(cur_arg); matched {
				data["fields"] = append(data["fields"], cur_arg)
				// option f requires short raw listing format
				data["format"] = append(data["format"], "s")
				rawprint = true
			} else {
				fmt.Fprintf(os.Stderr, "rt: show: No valid field list in '-f %s'.\n", cur_arg)
				bad = true
				break
			}

		} else if matched, _ = regexp.MatchString(`^\d+$`, cur_arg); matched {
			var tmp_arg string = "ticket/" + cur_arg

			tmp_arg = is_object_spec(tmp_arg, type_name)
			if tmp_arg != "" {
				objects = append(objects, tmp_arg)
			}
			histspec = is_object_spec("ticket/"+cur_arg+"/history", type_name)

		} else if matched, _ = regexp.MatchString(`^\d+\/`, cur_arg); matched {
			var tmp_arg string = "ticket/" + cur_arg

			if is_object_spec(tmp_arg, type_name) != "" {
				objects = append(objects, tmp_arg)
			}

			if matched, _ = regexp.MatchString("content$", cur_arg); matched {
				rawprint = true
			}

		} else if spec = is_object_spec(cur_arg, type_name); spec != "" {
			objects = append(objects, spec)

			// I'm using the first re to show whether they all need to be compiled.
			if content_re == nil {
				content_re = regexp.MustCompile("/content$")
				links_re = regexp.MustCompile("/links")
				ticket_re = regexp.MustCompile("^ticket")
			}

			if content_re.MatchString(spec) || links_re.MatchString(spec) || !ticket_re.MatchString(spec) {
				rawprint = true
			}
		} else {
			var arg string

			if matched, _ = regexp.MatchString(`^-`, cur_arg); matched {
				arg = "option"
			} else {
				arg = "argument"
			}

			fmt.Fprintf(os.Stderr, "rt: show Unrecognized %s '%s'.\n", arg, cur_arg)
			bad = true
			break
		}
	}

	if !rawprint {
		var ok bool

		if histspec != "" {
			objects = append(objects, histspec)
		}

		_, ok = data["format"]
		if !ok {
			data["format"] = append(data["format"], "l")
		}
	}

	if len(objects) == 0 {
		fmt.Fprintf(os.Stderr, "rt: %s: No objects specified.\n", action)
		bad = true
	}

	if bad {
		suggest_help(action, type_name)
		return
	}

	for _, str := range objects {
		data["id"] = append(data["id"], str)
	}
	fake_res = submit(rest+`/show`, data)

	if !strings.HasPrefix((*fake_res).frs_content_type, "text") {
		// if this isn't a text reply, remove the trailing newline so we
		// don't corrupt things like tarballs when people do
		// show ticket/id/attachments/id/content > foo.tar.gz
		(*fake_res).frs_content = strings.TrimRight((*fake_res).frs_content, "\n")
		rawprint = true
	}

	if rawprint {
		fmt.Printf("%s", (*fake_res).frs_content)
	} else {
		// The "(?m)" flag is critical so that the pattern looks at each line
		// rather than just the end of the whole string.
		var pattern string = `(?m)^RT\/[\d\.]+ 200 Ok$`
		var forms []form_t

		if line_re == nil {
			line_re = regexp.MustCompile(pattern)
		}

		(*fake_res).frs_content = line_re.ReplaceAllString((*fake_res).frs_content, `--`)
		forms = form_parse((*fake_res).frs_content)
		prettyshow(forms)
	}
}

func fn_take() {
	var cmd string
	var id string
	var bad bool
	var data = data_t{}
	var err error
	var form form_t
	var forms []form_t
	var text string
	var fake_res *fake_res_t
	var word string

	if argc == 3 {
		cmd = os.Args[cmd_line_ndx-1]
		id = os.Args[cmd_line_ndx]
		_, err = strconv.Atoi(id)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Invalid ticket ID %s specified\n", id)
			bad = true
		}
		form.o = append(form.o, "Ticket", "Action")
		form.k = make(key_t)
		form.k["Ticket"] = append(form.k["Ticket"], id)
		form.k["Action"] = append(form.k["Action"], cmd)
		form.k["Status"] = append(form.k["Status"], "")

		forms = append(forms, form)
		text = form_compose(forms)
		data["content"] = append(data["content"], text)
	} else {
		if argc == 2 {
			word = "few"
		} else {
			word = "many"
		}
		fmt.Fprintf(os.Stderr, "Too %s arguments specified.\n", word)
		bad = true
	}

	if bad {
		suggest_help("take", "ticket")
		return
	}

	fake_res = submit(rest+`/ticket/`+id+`/take`, data)
	fmt.Printf("%s", (*fake_res).frs_content)
}

func fn_version() {
	fmt.Printf("rt %s\n", version)
}
