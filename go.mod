module github.com/scottshotgg/express2

require (
	github.com/BurntSushi/toml v0.3.1 // indirect
	github.com/etcd-io/bbolt v1.3.2 // indirect
	github.com/inconshreveable/mousetrap v1.0.0 // indirect
	github.com/mitchellh/go-homedir v1.0.0 // indirect
	github.com/pkg/errors v0.8.1
	github.com/scottshotgg/express-ast v0.0.0-00010101000000-000000000000
	github.com/scottshotgg/express-lex v0.0.0-20190816233540-eb22f05bde7e
	github.com/scottshotgg/express-token v0.0.0-20190816231727-78f862b0ae0d
	github.com/spf13/cobra v0.0.3
	github.com/spf13/viper v1.3.1
	golang.org/x/sys v0.0.0-20190105165716-badf5585203e // indirect
)

replace (
	github.com/scottshotgg/express-ast => ../express-ast
	github.com/scottshotgg/express-lex => ../express-lex
	github.com/scottshotgg/express-token => ../express-token
)
