package etc

/*
#
# RT was configured with:
#
#   $ ./configure --prefix=/usr/local
#

package RT;

#############################  WARNING  #############################
#                                                                   #
#                     NEVER EDIT RT_Config.pm !                     #
#                                                                   #
#         Instead, copy any sections you want to change to          #
#         RT_SiteConfig.pm and edit them there.  Otherwise,         #
#         your changes will be lost when you upgrade RT.            #
#                                                                   #
#############################  WARNING  #############################

=head1 NAME

RT::Config

=head1 DESCRIPTION

RT has dozens of configuration options to customize how RT behaves for
different situations. The available options and valid values are
described below.

=head2 Server Configuration Files

This file, etc/RT_Config.pm, defines the available configuration options
and sets the defaults for RT. You should B<never> edit this file directly.
If you do, your changes will be lost when you upgrade and RT installs the
newest version of this file on your system.

The correct place to set site-specific options is in etc/RT_SiteConfig.pm.
If you have many customizations to manage, you can break your configuration
into multiple files and put them in a directory etc/RT_SiteConfig.d/. More
information about this option is available in the L<RT::Config>
documentation.

=head2 Web Configuration

RT also allows you to set configuration via the RT web interface at
Admin > Tools > System Configuration. Any configuration options set there take
precedence over values set in etc/RT_SiteConfig.pm. If you provide a
custom setting in both places, RT will issue a warning in the log as a
reminder to consider removing the setting from etc/RT_SiteConfig.pm
to avoid confusion.

Some settings that are core to RT cannot be changed via the web interface.
This prevents changes that could make your RT inoperable, leaving you
unable to restore the system via the web UI.

=head1 System

=head2 Base configuration

=over 4


$rtname is the string that RT will look for in mail messages to
figure out what ticket a new piece of mail belongs to.

Your domain name is recommended, so as not to pollute the namespace.
Once you start using a given tag, you should probably never change it;
otherwise, mail for existing tickets won't get put in the right place.

*/

var Rtname = "example.com"

/*

You should set this to your organization's DNS domain. For example,
I<fsck.com> or I<asylum.arkham.ma.us>. It is used by the linking
interface to guarantee that ticket URIs are unique and easy to
construct.  Changing it after you have created tickets in the system
will B<break> all existing ticket links!

*/

var Organization = "example.com"

/*

RT is designed such that any mail which already has a ticket-id
associated with it will get to the right place automatically.

$CorrespondAddress and $CommentAddress are the default addresses
that will be listed in From: and Reply-To: headers of correspondence
and comment mail tracked by RT, unless overridden by a queue-specific
address.  They should be set to email addresses which have been
configured as aliases for F<rt-mailgate>.

*/

var CorrespondAddress = ""

var CommentAddress = ""

/*

Domain name of the RT server, e.g. 'www.example.com'. It should not
contain anything except the server name.

*/

var WebDomain = "localhost"

/*

If we're running as a superuser, run on port 80.  Otherwise, pick a
high port for this user.

443 is default port for https protocol.

*/

var WebPort = 8080

/*

If you're putting the web UI somewhere other than at the root of your
server, you should set $WebPath to the path you'll be serving RT
at.

$WebPath requires a leading / but no trailing /, or it can be
blank.

In most cases, you should leave $WebPath set to "" (an empty
value).

*/

var WebPath = ""

/*

$Timezone is the default timezone, used to convert times entered by
users into GMT, as they are stored in the database, and back again;
users can override this.  It should be set to a timezone recognized by
your server.

*/

var Timezone = "US/Eastern"

/*

Once a plugin has been downloaded and installed, use Plugin() to add
to the enabled @Plugins list:

    Plugin( "RT::Extension::JSGantt" );

RT will also accept the distribution name (i.e. RT-Extension-JSGantt)
instead of the package name (RT::Extension::JSGantt).

*/

var Plugins = []string{}

/*

Set @StaticRoots to serve extra paths with a static handler.  The
contents of each hashref should be the the same arguments as
L<Plack::Middleware::Static> takes.  These paths will be checked before
any plugin or core static paths.

Example:

*/

/*
var @StaticRoots,
        {
            path => qr{^/static/},
            root => '/local/path/to/static/parent',
        },
    );

*/

var StaticRoots = []string{}

/*

=back




=head2 Database connection

=over 4


Database driver being used; case matters.  Valid types are "mysql",
"Oracle", and "Pg".  "SQLite" is also available for non-production use.

*/

var DatabaseType = "mysql"

/*

The domain name of your database server.  If you're running MySQL and
on localhost, leave it blank for enhanced performance.

DatabaseRTHost is the fully-qualified hostname of your RT server,
for use in granting ACL rights on MySQL.

*/

var DatabaseHost = "localhost"

var DatabaseRTHost = "localhost"

/*

The port that your database server is running on.  Ignored unless it's
a positive integer. It's usually safe to leave this blank; RT will
choose the correct default.

*/

var DatabasePort = ""

/*

The name of the user to connect to the database as.

*/

var DatabaseUser = "rt_user"

/*

The password the $DatabaseUser should use to access the database.

*/

var DatabasePassword = "rt_pass"

/*

The name of the RT database on your database server. For Oracle, the
SID and database objects are created in $DatabaseUser's schema.

*/

// var DatabaseName = "rt5"
var DatabaseName = "jon" // So that we don't override the real db while testing.

/*

Allows additional properties to be passed to the database connection
step.  Possible properties are specific to the database-type; see
https://metacpan.org/pod/DBI#connect

For PostgreSQL, for instance, the following enables SSL (but does no
certificate checking, providing data hiding but no MITM protection):

   # See https://metacpan.org/pod/DBD::Pg#connect
   # and http://www.postgresql.org/docs/8.4/static/libpq-ssl.html
*/

var DatabaseExtraDSN = map[string]string{"sslmode": "require"}

/*

For MySQL, the following acts similarly if the server has enabled SSL.
Otherwise, it provides no protection; MySQL provides no way to I<force>
SSL connections:

   # See https://metacpan.org/pod/DBD::mysql#connect
   # and http://dev.mysql.com/doc/refman/5.1/en/ssl-options.html
*/

// var DatabaseExtraDSN = map[string]bool{"mysql_ssl": true}

// var %DatabaseExtraDSN = ()

/*

The name of the database administrator to connect to the database as
during upgrades.

*/

var DatabaseAdmin = "root"

/*

=back




=head2 Logging

The default is to log anything except debugging information to syslog.
Check the L<Log::Dispatch> POD for information about how to get things
by syslog, mail or anything else, get debugging info in the log, etc.

It might generally make sense to send error and higher by email to
some administrator.  If you do this, be careful that this email isn't
sent to this RT instance.  Mail loops will generate a critical log
message.

=over 4


The minimum level error that will be logged to the specific device.
From lowest to highest priority, the levels are:

    debug info notice warning error critical alert emergency

Many syslogds are configured to discard or file debug messages away, so
if you're attempting to debug RT you may need to reconfigure your
syslogd or use one of the other logging options.

Logging to your screen affects scripts run from the command line as well
as the STDERR sent to your webserver (so these logs will usually show up
in your web server's error logs).

*/

var LogToSyslog = "info"

var LogToSTDERR = "info"

/*

Logging to a standalone file is also possible. The file needs to both
exist and be writable by all direct users of the RT API. This generally
includes the web server and whoever rt-crontool runs as. Note that
rt-mailgate and the RT CLI go through the webserver, so their users do
not need to have write permissions to this file. If you expect to have
multiple users of the direct API, Best Practical recommends using syslog
instead of direct file logging.

You should set $LogToFile to one of the levels documented above.

*/

var LogToFile = ""

var LogDir = "var/log"

var LogToFileNamed = "rt.log"

/*

If set to a log level then logging will include stack traces for
messages with level equal to or greater than specified.

NOTICE: Stack traces include parameters supplied to functions or
methods. It is possible for stack trace logging to reveal sensitive
information such as passwords or ticket content in your logs.

*/

var LogStackTraces = ""

/*

Additional options to pass to L<Log::Dispatch::Syslog>; the most
interesting flags include facility, logopt, and possibly ident.
See the L<Log::Dispatch::Syslog> documentation for more information.

*/

var LogToSyslogConf = []string{}

/*

=back



=head2 Incoming mail gateway

=over 4


This regexp controls what subject tags RT recognizes as its own.  If
you're not dealing with historical $rtname values, or historical
queue-specific subject tags, you'll likely never have to change this
configuration.

Be B<very careful> with it. Note that it overrides $rtname for
subject token matching.

The setting below would make RT behave exactly as it does without the
setting enabled.

*/

var EmailSubjectTagRegex = `/\Q` + Rtname + `\E/i `

/*

$OwnerEmail is the address of a human who manages RT. RT will send
errors generated by the mail gateway to this address. Because RT
sends errors to this address, it should I<not> be an address that's
managed by your RT instance, to avoid mail loops.

*/

var OwnerEmail = "root"

/*

If $LoopsToRTOwner is defined, RT will send mail that it believes
might be a loop to $OwnerEmail.

*/

var LoopsToRTOwner = 1

/*

If $StoreLoops is defined, RT will record messages that it believes
to be part of mail loops.  As it does this, it will try to be careful
not to send mail to the sender of these messages.

*/

var StoreLoops = -1

/*

$MaxAttachmentSize sets the maximum size (in bytes) of attachments
stored in the database.  This setting is irrelevant unless one of
$TruncateLongAttachments or $DropLongAttachments (below) are set, B<OR>
the database is stored in Oracle.  On Oracle, attachments larger than
this can be fully stored, but will be truncated to this length when
read.

*/

var MaxAttachmentSize = 10 * 1024 * 1024

/*

If this is set to a non-undef value, RT will truncate attachments
longer than $MaxAttachmentSize.

*/

var TruncateLongAttachments = 0

/*

If this is set to a non-undef value, RT will silently drop attachments
longer than MaxAttachmentSize.  $TruncateLongAttachments, above,
takes priority over this.

*/

var DropLongAttachments = 0

/*

$RTAddressRegexp is used to make sure RT doesn't add itself as a
ticket CC if $ParseNewMessageForTicketCcs, above, is enabled.  It
is important that you set this to a regular expression that matches
all addresses used by your RT.  This lets RT avoid sending mail to
itself.  It will also hide RT addresses from the list of "One-time Cc"
and Bcc lists on ticket reply.

If you have a number of addresses configured in your RT database
already, you can generate a naive first pass regexp by using:

    perl etc/upgrade/generate-rtaddressregexp

If left blank, RT will compare each address to your configured
$CorrespondAddress and $CommentAddress before searching for a
Queue configured with a matching "Reply Address" or "Comment Address"
on the Queue Admin page.

*/

var RTAddressRegexp = ""

/*

RT provides functionality which allows the system to rewrite incoming
email addresses, using L<RT::User/CanonicalizeEmailAddress>.  The
default implementation replaces all occurrences of the regular
expression in CanonicalizeEmailAddressMatch with
CanonicalizeEmailAddressReplace, via s/$Match/$Replace/gi.  The
most common use of this is to replace @something.example.com with
@example.com.  If more complex noramlization is required,
L<RT::User/CanonicalizeEmailAddress> can be overridden to provide it.

*/

var CanonicalizeEmailAddressMatch = "@subdomain.example.com$"

var CanonicalizeEmailAddressReplace = "@example.com"

/*

By default $ValidateUserEmailAddresses is 1, and RT will refuse to create
users with an invalid email address (as specified in RFC 2822) or with
an email address made of multiple email addresses.

Set this to false to skip any email address validation.  Doing so may open up
vulnerabilities.

*/

var ValidateUserEmailAddresses = true

/*

@MailPlugins is a list of authentication plugins for
L<RT::Interface::Email> to use; see L<rt-mailgate>

The default "extract remote tracking tags" scrip settings; these
detect when your RT is talking to another RT, and adjust the subject
accordingly.

*/

var ExtractSubjectTagMatch = `/\[[^\]]+? #\d+\]/`

// var $ExtractSubjectTagNoMatch = ( ${RT::EmailSubjectTagRegex}
//        ? qr{\[(?:https?://)?(?:${RT::EmailSubjectTagRegex}) #\d+\]}
//        : qr{\[(?:https?://)?\Q$RT::rtname\E #\d+\]}));

/*

Some email clients create a plain text version of HTML-formatted
email to help other clients that read only plain text.
Unfortunately, the plain text parts sometimes end up with
doubled newlines and these can then end up in RT. This
is most often seen in MS Outlook.

Enable this option to have RT check for additional mail headers
and attempt to identify email from MS Outlook. When detected,
RT will then clean up double newlines. Note that it may
clean up intentional double newlines as well.

*/

var CheckMoreMSMailHeaders = false

/*

Email coming into RT can sometimes have another email attached, when
an email is forwarded as an attachment, for example. By default, RT recognizes
that the attached content is an email and does some processing, including
some parsing the headers of the attached email. You can see this in suggested
email addresses on the People page and One-time Cc on reply.

If you want RT to treat attached email files as regular file attachments,
set this option to true (1). With this option enabled, attached email will
show up in the "Attachments" section like other types of file attachments
and content like headers will not be processed.

*/

var TreatAttachedEmailAsFiles = false

/*

=back



=head2 Outgoing mail

=over 4


$MailCommand defines which method RT will use to try to send mail.
We know that 'sendmailpipe' works fairly well.  If 'sendmailpipe'
doesn't work well for you, try 'sendmail'.  'qmail' is also a supported
value.

For testing purposes, or to simply disable sending mail out into the
world, you can set $MailCommand to 'mbox' which logs all mail, in
mbox format, to files in F</opt/rt5/var/> based in the process start
time.  The 'testfile' option is similar, but the files that it creates
(under /tmp) are temporary, and removed upon process completion; the
format is also not mbox-compatable.

*/

var MailCommand = "sendmailpipe"

/*

SetOutgoingMailFrom> tells RT to set the sender envelope to the
Correspond mail address of the ticket's queue.

Warning: If you use this setting, bounced mails will appear to be
incoming mail to the system, thus creating new tickets.

If the value contains an @, it is assumed to be an email address and used as
a global envelope sender.  Expected usage in this case is to simply set the
same envelope sender on all mail from RT, without defining
$OverrideOutgoingMailFrom.  If you do define $OverrideOutgoingMailFrom,
anything specified there overrides the global value (including Default).

This option only works if $MailCommand is set to 'sendmailpipe'.

*/

var SetOutgoingMailFrom = false

/*

$OverrideOutgoingMailFrom is used for overwriting the Correspond
address of the queue as it is handed to sendmail -f. This helps force
the From_ header away from www-data or other email addresses that show
up in the "Sent by" line in Outlook.

The option is a hash reference of queue id/name to email address. If
there is no ticket involved, then the value of the Default key will
be used.

SetOutgoingMailFrom> is enabled and
$MailCommand is set to 'sendmailpipe'.

*/

var OverrideOutgoingMailFrom = map[string]string{"Default": "admin@rt.example.com", "General": "general@rt.example.com"}

/*

$DefaultMailPrecedence is used to control the default Precedence
level of outgoing mail where none is specified.  By default it is
bulk, but if you only send mail to your staff, you may wish to
change it.

Note that you can set the precedence of individual templates by
including an explicit Precedence header.

If you set this value to undef then we do not set a default
Precedence header to outgoing mail. However, if there already is a
Precedence header, it will be preserved.

*/

var DefaultMailPrecedence = "bulk"

/*

$OverrideMailPrecedence is used for overwriting the $DefaultMailPrecedence
value for a queue.

The option is a hash reference of queue id/name to precedence. If you set the
precedence to undef, a Precedence header will not be added to the mail.

This option only works if $DefaultMailPrecedence is enabled.

*/

