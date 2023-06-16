package db

import (
	"github.com/andibalo/ramein/astra/internal/config"
	"github.com/gocql/gocql"
	"github.com/scylladb/gocqlx/v2"
	"go.uber.org/zap"
	"time"
)

func InitDB(cfg config.Config) (gocqlx.Session, error) {
	cluster := CreateCluster(gocql.Quorum, cfg.DBKeyspace(), cfg.DBHosts()...)

	session, err := gocqlx.WrapSession(cluster.CreateSession())
	if err != nil {
		cfg.Logger().Fatal("unable to connect to scylla", zap.Error(err))

		return session, err
	}

	cfg.Logger().Info("connected to scylla!")

	return session, nil
}

func CreateCluster(consistency gocql.Consistency, keyspace string, hosts ...string) *gocql.ClusterConfig {
	retryPolicy := &gocql.ExponentialBackoffRetryPolicy{
		Min:        time.Second,
		Max:        10 * time.Second,
		NumRetries: 5,
	}
	cluster := gocql.NewCluster(hosts...)
	cluster.Keyspace = keyspace
	cluster.Timeout = 5 * time.Second
	cluster.RetryPolicy = retryPolicy
	cluster.Consistency = consistency
	cluster.PoolConfig.HostSelectionPolicy = gocql.TokenAwareHostPolicy(gocql.RoundRobinHostPolicy())
	return cluster
}
