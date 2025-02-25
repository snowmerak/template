package postgres

import (
	"strconv"
	"strings"
)

const (
	SSLModeDisable    = "disable"
	SSLModeAllow      = "allow"
	SSLModePrefer     = "prefer"
	SSLModeRequire    = "require"
	SSLModeVerifyCA   = "verify-ca"
	SSLModeVerifyFull = "verify-full"
)

type ConnectionStringBuilder struct {
	host     string
	port     int
	user     string
	password string
	database string

	sslMode        string
	connectTimeout int
}

func NewConnectionStringBuilder(host string, port int, user string, password string, database string) *ConnectionStringBuilder {
	return &ConnectionStringBuilder{
		host:     host,
		port:     port,
		user:     user,
		password: password,
		database: database,
	}
}

func (b *ConnectionStringBuilder) SSLMode(sslMode string) *ConnectionStringBuilder {
	b.sslMode = sslMode
	return b
}

func (b *ConnectionStringBuilder) ConnectTimeout(connectTimeout int) *ConnectionStringBuilder {
	b.connectTimeout = connectTimeout
	return b
}

func (b *ConnectionStringBuilder) Build() string {
	builder := strings.Builder{}
	builder.WriteString("postgres://")
	builder.WriteString(b.user)
	builder.WriteString(":")
	builder.WriteString(b.password)
	builder.WriteString("@")
	builder.WriteString(b.host)
	builder.WriteString(":")
	builder.WriteString(strconv.Itoa(b.port))
	builder.WriteString("/")
	builder.WriteString(b.database)
	isFirstArg := true
	needAmpere := false
	if b.sslMode != "" {
		if isFirstArg {
			builder.WriteString("?")
			isFirstArg = false
		}
		if needAmpere {
			builder.WriteString("&")
		}
		builder.WriteString("sslmode=")
		builder.WriteString(b.sslMode)
		needAmpere = true
	}
	if b.connectTimeout > 0 {
		if isFirstArg {
			builder.WriteString("?")
			isFirstArg = false
		}
		if needAmpere {
			builder.WriteString("&")
		}
		builder.WriteString("connect_timeout=")
		builder.WriteString(strconv.Itoa(b.connectTimeout))
	}
	return builder.String()
}