// var $OverrideMailPrecedence = { "Queue 1" => "list", "Queue 2" => undef, }

/*

$DefaultErrorMailPrecedence is used to control the default
Precedence level of outgoing mail that indicates some kind of error
condition. By default it is bulk, but if you only send mail to your
staff, you may wish to change it.

If you set this value to undef then we do not add a Precedence
header to error mail.

*/

var DefaultErrorMailPrecedence = "bulk"

/*

$UseOriginatorHeader is used to control the insertion of an
RT-Originator Header in every outgoing mail, containing the mail
address of the transaction creator.

*/

var UseOriginatorHeader = true

/*

By default, RT sets the outgoing mail's "From:" header to "SenderName
Setting $UseFriendlyFromLine to false disables it.

*/

var UseFriendlyFromLine = true

/*

sprintf() format of the friendly 'From:' header; its arguments are
SenderName and SenderEmailAddress.

*/

var FriendlyFromLineFormat = `%s via RT" <%s>`

/*

RT can optionally set a "Friendly" 'To:' header when sending messages
to Ccs or AdminCcs (rather than having a blank 'To:' header.

This feature DOES NOT WORK WITH SENDMAIL[tm] BRAND SENDMAIL.  If you
are using sendmail, rather than postfix, qmail, exim or some other
MTA, you _must_ disable this option.

*/

var UseFriendlyToLine = false

/*

sprintf() format of the friendly 'To:' header; its arguments are
WatcherType and TicketId.

*/

// var FriendlyToLineFormat = "\"%s of " + RT->Config->Get('rtname') +" Ticket #%s\:"

/*

By default, RT doesn't notify the person who performs an update, as
they already know what they've done. If you'd like to change this
Set $NotifyActor to true.

*/

var NotifyActor = true

/*

By default, RT records each message it sends out to its own internal
database.  To change this behavior, set $RecordOutgoingEmail to false

If this is disabled, users' digest mail delivery preferences
(i.e. EmailFrequency) will also be ignored.

*/

var RecordOutgoingEmail = true

/*

Setting these options enables VERP support
L<http://cr.yp.to/proto/verp.txt>.

Uncomment the following two directives to generate envelope senders
of the form ${VERPPrefix}${originaladdress}@${VERPDomain}
(i.e. rt-jesse=fsck.com@rt.example.com ).

This currently only works with sendmail and sendmailpipe.

*/

var VERPPrefix = "rt-"

// var $VERPDomain = $RT::Organization

/*

By default, RT forwards a message using queue's address and adds RT's
tag into subject of the outgoing message, so recipients' replies go
into RT as correspondents.

To change this behavior, set $ForwardFromUser to true and RT
will use the address of the current user and remove RT's subject tag.

*/

var ForwardFromUser = false

/*

RT's default pure-perl formatter may fail to successfully convert even
on some relatively simple HTML; this will result in blank text/plain
parts, which is particuarly unfortunate if HTML templates are not in
use.

If the optional dependency L<HTML::FormatExternal> is installed, RT will
use external programs to render HTML to plain text.  The default is to
try, in order, w3m, elinks, html2text, links, lynx, and
then fall back to the core pure-perl formatter if none are installed.

Set $HTMLFormatter to one of the above programs (or the full path to
such) to use a different program than the above would choose by default.
Setting this requires that L<HTML::FormatExternal> be installed.

If the chosen formatter is not in the webserver's $PATH, you may set
this option the full path to one of the aforementioned executables.

*/

var HTMLFormatter = ""

/*

=back

=head3 Email dashboards

=over 4


The email address from which RT will send dashboards. If none is set,
then $OwnerEmail will be used.

*/

var DashboardAddress = ""

/*

Lets you set the subject of dashboards. Arguments are the frequency (Daily,
Weekly, Monthly) of the dashboard and the dashboard's name.

*/

var DashboardSubject = "%s Dashboard: %s"

/*

A list of regular expressions that will be used to remove content from
mailed dashboards.

*/

// var @EmailDashboardRemove = ()

/*

A list that specifies which language to use for dashboard subscription email.
There are several special keys:

* _subscription: the language chosen on the dashboard subscription page
* _recipient: the recipient's language, as chosen on their "About Me" page
* _subscriber: the subscriber's language, as chosen on their "About Me" page

The first key that produces a value is used for the email. Be aware that users
may not actually have a language set on their "About Me" page, since RT falls
back to the language their web browser specifies (and of course in a scheduled
email dashboard, there is no web browser).

You may also include a specific language as a fallback when there is no
language specified otherwise. Using a specific language never fails to produce
a value, so subsequent values in the list will never be considered.

By default, RT examines the subscription, then the recipient, then subscriber,
then finally falls back to English.

See also L</@LexiconLanguages>.


*/

// var @EmailDashboardLanguageOrder = qw(_subscription _recipient _subscriber en)

/*

=back



=head3 Sendmail configuration

These options only take effect if $MailCommand is 'sendmail' or
'sendmailpipe'

=over 4


$SendmailArguments> defines what flags to pass to C<$SendmailPath
These options are good for most sendmail wrappers and work-a-likes.

These arguments are good for sendmail brand sendmail 8 and newer:

var SendmailArguments,"-oi -ODeliveryMode=b -OErrorMode=m"

*/

var SendmailArguments = "-oi"

/*

$SendmailBounceArguments defines what flags to pass to $Sendmail
assuming RT needs to send an error (i.e. bounce).

*/

var SendmailBounceArguments = `-f "<>"`

/*

If you selected 'sendmailpipe' above, you MUST specify the path to
your sendmail binary in $SendmailPath.

*/

var SendmailPath = "/usr/sbin/sendmail"

/*

=back

=head3 Other mailers

=over 4

@MailParams defines a list of options passed to $MailCommand if it
is not "sendmailpipe" or "sendmail";

*/

// var @MailParams = ()

/*

=back

=head2 Application logic

=over 4

If $ParseNewMessageForTicketCcs is set to true, RT will attempt to
divine Ticket 'Cc' watchers from the To and Cc lines of incoming
messages that create new tickets. This option does not apply to replies
or comments on existing tickets. Be forewarned that if you have I<any>
addresses which forward mail to RT automatically and you enable this
option without modifying $RTAddressRegexp below, you will get
yourself into a heap of trouble.

See also the L<RT::Action::AutoAddWatchers> extension which adds
watchers from ticket replies on existing tickets.

*/

var ParseNewMessageForTicketCcs = true

/*

Set $UseTransactionBatch to true to execute transactions in batches,
such that a resolve and comment (for example) would happen
simultaneously, instead of as two transactions, unaware of each
others' existence.

*/

var UseTransactionBatch = true

/*

When this feature is enabled a user needs I<ModifyTicket> rights on
both tickets to link them together; otherwise, I<ModifyTicket> rights
on either of them is sufficient.

*/

var StrictLinkACL = true

/*

Should RT redistribute correspondence that it identifies as machine
generated?  A 1 will do so; setting this to false will cause no
such messages to be redistributed.  You can also use "privileged" (the
default), which will redistribute only to privileged users. This helps
to protect against malformed bounces and loops caused by auto-created
requestors with bogus addresses.

*/

var RedistributeAutoGeneratedMessages = "privileged"

/*

Should rejection notes from approvals be sent to the requestors?

*/

var ApprovalRejectionNotes = true

/*

Should approval tickets only be viewed and modified through the standard
approval interface?  With this setting enabled (by default), any attempt to use
the normal ticket display and modify page for approval tickets will be
redirected.

For example, with this option set to true and an approval ticket #123:

    /Ticket/Display.html?id=123

is redirected to

    /Approval/Display.html?id=123

With this option set to false, the redirect won't happen.

=back

*/

var ForceApprovalsView = true

/*

=head2 Extra security

This is a list of extra security measures to enable that help keep your RT
safe.  If you don't know what these mean, you should almost certainly leave the
defaults alone.

=over 4


If set to true, the ExecuteCode right will be removed from
all users, B<including> the superuser.  This is intended for when RT is
installed into a shared environment where even the superuser should not
be allowed to run arbitrary Perl code on the server via scrips.

*/

var DisallowExecuteCode = false

/*

If set to false, framekiller javascript will be disabled and the
X-Frame-Options: DENY header will be suppressed from all responses.
This disables RT's clickjacking protection.

*/

var Framebusting = true

/*

If set to false, the HTTP Referer (sic) header will not be
checked to ensure that requests come from RT's own domain.  As RT allows
for GET requests to alter state, disabling this opens RT up to
cross-site request forgery (CSRF) attacks.

*/

var RestrictReferrer = true

/*

If set to false, RT will allow the user to log in from any link
or request, merely by passing in user and pass parameters; setting
it to true forces all logins to come from the login box, so the
user is aware that they are being logged in.  The default is off, for
backwards compatability.

*/

var RestrictLoginReferrer = false

/*

This is a list of hostname:port combinations that RT will treat as being
part of RT's domain. This is particularly useful if you access RT as
multiple hostnames or have an external auth system that needs to
redirect back to RT once authentication is complete.

*/

// var @ReferrerWhitelist = qw(www.example.com:443  www3.example.com:80)

/*

If the "RT has detected a possible cross-site request forgery" error is triggered
by a host:port sent by your browser that you believe should be valid, you can copy
the host:port from the error message into this list.

Simple wildcards, similar to SSL certificates, are allowed.  For example:

    *.example.com:80    # matches foo.example.com
                        # but not example.com
                        #      or foo.bar.example.com

    www*.example.com:80 # matches www3.example.com
                        #     and www-test.example.com
                        #     and www.example.com


*/

// var @ReferrerWhitelist = qw()

/*

%ReferrerComponents is the hash to customize referrer checking behavior when
$RestrictReferrer is enabled, where you can whitelist or blacklist the
components along with their query args. e.g.

*/

// var %ReferrerComponents,
//         ( "/Foo.html" => 1, "/Bar.html" => 0, "/Baz.html" => [ "id", "results" ] )
//     );

/*

With this, '/Foo.html' will be whitelisted, and '/Bar.html' will be blacklisted.
'/Baz.html' with id/results query arguments will be whitelisted but blacklisted
if there are other query arguments.

*/

// var %ReferrerComponents

/*

If set to false, the C<X-Content-Type-Options: nosniff> header will be omitted on
attachments.  Because RT does not filter HTML content in unknown content types,
disabling this opens RT up to cross-site scripting (XSS) attacks by allowing
the execution of arbitrary Javascript when the browser detects HTML-looking
data in an attachment with an unknown content type.

*/

var StrictContentTypes = true

/*

This sets the default cost parameter used for the bcrypt key
derivation function.  Valid values range from 4 to 31, inclusive, with
higher numbers denoting greater effort.

*/

var BcryptCost = 12

/*

=back

=head2 Internationalization

=over 4

An array that contains languages supported by RT's
internationalization interface.  Defaults to all *.po lexicons;
setting it to C<qw(en ja)> will make RT bilingual instead of
multilingual, but will save some memory.

*/

var LexiconLanguages = []string{`*`}

/*

An array that contains default encodings used to guess which charset
an attachment uses, if it does not specify one explicitly.  All
options must be recognized by L<Encode::Guess>.

The first element may also be *, which attempts encoding detection
using L<Encode::Detect::Detector>.  This uses Mozilla's character
detection library to examine the bytes, and use frequency metrics to
rank the options.  This detection may fail (and fall back to other
options in the @EmailInputEncodings list) if no decoding has high
enough confidence metrics.  As of L<Encode::Detect::Detector> version
1.01, it knows the following encodings:

    big5-eten
    cp1250
    cp1251
    cp1253
    cp1255
    cp855
    cp866
    euc-jp
    euc-kr
    euc-tw
    gb18030
    iso-8859-2
    iso-8859-5
    iso-8859-7
    iso-8859-11
    koi8-r
    MacCyrillic
    shiftjis
    utf-8

*/

// var @EmailInputEncodings = qw(utf-8 iso-8859-1 us-ascii)

/*

The charset for localized email.  Must be recognized by Encode.

*/

var EmailOutputEncoding = "utf-8"

/*

=back

=head2 Date and time handling

=over 4


You can choose date and time format.  See the "Output formatters"
section in perldoc F<lib/RT/Date.pm> for more options.  This option
can be overridden by users in their preferences.

Some examples:

*/

var DateTimeFormat = "LocalizedDateTime"

/*
var $DateTimeFormat = { Format => "ISO", Seconds => 0 });>

var $DateTimeFormat = "RFC2822");>

var $DateTimeFormat = { Format => "RFC2822", Seconds => 0, DayOfWeek => 0 });>

var $DateTimeFormat = "DefaultFormat"
*/

/*

# Next two options are for Time::ParseDate

Set this to true if your local date convention looks like "dd/mm/yy"
instead of "mm/dd/yy". Used only for parsing, not for displaying
dates.

*/

var DateDayBeforeMonth = true

/*

Should an unspecified day or year in a date refer to a future or a
past value? For example, should a date of "Tuesday" default to mean
the date for next Tuesday or last Tuesday? Should the date "March 1"
default to the date for next March or last March?

Set $AmbiguousDayInPast for the last date, or
$AmbiguousDayInFuture for the next date; the default is usually
correct.  If both are set, $AmbiguousDayInPast takes precedence.

*/

var AmbiguousDayInPast = false

var AmbiguousDayInFuture = false

/*

Use this to set the default units for time entry to hours instead of
minutes.  Note that this only effects entry, not display.

*/

var DefaultTimeUnitsToHours = 0

/*

By default, events in the iCal feed on the ticket search page
Set
$TimeInICal if you have start or due dates on tickets that
have significant time values and you want those times to be
included in the events in the iCal feed.

This option can also be set as an individual user preference.

*/

var TimeInICal = 0

/*

By default, RT parses an unknown date first with L<Time::ParseDate>, and if
this fails with L<DateTime::Format::Natural>.
$PreferDateTimeFormatNatural changes this behavior to first parse with
L<DateTime::Format::Natural>, and if this fails with L<Time::ParseDate>.
This gives you the possibility to use the more advanced features of
L<DateTime::Format::Natural>.
For example with L<Time::ParseDate> it isn't possible to get the
'first day of the last month', where L<DateTime::Format::Natural> supports
this with 'last month'.

Be aware that L<Time::ParseDate> and L<DateTime::Format::Natural> have
different definitions for the relative date and time syntax.
L<Time::ParseDate> returns for 'last month' this DayOfMonth from the last month.
L<DateTime::Format::Natural> returns for 'last month' the first day of the last
month. So changing this config option maybe changes the results of your saved
searches.

*/

var PreferDateTimeFormatNatural = 0

/*

=back

=head2 Authorization and user configuration

=over 4


If $WebRemoteUserAuth is defined, RT will defer to the environment's
REMOTE_USER variable, which should be set by the webserver's
authentication layer.

*/

var WebRemoteUserAuth = ""

/*

If $WebRemoteUserContinuous is defined, RT will check for the
REMOTE_USER on each access.  If you would prefer this to only happen
once (at initial login) set this to false.  The default
setting will help ensure that if your webserver's authentication layer
deauthenticates a user, RT notices as soon as possible.

*/

var WebRemoteUserContinuous = true

/*

If $WebFallbackToRTLogin is defined, the user is allowed a
chance of fallback to the login screen, even if REMOTE_USER failed.

*/

var WebFallbackToRTLogin = ""

/*

By default, $LogoutURL is set to RT's logout page. When an
external service is used to log into RT, $LogoutURL can be set
to the identity provider's logout URL. Include the full path to the
logout endpoint, for example: 'https://www.example.com/logout'.

*/

