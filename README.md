# Firesword

An apachetop-style project + log parsing. Provides CLI piping, ncurses and SQL queries.

There's still a few bugs to be fixed and features I'm planning to add, but it's already a fairly stable project and very much usable.

_About nurses support:_

Ncurses can be enabled or disabled via '`// +build ignore`' at the top of ncurses.go

Ncurses mode also requires sqlite and readline, so compiling with ncurses breaks cross-compiling portability (which is why I disable it by default)

### Required imports:

    go get github.com/lmorg/apachelogs    # my apache log parsing package
    go get github.com/ActiveState/tail    # tail -f like library

_Additionally required for ncurses support:_

    go get github.com/nsf/termbox-go      # pretty terminal APIs
    go get github.com/shavac/readline     # readline.c support for inputting SQL
    go get github.com/mattn/go-sqlite3    # sqlite engine


### firesward -h

	Usage: firesword -n [-r int] [-t int]     -f str | *
	                 [--fmt str] [--grep str] --stdin | -f str | *
	                 -h | -hf | -hg | -v

	Global preferences:
	-------------------
	  --no-smp       Disable multi-processor support (SMP enabled by default)
	  --no-errors    Surpress error messages, don't fail on unless fatal

	Ncurses interface:
	------------------
	  -n             Start ncurses mode. The real time user interface
	  -r int         Refresh rate in seconds (default is 1 second)
	  --sql str      SQL to start with (default is "SELECT * FROM default_view")

	Command line interface:
	-----------------------
	  --fmt str      Output format (default: "{ip} {uri} {status} {stitle}")
	                     (-hf for field names and how to declare field lengths)
	  --grep str     Filter results (-hg for patterns)
	  --trim-slash   Trims trailing slash from URI (useful for "sort | uniq -c")

	Input streams:
	--------------
	  --stdin        Read from STDIN (not available in ncurses mode)
	  -f str         Read from text stream, equivalent to tail -f (file name as string)
	  *              Read from text / gzip file (multiple files space delimited)

	Help:
	-----
	  -h | -?        Prints this usage guide
	  -hf            Prints format field names
	  -hg            Prints grep pattern guide
	  -v             Prints version number


### firesward -hf

	Field names:
	------------
	  date           Date of request
	  datetime       Date and time of request
	  epoch / unix   EPOCH (UNIX) timestamp (seconds)
	  file           File name of log (not available to --grep)
	  ip             IP address
	  method         HTTP method (eg GET, POST, HEAD, etc)
	  proc           Processing time
	  proto          Protocol (eg HTTP, HTTPS)
	  qs             Query string
	  ref            Referrer
	  sdesc          HTTP status description (not available to --grep)
	  size           HTTP response body size
	  status         HTTP status (--grep as string)
	  stitle         HTTP status title (not available to --grep)
	  time           Time of request
	  ua             User agent
	  uid            User ID (if applicable)
	  uri            URI

	--fmt field lengths are comma-separated printf-style values. eg
	  "{ref,-30}" == Referrer,  30 characters padding, left justified
	  "{file,40}" == File name, 40 characters padding, right justified


### firesward -hg

	The following is the grep format (braces for illustration purposes only):
	--grep '(field name)(operator)(comparison);(field name)(operator)(comparison);'
	eg: firesword --grep 'time>12:00'
	    firesword --grep 'status=500;time<14:35'
	    firesword --grep 'uri{}{foo}{bar}'

	Operators:
	----------
	  <   less than                             (date and numeric fields only)
	  >   greater than                          (date and numeric fields only)
	  ==  equals (= is also valid)              (all fields)
	  !=  does not equal (<> is also valid)     (all fields)
	  =+  contains                              (string fields only)
	  !+  does not contain                      (string fields only)
	  =~  regex matches                         (string fields only)
	  !~  regex does not match                  (string fields only)
	  ~<  round field down to the nearest n     (numeric fields only)
	  ~>  round field up to the nearest n       (numeric fields only)
	  {}  regex substitution: {search}{replace} (string fields only)
	  /   divide                                (numeric fields only)
	  *   multiply                              (numeric fields only)

	Date / Time Formats:
	--------------------
	  date: dd-mm-yyyy
	  time: hh:mm
	  Date and time fields are entered as strings but processed as numeric fields.

	String Formats:
	--------------------
	  Regex operators are case insensitive, all other operators are case sensitive.
