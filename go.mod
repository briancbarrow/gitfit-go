module github.com/briancbarrow/gitfit-go

go 1.21.3

require (
	github.com/a-h/templ v0.2.432
	github.com/alexedwards/scs/libsqlstore v0.0.0-00010101000000-000000000000
	github.com/alexedwards/scs/v2 v2.7.0
	github.com/go-playground/form/v4 v4.2.1
	github.com/joho/godotenv v1.5.1
	github.com/julienschmidt/httprouter v1.3.0
	github.com/justinas/alice v1.2.0
	github.com/justinas/nosurf v1.1.1
	github.com/mattn/go-sqlite3 v1.14.18
	github.com/pressly/goose/v3 v3.17.0
	github.com/stytchauth/stytch-go/v11 v11.5.2
	github.com/tursodatabase/libsql-client-go v0.0.0-20231216154754-8383a53d618f
	golang.org/x/crypto v0.16.0
)

require (
	github.com/MicahParks/keyfunc/v2 v2.0.1 // indirect
	github.com/antlr/antlr4/runtime/Go/antlr/v4 v4.0.0-20230512164433-5d1fd1a340c9 // indirect
	github.com/golang-jwt/jwt/v5 v5.0.0 // indirect
	github.com/klauspost/compress v1.17.2 // indirect
	github.com/libsql/sqlite-antlr4-parser v0.0.0-20230802215326-5cb5bb604475 // indirect
	github.com/mitchellh/mapstructure v1.5.0 // indirect
	github.com/sethvargo/go-retry v0.2.4 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	golang.org/x/exp v0.0.0-20220722155223-a9213eeb770e // indirect
	golang.org/x/sync v0.5.0 // indirect
	nhooyr.io/websocket v1.8.7 // indirect
)

// replace github.com/alexedwards/scs/sqlite3store => ../scs/sqlite3store

replace github.com/alexedwards/scs/libsqlstore => ../scs/libsqlstore
