package cache

import (
	"github.com/google/uuid"
	"github.com/ory/dockertest/v3"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
	"time"
)

var c Cache

func TestMain(m *testing.M) {
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	resource, err := pool.Run("launcher.gcr.io/google/memcached1", "latest", []string{})
	if err != nil {
		log.Fatalf("Could not start resource: %s", err)
	}
	<-time.After(time.Second)
	err = pool.Retry(func() error {
		hostAndPort := resource.GetHostPort("11211/tcp")
		c, err = New(hostAndPort)
		if err == nil {
			<-time.After(time.Second)
		}
		return err
	})
	if err != nil {
		log.Fatalf("Could not connect to database: %s", err)
	}

	code := m.Run()

	if err = pool.Purge(resource); err != nil {
		log.Errorf("Could not purge resource: %s", err)
	}

	err = resource.Expire(1)
	if err != nil {
		log.Error(err)
	}

	os.Exit(code)
}

func TestGetNotFound(t *testing.T) {
	key := uuid.New().String()
	res, err := c.Get(key)
	require.Nil(t, res)
	require.Nil(t, err)
}

func TestGet(t *testing.T) {
	key := uuid.New().String()
	err := c.Set(key, []byte(key), time.Minute)
	require.Nil(t, err)
	res, err := c.Get(key)
	require.Equal(t, key, string(res))
	require.Nil(t, err)
}

func TestExpiration(t *testing.T) {
	key := uuid.New()
	err := c.Set(key.String(), []byte(key.String()), time.Second)
	require.Nil(t, err)
	res, err := c.Get(key.String())
	require.Equal(t, key.String(), string(res))
	require.Nil(t, err)
	<-time.After(time.Second)
	res, err = c.Get(key.String())
	require.Nil(t, res)
	require.Nil(t, err)
}

func TestDelete(t *testing.T) {
	key := uuid.New().String()
	err := c.Set(key, []byte(key), time.Hour)
	require.Nil(t, err)
	res, err := c.Get(key)
	require.Equal(t, key, string(res))
	require.Nil(t, err)
	err = c.Delete(key)
	require.Nil(t, err)
	res, err = c.Get(key)
	require.Nil(t, res)
	require.Nil(t, err)
}
