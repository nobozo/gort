package main

import (
	// "fmt"
	"regexp"
	"strings"
)

type key_t map[string][]string

type form_t struct {
	c string   // form comments
	e string   // optional error text
	o []string // order of keys
	k key_t    // key/value pairs
}

var (
	whitespace_re,
	all_whitespace_re,
	kv_re,
	beg_line_re,
	leading_spaces_re *regexp.Regexp
)

// Forms are RFC822-style sets of (field, value) specifications with some
// initial comments and interspersed blank lines allowed for convenience.
// Sets of forms are separated by --\n (in a cheap parody of MIME).
//
// Each form is parsed into a struct with four fields: commented text
// at the start of the form, an array with the order of keys, a hash with
// key/value pairs, and optional error text if the form syntax was wrong.
//
// Returns a slice containing all the parsed form structures.
func form_parse(body string) []form_t {
	var lines []string
	var line string
	var state int = 0 // -1 = syntax error, 0 = start of form, 1 = processing form.
	var line_cnt, cur_line int
	var kv_pattern string = `^(` + field + `):(?:\s+(.*))?$`
	var form form_t
	var forms []form_t

	// Create a slice with one line of the form body per element.
	// This might not be very efficient.
	lines = strings.Split(body, "\n")

	// How many lines are in the body?
	line_cnt = len(lines)

	// A pattern for breaking apart a key:value line.
	if kv_re == nil {
		kv_re = regexp.MustCompile(kv_pattern)
	}

	// Create the map that will hold the key:value pairs in the form.
	if form.k == nil {
		form.k = make(key_t)
	}

	// This is the slice that will hold each form struct.
	forms = make([]form_t, 0)

	if all_whitespace_re == nil {
		all_whitespace_re = regexp.MustCompile(`^\s+$`)
	}

	// Go through each line in the body.
	for cur_line = 0; cur_line < line_cnt; cur_line++ {

		line = lines[cur_line]
		// fmt.Printf("forms = %d\nforms = %v\n", forms)

		// We can safely ignore blank lines or lines with only whitespace.
		if len(line) == 0 || all_whitespace_re.MatchString(line) {
			continue
		}

		// If at the beginning of a form.
		if line == "--" {
			// We're either at the beginning of the body or we've reached the end of one form.
			// We'll ignore it if it was empty. Otherwise store it, errors and all.
			if (form.e != "") || (form.c != "") || (len(form.o) != 0) {
				// The form isn't empty so append it to forms.
				forms = append(forms, form)

				// Empty out the form for next form.
				form.e = ""
				form.c = ""
				form.o = nil
				form.k = make(key_t)
			}
			state = 0

			// If there hasn't been a syntax error.
		} else if state != -1 {
			// Read an optional block of comments only at the start of a form.
			if (state == 0) && (line[0:1] == "#") {
				state = 1

				// Go through the body looking for all comment lines.
				// Append them all into form.c.
				for cur_line < line_cnt {
					form.c += line
					cur_line++
					line = lines[cur_line]
					if (len(line) != 0) && (line[0:1] == "#") {
						continue
					} else {
						// No more comments.
						cur_line--
						break
					}
				}
				// We found a key:value line.
			} else if matches := kv_re.FindAllStringSubmatch(line, -1); len(matches) > 0 {
				var key string
				var value string
				var cont_lines []string
				var whitespace string

				if whitespace_re == nil {
					whitespace_re = regexp.MustCompile(`^\s+`)
				}

				key = matches[0][1]
				value = matches[0][2]

				// At this point we know what the key is.
				form.o = append(form.o, key)

				// The value could be in one of three places:
				// 1) after the ":" on the current line.
				// 2) after the ":" on the current line and on
				//     following continuation lines.
				// 3) only on following continuation lines.
				//
				// (A continuation line is a line starting with whitespace and
				// followed by any text.) We already know the current
				// line isn't all whitespace due to a check above.
				// So now look to see which it is.
				//
				// I'm assuming there can be any number of continuation lines.

				if value != "" {
					form.k[key] = append(form.k[key], value)
				}

				// If the next line is empty, or consists of all whitespace, or a regular
				// key:value line, then there are no more continuation lines.
				for cur_line++; cur_line < (line_cnt - 1); cur_line++ {
					line = lines[cur_line]

					if len(line) == 0 || all_whitespace_re.MatchString(line) {
						break
					}

					// If the line has leading whitespace.
					if whitespace_re.MatchString(line) {
						// It's a continuation line.
						cont_lines = append(cont_lines, line)
					} else {
						// It's not a continuation line. Put it back.
						cur_line--
						break
					}
				}

				// If no continuation lines were found then the current line
				// is just a key with an empty value. There's nothing
				// more to do. If there's one continuation line then it will become
				// the value of the current key. If there is more than one, then
				// look through them for a common length of leading whitespace.
				if len(cont_lines) == 1 {
					form.k[key] = append(form.k[key], cont_lines[0])
					continue
				} else if len(cont_lines) > 1 {
					// Find longest common leading whitespace indent.
					var cont_line_whitespace string
					var cont_line string

					// Go through all the continuation lines looking for the shortest
					// common leading whitespace.
					for _, cont_line = range cont_lines {
						cont_line_whitespace = whitespace_re.FindString(cont_line)
						if len(cont_line_whitespace) > len(whitespace) {
							whitespace = cont_line_whitespace
						}
					}

					// If a common length of whitespace was found, strip it from each line.
					if len(whitespace) > 0 {
						var i int

						for i, cont_line = range cont_lines {

							cont_lines[i] = strings.Replace(cont_lines[i], whitespace, "", -1)
							value += cont_lines[i] + "\n"
						}
						value = strings.TrimSuffix(value, "\n")
						form.k[key] = append(form.k[key], value)
					}
				}
				state = 1
			} else if line[0:1] != "#" {
				// We've found a syntax error, so we'll reconstruct the
				// form parsed thus far, and add an error marker. (>>)
				state = -1
				form.e = form_compose(forms)
				if strings.HasPrefix(line, ">>") {
					form.e += line
				} else {
					form.e += ">> " + line
				}
			}
		} else {
			// We saw a syntax error earlier, so we'll accumulate the
			// contents of this form until the end.
			form.e += line
		}
	}

	if (form.e != "") || (form.c != "") || (form.o != nil) {
		forms = append(forms, form)
	}

	// fmt.Printf("forms len = %d\n", len(forms))
	return forms
}