var LogoutURL = "/NoAuth/Logout.html"

/*

$WebRemoteUserGecos means to match 'gecos' field as the user
identity; useful with mod_auth_external.

*/

var WebRemoteUserGecos = false

/*

$WebRemoteUserAutocreate will create users under the same name as
REMOTE_USER upon login, if they are missing from the Users table.

*/

var WebRemoteUserAutocreate = false

/*

If $WebRemoteUserAutocreate is set to true, $UserAutocreateDefaultsOnLogin
will be passed to L<RT::User/Create> when the user is created.
Use it to set default settings for the new user account, such
as creating users as unprivileged with:

*/

// var $UserAutocreateDefaultsOnLogin = { Privileged => 0 }

/*

or privileged:

*/

// var $UserAutocreateDefaultsOnLogin = { Privileged => 1 }

/*

The settings must be in a hashref as shown.

This option is also used if you have External Auth configured.
See L</External Authentication and Authorization> for details.

*/

var UserAutocreateDefaultsOnLogin = false

/*

$WebSessionClass is the class you wish to use for storing sessions.  On
MySQL, Pg, and Oracle it defaults to using your database, in other cases
sessions are stored in files using L<Apache::Session::File>. Other installed
Apache::Session::* modules can be used to store sessions.

*/

var WebSessionClass = ""

/*

%WebSessionProperties is the hash to configure class L</$WebSessionClass>
in case custom class is used. By default it's empty and values are picked
depending on the class. Make sure that it's empty if you're using DB as session
backend.

*/

// var %WebSessionProperties

/*

By default, RT's user sessions persist until a user closes his or her
browser. With the $AutoLogoff option you can setup session lifetime
in minutes. A user will be logged out if he or she doesn't send any
requests to RT for the defined time.

*/

var AutoLogoff = false

/*

The number of seconds to wait after logout before sending the user to
the login page. By default, 1 second, though you may want to increase
this if you display additional information on the logout page.

*/

var LogoutRefresh = 1

/*

By default, RT's session cookie isn't marked as "secure". Some web
browsers will treat secure cookies more carefully than non-secure
ones, being careful not to write them to disk, only sending them over
an SSL secured connection, and so on. To enable this behavior, set
$WebSecureCookies to true.  NOTE: You probably don't want to turn this
on I<unless> users are only connecting via SSL encrypted HTTPS
connections.

*/

var WebSecureCookies = false

/*

Default RT's session cookie to not being directly accessible to
javascript.  The content is still sent during regular and AJAX requests,
and other cookies are unaffected, but the session-id is less
programmatically accessible to javascript.  Turning this off should only
be necessary in situations with odd client-side authentication
requirements.

*/

var WebHttpOnlyCookies = true

/*

$MinimumPasswordLength defines the minimum length for user
Setting it to false disables this check.

*/

var MinimumPasswordLength = 5

/*

=back

=head3 External Authentication and Authorization

RT has a built-in module for integrating with a directory service like
LDAP or Active Directory for authentication (login) and authorization
(enabling/disabling users and setting user attributes). The core configuration
settings for the service are listed here. Additional details are available
in the L<RT::Authen::ExternalAuth> module documentation.

See also L</$UserAutocreateDefaultsOnLogin> for configuring
defaults for autocreated users.

=over 4


This option, along with the following options, activate and configure authentication
via a resource external to RT. All of the configuration for your external authentication
service, like LDAP or Active Directory, are defined in a data structure in this option.
You can find full details on the configuration
options in the L<RT::Authen::ExternalAuth> documentation.


# No defaults are set for ExternalAuth because this is an optional feature.


Sets the priority of authentication resources if you have multiple configured.
RT will attempt authorization with each resource, in order, until one succeeds or
no more remain. See L<RT::Authen::ExternalAuth> for details.


Sets the order of resources for querying user information if you have multiple
configured. RT will query each resource, in order, until one succeeds or
no more remain. See L<RT::Authen::ExternalAuth> for details.


A hashref of options to set for users who are autocreated on login via
ExternalAuth. For example, you can automatically make "Privileged" users
who were authenticated and created from LDAP or Active Directory.
See L<RT::Authen::ExternalAuth> for details.


Users should still be autocreated by RT as internal users if they
fail to exist in an external service; this is so requestors who
are not in LDAP can still be created when they email in.
See L<RT::Authen::ExternalAuth> for details.


If you have a mix of RT and federated authentication, RT can't directly
verify a user's password against the federated IdP. You can explicitly
disable the password prompt when creating a token by setting this option
to true (1).

=back


*/

var DisablePasswordForAuthToken = false

/*

=head2 Initialdata Formats

RT supports pluggable data format parsers for F<initialdata> files.

If you add format handlers, note that you can remove the perl entry if you
don't want it available. B<Removing the default perl entry may cause problems
installing plugins and RT updates>. If so, re-enable it temporarily.

=over 4


Set the $InitialdataFormatHandlers to an arrayref containing a list of
format handler modules. The 'perl' entry is the system default, and handles
perl-style intialdata files.

The JSON format handler is also available in RT, but it is not loaded by
default. Add it to your configuration as shown below to enable it.

*/

// var $InitialdataFormatHandlers,
//          [
//             'perl',
//             'RT::Initialdata::JSON',
//             'RT::Extension::Initialdata::Foo',
//             ...
//          ]
//        );

/*

=back

*/

// var InitialdataFormatHandlers = [ "perl" ]

/*

=head2 Development options

=over 4


RT comes with a "Development mode" setting.  This setting, as a
convenience for developers, turns on several of development options
that you most likely don't want in production:

=over 4


Disables CSS and JS minification and concatenation.  Both CSS and JS
will be instead be served as a number of individual smaller files,
unchanged from how they are stored on disk.

Uses L<Module::Refresh> to reload changed Perl modules on each
request.

Turns off Mason's static_source directive; this causes Mason to
reload template files which have been modified on disk.

Turns on Mason's HTML error_format; this renders compilation errors
to the browser, along with a full stack trace.  It is possible for
stack traces to reveal sensitive information such as passwords or
ticket content.

Turns off caching of callbacks; this enables additional callbacks to
be added while the server is running.

=back

*/

var DevelMode = false

/*

What abstract base class should RT use for its records. You should
probably never change this.

Valid values are DBIx::SearchBuilder::Record or
DBIx::SearchBuilder::Record::Cachable

*/

var RecordBaseClass = "DBIx::SearchBuilder::Record::Cachable"

/*

@MasonParameters is the list of parameters for the constructor of
HTML::Mason's Apache or CGI Handler.  This is normally only useful for
debugging, e.g. profiling individual components with:

    use MasonX::Profiler; # available on CPAN
*/

// var @MasonParameters = (preamble => 'my $p = MasonX::Profiler->new($m, $r);')

// var @MasonParameters = ()

/*

RT has rudimentary SQL statement logging support; simply set
$StatementLog to be the level that you wish SQL statements to be
logged at.

Enabling this option will also expose the SQL Queries page in the
Admin -> Tools menu for SuperUsers.

*/

var StatementLog = false

/*

RT enables SQL bind parameters for all searches by default, which improves
performance especially for Oracle. If you need to disable this for some
reason, add this to config:

    $ENV{SB_PREFER_BIND} = 0;

See also L<DBIx::SearchBuilder/BuildSelectQuery>.

=back


=head1 Web interface

=head2 Base configuration

=over 4


This determines the default stylesheet the RT web interface will use.
RT ships with several themes by default:

  elevator-light  The default light theme for RT 5
  elevator-dark   The dark theme for RT 5

This value actually specifies a directory in F<share/static/css/>
from which RT will try to load the file main.css (which should @import
any other files the stylesheet needs).  This allows you to easily and
cleanly create your own stylesheets to apply to RT.  This option can
be overridden by users in their preferences.

*/

var WebDefaultStylesheet = "elevator-light"

/*

This is the email address of the person who administers and provides
support for RT itself. If set, it is displayed on the RT login page
and as such is likely to receive email from users who are unable to
log in.

*/

var RTSupportEmail = ""

/*

Starting with RT 5.0, RT's web interface is fully responsive and
will render correctly on most mobile devices. However, RT also has
a mobile-optimized mode that shows a limited feature set
focused on ticket updates. To default to this site when RT is accessed
from a mobile device, enable this option (set to true).

*/

var ShowMobileSite = false

/*

Use this to select the default queue name that will be used for
creating new tickets. You may use either the queue's name or its
ID. This only affects the queue selection boxes on the web interface.

*/

var DefaultQueue = "General"

/*

When a queue is selected in the new ticket dropdown, make it the new
default for the new ticket dropdown.

*/

var RememberDefaultQueue = true

/*

Hide all links and portlets related to Reminders by setting this to false

*/

var EnableReminders = true

/*

Set @CustomFieldValuesSources to a list of class names which extend
L<RT::CustomFieldValues::External>.  This can be used to pull lists of
custom field values from external sources at runtime.

*/

// var @CustomFieldValuesSources = ()

/*

Set @CustomFieldValuesCanonicalizers to a list of class names which extend
L<RT::CustomFieldValues::Canonicalizer>. This can be used to rewrite
(canonicalize) values entered by users to fit some defined format.

See the documentation in L<RT::CustomFieldValues::Canonicalizer> for adding
your own canonicalizers.

*/

// var @CustomFieldValuesCanonicalizers = qw(
//     RT::CustomFieldValues::Canonicalizer::Uppercase
//     RT::CustomFieldValues::Canonicalizer::Lowercase
// ));

/*

This option affects the display of ticket, user, group, and asset custom fields
in the web interface. It does not address the sorting of custom fields within
the groupings; that ordering is controlled by the Ticket Custom Fields tab in
Queue configuration in the Admin UI. Asset custom field ordering is
found in the Asset Custom Fields tab in Catalog configuration.

A nested data structure defines how to group together custom fields
under a mix of built-in and arbitrary headings ("groupings").

Set %CustomFieldGroupings to a nested structure similar to the following:

*/

//var %CustomFieldGroupings,
//        'RT::Ticket' => [
//            'Grouping Name'     => ['CF Name', 'Another CF'],
//            'Another Grouping'  => ['Some CF'],
//            'Dates'             => ['Shipped date'],
//        ],
//        'RT::User' => [
//            'Phones' => ['Fax number'],
//        ],
//        'RT::Asset' => [
//            'Asset Details' => ['Serial Number', 'Manufacturer', 'Type', 'Tracking Number'],
//            'Dates'         => ['Support Expiration', 'Issue Date'],
//        ],
//        'RT::Group' => [
//            'Basics' => ['Department'],
//        ],
//    );

/*

The first level keys are record types for which CFs may be used, and the
values are either hashrefs or arrayrefs -- if arrayrefs, then the
order of grouping entries is preserved during display, otherwise groupings
are displayed alphabetically. The second level keys are the grouping names
and the values are array refs containing a list of CF names.

For RT::Ticket and RT::Asset, you can specify global, and queue
or catalog level groupings. For example, if you wanted to diplay some
groupings only on tickets in the General queue, you can create an entry
for 'General'. Global configurations then go in 'Default' as shown below.

    'RT::Ticket' => {
        'Default' => [
            'Grouping Name'    => [ 'CF Name' ],
        ],
        'General' => [
            'Grouping Name'    => [ 'CF Name', 'Another CF' ],
            'Another Grouping' => ['Some CF'],
            'Dates'            => ['Shipped date'],
        ],
    },

There are several special built-in groupings which RT displays in
specific places (usually the collapsible box of the same title).  The
ordering of these standard groupings cannot be modified.  You may also
only append Custom Fields to the list in these boxes, not reorder or
remove core fields.

For RT::Ticket, these groupings are: Basics, Dates, Links, People

For RT::User: Identity, C<Access control>, Location, Phones

For RT::Group: Basics

For RT::Asset: Basics, Dates, People, Links

Extensions may also add their own built-in groupings, refer to the individual
extension documentation for those.


Set $CanonicalizeRedirectURLs to true to use $WebURL when
redirecting rather than the one we get from %ENV.

Apache's UseCanonicalName directive changes the hostname that RT
finds in %ENV.  You can read more about what turning it On or Off
means in the documentation for your version of Apache.

If you use RT behind a reverse proxy, you almost certainly want to
enable this option.

*/

var CanonicalizeRedirectURLs = false

/*

Set $CanonicalizeURLsInFeeds to true to use $WebURL in feeds
rather than the one we get from request.

If you use RT behind a reverse proxy, you almost certainly want to
enable this option.

*/

var CanonicalizeURLsInFeeds = false

/*

A list of additional JavaScript files to be included in head.

*/

// var @JSFiles = qw//

/*

A list of additional CSS files to be included in head.

If you're a plugin author, refer to RT->AddStyleSheets.

*/

// var @CSSFiles = qw//

/*

This determines how user info is displayed. 'concise' will show the
first of RealName, Name or EmailAddress that has a value. 'verbose' will
show EmailAddress, and the first of RealName or Name which is defined.
The default, 'role', uses 'verbose' for unprivileged users, and the Name
followed by the RealName for privileged users.

*/

var UsernameFormat = "role"

/*

This controls the display of lists of users returned from the User
Summary Search. The display of users in the Admin interface is
controlled by %AdminSearchResultFormat.

*/

// var $UserSearchResultFormat,
//          q{ '<a href="__WebPath__/User/Summary.html?id=__id__">__id__</a>/TITLE:#'}
//         .q{,'<a href="__WebPath__/User/Summary.html?id=__id__">__Name__</a>/TITLE:Name'}
//         .q{,__RealName__, __EmailAddress__}
// );

/*

A list of portlets to be displayed on the User Summary page.
By default, we show all of the available portlets.
Extensions may provide their own portlets for this page.

*/

// var @UserSummaryPortlets = (qw/ExtraInfo CreateTicket ActiveTickets InactiveTickets UserAssets /)

/*

This controls what information is displayed on the User Summary
portal. By default the user's Real Name, Email Address and Username
are displayed. You can remove these or add more as needed. This
expects a Format string of user attributes. Please note that not all
the attributes are supported in this display because we're not
building a table.

*/

var UserSummaryExtraInfo = "RealName, EmailAddress, Name"

/*

Control the appearance of the Active and Inactive ticket lists in the
User Summary.

*/

//var $UserSummaryTicketListFormat = q{
//       '<B><A HREF="__WebPath__/Ticket/Display.html?id=__id__">__id__</a></B>/TITLE:#',
//       '<B><A HREF="__WebPath__/Ticket/Display.html?id=__id__">__Subject__</a></B>/TITLE:Subject',
//       Status,
//       QueueName,
//       Owner,
//       Priority,
//       '__NEWLINE__',
//       "",
//       '<small>__Requestors__</small>',
//       '<small>__CreatedRelative__</small>',
//       '<small>__ToldRelative__</small>',
//       '<small>__LastUpdatedRelative__</small>',
//       '<small>__TimeLeft__</small>'
//});

/*

Usually you don't want to set these options. The only obvious reason
is if RT is accessible via https protocol on a non standard port, e.g.
'https://rt.example.com:9999'. In all other cases these options are
computed using $WebDomain, $WebPort and $WebPath.

$WebBaseURL is the scheme, server and port
(e.g. 'http://rt.example.com') for constructing URLs to the web
UI. $WebBaseURL doesn't need a trailing /.

$WebURL is the $WebBaseURL, $WebPath and trailing /, for
example: 'http://www.example.com/rt/'.


my $port = RT->Config->Get('WebPort');

*/

//var $WebBaseURL,
//    ($port == 443? 'https': 'http') .'://'
//    . RT->Config->Get('WebDomain')
//    . ($port != 80 && $port != 443? ":$port" : "")
//);

