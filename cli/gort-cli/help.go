package main

import (
	"fmt"
	"os"
)

func fn_help() {
	var arg string

	if argc < 3 {
		arg = "intro"
	} else {
		arg = os.Args[2]
	}

	switch arg {

	case "usage":
		fallthrough
	case "syntax":
		fmt.Fprintf(os.Stderr, `
    Syntax:

        rt <action> [options] [arguments]
      or
        rt shell

    Each invocation of this program must specify an action (e.g. "edit",
    "create"), options to modify behaviour, and other arguments required
    by the specified action. (For example, most actions expect a list of
    numeric object IDs to act upon.)

    The details of the syntax and arguments for each action are given by
    "rt help <action>". Some actions may be referred to by more than one
    name ("create" is the same as "new", for example).  

    You may also call "rt shell", which will give you an 'rt>' prompt at
    which you can issue commands of the form "<action> [options] 
    [arguments]".  See "rt help shell" for details.

    Objects are identified by a type and an ID (which can be a name or a
    number, depending on the type). For some actions, the object type is
    implied (you can only comment on tickets); for others, the user must
    specify it explicitly. See "rt help objects" for details.

    In syntax descriptions, mandatory arguments that must be replaced by
    appropriate value are enclosed in <>, and optional arguments are
    indicated by [] (for example, <action> and [options] above).

    For more information:

        - rt help objects       (how to specify objects)
        - rt help actions       (a list of actions)
        - rt help types         (a list of object types)
        - rt help shell         (how to use the shell)

`)

	case "conf":
		fallthrough
	case "config":
		fallthrough
	case "configuration":
		fmt.Fprintf(os.Stderr, `
    This program has two major sources of configuration information: its
    configuration files, and the environment.

    The program looks for configuration directives in a file named .rtrc
    (or $RTCONFIG; see below) in the current directory, and then in more
    distant ancestors, until it reaches /. If no suitable configuration
    files are found, it will also check for ~/.rtrc, local/etc/rt.conf
    and /etc/rt.conf.

    Configuration directives:

        The following directives may occur, one per line:

        - server <URL>           URL to RT server.
        - user <username>        RT username.
        - passwd <passwd>        RT user's password.
        - query <RT Query>       Default RT Query for list action
        - orderby <order>        Default RT order for list action
        - queue <queuename>      Default RT Queue for list action
        - auth <rt|basic|gssapi> Method to authenticate via; "basic"
                     means HTTP Basic authentication, "gssapi" means
                     Kerberos credentials, if your RT is configured
                     with $WebRemoteUserAuth.  For backwards
                     compatibility, "externalauth 1" means "auth basic"

        Blank and #-commented lines are ignored.

    Sample configuration file contents:

         server  https://rt.somewhere.com/
         # more than one queue can be given (by adding a query expression)
         queue helpdesk or queue=support
         query Status != resolved and Owner=myaccount


    Environment variables:

        The following environment variables override any corresponding
        values defined in configuration files:

        - RTUSER
        - RTPASSWD
        - RTAUTH
        - RTSERVER
        - RTDEBUG       Numeric debug level. (Set to 3 for full logs.)
        - RTCONFIG      Specifies a name other than ".rtrc" for the
                        configuration file.
        - RTQUERY       Default RT Query for rt list
        - RTORDERBY     Default order for rt list

`)

	case "objects":
		fmt.Fprintf(os.Stderr, `
    Syntax:

        <type>/<id>[/<attributes>]

    Every object in RT has a type (e.g. "ticket", "queue") and a numeric
    ID. Some types of objects can also be identified by name (like users
    and queues). Furthermore, objects may have named attributes (such as
    "ticket/1/history").

    An object specification is like a path in a virtual filesystem, with
    object types as top-level directories, object IDs as subdirectories,
    and named attributes as further subdirectories.

    A comma-separated list of names, numeric IDs, or numeric ranges can
    be used to specify more than one object of the same type. Note that
    the list must be a single argument (i.e., no spaces). For example,
    "user/root,1-3,5,7-10,ams" is a list of ten users; the same list
    can also be written as "user/ams,root,1,2,3,5,7,8-10".
    
    If just a number is given as object specification it will be
    interpreted as ticket/<number>

    Examples:

        1                   # the same as ticket/1
        ticket/1
        ticket/1/attachments
        ticket/1/attachments/3
        ticket/1/attachments/3/content
        ticket/1-3/links
        ticket/1-3,5-7/history

        user/ams

    For more information:

        - rt help <action>      (action-specific details)
        - rt help <type>        (type-specific details)

`)

	case "actions":
		fallthrough
	case "commands":
		fmt.Fprintf(os.Stderr, `
    You can currently perform the following actions on all objects:

        - list          (list objects matching some condition)
        - show          (display object details)
        - edit          (edit object details)
        - create        (create a new object)

    Each type may define actions specific to itself; these are listed in
    the help item about that type.

    For more information:

        - rt help <action>      (action-specific details)
        - rt help types         (a list of possible types)

    The following actions on tickets are also possible:

        - comment       Add comments to a ticket
        - correspond    Add comments to a ticket
        - merge         Merge one ticket into another
        - link          Link one ticket to another
        - take          Take a ticket (steal and untake are possible as well)

    For several edit set subcommands that are frequently used abbreviations
    have been introduced. These abbreviations are:

        - delete or del  delete a ticket           (edit set status=deleted)
        - resolve or res resolve a ticket          (edit set status=resolved)
        - subject        change subject of ticket  (edit set subject=string)
        - give           give a ticket to somebody (edit set owner=user)

`)

	case "types":
		fmt.Fprintf(os.Stderr, `
    You can currently operate on the following types of objects:

        - tickets
        - users
        - groups
        - queues

    For more information:

        - rt help <type>        (type-specific details)
        - rt help objects       (how to specify objects)
        - rt help actions       (a list of possible actions)

`)

	case "ticket":
		fmt.Fprintf(os.Stderr, `
    Tickets are identified by a numeric ID.

    The following generic operations may be performed upon tickets:

        - list
        - show
        - edit
        - create

    In addition, the following ticket-specific actions exist:

        - link
        - merge
        - comment
        - correspond
        - take
        - steal
        - untake
        - give
        - resolve
        - delete
        - subject

    Attributes:

        The following attributes can be used with "rt show" or "rt edit"
        to retrieve or edit other information associated with tickets:

        links                      A ticket's relationships with others.
        history                    All of a ticket's transactions.
        history/type/<type>        Only a particular type of transaction.
        history/id/<id>            Only the transaction of the specified id.
        attachments                A list of attachments.
        attachments/<id>           The metadata for an individual attachment.
        attachments/<id>/content   The content of an individual attachment.

`)

	case "user":
		fallthrough
	case "group":
		fmt.Fprintf(os.Stderr, `
    Users and groups are identified by name or numeric ID.

    The following generic operations may be performed upon them:

        - list
        - show
        - edit
        - create

`)

	case "queue":
		fmt.Fprintf(os.Stderr, `
    Queues are identified by name or numeric ID.

    Currently, they can be subjected to the following actions:

        - show
        - edit
        - create

`)

	case "subject":
		fmt.Fprintf(os.Stderr, `
    Syntax:

        rt subject <id> <new subject text>

    Change the subject of a ticket whose ticket id is given.

`)

	case "give":
		fmt.Fprintf(os.Stderr, `
    Syntax:

        rt give <id> <accountname>

    Give a ticket whose ticket id is given to another user.

`)

	case "resolve":
		fallthrough
	case "res":
		fmt.Fprintf(os.Stderr, `
    Syntax:

        rt resolve <id>

    Resolves a ticket whose ticket id is given.

`)

	case "delete":
		fallthrough
	case "del":
		fmt.Fprintf(os.Stderr, `
    Syntax:

        rt delete <id>

    Deletes a ticket whose ticket id is given.

`)

	case "logout":
		fmt.Fprintf(os.Stderr, `
    Syntax:

        rt logout

    Terminates the currently established login session. You will need to
    provide authentication credentials before you can continue using the
    server. (See "rt help config" for details about authentication.)

`)

	case "ls":
		fallthrough
	case "list":
		fallthrough
	case "search":
		fmt.Fprintf(os.Stderr, `
    Syntax:

        rt <ls|list|search> [options] "query string"

    Displays a list of objects matching the specified conditions.
    ("ls", "list", and "search" are synonyms.)

    The query string must be supplied as one argument.

    if on tickets, query is in the SQL-like syntax used internally by
    RT. (For more information, see "rt help query".), otherwise, query
    is plain string with format "FIELD OP VALUE", e.g. "Name = General".

    if query string is absent, we limit to privileged ones on users and
    user defined ones on groups automatically.

    Options:

        The following options control how much information is displayed
        about each matching object:

        -i             Numeric IDs only. (Useful for |rt edit -; see examples.)
        -s             Short description.
        -l             Longer description.
        -f <field[s]   Display only the fields listed and the ticket id

        In addition,
        
        -o +/-<field>  Orders the returned list by the specified field.
        -r             reversed order (useful if a default was given)
        -q queue[s]    restricts the query to the queue[s] given
                       multiple queues are separated by comma
        -S var=val     Submits the specified variable with the request.
        -t type        Specifies the type of object to look for. (The
                       default is "ticket".)

    Examples:

        rt ls "Priority > 5 and Status=new"
        rt ls -o +Subject "Priority > 5 and Status=new"
        rt ls -o -Created "Priority > 5 and Status=new"
        rt ls -i "Priority > 5"|rt edit - set status=resolved
        rt ls -t ticket "Subject like '[PATCH]%%'"
        rt ls -q systems
        rt ls -f owner,subject
        rt ls -t queue 'Name = General'
        rt ls -t user 'EmailAddress like foo@bar.com'
        rt ls -t group 'Name like foo'

`)

	case "show":
		fmt.Fprintf(os.Stderr, `
    Syntax:

        rt show [options] <object-ids>

    Displays details of the specified objects.

    For some types, object information is further classified into named
    attributes (for example, "1-3/links" is a valid ticket specification
    that refers to the links for tickets 1-3). Consult "rt help <type>"
    and "rt help objects" for further details.

    If only a number is given it will be interpreted as the objects
    ticket/number and ticket/number/history

    This command writes a set of forms representing the requested object
    data to STDOUT.

    Options:

        The following options control how much information is displayed
        about each matching object:

        Without any formatting options prettyprinted output is generated.
        Giving any of the two options below reverts to raw output.
        -s      Short description (history and attachments only).
        -l      Longer description (history and attachments only).

        In addition,
        -               Read IDs from STDIN instead of the command-line.
        -t type         Specifies object type.
        -f a,b,c        Restrict the display to the specified fields.
        -S var=val      Submits the specified variable with the request.

    Examples:

        rt show -t ticket -f id,subject,status 1-3
        rt show ticket/3/attachments/29
        rt show ticket/3/attachments/29/content
        rt show ticket/1-3/links
        rt show ticket/3/history
        rt show -l ticket/3/history
        rt show -t user 2
        rt show 2

`)

	case "new":
		fallthrough
	case "edit":
		fallthrough
	case "create":
		fmt.Fprintf(os.Stderr, `
    Syntax:

        rt edit [options] <object-ids> set field=value [field=value] ...
                                       add field=value [field=value] ...
                                       del field=value [field=value] ...

    Edits information corresponding to the specified objects.

    A purely numeric object id nnn is translated into ticket/nnn

    If, instead of "edit", an action of "new" or "create" is specified,
    then a new object is created. In this case, no numeric object IDs
    may be specified, but the syntax and behaviour remain otherwise
    unchanged.

    This command typically starts an editor to allow you to edit object
    data in a form for submission. If you specified enough information
    on the command-line, however, it will make the submission directly.

    The command line may specify field-values in three different ways.
    "set" sets the named field to the given value, "add" adds a value
    to a multi-valued field, and "del" deletes the corresponding value.
    Each "field=value" specification must be given as a single argument.

    For some types, object information is further classified into named
    attributes (for example, "1-3/links" is a valid ticket specification
    that refers to the links for tickets 1-3). These attributes may also
    be edited. Consult "rt help <type>" and "rt help object" for further
    details.

    Options:

        -       Read numeric IDs from STDIN instead of the command-line.
                (Useful with rt ls ... | rt edit -; see examples below.)
        -i      Read a completed form from STDIN before submitting.
        -o      Dump the completed form to STDOUT instead of submitting.
        -e      Allows you to edit the form even if the command-line has
                enough information to make a submission directly.
        -S var=val
                Submits the specified variable with the request.
        -t type Specifies object type.
        -ct content-type Specifies content type of message(tickets only).

    Examples:

        # Interactive (starts $EDITOR with a form).
        rt edit ticket/3
        rt create -t ticket
        rt create -t ticket -ct text/html

        # Non-interactive.
        rt edit ticket/1-3 add cc=foo@example.com set priority=3 due=tomorrow
        rt ls -t tickets -i 'Priority > 5' | rt edit - set status=resolved
        rt edit ticket/4 set priority=3 owner=bar@example.com \
                         add cc=foo@example.com bcc=quux@example.net
        rt create -t ticket set subject='new ticket' priority=10 \
                            add cc=foo@example.com

`)

	case "comment":
		fallthrough
	case "correspond":
		fmt.Fprintf(os.Stderr, `
    Syntax:

        rt <comment|correspond> [options] <ticket-id>

    Adds a comment (or correspondence) to the specified ticket (the only
    difference being that comments aren't sent to the requestors.)

    This command will typically start an editor and allow you to type a
    comment into a form. If, however, you specified all the necessary
    information on the command line, it submits the comment directly.

    (See "rt help forms" for more information about forms.)

    Options:

        -m <text>       Specify comment text.
        -ct <content-type> Specify content-type of comment text.
        -a <file>       Attach a file to the comment. (May be used more
                        than once to attach multiple files.)
        -c <addrs>      A comma-separated list of Cc addresses.
        -b <addrs>      A comma-separated list of Bcc addresses.
        -s <status>     Set a new status for the ticket (default will
                        leave the status unchanged)
        -w <time>       Specify the time spent working on this ticket.
        -e              Starts an editor before the submission, even if
                        arguments from the command line were sufficient.

    Examples:

        rt comment -m 'Not worth fixing.' -a stddisclaimer.h 23

`)

	case "merge":
		fmt.Fprintf(os.Stderr, `
    Syntax:

        rt merge <from-id> <to-id>

    Merges the first ticket specified into the second ticket specified.

`)

	case "link":
		fmt.Fprintf(os.Stderr, `
    Syntax:

        rt link [-d] <id-A> <link> <id-B>

    Creates (or, with -d, deletes) a link between the specified tickets.
    The link can (irrespective of case) be any of:

        DependsOn/DependedOnBy:     A depends upon B (or vice versa).
        RefersTo/ReferredToBy:      A refers to B (or vice versa).
        MemberOf/HasMember:         A is a member of B (or vice versa).

    To view a ticket's links, use "rt show ticket/3/links". (See
    "rt help ticket" and "rt help show".)

    Options:

        -d      Deletes the specified link.

    Examples:

        rt link 2 dependson 3
        rt link -d 4 referredtoby 6     # 6 no longer refers to 4

`)

	case "query":
		fmt.Fprintf(os.Stderr, `
    RT uses an SQL-like syntax to specify object selection constraints.
    See the <RT:...> documentation for details.
    
    (XXX: I'm going to have to write it, aren't I?)

    Until it exists here a short description of important constructs:

    The two simple forms of query expressions are the constructs
    Attribute like Value and
    Attribute = Value or Attribute != Value

    Whether attributes can be matched using like or using = is built into RT.
    The attributes id, Queue, Owner Priority and Status require the = or !=
    tests.

    If Value is a string it must be quoted and may contain the wildcard
    character %%. If the string does not contain white space, the quoting
    may however be omitted, it will be added automatically when parsing
    the input.

    Simple query expressions can be combined using and, or and parentheses
    can be used to group expressions.

    As a special case a standalone string (which would not form a correct
    query) is transformed into (Owner='string' or Requestor like 'string%%')
    and added to the default query, i.e. the query is narrowed down.

    If no Queue=name clause is contained in the query, a default clause
    Queue=$config{queue} is added.

    Examples:
    Status!='resolved' and Status!='rejected'
    (Owner='myaccount' or Requestor like 'myaccount%%') and Status!='resolved'

`)

	case "form":
		fallthrough
	case "forms":
		fmt.Fprintf(os.Stderr, `
    This program uses RFC822 header-style forms to represent object data
    in a form that's suitable for processing both by humans and scripts.

    A form is a set of (field, value) specifications, with some initial
    commented text and interspersed blank lines allowed for convenience.
    Field names may appear more than once in a form; a comma-separated
    list of multiple field values may also be specified directly.
    
    Field values can be wrapped as in RFC822, with leading whitespace.
    The longest sequence of leading whitespace common to all the lines
    is removed (preserving further indentation). There is no limit on
    the length of a value.

    Multiple forms are separated by a line containing only "--\n".

    (XXX: A more detailed specification will be provided soon. For now,
    the server-side syntax checking will suffice.)

`)

	case "topics":
		fmt.Fprintf(os.Stderr, `
    Syntax:

        rt help <topic>

    Get help on any of the following subjects:

        - tickets, users, groups, queues.
        - show, edit, ls/list/search, new/create.

        - query                                 (search query syntax)
        - forms                                 (form specification)

        - objects                               (how to specify objects)
        - types                                 (a list of object types)
        - actions/commands                      (a list of actions)
        - usage/syntax                          (syntax details)
        - conf/config/configuration             (configuration details)
        - examples                              (a few useful examples)

`)

	case "example":
		fallthrough
	case "examples":
		fmt.Fprintf(os.Stderr, `
    some useful examples

    All the following list requests will be restricted to the default queue.
    That can be changed by adding the option -q queuename

    List all tickets that are not rejected/resolved
        rt ls
    List all tickets that are new and do not have an owner
        rt ls "status=new and owner=nobody"
    List all tickets which I have sent or of which I am the owner
        rt ls myaccount
    List all attributes for the ticket 6977 (ls -l instead of ls)
        rt ls -l 6977
    Show the content of ticket 6977
        rt show 6977
    Show all attributes in the ticket and in the history of the ticket
        rt show -l 6977
    Comment a ticket (mail is sent to all queue watchers, i.e. AdminCc's)
        rt comment 6977
        This will open an editor and lets you add text (attribute Text:)
        Other attributes may be changed as well, but usually don't do that.
    Correspond a ticket (like comment, but mail is also sent to requestors)
        rt correspond 6977
    Edit a ticket (generic change, interactive using the editor)
        rt edit 6977
    Change the owner of a ticket non interactively
        rt edit 6977 set owner=myaccount
        or
        rt give 6977 account
        or
        rt take 6977
    Change the status of a ticket
        rt edit 6977 set status=resolved
        or
        rt resolve 6977
    Change the status of all tickets I own to resolved !!!
        rt ls -i owner=myaccount | rt edit - set status=resolved

`)

	case "shell":
		fmt.Fprintf(os.Stderr, `
    Syntax:

        rt shell

    Opens an interactive shell, at which you can issue commands of 
    the form "<action> [options] [arguments]".

    To exit the shell, type "quit" or "exit".

    Commands can be given at the shell in the same form as they would 
    be given at the command line without the leading 'rt' invocation.

    Example:
        $ rt shell
        rt> create -t ticket set subject='new' add cc=foo@example.com
        # Ticket 8 created.
        rt> quit
        $

`)

	case "take":
		fallthrough
	case "untake":
		fallthrough
	case "steal":
		fmt.Fprintf(os.Stderr, `
    Syntax:

        rt <take|untake|steal> <ticket-id>

    Sets the owner of the specified ticket to the current user, 
    assuming said user has the bits to do so, or releases the 
    ticket.  
    
    'Take' is used on tickets which are not currently owned 
    (Owner: Nobody), 'steal' is used on tickets which *are* 
    currently owned, and 'untake' is used to "release" a ticket 
    (reset its Owner to Nobody).  'Take' cannot be used on
    tickets which are currently owned.

    Example:
        alice$ rt create -t ticket set subject="New ticket"
        # Ticket 7 created.
        alice$ rt take 7
        # Owner changed from Nobody to alice
        alice$ su bob
        bob$ rt steal 7
        # Owner changed from alice to bob
        bob$ rt untake 7
        # Owner changed from bob to Nobody

`)

	case "quit":
		fallthrough
	case "exit":
		fmt.Fprintf(os.Stderr, `
    Use "quit" or "exit" to leave the shell.  Only valid within shell 
    mode.

    Example:
        $ rt shell
        rt> quit
        $

`)
	case "intro":
		fallthrough
	case "introduction":
		fallthrough
	default:
		fmt.Fprintf(os.Stderr, `
    This is a command-line interface to RT 3.0 or newer.

    It allows you to interact with an RT server over HTTP, and offers an
    interface to RT's functionality that is better-suited to automation
    and integration with other tools.

    In general, each invocation of this program should specify an action
    to perform on one or more objects, and any other arguments required
    to complete the desired action.

    For more information:

        - rt help usage         (syntax information)
        - rt help objects       (how to specify objects)
        - rt help actions       (a list of possible actions)
        - rt help types         (a list of object types)

        - rt help config        (configuration details)
        - rt help examples      (a few useful examples)
        - rt help topics        (a list of help topics)

`)
	}
}

func usage() {
	fmt.Fprintf(os.Stderr, `

NAME

rt - command-line interface to RT 3.0 or newer

SYNOPSIS

    rt help

DESCRIPTION

This script allows you to interact with an RT server over HTTP, and offers an
interface to RT's functionality that is better-suited to automation and
integration with other tools.

In general, each invocation of this program should specify an action to
perform on one or more objects, and any other arguments required to complete
the desired action.

`)

}
