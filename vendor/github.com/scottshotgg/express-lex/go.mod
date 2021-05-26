module github.com/scottshotgg/express-lex

go 1.12

require (
	github.com/pkg/errors v0.8.1
	github.com/scottshotgg/express-token v0.0.0-20200121235105-7ab119fd3a82
)


replace (
	github.com/scottshotgg/express-token => ../express-token
)