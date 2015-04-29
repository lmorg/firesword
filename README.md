# firesword
An apachetop-style project + log parsing. Provides CLI piping, ncurses and SQL queries

_Required imports:_

    go get github.com/lmorg/apachelogs    # my apache log parsing package
    go get github.com/nsf/termbox-go      # pretty terminal APIs
    go get github.com/gizak/termui        # required by termbox-go
    go get github.com/shavac/readline     # readline.c support for inputting SQL
    go get github.com/mattn/go-sqlite3    # sqlite engine
    go get github.com/mattn/go-runewidth  # required by go-sqlite3
