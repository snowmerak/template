//go:build windows

package postgres

// when you are working on windows
//go:generate docker run --rm -v "${PWD}:/src" -w /src sqlc/sqlc generate