// var WebURL = WebBaseURL + WebPath' + "/"

/*

$WebImagesURL points to the base URL where RT can find its images.
Define the directory name to be used for images in RT web documents.


*/

var WebImagesURL = WebPath + "/static/images/"

/*

$LogoURL points to the URL of the RT logo displayed in the web UI.
This can also be configured via the web UI.

*/

var LogoURL = WebImagesURL + "request-tracker-logo.svg"

/*

$LogoLinkURL is the URL that the RT logo hyperlinks to.

*/

var LogoLinkURL = "http://bestpractical.com"

/*

$LogoAltText is a string of text for the alt-text of the logo. It
will be passed through loc for localization.

*/

var LogoAltText = "Request Tracker logo"

/*

What portion of RT's URL space should not require authentication.  The
default is almost certainly correct, and should only be changed if you
are extending RT.

*/

// var $WebNoAuthRegex = qr{^ (?:/+NoAuth/ | /+REST/\d+\.\d+/NoAuth/) }x

/*

By default, RT clears its database cache after every page view.  This
ensures that you've always got the most current information when
Setting
$WebFlushDbCacheEveryRequest to false will turn this off, which will
speed RT up a bit, at the expense of a tiny bit of data accuracy.

*/

var WebFlushDbCacheEveryRequest = true

/*

The L<GD> module (which RT uses for graphs) ships with a built-in font
that doesn't have full Unicode support. You can use a given TrueType
font for a specific language by setting %ChartFont to (language =E<gt>
the absolute path of a font) pairs. Your GD library must have support
for TrueType fonts to use this option. If there is no entry for a
language in the hash then font with 'others' key is used.

RT comes with two TrueType fonts covering most available languages.

*/

//var %ChartFont,
//    'zh-cn'  => "$RT::FontPath/DroidSansFallback.ttf",
//    'zh-tw'  => "$RT::FontPath/DroidSansFallback.ttf",
//    'ja'     => "$RT::FontPath/DroidSansFallback.ttf",
//    'others' => "$RT::FontPath/DroidSans.ttf",
//);

/*

RT stores dates using the UTC timezone in the DB, so charts grouped by
Set $ChartsTimezonesInDB to true
to enable timezone conversions using your DB's capabilities. You may
need to do some work on the DB side to use this feature, read more in
F<docs/customizing/timezones_in_charts.pod>.

At this time, this feature only applies to MySQL and PostgreSQL.

*/

var ChartsTimezonesInDB = false

/*

An array of 6-digit hexadecimal RGB color values used for chart series.  By
default there are 12 distinct colors.

*/

//var @ChartColors = qw(
//    66cc66 ff6666 ffcc66 663399
//    3333cc 339933 993333 996633
//    33cc33 cc3333 cc9933 6633cc
//));

/*

Set this to false to disable Chart in JavaScript.

*/

var EnableJSChart = true

/*

The color scheme to use for Chart in Javascript. By default it's
I<brewer.Paired12>.  The full list is:
L<https://nagix.github.io/chartjs-plugin-colorschemes/colorchart.html>

*/

var JSChartColorScheme = "brewer.Paired12"

/*

=back



=head2 Home page

=over 4

$DefaultSummaryRows is default number of rows displayed in for
search results on the front page.

*/

var DefaultSummaryRows = 10

/*

This setting defines the possible homepage and search result refresh
options. Each value is a number of seconds. You should not include a value
of 0, as that is always provided as an option.

See also L</$HomePageRefreshInterval> and L</$SearchResultsRefreshInterval>.

*/

// var @RefreshIntervals = qw(120 300 600 1200 3600 7200)

/*

$HomePageRefreshInterval is default number of seconds to refresh
the RT home page. Choose from any value in L</@RefreshIntervals>,
or the default of 0 for no automatic refresh.

*/

var HomePageRefreshInterval = 0

/*

$HomepageComponents is an arrayref of allowed components on a
user's customized homepage ("RT at a glance").

*/

//var $HomepageComponents,
//    [
//        qw(QuickCreate QueueList QueueListAllStatuses MyAdminQueues MySupportQueues MyReminders RefreshHomepage Dashboards SavedSearches FindUser MyAssets FindAsset FindGroup) # loc_qw
//    ]
//);

/*

=back




=head2 Ticket search

=over 4

Historically, ACLs were checked on display, which could lead to empty
Set $UseSQLForACLChecks to false
to go back to this method; this will reduce the complexity of the
generated SQL statements, at the cost of the aforementioned bugs.

*/

// var $UseSQLForACLChecks = true

/*

On the display page of a ticket from search results, RT provides links
to the first, next, previous and last ticket from the results.  In
order to build these links, RT needs to re-run the original search
and fetch the full result set from the database. If the original
search was resource-intensive, this will then slow down diplay of the
ticket page.

Set $TicketsItemMapSize to the number of tickets you want RT to examine
to build these links. If the full result set is larger than this
Set this to zero to
always examine all results. This can improve performance for searches
with large result sets.

Set $ShowSearchNavigation to false to not build these links at all and
completely avoid re-running the original search query.

*/

var TicketsItemMapSize = 1000

var ShowSearchNavigation = true

/*

$SearchResultsRefreshInterval is default number of seconds to refresh
search results in RT. Choose from any value in L</@RefreshIntervals>, or
the default of 0 for no automatic refresh.

*/

var SearchResultsRefreshInterval = 0

/*

$DefaultSearchResultFormat is the default format for RT search results

*/

//Set ($DefaultSearchResultFormat, qq{
//   '<B><A HREF="__WebPath__/Ticket/Display.html?id=__id__">__id__</a></B>/TITLE:#',
//   '<B><A HREF="__WebPath__/Ticket/Display.html?id=__id__">__Subject__</a></B>/TITLE:Subject',
//   Status,
//   QueueName,
//   Owner,
//   Priority,
//   '__NEWLINE__',
//   '__NBSP__',
//   '<small>__Requestors__</small>',
//   '<small>__CreatedRelative__</small>',
//   '<small>__ToldRelative__</small>',
//   '<small>__LastUpdatedRelative__</small>',
//   '<small>__TimeLeft__</small>'});

/*

This is the format of ticket search result for "Download User Tickets" links. It
defaults to DefaultSearchResultFormat for privileged users and DefaultSelfServiceSearchResultFormat
for unprivileged users if it's not set.

*/

var UserTicketDataResultFormat = ""

/*

This is the format of the user search result for "Download User Data" links.

*/

//var $UserDataResultFormat = "'__id__', '__Name__', '__EmailAddress__', '__RealName__',\
//                            '__NickName__', '__Organization__', '__HomePhone__', '__WorkPhone__',\
//                            '__MobilePhone__', '__PagerPhone__', '__Address1__', '__Address2__',\
//                            '__City__', '__State__','__Zip__', '__Country__', '__Gecos__', '__Lang__',\
//                            '__Timezone__', '__FreeFormContactInfo__'");

/*

This is the format of the user transaction search result for "Download User Transaction Data" links.

*/

// var $UserTransactionDataResultFormat = "'__ObjectId__/TITLE:Ticket Id', '__id__', '__Created__', '__Description__',\
//                                         '__OldValue__', '__NewValue__', '__Content__'");

/*

What Tickets column should we order by for RT Ticket search results.

*/

var DefaultSearchResultOrderBy = "id"

/*

When ordering RT Ticket search results by $DefaultSearchResultOrderBy,
should the sort be ascending (ASC) or descending (DESC).

*/

var DefaultSearchResultOrder = "ASC"

/*

Display search result count on ticket lists. Defaults to true (show them).

*/

var ShowSearchResultCount = true

/*

Full text search (FTS) without database indexing is a very slow
operation, and is thus disabled by default.

Before setting Indexed to true, read F<docs/full_text_indexing.pod> for
the full details of FTS on your particular database.

It is possible to enable FTS without database indexing support, simply
by setting the Enable key to true, while leaving Indexed set to false.
This is not generally suggested, as unindexed full-text searching can
cause severe performance problems.

*/

// var %FullTextSearch = Enable  => 0, Indexed => 0,

/*

On some systems, very large attachments can cause memory and other
performance issues for the indexer making it unable to complete
indexing. Adding resources like memory and CPU will solve this
issue, but in cases where that isn't possible, this option
sets a maximum size in bytes on attachments to index. Attachments
larger than this limit are skipped and will not be available to
full text searches.

# Default 0 means no limit
*/

var MaxFulltextAttachmentSize = 0

/*

If $DontSearchFileAttachments is set to true, then uploaded files
(attachments with file names) are not searched during content
search.

Note that if you use indexed FTS then named attachments are still
indexed by default regardless of this option.

*/

var DontSearchFileAttachments = false

/*

When query in simple search doesn't have status info, use this to only
search active ones.

*/

var OnlySearchActiveTicketsInSimpleSearch = true

/*

When only one ticket is found in search, use this to redirect to the
ticket display page automatically.

*/

var SearchResultsAutoRedirect = false

/*

Allows users to update tickets directly on search results and ticket display
pages.

*/

var InlineEdit = true

/*

This setting allows you to control which panels on display pages participate
in inline edit, as well as fine-tuning their specific behavior.

Much like L</%CustomFieldGroupings>, you first specify a record type you
want to configure (though since currently only ticket display supports inline
edit, keys besides RT::Ticket are ignored). Then, for each panel, you
specify its behavior. The valid behaviors are:

=over 4


The panel will have an "Edit" link in the top right, which when clicked
immediately activates inline edit. The "Edit" link will change to
"Cancel" to restore the readonly display.

Much like link, except you may click anywhere inside the panel to
activate inline edit.

Turns off inline edit entirely for this panel.

Turns off the readonly display for this panel, providing I<only> inline
edit capabilities.

=back

You may also provide the special key _default inside a record type to
specify a default behavior for all panels.

This sample configuration will provide a default inline edit behavior of
click, but also specifies different behaviors for several other panels.
Note that the non-standard panel names "Grouping Name" and "Another
Grouping" are created by the L</%CustomFieldGroupings> setting.

*/

//var %InlineEditPanelBehavior,
//        'RT::Ticket' => {
//            '_default'          => 'click',
//
//            'Grouping Name'     => 'link',
//            'Another Grouping'  => 'click',
//            'Dates'             => 'always',
//            'Links'             => 'hide',
//            'People'            => 'link',
//        },
//    );

/*

=back

=head2 Ticket options

=over 4

Set to true to display Total Time Worked in the Basics section.

Total Time Worked is a dynamic value containing a sum of Time Worked for
the parent and all child tickets.  This value is generated when displaying
the ticket and automatically updates when a child ticket is added or removed.
Total Time Worked follows only parent/child link relationships.  Tickets
linked with depends-on or refers-to links are not included.

Total Time Worked is also available as a column for reports generated with
the Query Builder.

*/

var DisplayTotalTimeWorked = false

/*

This determines if the 'More about requestor' box on
Ticket/Display.html is shown for Privileged Users.

*/

var ShowMoreAboutPrivilegedUsers = false

/*

This can be set to Active, Inactive, All or None.  It controls what
ticket list will be displayed in the 'More about requestor' box on
Ticket/Display.html.  This option can be controlled by users also.

*/

var MoreAboutRequestorTicketList = "Active"

/*

Control the appearance of the ticket lists in the 'More About Requestors' box.

*/

//var $MoreAboutRequestorTicketListFormat = q{
//       '<a href="__WebPath__/Ticket/Display.html?id=__id__">__id__</a>',
//       '__Owner__',
//       '<a href="__WebPath__/Ticket/Display.html?id=__id__">__Subject__</a>',
//       '__Status__',
//});

/*

By default, the 'More about requestor' box on Ticket/Display.html
shows the Requestor's name and ticket list.  If you would like to see
extra information about the user, this expects a Format string of user
attributes.  Please note that not all the attributes are supported in
this display because we're not building a table.

Example:
*/

var MoreAboutRequestorExtraInfo = "Organization = Address1"

// var $MoreAboutRequestorExtraInfo = ""

/*

By default, the 'More about requestor' box on Ticket/Display.html
shows all the groups of the Requestor.  Use this to limit the number
of groups; a value of undef removes the group display entirely.

*/

var MoreAboutRequestorGroupsLimit = 0

/*

Should the ticket create and update forms use a more space efficient
two column layout.  This layout may not work in narrow browsers if you
set a MessageBoxWidth (below).

*/

var UseSideBySideLayout = true

/*

When displaying a list of Ticket Custom Fields for editing, RT
defaults to a 2 column list.  If you set this to true, it will instead
display the Custom Fields in a single column.

*/

var EditCustomFieldsSingleColumn = false

/*

If set to true, RT will prompt users when there are new,
unread messages on tickets they are viewing.

*/

var ShowUnreadMessageNotifications = false

/*

If set to true, the owner drop-downs for ticket update/modify and the query
builder are replaced by text fields that autocomplete.  This can
alleviate the sometimes huge owner list for installations where many
users have the OwnTicket right.

The Owner entry is automatically converted to an autocomplete box if the list
of owners exceeds $DropdownMenuLimit items. However, the query to generate
the list of owners is still run and this can increase page load times. If
your owner lists exceed the limit and you are using the autocomplete box, you
can improve performance by explicitly setting $AutocompleteOwners.

Drop down doesn't show unprivileged users. If your setup allows unprivileged
to own ticket then you have to enable autocompleting.

*/

var AutocompleteOwners = false

/*

The Owner dropdown menu, used in various places in RT including the Query
Builder and ticket edit pages, automatically changes from a dropdown menu to
an autocomplete field once the menu holds more than the $DropdownMenuLimit
owners. Dropdown menus become more difficult to use when they contain a large
number of values and the autocomplete textbox can be more usable.

If you have very large numbers of users who can be owners, this can cause
slow page loads on pages with an Owner selection. See L</$AutocompleteOwners>
for a way to potentially speed up page loads.

*/

var DropdownMenuLimit = 50

/*

If set to true, the owner drop-downs for the query builder are always
replaced by text field that autocomplete and $AutocompleteOwners
is ignored. Helpful when owners list is huge in the query builder.

*/

var AutocompleteOwnersForSearch = false

/*

If set to true, any queue drop-downs are replaced by text fields that
autocomplete. This can alleviate the sometimes huge queue list for
installations with many queues, and can also increase page load
times in some cases. A user can override this setting as a personal
preference.

*/

var AutocompleteQueues = false

/*

Used when searching for an Article to Include.

Specifies which fields of L<RT::Article> to match against and how to match
each field when autocompleting articles.  Valid match methods are LIKE,
STARTSWITH, ENDSWITH, =, and !=.  Valid search fields are the core Article
fields, as well as custom fields, including Content, which are specified as
"CF.1234" or "CF.Name"


*/

//var $ArticleSearchFields = {
//    Name         => 'STARTSWITH',
//    Summary      => 'LIKE',
//})

/*

Used by the User Autocompleter as well as the User Search.

Specifies which fields of L<RT::User> to match against and how to match
each field when autocompleting users.  Valid match methods are LIKE,
STARTSWITH, ENDSWITH, =, and !=.  Valid search fields are the core User
fields, as well as custom fields, which are specified as "CF.1234" or
"CF.Name"

*/

//var $UserSearchFields = {
//    EmailAddress => 'STARTSWITH',
//    Name         => 'STARTSWITH',
//    RealName     => 'LIKE',
//})

/*

Specifies which fields of L<RT::Ticket> to match against and how to match each
field when autocompleting users.  Valid match methods are LIKE, STARTSWITH,
ENDSWITH, =, and !=.

Not all Ticket fields are publically accessible and hence won't work for
autocomplete unless you override their accessibility using a local overlay or a
plugin.  Out of the box the following fields are public: id, Subject.

*/

