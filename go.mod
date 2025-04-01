module go-caro

go 1.23.3

require (
	github.com/jackc/pgx/v5 v5.7.3
	gopkg.in/telebot.v4 v4.0.0-beta.4
	gopkg.in/yaml.v3 v3.0.1
)

require (
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20240606120523-5a60cdf6a761 // indirect
	github.com/jackc/puddle/v2 v2.2.2 // indirect
	github.com/kr/pretty v0.3.1 // indirect
	github.com/stretchr/testify v1.8.4 // indirect
	golang.org/x/crypto v0.36.0 // indirect
	golang.org/x/sync v0.12.0 // indirect
	golang.org/x/text v0.23.0 // indirect
)

replace gopkg.in/telebot.v4 => github.com/chmzy/telebot v0.0.0-20250330133046-27de83b2157f
