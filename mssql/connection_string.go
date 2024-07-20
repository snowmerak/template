package mssql

import (
	"github.com/microsoft/go-mssqldb/msdsn"
	"net/url"
	"strconv"
	"strings"
	"time"
)

const (
	scheme      = "sqlserver://"
	defaultPort = 1433
)

type ConnectionStringBuilder struct {
	host                   string
	port                   int
	username               *string
	password               *string
	instance               *string
	database               *string
	ConnectionTimeout      *time.Duration
	TrustServerCertificate *bool
	encrypt                *string
}

func NewConnectionStringBuilder(host string, port int) *ConnectionStringBuilder {
	if port == 0 {
		port = defaultPort
	}

	return &ConnectionStringBuilder{
		host: host,
		port: port,
	}
}

func (b *ConnectionStringBuilder) WithUserAuth(username string, password string) *ConnectionStringBuilder {
	b.username = &username
	b.password = &password
	return b
}

func (b *ConnectionStringBuilder) WithInstance(instance string) *ConnectionStringBuilder {
	b.instance = &instance
	return b
}

func (b *ConnectionStringBuilder) WithDatabase(database string) *ConnectionStringBuilder {
	b.database = &database
	return b
}

func (b *ConnectionStringBuilder) WithConnectionTimeout(connectionTimeout time.Duration) *ConnectionStringBuilder {
	b.ConnectionTimeout = &connectionTimeout
	return b
}

func (b *ConnectionStringBuilder) WithTrustServerCertificate(trustServerCertificate bool) *ConnectionStringBuilder {
	b.TrustServerCertificate = &trustServerCertificate
	return b
}

func (b *ConnectionStringBuilder) WithEncrypt(encrypt string) *ConnectionStringBuilder {
	b.encrypt = &encrypt
	return b
}

func (b *ConnectionStringBuilder) BuildToURL() string {
	q := url.Values{}
	if b.database != nil {
		q.Add(msdsn.Database, *b.database)
	}
	if b.ConnectionTimeout != nil {
		q.Add(msdsn.ConnectionTimeout, strconv.FormatInt(int64(*b.ConnectionTimeout/time.Second), 10))
	}
	if b.TrustServerCertificate != nil {
		q.Add(msdsn.TrustServerCertificate, strconv.FormatBool(*b.TrustServerCertificate))
	}
	if b.encrypt != nil {
		q.Add(msdsn.Encrypt, *b.encrypt)
	}

	u := url.URL{
		Scheme: scheme,
		Host:   b.host + ":" + strconv.FormatInt(int64(b.port), 10),
	}

	if b.username != nil && b.password != nil {
		u.User = url.UserPassword(*b.username, *b.password)
	}

	if b.instance != nil {
		u.Path = *b.instance
	}

	u.RawQuery = q.Encode()

	return u.String()
}

func (b *ConnectionStringBuilder) BuildToString() string {
	builder := strings.Builder{}
	builder.WriteString("Server=")
	builder.WriteString(b.host)
	builder.WriteString(",")
	builder.WriteString(strconv.FormatInt(int64(b.port), 10))
	if b.instance != nil {
		builder.WriteString("\\")
		builder.WriteString(*b.instance)
	}
	if b.username != nil && b.password != nil {
		builder.WriteString(";User Id=")
		builder.WriteString(*b.username)
		builder.WriteString(";Password=")
		builder.WriteString(*b.password)
	}
	if b.database != nil {
		builder.WriteString(";Database=")
		builder.WriteString(*b.database)
	}
	if b.ConnectionTimeout != nil {
		builder.WriteString(";Connect Timeout=")
		builder.WriteString(strconv.FormatInt(int64(*b.ConnectionTimeout/time.Second), 10))
	}
	if b.TrustServerCertificate != nil {
		builder.WriteString(";Trust Server Certificate=")
		switch *b.TrustServerCertificate {
		case true:
			builder.WriteString("True")
		case false:
			builder.WriteString("False")
		}
	}
	if b.encrypt != nil {
		builder.WriteString(";Encrypt=")
		builder.WriteString(*b.encrypt)
	}

	return builder.String()
}