//var $TicketAutocompleteFields = {
//    id      => 'STARTSWITH',
//    Subject => 'LIKE',
//})

/*

Enable this to redirect to the created ticket display page
automatically when using QuickCreate.

*/

var DisplayTicketAfterQuickCreate = false

/*

Setting this to true
causes InterCapped or ALLCAPS words in WikiText fields to automatically
become links to searches for those words.  If used on Articles, it links
to the Article with that name.

*/

var WikiImplicitLinks = false

/*

%LinkedQueuePortlets allows you to display links to tickets in
another queue in a stand-alone portlet on the ticket display page.
This makes it easier to highlight specific ticket links separate from
the standard Links portlet.

For example, you might have a Sales queue that tracks incoming product
requests, and for each ticket you create a linked ticket in the Shipping
queue for each outgoing shipment. You could add the configuration below
to create a stand-alone Shipping portlet on tickets in the Sales queue,
making it easier to see those linked tickets. You might have a Returns
queue to show as well.

*/

//var %LinkedQueuePortlets = (
//        'Sales' => [
//            { 'Shipping'   => [ 'All' ] },
//            { 'Returns'   => [ 'RefersTo' ] },
//        ],
//        'Shipping'   => [
//            { 'Postage' => [ 'DependsOn', 'HasMember' ] },
//        ],
//    ));

/*

You can include multiple linked queues in each ticket and they are
displayed in the order you define them in the configuration. The values
are RT link types: 'DependsOn', 'DependedOnBy', 'HasMember'
(children), 'MemberOf' (parents), 'RefersTo', and 'ReferredToBy'.
'All' lists all linked tickets. You can include multiple link types for
each as shown above.

*/

// var %LinkedQueuePortlets = ()

/*

%LinkedQueuePortletFormats defines the format for displaying
linked tickets in each linked queue portlet defined by %LinkedQueuePortlets.

To change just the General list you would do:

*/

// var %LinkedQueuePortletFormats = General => 'modified configuration'

//var %LinkedQueuePortletFormats,
//    Default =>
//        q{'<b><a href="__WebPath__/Ticket/Display.html?id=__id__">__id__</a></b>/TITLE:#',}.
//        q{'<b><a href="__WebPath__/Ticket/Display.html?id=__id__">__Subject__</a></b>/TITLE:Subject',}.
//        q{Status},
//);

/*

Set $PreviewScripMessages to true if the scrips preview on the ticket
reply page should include the content of the messages to be sent.

*/

var PreviewScripMessages = false

/*

If $SimplifiedRecipients is set, a simple list of who will receive
B<any> kind of mail will be shown on the ticket reply page, instead of a
detailed breakdown by scrip.

*/

var SimplifiedRecipients = false

/*

If $SquelchedRecipients is set, the checkbox list of who will receive
B<any> kind of mail on the ticket reply page are displayed initially as
B<un>checked - which means nobody in that list would get any mail. It
does not affect correspondence done via email yet.

*/

var SquelchedRecipients = false

/*

If set to true, this option will skip ticket menu actions which can't be
completed successfully because of outstanding active Depends On tickets.

By default, all ticket actions are displayed in the menu even if some of
them can't be successful until all Depends On links are resolved or
transitioned to another inactive status.

*/

var HideResolveActionsWithDependencies = false

/*

This determines if we should hide unset fields on ticket display page.
Set this to true to hide unset fields.

*/

var HideUnsetFieldsOnDisplay = false

/*

On ticket comment and correspond there are "One-time Cc" and "One-time Bcc"
fields. As part of this section, RT includes a list of suggested email
addresses based on the correspondence history for that ticket. This list may
grow quite large over time.

Enabling this option will hide the list behind a "(show suggestions)" link to
cut down on page clutter. Once this option is clicked the link will change to
"(hide suggestions)" and the full list of email addresses will be shown.

*/

var HideOneTimeSuggestions = false

/*

Priority is stored as a number internally. This determines whether
Priority is displayed to users as a number or using configured
labels like Low, Medium, High. See L<%PriorityAsString> for details
on this configuration.

Set to false to display
numbers, which was the previous default for RT.

*/

var EnablePriorityAsString = true

/*

This setting allows you to define labels for priority values
available on tickets. RT stores these values internally as a number,
but this number will be hidden if $EnablePriorityAsString is true.
For the configuration, link the labels to numbers as shown below. If
you have more or less priority settings, you can adjust the numbers,
giving a unique number to each.

*/

//var %PriorityAsString,
//        Default => { Low => 0, Medium => 50, High => 100 },
//        General => [ Medium => 50, Low => 0, High => 80, 'On Fire' => 100],
//        Support => 0,
//    );

/*

The key is queue name or "Default", which is the fallback for unspecified
queues. Values can be an ArrayRef, HashRef, or 0.

=over


This is the ordered String => Number map list. Pririty options will be
rendered in the order they are listed in the list.


This is the unordered String => Number map list. Priority options will be
rendered in numerical ascending order.


Priority is rendered as a number.

=back


*/

//var %PriorityAsString,
//    Default => { Low => 0, Medium => 50, High => 100 },
//);

/*

=back

=head2 Group Summary Configuration

Below are configuration options for the Group Summary page.

=over

This controls the display of lists of groups returned from the Group
Summary Search. The display of groups in the Admin interface is
controlled by %AdminSearchResultFormat.


*/

//var $GroupSearchResultFormat,
//         q{ '<a href="__WebPath__/Group/Summary.html?id=__id__">__id__</a>/TITLE:#'}
//        .q{,'<a href="__WebPath__/Group/Summary.html?id=__id__">__Name__</a>/TITLE:Name'}
//);

/*

A list of portlets to be displayed on the Group Summary page.
By default, we show all of the available portlets.
Extensions may provide their own portlets for this page.

*/

// var @GroupSummaryPortlets = (qw/ExtraInfo CreateTicket ActiveTickets InactiveTickets GroupAssets /)

/*

This controls what information is displayed on the Group Summary
portal. By default the group Name and Description are displayed.

*/

var GroupSummaryExtraInfo = "id, Name, Description"

/*

Control the appearance of the Active and Inactive ticket lists in the
Group Summary.

*/

//var $GroupSummaryTicketListFormat = q{
//       '<B><A HREF="__WebPath__/Ticket/Display.html?id=__id__">__id__</a></B>/TITLE:#',
//       '<B><A HREF="__WebPath__/Ticket/Display.html?id=__id__">__Subject__</a></B>/TITLE:Subject',
//       Status,
//       QueueName,
//       Owner,
//       Priority,
//       '__NEWLINE__',
//       "",
//       '<small>__Requestors__</small>',
//       '<small>__CreatedRelative__</small>',
//       '<small>__ToldRelative__</small>',
//       '<small>__LastUpdatedRelative__</small>',
//       '<small>__TimeLeft__</small>'
//})

/*

Specifies which fields of L<RT::Group> to match against and how to match
each field when performing a quick search on groups.  Valid match
methods are LIKE, STARTSWITH, ENDSWITH, =, and !=.  Valid search fields
are id, Name, Description, or custom fields, which are specified as
"CF.1234" or "CF.Name"

*/

//var $GroupSearchFields = {
//    id          => '=',
//    Name        => 'LIKE',
//    Description => 'LIKE',
//})

/*

Defines whether unprivileged users (users of SelfService) are allowed
Setting this option to true means
unprivileged users will be able to search all your user created
group names. Users will also need the SeeGroup privilege to use
this feature.

*/

var AllowGroupAutocompleteForUnprivileged = false

/*

=back

=head2 Self Service Interface

The Self Service Interface is a view automatically presented to Unprivileged
users who have a password and log into the web UI. The following options
modify the default behavior of the Self Service pages.

=over 4


On the ticket display page, show only correspondence transactions in the
ticket history. This hides all ticket update transactions like status changes,
custom field updates, updates to watchers, etc.


*/

var SelfServiceCorrespondenceOnly = false

/*

This determines if we should hide Time Worked, Time Estimated, and
Time Left for unprivileged users.
Set this to true to hide those fields.

*/

var HideTimeFieldsFromUnprivilegedUsers = false

/*

Should unprivileged users (users of SelfService) be allowed to
Setting this option to true means unprivileged users
will be able to search all your users.

*/

var AllowUserAutocompleteForUnprivileged = false

/*

$DefaultSelfServiceSearchResultFormat is the default format of
searches displayed in the SelfService interface.

*/

//var $DefaultSelfServiceSearchResultFormat = qq{
//   '<B><A HREF="__WebPath__/SelfService/Display.html?id=__id__">__id__</a></B>/TITLE:#',
//   '<B><A HREF="__WebPath__/SelfService/Display.html?id=__id__">__Subject__</a></B>/TITLE:Subject',
//   Status,
//   Requestors,
//   Owner});

/*

What portion of RT's URLspace should be accessible to Unprivileged
users This does not override the redirect from F</Ticket/Display.html>
to F</SelfService/Display.html> when Unprivileged users attempt to
access ticked displays.

*/

// var $SelfServiceRegex = qr!^(?:/+SelfService/)!x

/*

This option controls how the SelfService user preferences page is
displayed. It accepts a string from one of the four possible modes
below.

=over

When set to edit-prefs, self service users will be able to update
their Timezone and Language preference and update their password.
This is the default behavior of RT.


When set to view-info, users will have full access to all their
user information stored in RT on a read-only page.


When set to edit-prefs-view-info, users will have full access as in
the view-info option, but also will be able to update their Locale
and password as in the default edit-prefs option.


When set to full-edit, users will be able to fully view and update
all of their stored RT user information.

=back

*/

var SelfServiceUserPrefs = "edit-prefs"

/*

Set this to the name of the queue to use for tickets requesting updates
to user infomation from Self Service users. Once it's set, a quick
ticket create portlet will show up on Preferences page for self service
users. This option is only available when $SelfServiceUserPrefs is set
to 'view-info' or 'edit-prefs-view-info'.

Self service users need the CreateTicket right on this queue to create
a ticket.

*/

var SelfServiceRequestUpdateQueue = ""

/*

Allow Self Service users to download their user information, ticket data,
and transaction data as a .tsv file. When enabled, these options
will appear in the self service interface at Logged in as > Preferences.
Users also need the ModifySelf right to have access to this page.

*/

var SelfServiceDownloadUserData = false

/*

Set this option to true to show a section with group tickets
on self service pages.

*/

var SelfServiceShowGroupTickets = false

/*

$SelfServiceUseDashboard is a flag indicating whether or not to use
a dashboard for the Self Service home page.  If it is set to false,
then the normal Open Tickets / Closed Tickets menu is shown rather
than a dashboard.

*/

var SelfServiceUseDashboard = false

/*

$SelfServicePageComponents is an arrayref of allowed components on
the SelfService page, if you have set $SelfServiceUseDashboard to true.

*/

// var $SelfServicePageComponents = [qw(SelfServiceTopArticles SelfServiceNewestArticles)]
// );

/*

If enabled, $SelfServiceShowArticleSearch displays a "Search Articles"
box in the menu bar in the self service interface. This option controls
only showing or hiding the search box. Users still need appropriate rights
to see article search results and view articles.

*/

var SelfServiceShowArticleSearch = false

/*

=back

=head2 Articles

=over 4

Set this to true to display the Articles interface on the Ticket Create
page in addition to the Reply/Comment page.

*/

var ArticleOnTicketCreate = false

/*

Set this to true to hide "Include Article" box on the ticket update page.

*/

var HideArticleSearchOnReplyCreate = false

/*

Set this to false to suppress the default behavior of automatically linking
to Articles when they are included in a message.

*/

var LinkArticlesOnInclude = true

/*

=back

=head2 Assets

=over 4

This should be a list of names of queues whose tickets should always
display the "Assets" box.  This is useful for queues which deal
primarily with assets, as it provides a ready box to link an asset to
the ticket, even when the ticket has no related assets yet.

*/

// var @AssetQueues = ()

/*

This provides the default catalog after a user initially logs in.
However, the default catalog is "sticky," and so will remember the
last-selected catalog thereafter.

*/

var DefaultCatalog = "General assets"

/*

Specifies which fields of L<RT::Asset> to match against and how to match
each field when performing a quick search on assets.  Valid match
methods are LIKE, STARTSWITH, ENDSWITH, =, and !=.  Valid search fields
are id, Name, Description, or custom fields, which are specified as
"CF.1234" or "CF.Name"

*/

//var $AssetSearchFields = {
//    id          => '=',
//    Name        => 'LIKE',
//    Description => 'LIKE',
//})

/*

The format that results of the asset search are displayed with.

 loc('Related tickets')
*/

//var $AssetDefaultSearchResultFormat = q[
//    '<a href="__WebPath__/Asset/Display.html?id=__id__">__id__</a>/TITLE:#',
//    '<a href="__WebHomePath__/Asset/Display.html?id=__id__">__Name__</a>/TITLE:Name',
//    Status,
//    Catalog,
//    Owner,
//    '__ActiveTickets__ __InactiveTickets__/TITLE:Related tickets',
//    '__NEWLINE__',
//    '__NBSP__',
//    '<small>__Description__</small>',
//    '<small>__CreatedRelative__</small>',
//    '<small>__LastUpdatedRelative__</small>',
//    '<small>__Contacts__</small>',
//])

/*

What column should we order by for RT asset search results.

*/

var AssetDefaultSearchResultOrderBy = "Name"

/*

When ordering asset search results by $AssetDefaultSearchResultOrderBy,
should the sort be ascending (ASC) or descending (DESC).

*/

var AssetDefaultSearchResultOrder = "ASC"

/*

Display search result count on asset lists. Defaults to true (show them).

*/

var AssetShowSearchResultCount = true

/*

The format that results of the asset simple search are displayed with.  This
is either a string, which will be used for all catalogs, or a hash
reference, keyed by catalog's name/id.  If a hashref and neither name or id
is found therein, falls back to the key "".

If you wish to use the multiple catalog format, your configuration would look
something like:

*/

//var $AssetSimpleSearchFormat = {
//        'General assets' => q[Format String for the General Assets Catalog],
//        8                => q[Format String for Catalog 8],
//        ""               => q[Format String for any catalogs not listed explicitly],
//    });

/*

 loc('Related tickets')

*/

//var $AssetSimpleSearchFormat = q[
//    '<a href="__WebPath__/Asset/Display.html?id=__id__">__id__</a>/TITLE:#',
//    '<a href="__WebHomePath__/Asset/Display.html?id=__id__">__Name__</a>/TITLE:Name',
//    Status,
//    Catalog,
//    Owner,
//    '__ActiveTickets__ __InactiveTickets__/TITLE:Related tickets',
//    '__NEWLINE__',
//    '__NBSP__',
//    '<small>__Description__</small>',
//    '<small>__CreatedRelative__</small>',
//    '<small>__LastUpdatedRelative__</small>',
//    '<small>__Contacts__</small>',
//]);

/*

The information that is displayed on ticket display pages about assets
related to the ticket.  This is displayed in a table beneath the asset
name.

*/

//var $AssetSummaryFormat = q[
//    '<a href="__WebHomePath__/Asset/Display.html?id=__id__">__Name__</a>/TITLE:Name',
//    Description,
//    '__Status__ (__Catalog__)/TITLE:Status',
//    Owner,
//    HeldBy,
//    Contacts,
//    '__ActiveTickets__ __InactiveTickets__/TITLE:Related tickets',
//]);

/*

The information that is displayed on ticket display pages about tickets
related to assets related to the ticket.  This is displayed as a list of
tickets underneath the asset properties.

*/

