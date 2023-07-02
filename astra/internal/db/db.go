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

func InitKeyspaceAndTables(cfg config.Config, session gocqlx.Session) error {

	err := session.ExecStmt(`CREATE KEYSPACE IF NOT EXISTS astra_ks WITH REPLICATION = { 'class' : 'SimpleStrategy', 'replication_factor' : 3};`)
	if err != nil {
		cfg.Logger().Fatal("error create keyspace", zap.Error(err))

		return err
	}

	err = session.ExecStmt(`CREATE TABLE IF NOT EXISTS message_by_conversation_id(
		conversation_id uuid,
		message_id uuid,
		conversation_name varchar,
		from_user_id varchar,
		from_user_number varchar,
		from_user_first_name varchar,
		from_user_last_name varchar,
		from_user_email varchar,
		text_content varchar,
		sent_at timestamp,
		seen_at timestamp,
		created_by varchar,
		created_at timestamp,
		updated_by varchar,
		updated_at timestamp,
		deleted_by varchar,
		deleted_at timestamp,
		PRIMARY KEY (conversation_id, sent_at)
	) WITH CLUSTERING ORDER BY (sent_at DESC)`)
	if err != nil {
		cfg.Logger().Fatal("error create message table", zap.Error(err))

		return err
	}

	return nil
}
