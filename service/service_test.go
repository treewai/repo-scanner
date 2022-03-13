package service

import (
	"reflect"
	"runtime"
	"secret-scanner/models"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const timeout = 3

type mockDB struct {
	req  interface{}
	resp interface{}
	err  error

	scanResp interface{}
	scanErr  error
}

func getService() *service {
	return &service{
		db: &mockDB{},
	}
}

func (s *service) getServiceDB() *mockDB {
	return s.db.(*mockDB)
}

func TestServiceNew(t *testing.T) {
	f := func(*models.Result) error {
		return nil
	}
	s := NewService(nil, f)
	switch s.(type) {
	case Service:
		sv, ok := s.(*service)
		require.True(t, ok)
		assert.Equal(t,
			runtime.FuncForPC(reflect.ValueOf(f).Pointer()).Name(),
			runtime.FuncForPC(reflect.ValueOf(sv.register).Pointer()).Name())
	default:
		t.Errorf("invalid interface type")
	}
}

func TestServiceRegister(t *testing.T) {
	scanc := make(chan *models.Result, 1)
	f := Register(scanc)
	result := &models.Result{
		ID:       "12345",
		QueuedAt: time.Now(),
	}

	f(result)

	select {
	case res := <-scanc:
		assert.Equal(t, result, res)
	default:
		t.Errorf("result not found")
	}
}
