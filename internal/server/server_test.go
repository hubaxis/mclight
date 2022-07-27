package server

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/hubaxis/mclight/internal/cache/mocks"
	"github.com/hubaxis/mclight/protocol/mclight"
	m "github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestGet(t *testing.T) {
	c := mocks.NewCache(t)
	key := uuid.New().String()
	c.On("Get", key).Return([]byte(key), nil)
	s := New(c)
	res, err := s.Get(context.Background(), &mclight.GetRequest{Key: key})
	require.Nil(t, err)
	require.Equal(t, key, string(res.Value))
}

func TestGetNotFound(t *testing.T) {
	c := mocks.NewCache(t)
	key := uuid.New().String()
	c.On("Get", key).Return(nil, nil)
	s := New(c)
	res, err := s.Get(context.Background(), &mclight.GetRequest{Key: key})
	require.Nil(t, err)
	require.Nil(t, res.Value)
}

func TestSet(t *testing.T) {
	c := mocks.NewCache(t)
	key := uuid.New().String()
	c.On("Set", key, []byte(key), m.AnythingOfType("Duration")).Return(nil)
	s := New(c)
	res, err := s.Set(context.Background(), &mclight.SetRequest{Key: key, Value: []byte(key)})
	require.Nil(t, err)
	require.NotNil(t, res)
}

func TestSetError(t *testing.T) {
	c := mocks.NewCache(t)
	key := uuid.New().String()
	c.On("Set", key, []byte(key), m.AnythingOfType("Duration")).Return(fmt.Errorf("test"))
	s := New(c)
	res, err := s.Set(context.Background(), &mclight.SetRequest{Key: key, Value: []byte(key)})
	require.NotNil(t, err)
	require.Nil(t, res)
}

func TestDeleteError(t *testing.T) {
	c := mocks.NewCache(t)
	key := uuid.New().String()
	c.On("Delete", key).Return(fmt.Errorf("test"))
	s := New(c)
	res, err := s.Delete(context.Background(), &mclight.DeleteRequest{Key: key})
	require.NotNil(t, err)
	require.Nil(t, res)
}

func TestDelete(t *testing.T) {
	c := mocks.NewCache(t)
	key := uuid.New().String()
	c.On("Delete", key).Return(nil)
	s := New(c)
	res, err := s.Delete(context.Background(), &mclight.DeleteRequest{Key: key})
	require.Nil(t, err)
	require.NotNil(t, res)
}
