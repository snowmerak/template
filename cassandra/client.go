package cassandra

import (
	"github.com/gocql/gocql"
	"github.com/scylladb/gocqlx/v2"
)

type Client struct {
	clusterConfig gocql.ClusterConfig
}

func NewClient(host ...string) *Client {
	cfg := gocql.NewCluster(host...)

	return &Client{
		clusterConfig: *cfg,
	}
}

func (c *Client) SetKeyspace(keyspace string) {
	c.clusterConfig.Keyspace = keyspace
}

func (c *Client) SetConfig(f func(*gocql.ClusterConfig)) {
	f(&c.clusterConfig)
}

func (c *Client) Connect() (*gocql.Session, error) {
	return c.clusterConfig.CreateSession()
}

func (c *Client) ConnectX() (gocqlx.Session, error) {
	session, err := gocqlx.WrapSession(c.clusterConfig.CreateSession())
	if err != nil {
		return session, err
	}
	return session, nil
}
