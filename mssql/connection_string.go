package mssql

import (
	"github.com/microsoft/go-mssqldb/msdsn"
	"net/url"
	"strconv"
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
	DialTimeout            *time.Duration
	TrustServerCertificate *bool
	logFlag                *int
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

func (b *ConnectionStringBuilder) WithDialTimeout(dialTimeout time.Duration) *ConnectionStringBuilder {
	b.DialTimeout = &dialTimeout
	return b
}

func (b *ConnectionStringBuilder) WithTrustServerCertificate(trustServerCertificate bool) *ConnectionStringBuilder {
	b.TrustServerCertificate = &trustServerCertificate
	return b
}

func (b *ConnectionStringBuilder) WithLogFlag(logFlag int) *ConnectionStringBuilder {
	b.logFlag = &logFlag
	return b
}

func (b *ConnectionStringBuilder) WithEncrypt(encrypt string) *ConnectionStringBuilder {
	b.encrypt = &encrypt
	return b
}

func (b *ConnectionStringBuilder) Build() string {
	q := url.Values{}
	if b.database != nil {
		q.Add(msdsn.Database, *b.database)
	}
	if b.ConnectionTimeout != nil {
		q.Add(msdsn.ConnectionTimeout, strconv.FormatInt(int64(*b.ConnectionTimeout/time.Second), 10))
	}
	if b.DialTimeout != nil {
		q.Add(msdsn.DialTimeout, strconv.FormatInt(int64(*b.DialTimeout/time.Second), 10))
	}
	if b.TrustServerCertificate != nil {
		q.Add(msdsn.TrustServerCertificate, strconv.FormatBool(*b.TrustServerCertificate))
	}
	if b.logFlag != nil {
		q.Add(msdsn.LogParam, strconv.FormatInt(int64(*b.logFlag), 10))
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