//var $AssetSummaryRelatedTicketsFormat = q[
//    '<a href="__WebPath__/Ticket/Display.html?id=__id__">__id__</a>',
//    '(__OwnerName__)',
//    '<a href="__WebPath__/Ticket/Display.html?id=__id__">__Subject__</a>',
//    QueueName,
//    Status,
//]);

/*

Specify a list of Asset custom fields to show in "Basics" widget on create.

e.g.

*/

// var $AssetBasicCustomFieldsOnCreate = [ 'foo', 'bar' ]

// var $AssetBasicCustomFieldsOnCreate = undef

/*


Set to a true value to hide the legacy Asset Simple Search in favor of AssetSQL
added in RT 5.0.

When hidden, the Asset search menu shows the Current Search menu like tickets,
giving quick access back to a search after clicking on an asset.


*/

var AssetHideSimpleSearch = false

/*


By default an asset is limited to a single user as an owner. By setting
this to a true value, you can allow multiple users and groups as owner.
If you change this back to a false value while having multiple owners
set on any assets, RT's behavior may be inconsistent.


*/

var AssetMultipleOwner = false

/*

By default the People portlet on the asset display page shows the Name field
Set this value with additional fields from the user record to
show more information about the users. The value is a Format string of user
attributes or custom fields. For example, to show the user's email, city, state,
zip, and the user custom field Primary Office:

*/

var UserAssetExtraInfo = `EmailAddress, City, State, Zip, "CF.{Primary Office}"`

// var $UserAssetExtraInfo = ""

/*

=back

=head2 Message box properties

=over 4


For message boxes, set the entry box width, height and what type of
wrapping to use.  These options can be overridden by users in their
preferences.

When the width is set to undef, no column count is specified and the
message box will take up 100% of the available width.  Combining this
with HARD messagebox wrapping (below) is not recommended, as it will
lead to inconsistent width in transactions between browsers.

These settings only apply to the non-RichText message box.  See below
for Rich Text settings.

*/

var MessageBoxWidth = 0

var MessageBoxHeight = 15

/*

Should "rich text" editing be enabled? This option lets your users
send HTML email messages from the web interface.

*/

var MessageBoxRichText = true

/*

Height of rich text JavaScript enabled editing boxes (in pixels)

*/

var MessageBoxRichTextHeight = 300

/*

Should your users' signatures (from their Preferences page) be
included in Comments and Replies.

*/

var MessageBoxIncludeSignature = true

/*

Should your users' signatures (from their Preferences page) be
Setting this to false overrides
$MessageBoxIncludeSignature.

*/

var MessageBoxIncludeSignatureOnComment = true

/*

By default RT places the signature at the bottom of the quoted text in
Set this to true to place the signature
above the quoted text.

*/

var SignatureAboveQuote = false

/*

=back

=head2 Attach Files

=over 4


By default, RT uses Dropzone to attach files if possible. If
$PreferDropzone is set to false, RT will always use plain file inputs.

*/

var PreferDropzone = true

/*

=back

=head2 Transaction search

=over 4


%TransactionDefaultSearchResultFormat is the default format for RT
transaction search results for various object types. Keys are object types
like RT::Ticket, values are the format string.

*/

//var %TransactionDefaultSearchResultFormat,
//    'RT::Ticket' => qq{
//        '<B><A HREF="__WebPath__/Transaction/Display.html?id=__id__">__id__</a></B>/TITLE:ID',
//        '<B><A HREF="__WebPath__/Ticket/Display.html?id=__ObjectId__">__ObjectId__</a></B>/TITLE:Ticket',
//        '__Description__',
//        '<small>__OldValue__</small>',
//        '<small>__NewValue__</small>',
//        '<small>__Content__</small>',
//        '<small>__CreatedRelative__</small>',
//   },
//);

/*

What Transactions column should we order by for RT Transaction search
results for various object types.  Keys are object types like RT::Ticket,
values are the column names.

Defaults to I<id>.


*/

// var %TransactionDefaultSearchResultOrderBy = 'RT::Ticket' => 'id'

/*

When ordering RT Transaction search results by
%TransactionDefaultSearchResultOrderBy, should the sort be ascending
(ASC) or descending (DESC).  Keys are object types like RT::Ticket,
values are either "ASC" or "DESC".

Defaults to I<ASC>.

*/

// var %TransactionDefaultSearchResultOrder = 'RT::Ticket' => 'ASC'

/*

Display search result count on transaction lists.  Keys are object types
like RT::Ticket, values are either 1 or 0.

Defaults to true (show them).

*/

// var %TransactionShowSearchResultCount = 'RT::Ticket' => 1

/*

=back

=head2 Transaction display

=over 4


By default, RT shows newest transactions at the bottom of the ticket
history page, if you want see them at the top set this to false.  This
option can be overridden by users in their preferences.

*/

var OldestTransactionsFirst = true

/*

This option controls how history is shown on the ticket display page.  It
accepts one of three possible modes and is overrideable on a per-user
preference level.  If you regularly deal with long tickets and don't care much
about the history, you may wish to change this option to click.

=over

When set to delay, history is loaded via javascript after the rest of the
page has been loaded.  This speeds up apparent page load times and generally
provides a smoother experience.  You may notice slight delays before the ticket
history appears on very long tickets.

When set to click, history is loaded on demand when a placeholder link is
clicked.  This speeds up ticket display page loads and history is never loaded
if not requested.

When set to always, history is loaded before showing the page.  This ensures
history is always available immediately, but at the expense of longer page load
times.  This behaviour was the default in RT 4.0.

When set to scroll, history is loaded via javascript after the rest of the
page has been loaded, as you scroll down the page. Ten transactions are loaded
initially, and then more are loaded ten at a time. This can dramatically speed
up initial page load times on very long tickets.

=back

*/

var ShowHistory = "delay"

/*

By default, RT hides from the web UI information about blind copies
user sent on reply or comment.

*/

var ShowBccHeader = false

/*

If TrustHTMLAttachments is not defined, we will display them as
text. This prevents malicious HTML and JavaScript from being sent in a
request (although there is probably more to it than that)

*/

var TrustHTMLAttachments = false

/*

Always download attachments, regardless of content type. If set, this
overrides TrustHTMLAttachments.

*/

var AlwaysDownloadAttachments = false

/*

Sets the number of attachments to display on ticket display and ticket
update pages (default is 5). Attachments beyond this number are displayed
Set to undef to always show
all attachments. A value of 0 means show no attachments by default.

*/

var AttachmentListCount = 5

/*

By default, RT shows rich text (HTML) messages if possible.  If
$PreferRichText is set to false, RT will show plain text messages in
preference to any rich text alternatives.

As a security precaution, RT limits the HTML that is displayed to a
known-good subset -- as allowing arbitrary HTML to be displayed exposes
multiple vectors for XSS and phishing attacks.  If
L</$TrustHTMLAttachments> is enabled, the original HTML is available for
viewing via the "Download" link.

*/

var PreferRichText = true

/*

$MaxInlineBody is the maximum textual attachment size that we want to
see inline when viewing a transaction.  RT will inline any text if the
value is undefined or 0.  This option can be overridden by users in
their preferences.  The default is 25k.

*/

var MaxInlineBody = 25 * 1024

/*

By default, RT shows images attached to incoming (and outgoing) ticket
Set this variable to false if you'd like to disable that
behavior.

*/

var ShowTransactionImages = true

/*

By default, RT doesn't show remote images attached to incoming (and outgoing)
Set this variable to true if you'd like to enable remote
image display.  Showing remote images may allow spammers and other senders to
track when messages are viewed and see referer information.

Note that this setting is independent of L</$ShowTransactionImages> above.

*/

var ShowRemoteImages = false

/*

Normally plaintext attachments are displayed as HTML with line breaks
preserved.  This causes space- and tab-based formatting not to be
Set $PlainTextMono to true to use a monospaced
font and preserve formatting.

*/

var PlainTextMono = false

/*

If $SuppressInlineTextFiles is set to true, then uploaded text files
(text-type attachments with file names) are prevented from being
displayed in-line when viewing a ticket's history.

*/

var SuppressInlineTextFiles = false

/*

MakeClicky detects various formats of data in headers and email
messages, and extends them with supporting links.  By default, RT
provides two formats:

* 'httpurl': detects http:// and https:// URLs and adds '[Open URL]'
  link after the URL.

* 'httpurl_overwrite': also detects URLs as 'httpurl' format, but
  replaces the URL with a link.  Enabled by default.

See F<share/html/Elements/MakeClicky> for documentation on how to add
your own styles of link detection.

*/

// var @Active_MakeClicky = qw(httpurl_overwrite)

/*

Quote folding is the hiding of old replies in transaction history.
Set this to false to disable it.

*/

var QuoteFolding = true

/*

$QuoteWrapWidth controls the number of columns to use when wrapping
quoted text within transactions.

*/

var QuoteWrapWidth = 70

/*

=back


=head2 Administrative interface

=over 4


RT can show administrators a feed of recent RT releases and other
related announcements and information from Best Practical on the top
level Admin page.  This feature helps you stay up to date on
RT security announcements and version updates.

RT provides this feature using an "iframe" on /Admin/index.html
which asks the administrator's browser to show an inline page from
Best Practical's website.

If you'd rather not make this feature available to your
administrators, set $ShowRTPortal to false.

*/

var ShowRTPortal = true

/*

In the admin interface, format strings similar to tickets result
formats are used. Use %AdminSearchResultFormat to define the format
strings used in the admin interface on a per-RT-class basis.

*/

//var %AdminSearchResultFormat,
//    Queues =>
//        q{'<a href="__WebPath__/Admin/Queues/Modify.html?id=__id__">__id__</a>/TITLE:#'}
//        .q{,'<a href="__WebPath__/Admin/Queues/Modify.html?id=__id__">__Name__</a>/TITLE:Name'}
//        .q{,__Description__,__Address__,__Priority__,__DefaultDueIn__,__Lifecycle__,__SubjectTag__,__Disabled__,__SortOrder__},
//
//    Groups =>
//        q{'<a href="__WebPath__/Admin/Groups/Modify.html?id=__id__">__id__</a>/TITLE:#'}
//        .q{,'<a href="__WebPath__/Admin/Groups/Modify.html?id=__id__">__Name__</a>/TITLE:Name'}
//        .q{,'__Description__',__Disabled__},
//
//    Users =>
//        q{'<a href="__WebPath__/Admin/Users/Modify.html?id=__id__">__id__</a>/TITLE:#'}
//        .q{,'<a href="__WebPath__/Admin/Users/Modify.html?id=__id__">__Name__</a>/TITLE:Name'}
//        .q{,__RealName__, __EmailAddress__,__SystemGroup__,__Disabled__},
//
//    CustomFields =>
//        q{'<a href="__WebPath__/Admin/CustomFields/Modify.html?id=__id__">__id__</a>/TITLE:#'}
//        .q{,'<a href="__WebPath__/Admin/CustomFields/Modify.html?id=__id__">__Name__</a>/TITLE:Name'}
//        .q{,__AddedTo__, __EntryHint__, __FriendlyPattern__,__Disabled__},
//
//    CustomRoles =>
//        q{'<a href="__WebPath__/Admin/CustomRoles/Modify.html?id=__id__">__id__</a>/TITLE:#'}
//        .q{,'<a href="__WebPath__/Admin/CustomRoles/Modify.html?id=__id__">__Name__</a>/TITLE:Name'}
//        .q{,__Description__,__MaxValues__,__Disabled__},
//
//    Scrips =>
//        q{'<a href="__WebPath__/Admin/Scrips/Modify.html?id=__id____From__">__id__</a>/TITLE:#'}
//        .q{,'<a href="__WebPath__/Admin/Scrips/Modify.html?id=__id____From__">__Description__</a>/TITLE:Description'}
//        .q{,__Condition__, __Action__, __Template__, __Disabled__},
//
//    Templates =>
//        q{'<a href="__WebPath__/__WebRequestPathDir__/Template.html?Queue=__QueueId__&Template=__id__">__id__</a>/TITLE:#'}
//        .q{,'<a href="__WebPath__/__WebRequestPathDir__/Template.html?Queue=__QueueId__&Template=__id__">__Name__</a>/TITLE:Name'}
//        .q{,'__Description__','__UsedBy__','__IsEmpty__'},
//    Classes =>
//        q{ '<a href="__WebPath__/Admin/Articles/Classes/Modify.html?id=__id__">__id__</a>/TITLE:#'}
//        .q{,'<a href="__WebPath__/Admin/Articles/Classes/Modify.html?id=__id__">__Name__</a>/TITLE:Name'}
//        .q{,__Description__,__Disabled__},
//
//    Catalogs =>
//        q{'<a href="__WebPath__/Admin/Assets/Catalogs/Modify.html?id=__id__">__id__</a>/TITLE:#'}
//        .q{,'<a href="__WebPath__/Admin/Assets/Catalogs/Modify.html?id=__id__">__Name__</a>/TITLE:Name'}
//        .q{,__Description__,__Lifecycle__,__Disabled__},
//);

/*

Use %AdminSearchResultRows to define the search result rows in the admin
interface on a per-RT-class basis.

*/

//var %AdminSearchResultRows,
//    Queues       => 50,
//    Groups       => 50,
//    Users        => 50,
//    CustomFields => 50,
//    CustomRoles  => 50,
//    Scrips       => 50,
//    Templates    => 50,
//    Classes      => 50,
//    Catalogs     => 50,
//    Assets       => 50,
//);

/*

Starting in RT 5.0, SuperUsers can edit RT system configuration via the web UI.
Options set in the web UI take precedence over those set in configuration files.

If you prefer to set configuration only via files, set $ShowEditSystemConfig
to false to disable the web UI editing interface.

*/

var ShowEditSystemConfig = true

/*

Starting in RT 5.0, SuperUsers can edit lifecycle configuration via the web UI.
Options set in the web UI take precedence over those set in configuration files.

Set $ShowEditLifecycleConfig to false to disable the web UI editing interface.

*/

var ShowEditLifecycleConfig = true

/*

=back


=head1 Features

=head2 Syncing Users and Groups with LDAP or AD

In addition to the authentication services described above, RT also
has a utility you can schedule to periodicially sync from your
directory service additional user attributes, new users,
disabled users, and group membership. Options for the
LDAPImport tool are listed here. Additional information is
available in the L<RT::LDAPImport> documentation.

=over 4


Your LDAP server hostname.


The LDAP user to log in with.


LDAP options that are supported by the L<Net::LDAP> new method.

Example:

*/

// var $LDAPOptions = [ port => 636 ]

/*

Password for LDAPUser.

The LDAP search base.

Example:

*/

var LDAPBase = "ou=Organisational Unit,dc=domain,dc=TLD"

/*

The filter to use when querying LDAP for the set of users to sync.

Mapping to apply between LDAP attributes retrieved and RT user
record attributes. See the L<RT::LDAPImport> documentation
for details.

The base for the LDAP group search.

The filter to use when querying LDAP for groups to sync.

Mapping to apply between LDAP group member attributes retrieved and
RT groups. See the L<RT::LDAPImport> documentation
for details.

=back

=head2 Cryptography

A complete description of RT's cryptography capabilities can be found in
L<RT::Crypt>. At this moment, GnuPG (PGP) and SMIME security protocols are
supported.

=over 4


The following options apply to all cryptography protocols.

By default, all enabled security protocols will analyze each incoming
email. You may set Incoming to a subset of this list, if some enabled
protocols do not apply to incoming mail; however, this is usually
unnecessary.

For outgoing emails, the first security protocol from the above list is
used. Use the Outgoing option to set a security protocol that should
be used in outgoing emails.  At this moment, only one protocol can be
used to protect outgoing emails.

Set RejectOnMissingPrivateKey to false if you don't want to reject
emails encrypted for key RT doesn't have and can not decrypt.

Set RejectOnBadData to false if you don't want to reject letters
with incorrect data.

If you want to allow people to encrypt attachments inside the DB then
set AllowEncryptDataInDB to true.

Set Dashboards to a hash with Encrypt and Sign keys to control
whether dashboards should be encrypted and/or signed correspondingly.
By default they are not encrypted or signed.

Similarly, set DigestEmail to a hash with Encrypt and Sign keys
to control whether digest email should be encrypted and/or signed.
By default they are not encrypted or signed.

=back

*/

