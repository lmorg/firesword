# firesword
An apachetop-style project + log parsing. Provides CLI piping, ncurses and SQL queries.

There's still a few bugs to be fixed and features I'm planning to add, but it's already a fairly stable project and very much usable.

_Required imports:_

    go get github.com/lmorg/apachelogs    # my apache log parsing package
    go get github.com/nsf/termbox-go      # pretty terminal APIs
    go get github.com/gizak/termui        # required by termbox-go
    go get github.com/shavac/readline     # readline.c support for inputting SQL
    go get github.com/mattn/go-sqlite3    # sqlite engine
    go get github.com/mattn/go-runewidth  # required by go-sqlite3