// Given a slice of form_t structs, return their text representation as a string.
func form_compose(forms []form_t) string {
	var form form_t
	var text_sl []string

	if beg_line_re == nil {
		beg_line_re = regexp.MustCompile(`(?m)^`)
	}

	// Go through each form_t struct passed as an argument.
	for _, form = range forms {
		var text string

		// A comment can have any number of newlines at the end. So, normalize them
		// by getting rid of all of them, and then adding just one.
		if form.c != "" {
			form.c = strings.TrimSuffix(form.c, "\n")
			text = form.c + "\n"
		}

		// Collect any error strings.
		if form.e != "" {
			text += form.e
		} else if form.o != nil {
			// Go through all the keys in form.o. Since we want our output to be in
			// the correct order, we'll get the keys from form.o, which is in the
			// order they were originally found. We can't just go through form.k because
			// there's no guaranty what order its elements are in.
			var lines []string

			for _, key := range form.o {
				var spaces, line, value string
				var values []string

				// Get the corresponding value.
				if len(form.k[key]) > 0 {
					values = form.k[key]
				} else {
					values = nil
				}

				spaces = strings.Repeat(" ", len(key+":  "))
				if len(spaces) > 16 {
					spaces = strings.Repeat(" ", 4)
				}

				for _, value = range values {
					if strings.Contains(value, "\n") {
						value = beg_line_re.ReplaceAllString(value, spaces)
						if leading_spaces_re == nil {
							leading_spaces_re = regexp.MustCompile("^" + spaces)
						}
						value = leading_spaces_re.ReplaceAllString(value, "")

						if line != "" {
							lines = append(lines, line+"\n\n")
							line = ""
						} else if (len(lines) > 0) && (lines[len(lines)-1] != "\n\n") {
							lines[len(lines)-1] += "\n"
						}

						lines = append(lines, key+": "+value+"\n\n")
					} else if (line != "") && (len(line)+len(value) >= 70) {
						line += ",\n" + spaces + value
					} else if line != "" {
						line = line + "," + value
					} else {
						line = key + ": " + value
					}
				}

				if len(values) == 0 {
					line = key + ":"
				}

				if line != "" {
					if strings.Contains(line, "\n") {
						var lines_len int

						lines_len = len(lines)
						if (lines_len != 0) && !strings.HasSuffix(lines[lines_len-1], "\n\n") {
							lines[lines_len-1] += "\n"
						} else {
							line += "\n"
						}
					}
					lines = append(lines, line+"\n")
				}
			}

			text += strings.Join(lines, "")
		} else {
			text = strings.TrimSuffix(text, "\n")
		}
		text_sl = append(text_sl, text)
	}

	return strings.Join(text_sl, "\n--\n\n")
}