//var %Crypt,
//    Incoming                  => undef, # ['GnuPG', 'SMIME']
//    Outgoing                  => undef, # 'SMIME'
//
//    RejectOnMissingPrivateKey => 1,
//    RejectOnBadData           => 1,
//
//    AllowEncryptDataInDB      => 0,
//
//    Dashboards => {
//        Encrypt => 0,
//        Sign    => 0,
//    },
//    DigestEmail => {
//        Encrypt => 0,
//        Sign    => 0,
//    },
//);

/*

=head3 SMIME configuration

A full description of the SMIME integration can be found in
L<RT::Crypt::SMIME>.

=over 4


Set Enable to false or true to disable or enable SMIME for
encrypting and signing messages.

Set OpenSSL to path to F<openssl> executable.

Set Keyring to directory with key files.  Key and certificates should
be stored in a PEM file in this directory named named, e.g.,
F<email.address@example.com.pem>.

Set CAPath to either a PEM-formatted certificate of a single signing
certificate authority, or a directory of such (including hash symlinks
as created by the openssl tool c_rehash).  Only SMIME certificates
signed by these certificate authorities will be treated as valid
signatures.  If left unset (and AcceptUntrustedCAs is unset, as it is
by default), no signatures will be marked as valid!

Set AcceptUntrustedCAs to allow arbitrary SMIME certificates, no
matter their signing entities.  Such mails will be marked as untrusted,
but signed; CAPath will be used to mark which mails are signed by
trusted certificate authorities.  This configuration is generally
insecure, as it allows the possibility of accepting forged mail signed
by an untrusted certificate authority.

Setting AcceptUntrustedCAs also allows encryption to users with
certificates created by untrusted CAs.

Set Passphrase to a scalar (to use for all keys), an anonymous
function, or a hash (to look up by address).  If the hash is used, the
"" key is used as a default.

Set OtherCertificatesToSend to path to a PEM-formatted certificate file.
Certificates in the file will be include in outgoing signed emails.

Set CheckCRL to a true value to have RT check for revoked certificates
by downloading a CRL. By default, CheckCRL is disabled.

Set CheckOCSP to a true value to have RT check for revoked certificates
against an OCSP server if possible.  By default, CheckOCSP is disabled.

Set CheckRevocationDownloadTimeout to the timeout in seconds for
downloading a CRL or an issuer certificate (the latter is used when
checking against OCSP).  The default timeout is 30 seconds.

See L<RT::Crypt::SMIME> for details.

=back

*/

//var %SMIME,
//    Enable => 0,
//    OpenSSL => 'openssl',
//    Keyring => q{var/data/smime},
//    CAPath => undef,
//    AcceptUntrustedCAs => undef,
//    Passphrase => undef,
//    OtherCertificatesToSend => undef,
//    CheckCRL => 0,
//    CheckOCSP => 0,
//    CheckRevocationDownloadTimeout => 30,
//);

/*

=head3 GnuPG configuration

A full description of the (somewhat extensive) GnuPG integration can
be found by running the command `perldoc L<RT::Crypt::GnuPG>` (or
`perldoc lib/RT/Crypt/GnuPG.pm` from your RT install directory).

=over 4


Set Enable to false or true to disable or enable GnuPG interfaces
for encrypting and signing outgoing messages.

Set GnuPG to the name or path of the gpg binary to use.

Set Passphrase to a scalar (to use for all keys), an anonymous
function, or a hash (to look up by address).  If the hash is used, the
"" key is used as a default.

Set OutgoingMessagesFormat to 'inline' to use inline encryption and
signatures instead of 'RFC' (GPG/MIME: RFC3156 and RFC1847) format.

Plain (not encrypted) email sent to RT can have attached files that
Set FileExtensions to any file extensions on
attachments that RT should treat as encrypted and attempt to
decrypt with gpg.


*/

//var %GnuPG,
//    Enable                 => 0,
//    GnuPG                  => 'gpg',
//    Passphrase             => undef,
//    OutgoingMessagesFormat => "RFC", # Inline
//    FileExtensions         => [ 'pgp', 'gpg', 'asc' ],
//);

/*

Options to pass to the GnuPG program.

If you override this in your RT_SiteConfig, you should be sure to
include a homedir setting.

Note that options with '-' character MUST be quoted.

*/

//var %GnuPGOptions,
//    homedir => q{var/data/gpg},
//
//# URL of a keyserver
//    keyserver => 'hkp://subkeys.pgp.net',
//
//# enables the automatic retrieving of keys when verifying signatures
//    'keyserver-options' => 'auto-key-retrieve',
//)

/*

=back

=head2 External storage

By default, RT stores attachments in the database.  ExternalStorage moves
all attachments that RT does not need efficient access to (which include
textual content and images) to outside of the database.  This may either
be on local disk, or to a cloud storage solution.  This decreases the
size of RT's database, in turn decreasing the burden of backing up RT's
database, at the cost of adding additional locations which must be
configured or backed up.  Attachment storage paths are calculated based
on file contents; this provides de-duplication.

A full description of external storage can be found by running the command
`perldoc L<RT::ExternalStorage>` (or `perldoc lib/RT/ExternalStorage.pm`
from your RT install directory).

Note that simply configuring L<RT::ExternalStorage> is insufficient; there
are additional steps required (including setup of a regularly-scheduled
upload job) to enable RT to make use of external storage.

=over 4


This selects which storage engine is used, as well as options for
configuring that specific storage engine. RT ships with the following
storage engines:

L<RT::ExternalStorage::Disk>, which stores files on directly onto disk.

L<RT::ExternalStorage::AmazonS3>, which stores files on Amazon's S3 service.

L<RT::ExternalStorage::Dropbox>, which stores files in Dropbox.

See each storage engine's documentation for the configuration it requires
and accepts.

*/

//var %ExternalStorage,
//        Type => 'Disk',
//        Path => '/opt/rt5/var/attachments',
//    );
//
//var %ExternalStorage = ()

/*

Certain object types, like values for Binary (aka file upload) custom
fields, are always put into external storage. However, for other
object types, like images and text, there is a line in the sand where
you want small objects in the database but large objects in external
storage. By default, objects larger than 10 MiB will be put into external
storage. $ExternalStorageCutoffSize adjusts that line in the sand.

Note that changing this setting does not affect existing attachments, only
the new ones that sbin/rt-externalize-attachments hasn't seen yet.

*/

var ExternalStorageCutoffSize = 10 * 1024 * 1024

/*

Certain ExternalStorage backends can serve files over HTTP.  For such
backends, RT can link directly to those files in external storage.  This
cuts down download time and relieves resource pressure because RT's web
server is no longer involved in retrieving and then immediately serving
each attachment.

Of the storage engines that RT ships, only
L<RT::ExternalStorage::AmazonS3> supports this feature, and you must
manually configure it to allow direct linking.

Set this to true to link directly to files in external storage.

*/

var ExternalStorageDirectLink = false

/*

=back


=head2 Lifecycles

=head3 Lifecycle definitions

Each lifecycle is a list of possible statuses split into three logic
sets: B<initial>, B<active> and B<inactive>. Each status in a
lifecycle must be unique. (Statuses may not be repeated across sets.)
Each set may have any number of statuses.

For example:

    default => {
        initial  => ['new'],
        active   => ['open', 'stalled'],
        inactive => ['resolved', 'rejected', 'deleted'],
        ...
    },

Status names can be from 1 to 64 ASCII characters.  Statuses are
localized using RT's standard internationalization and localization
system.

=over 4


You can define multiple B<initial> statuses for tickets in a given
lifecycle.

RT will automatically set its B<Started> date when you change a
ticket's status from an B<initial> state to an B<active> or
B<inactive> status.


B<Active> tickets are "currently in play" - they're things that are
being worked on and not yet complete.


B<Inactive> tickets are typically in their "final resting state".

While you're free to implement a workflow that ignores that
description, typically once a ticket enters an inactive state, it will
never again enter an active state.

RT will automatically set the B<Resolved> date when a ticket's status
is changed from an B<Initial> or B<Active> status to an B<Inactive>
status.

B<deleted> is still a special status and protected by the
B<DeleteTicket> right, unless you re-defined rights (read below). If
you don't want to allow ticket deletion at any time simply don't
include it in your lifecycle.

=back

Statuses in each set are ordered and listed in the UI in the defined
order.

Changes between statuses are constrained by transition rules, as
described below.

=head3 Default values

In some cases a default value is used to display in UI or in API when
value is not provided. You can configure defaults using the following
syntax:

    default => {
        ...
        defaults => {
            on_create => 'new',
            on_resolve => 'resolved',
            ...
        },
    },

The following defaults are used.

=over 4


If you (or your code) doesn't specify a status when creating a ticket,
RT will use the this status. See also L</Statuses available during
ticket creation>.


When an approval is accepted, the status of depending tickets will
be changed to this value.


When an approval is denied, the status of depending tickets will
be changed to this value.


When a reminder is opened, the status will be changed to this value.


When a reminder is resolved, the status will be changed to this value.

=back

=head3 Transitions between statuses and UI actions

A B<Transition> is a change of status from A to B. You should define
all possible transitions in each lifecycle using the following format:

    default => {
        ...
        transitions => {
            ""       => [qw(new open resolved)],
            new      => [qw(open resolved rejected deleted)],
            open     => [qw(stalled resolved rejected deleted)],
            stalled  => [qw(open)],
            resolved => [qw(open)],
            rejected => [qw(open)],
            deleted  => [qw(open)],
        },
        ...
    },

The order of items in the listing for each transition line affects
the order they appear in the drop-down. If you change the config
for 'open' state listing to:

    open     => [qw(stalled rejected deleted resolved)],

then the 'resolved' status will appear as the last item in the drop-down.

=head4 Statuses available during ticket creation

By default users can create tickets with a status of new,
open, or resolved, but cannot create tickets with a status of
rejected, stalled, or deleted. If you want to change the statuses
available during creation, update the transition from "" (empty
string), like in the example above.

=head4 Protecting status changes with rights

A transition or group of transitions can be protected by a specific
right.  Additionally, you can name new right names, which will be added
to the system to control that transition.  For example, if you wished to
create a lesser right than ModifyTicket for rejecting tickets, you could
write:

    default => {
        ...
        rights => {
            '* -> deleted'  => 'DeleteTicket',
            '* -> rejected' => 'RejectTicket',
            '* -> *'        => 'ModifyTicket',
        },
        ...
    },

This would create a new RejectTicket right in the system which you
could assign to whatever groups you choose.

On the left hand side you can have the following variants:

    '<from> -> <to>'
    '* -> <to>'
    '<from> -> *'
    '* -> *'

Valid transitions are listed in order of priority. If a user attempts
to change a ticket's status from B<new> to B<open> then the lifecycle
is checked for presence of an exact match, then for 'any to B<open>',
'B<new> to any' and finally 'any to any'.

If you don't define any rights, or there is no match for a transition,
RT will use the B<DeleteTicket> or B<ModifyTicket> as appropriate.

=head4 Labeling and defining actions

For each transition you can define an action that will be shown in the
UI; each action annotated with a label and an update type.

Each action may provide a default update type, which can be
B<Comment>, B<Respond>, or absent. For example, you may want your
staff to write a reply to the end user when they change status from
B<new> to B<open>, and thus set the update to B<Respond>.  Neither
B<Comment> nor B<Respond> are mandatory, and user may leave the
message empty, regardless of the update type.

This configuration can be used to accomplish what
$ResolveDefaultUpdateType was used for in RT 3.8.

Use the following format to define labels and actions of transitions:

    default => {
        ...
        actions => [
            'new -> open'     => { label => 'Open it', update => 'Respond' },
            'new -> resolved' => { label => 'Resolve', update => 'Comment' },
            'new -> rejected' => { label => 'Reject',  update => 'Respond' },
            'new -> deleted'  => { label => 'Delete' },

            'open -> stalled'  => { label => 'Stall',   update => 'Comment' },
            'open -> resolved' => { label => 'Resolve', update => 'Comment' },
            'open -> rejected' => { label => 'Reject',  update => 'Respond' },

            'stalled -> open'  => { label => 'Open it' },
            'resolved -> open' => { label => 'Re-open', update => 'Comment' },
            'rejected -> open' => { label => 'Re-open', update => 'Comment' },
            'deleted -> open'  => { label => 'Undelete' },
        ],
        ...
    },

In addition, you may define multiple actions for the same transition.
Alternately, you may use '* -> x' to match more than one transition.
For example:

    default => {
        ...
        actions => [
            ...
            'new -> rejected' => { label => 'Reject', update => 'Respond' },
            'new -> rejected' => { label => 'Quick Reject' },
            ...
            '* -> deleted' => { label => 'Delete' },
            ...
        ],
        ...
    },

=head3 Moving tickets between queues with different lifecycles

Unless there is an explicit mapping between statuses in two different
lifecycles, you can not move tickets between queues with these
lifecycles -- even if both use the exact same set of statuses.
Such a mapping is defined as follows:

    __maps__ => {
        'from lifecycle -> to lifecycle' => {
            'status in left lifecycle' => 'status in right lifecycle',
            ...
        },
        ...
    },

*/

//var %Lifecycles,
//    default => {
//        initial         => [qw(new)], # loc_qw
//        active          => [qw(open stalled)], # loc_qw
//        inactive        => [qw(resolved rejected deleted)], # loc_qw
//
//        defaults => {
//            on_create => 'new',
//            approved  => 'open',
//            denied    => 'rejected',
//            reminder_on_open     => 'open',
//            reminder_on_resolve  => 'resolved',
//        },
//
//        transitions => {
//            ""       => [qw(new open resolved)],
//
//            # from   => [ to list ],
//            new      => [qw(    open stalled resolved rejected deleted)],
//            open     => [qw(new      stalled resolved rejected deleted)],
//            stalled  => [qw(new open         resolved rejected deleted)],
//            resolved => [qw(new open stalled          rejected deleted)],
//            rejected => [qw(new open stalled resolved          deleted)],
//            deleted  => [qw(new open stalled resolved rejected        )],
//        },
//        rights => {
//            '* -> deleted'  => 'DeleteTicket',
//            '* -> *'        => 'ModifyTicket',
//        },
//        actions => [
//            'new -> open'      => { label  => 'Open It', update => 'Respond' }, # loc{label}
//            'new -> resolved'  => { label  => 'Resolve', update => 'Comment' }, # loc{label}
//            'new -> rejected'  => { label  => 'Reject',  update => 'Respond' }, # loc{label}
//            'new -> deleted'   => { label  => 'Delete',                      }, # loc{label}
//            'open -> stalled'  => { label  => 'Stall',   update => 'Comment' }, # loc{label}
//            'open -> resolved' => { label  => 'Resolve', update => 'Comment' }, # loc{label}
//            'open -> rejected' => { label  => 'Reject',  update => 'Respond' }, # loc{label}
//            'stalled -> open'  => { label  => 'Open It',                     }, # loc{label}
//            'resolved -> open' => { label  => 'Re-open', update => 'Comment' }, # loc{label}
//            'rejected -> open' => { label  => 'Re-open', update => 'Comment' }, # loc{label}
//            'deleted -> open'  => { label  => 'Undelete',                    }, # loc{label}
//        ],
//    },
//    assets => {
//        type     => "asset",
//        initial  => [
//            'new' # loc
//        ],
//        active   => [
//            'allocated', # loc
//            'in-use' # loc
//        ],
//        inactive => [
//            'recycled', # loc
//            'stolen', # loc
//            'deleted' # loc
//        ],
//
//        defaults => {
//            on_create => 'new',
//        },
//
//        transitions => {
//            ""        => [qw(new allocated in-use)],
//            new       => [qw(allocated in-use stolen deleted)],
//            allocated => [qw(in-use recycled stolen deleted)],
//            "in-use"  => [qw(allocated recycled stolen deleted)],
//            recycled  => [qw(allocated)],
//            stolen    => [qw(allocated)],
//            deleted   => [qw(allocated)],
//        },
//        rights => {
//            '* -> *'        => 'ModifyAsset',
//        },
//        actions => {
//            '* -> allocated' => {
//                label => "Allocate" # loc
//            },
//            '* -> in-use'    => {
//                label => "Now in-use" # loc
//            },
//            '* -> recycled'  => {
//                label => "Recycle" # loc
//            },
//            '* -> stolen'    => {
//                label => "Report stolen" # loc
//            },
//        },
//    },

