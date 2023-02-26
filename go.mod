module github.com/scottshotgg/express2

go 1.16

require (
	github.com/BurntSushi/toml v0.3.1 // indirect
	github.com/inconshreveable/mousetrap v1.0.0 // indirect
	github.com/pkg/errors v0.8.1
	github.com/scottshotgg/express-ast v0.0.0-20190816231702-96cd278652f9
	github.com/scottshotgg/express-lex v0.0.0-20200121235042-29d2e6df5787
	github.com/scottshotgg/express-token v0.0.0-20200121235105-7ab119fd3a82
	github.com/spf13/cobra v0.0.3
	github.com/spf13/viper v1.3.1
	golang.org/x/sys v0.0.0-20190105165716-badf5585203e // indirect
)

replace (
	github.com/scottshotgg/express-ast => ../express-ast
	github.com/scottshotgg/express-lex => ../express-lex
	github.com/scottshotgg/express-token => ../express-token
)
