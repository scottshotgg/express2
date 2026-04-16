module github.com/scottshotgg/express2

go 1.26.1

require (
	github.com/pkg/errors v0.9.1
	github.com/scottshotgg/express-ast v0.0.0-20190816231702-96cd278652f9
	github.com/scottshotgg/express-lex v0.0.0-20200121235042-29d2e6df5787
	github.com/scottshotgg/express-token v0.0.0-20230327011102-da006d30a2eb
	github.com/spf13/cobra v1.10.2
	github.com/spf13/viper v1.21.0
)

require (
	github.com/fsnotify/fsnotify v1.9.0 // indirect
	github.com/go-viper/mapstructure/v2 v2.5.0 // indirect
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/pelletier/go-toml/v2 v2.3.0 // indirect
	github.com/sagikazarmark/locafero v0.12.0 // indirect
	github.com/spf13/afero v1.15.0 // indirect
	github.com/spf13/cast v1.10.0 // indirect
	github.com/spf13/pflag v1.0.10 // indirect
	github.com/subosito/gotenv v1.6.0 // indirect
	go.yaml.in/yaml/v3 v3.0.4 // indirect
	golang.org/x/sys v0.43.0 // indirect
	golang.org/x/text v0.36.0 // indirect
	gopkg.in/check.v1 v1.0.0-20190902080502-41f04d3bba15 // indirect
)

replace (
	github.com/scottshotgg/express-ast => ../express-ast
	github.com/scottshotgg/express-lex => ../express-lex
	github.com/scottshotgg/express-token => ../express-token
)