/*
 don't change lifecyle of the approvals, they are not capable to deal with
 custom statuses
*/

//    approvals => {
//        initial         => [ 'new' ],
//        active          => [ 'open', 'stalled' ],
//        inactive        => [ 'resolved', 'rejected', 'deleted' ],
//
//        defaults => {
//            on_create => 'new',
//            reminder_on_open     => 'open',
//            reminder_on_resolve  => 'resolved',
//        },
//
//        transitions => {
//            ""       => [qw(new open resolved)],
//
//            # from   => [ to list ],
//            new      => [qw(open stalled resolved rejected deleted)],
//            open     => [qw(new stalled resolved rejected deleted)],
//            stalled  => [qw(new open rejected resolved deleted)],
//            resolved => [qw(new open stalled rejected deleted)],
//            rejected => [qw(new open stalled resolved deleted)],
//            deleted  => [qw(new open stalled rejected resolved)],
//        },
//        rights => {
//            '* -> deleted'  => 'DeleteTicket',
//            '* -> rejected' => 'ModifyTicket',
//            '* -> *'        => 'ModifyTicket',
//        },
//        actions => [
//            'new -> open'      => { label  => 'Open It', update => 'Respond' }, # loc{label}
//            'new -> resolved'  => { label  => 'Resolve', update => 'Comment' }, # loc{label}
//            'new -> rejected'  => { label  => 'Reject',  update => 'Respond' }, # loc{label}
//            'new -> deleted'   => { label  => 'Delete',                      }, # loc{label}
//            'open -> stalled'  => { label  => 'Stall',   update => 'Comment' }, # loc{label}
//            'open -> resolved' => { label  => 'Resolve', update => 'Comment' }, # loc{label}
//            'open -> rejected' => { label  => 'Reject',  update => 'Respond' }, # loc{label}
//            'stalled -> open'  => { label  => 'Open It',                     }, # loc{label}
//            'resolved -> open' => { label  => 'Re-open', update => 'Comment' }, # loc{label}
//            'rejected -> open' => { label  => 'Re-open', update => 'Comment' }, # loc{label}
//            'deleted -> open'  => { label  => 'Undelete',                    }, # loc{label}
//        ],
//    },
//);

/*

=head2 SLA

=over 4

*/

//var %ServiceAgreements = (
//        Default => '4h',
//        QueueDefault => {
//            'Incident' => '2h',
//        },
//        Levels => {
//            '2h' => { Resolve => { RealMinutes => 60*2 } },
//            '4h' => { Resolve => { RealMinutes => 60*4 } },
//        },
//    ));

/*

In this example I<Incident> is the name of the queue, and I<2h> is the name of
the SLA which will be applied to this queue by default.

Each service level can be described using several options:
L<Starts|/"Starts (interval, first business minute)">,
L<Resolve|/"Resolve and Response (interval, no defaults)">,
L<Response|/"Resolve and Response (interval, no defaults)">,
L<KeepInLoop|/"Keep in loop (interval, no defaults)">,
L<OutOfHours|/"OutOfHours (struct, no default)">
and L<ServiceBusinessHours|/"Configuring business hours">.

=over 4


By default when a ticket is created Starts date is set to
first business minute after time of creation. In other
words if a ticket is created during business hours then
Starts will be equal to Created time, otherwise Starts will
be beginning of the next business day.

However, if you provide 24/7 support then you most
probably would be interested in Starts to be always equal
to Created time.

Starts option can be used to adjust behaviour. Format
of the option is the same as format for deadlines which
described later in details. RealMinutes, BusinessMinutes
options and OutOfHours modifiers can be used here like
for any other deadline. For example:

    'standard' => {
        # give people 15 minutes
        Starts   => { BusinessMinutes => 15  },
    },

You can still use old option StartImmediately to set
Starts date equal to Created date.

Example:

    '24/7' => {
        StartImmediately => 1,
        Response => { RealMinutes => 30 },
    },

But it's the same as:

    '24/7' => {
        Starts => { RealMinutes => 0 },
        Response => { RealMinutes => 30 },
    },


These two options define deadlines for resolve of a ticket
and reply to customer(requestors) questions accordingly.

You can define them using real time, business or both. Read more
about the latter L<below|/"Using both Resolve and Response in the same level">.

The Due date field is used to store calculated deadlines.

=over 4


Defines deadline when a ticket should be resolved. This option is
quite simple and straightforward when used without L</Response>.

Example:

    # 8 business hours
    'simple' => { Resolve => 60*8 },
    ...
    # one real week
    'hard' => { Resolve => { RealMinutes => 60*24*7 } },


In many companies providing support service(s) resolve time of a ticket
is less important than time of response to requestors from staff
members.

You can use Response option to define such deadlines.  The Due date is
set when a ticket is created, unset when a worker replies, and re-set
when the requestor replies again -- until the ticket is closed, when the
ticket's Due date is unset.

B<NOTE> that this behaviour changes when Resolve and Response options
are combined; see L</"Using both Resolve and Response in the same
level">.

Note that by default, only the requestors on the ticket are considered
"outside actors" and thus require a Response due date; all other email
addresses are treated as workers of the ticket, and thus count as
meeting the SLA.  If you'd like to invert this logic, so that the Owner
and AdminCcs are the only worker email addresses, and all others are
external, see the L</AssumeOutsideActor> configuration.

The owner is never treated as an outside actor; if they are also the
requestor of the ticket, it will have no SLA.

If an outside actor replies multiple times, their later replies are
ignored; the deadline is always calculated from the oldest
correspondence from the outside actor.

Resolve and Response can be combined. In such case due date is set
according to the earliest of two deadlines and never is dropped to
'not set'.

If a ticket met its Resolve deadline then due date stops "flipping",
is freezed and the ticket becomes overdue. Before that moment when
an inside actor replies to a ticket, due date is changed to Resolve
Set', as well this happens when a ticket
is closed. So all the time due date is defined.

Example:

    'standard delivery' => {
        Response => { RealMinutes => 60*1  }, # one hour
        Resolve  => { RealMinutes => 60*24 }, # 24 real hours
    },

A client orders goods and due date of the order is set to the next one
hour, you have this hour to process the order and write a reply.
As soon as goods are delivered you resolve tickets and usually meet
Resolve deadline, but if you don't resolve or user replies then most
probably there are problems with delivery of the goods. And if after
a week you keep replying to the client and always meeting one hour
response deadline that doesn't mean the ticket is not over due.
Due date was frozen 24 hours after creation of the order.

It's quite rare situation when people need it, but we've decided
that business is applied first and then real time when deadline
described using both types of time. For example:

    'delivery' => {
        Resolve => { BusinessMinutes => 0, RealMinutes => 60*8 },
    },
    'fast delivery' {
        StartImmediately => 1,
        Resolve => { RealMinutes => 60*8 },
    },

For delivery requests which come into the system during business
hours these levels define the same deadlines, otherwise the first
level set deadline to 8 real hours starting from the next business
day, when tickets with the second level should be resolved in the
next 8 hours after creation.

=back


If response deadline is used then Due date is changed to repsonse
Set" when staff replies to a ticket. In some
cases you want to keep requestors in loop and keed them up to date
every few hours. KeepInLoop option can be used to achieve this.

    'incident' => {
        Response   => { RealMinutes => 60*1  }, # one hour
        KeepInLoop => { RealMinutes => 60*2 }, # two hours
        Resolve    => { RealMinutes => 60*24 }, # 24 real hours
    },

In the above example Due is set to one hour after creation, reply
of a inside actor moves Due date two hours forward, outside actors'
replies move Due date to one hour and resolve deadine is 24 hours.

=over 4


Out of hours modifier. Adds more real or business minutes to resolve
and/or reply options if event happens out of business hours, read also
</"Configuring business hours"> below.

Example:

    'level x' => {
        OutOfHours => { Resolve => { RealMinutes => +60*24 } },
        Resolve    => { RealMinutes => 60*24 },
    },

If a request comes into the system during night then supporters have two
hours, otherwise only one.

    'level x' => {
        OutOfHours => { Response => { BusinessMinutes => +60*2 } },
        Resolve    => { BusinessMinutes => 60 },
    },

Supporters have two additional hours in the morning to deal with bunch
of requests that came into the system during the last night.


Allows you to ignore a deadline when ticket has certain status. Example:

    'level x' => {
        KeepInLoop => { BusinessMinutes => 60, IgnoreOnStatuses => ['stalled'] },
    },

In above example KeepInLoop deadline is ignored if ticket is stalled.

B<NOTE>: When a ticket goes from an ignored status to a normal status, the new
Due date is calculated from the last action (reply, SLA change, etc) which fits
the SLA type (Response, Starts, KeepInLoop, etc).  This means if a ticket in
the above example flips from stalled to open without a reply, the ticket will
probably be overdue.  In most cases this shouldn't be a problem since moving
out of stalled-like statuses is often the result of RT's auto-open on reply
scrip, therefore ensuring there's a new reply to calculate Due from.  The
overall effect is that ignored statuses don't let the Due date drift
arbitrarily, which could wreak havoc on your SLA performance.
ExcludeTimeOnIgnoredStatuses option could get around the "probably be
overdue" issue by excluding the time spent on ignored statuses, e.g.

    'level x' => {
        KeepInLoop => {
            BusinessMinutes => 60,
            ExcludeTimeOnIgnoredStatuses => 1,
            IgnoreOnStatuses => ['stalled'],
        },
    },

=back


In the config you can set per queue defaults, using:

*/

//var %ServiceAgreements = (
//        Default => 'global default level of service',
//        QueueDefault => {
//            'queue name' => 'default value for this queue',
//            ...
//        },
//        ...
//    ));

/*

When using a L<Response|/"Resolve and Response (interval, no defaults)">
configuration, the due date is unset when anyone who is not a requestor
replies.  If it is common for non-requestors to reply to tickets, and
this should I<not> satisfy the SLA, you may wish to set
AssumeOutsideActor.  This causes the extension to assume that the
Response SLA has only been met when the owner or AdminCc reply.

*/

//Set ( %ServiceAgreements = (
//        AssumeOutsideActor => 1,
//        ...
//    ));

/*

=back

*/

// var %ServiceAgreements =

/*

In the config you can set one or more work schedules, e.g.

*/

//var %ServiceBusinessHours = (
//        'Default' => {
//            ... description ...
//        },
//        'Support' => {
//            ... description ...
//        },
//        'Sales' => {
//            ... description ...
//        },
//    ));

/*

Read more about how to describe a schedule in L<Business::Hours>.

=over 4


Each level supports BusinessHours option to specify your own business
hours.

    'level x' => {
        BusinessHours => 'work just in Monday',
        Resolve    => { BusinessMinutes => 60 },
    },

then L<%ServiceBusinessHours> should have the corresponding definition:

*/

//var %ServiceBusinessHours = (
//        'work just in Monday' => {
//            1 => { Name => 'Monday', Start => '9:00', End => '18:00' },
//        },
//    ));

/*

Default Business Hours setting is in $ServiceBusinessHours{'Default'}.

=back


*/

// var %ServiceBusinessHours =

/*

=back

=head2 Custom Date Ranges

=over 4

This option lets you declare additional date ranges to be calculated
and displayed in search results. Durations between any two core date
fields, as well as date custom fields, are supported. Each custom
date range is added as an additional display column in the query builder.

You can create basic date calculations via the web UI. SuperUsers can
create them in the main System Configuration section. Individual users
can also create date ranges in the Search options section of user
preferences. More complicated configurations, such as those with
custom code, can be added in your RT_SiteConfig.pm file as described
below.

Business hours are also supported in calculations if you have
L<%ServiceBusinessHours> configured.

Set %CustomDateRanges to a nested structure similar to the following:

*/

//var %CustomDateRanges,
//        'RT::Ticket' => {
//            'Resolution Time' => 'Resolved - Created',
//
//            'Downtime' => {
//                value => 'CF.Recovered - CF.{First Alert}',
//                business_time => 1,
//            },
//
//            'Time To Beta' => {
//                value => 'CF.Beta - now',
//
//                format => sub {
//                    my ($duration, $beta, $now, $ticket) = @_;
//                    my $days = int($duration / (24*60*60));
//                    if ($days < 0) {
//                        $ticket->loc('[quant,_1,day,days] ago', -$days);
//                    }
//                    else {
//                        $ticket->loc('in [quant,_1,day,days]', $days);
//                    }
//                },
//            },
//        },
//    );

/*

The first level keys are record types. Each record type's value must be a
hash reference. Each pair in the second-level hash defines a new range. The
key is the range's name (which is displayed to users in the UI), and its
value describes the range and must be either a string or a hashref.

Values that are plain strings simply describe the calculation to be made.

Values that are hashrefs that could include:

=over 4


A string that describes the calculation to be made.

The calculation is expected to be of the format C<"field - field"> where each
field may be:

=over 4


For example, L<RT::Ticket> supports: Created, Starts, Started, LastUpdated,
Told or LastContact, Due and Resolved.


You may use either CF.Name or C<CF.{Longer Name}> syntax.


=back

Custom date range calculations are defined using typical math operators with
a space before and after. Subtraction (-) is currently supported.

If either field and its corresponding fallback field(see blow) are both unset,
then nothing will be displayed for that record (and the format code
reference will not be called).  If you need additional control over how
results are calculated, see L<docs/customizing/search_result_columns.pod>.


When value is not set, from/to will be used to calculate instead.
Technically, C<Resolved - Created"> is equal to:

    { from => 'Created', to => 'Resolved' }


Fallback fields when the main fields are unset, e.g.

    {   from        => 'CF.{First Alert}',
        to          => 'CF.Recovered',
        to_fallback => 'now',
    }

When CF.Recovered is unset, "now" will be used in the calculation.


A boolean value to indicate if it's a business time or not.

When the schedule can't be deducted from corresponding object, the
Default one defined in %ServiceBusinessHours will be used instead.


A code reference that allows customization of how the duration is displayed
to users.  This code reference receives four parameters: the duration (a
number of seconds), the end time (an L<RT::Date> object), the start time
(another L<RT::Date>), and the record itself (which corresponds to the
first-level key; in the example config above, it would be the L<RT::Ticket>
object). The code reference should return the string to be displayed to the
user.

=back

=back

*/
