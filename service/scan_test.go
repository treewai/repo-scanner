package service

import (
	"net/http/httptest"
	"secret-scanner/models"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func (m *mockDB) AddResult(res *models.Result) (*models.Result, error) {
	if m.scanErr != nil {
		return nil, m.scanErr
	}
	return m.scanResp.(*models.Result), nil
}

func TestServiceScanRepository(t *testing.T) {
	tests := []struct {
		name     string
		respRepo *models.Repository
		errRepo  error
		req      *models.Result
		respScan *models.Result
		errScan  error
		err      error
		code     int
	}{
		{
			name: "success",
			respRepo: &models.Repository{
				ID:   "1",
				Name: "Repo-1",
				Url:  "https://example.com/test",
			},
			respScan: &models.Result{
				Status: models.StatusQueued,
			},
			code: fiber.StatusOK,
		},
		{
			name: "register error",
			respRepo: &models.Repository{
				ID:   "1",
				Name: "Repo-1",
				Url:  "https://example.com/test",
			},
			respScan: &models.Result{
				Status: models.StatusQueued,
			},
			err:  fiber.ErrInternalServerError,
			code: fiber.StatusInternalServerError,
		},
		{
			name: "add scan error",
			respRepo: &models.Repository{
				ID:   "1",
				Name: "Repo-1",
				Url:  "https://example.com/test",
			},
			errScan: fiber.ErrInternalServerError,
			code:    fiber.StatusInternalServerError,
		},
		{
			name: "repo not found error",
			respScan: &models.Result{
				Status: models.StatusQueued,
			},
			errScan: fiber.ErrInternalServerError,
			code:    fiber.StatusNotFound,
		},
		{
			name: "get repo error",
			respRepo: &models.Repository{
				ID:   "1",
				Name: "Repo-1",
				Url:  "https://example.com/test",
			},
			errRepo: fiber.ErrInternalServerError,
			code:    fiber.StatusInternalServerError,
		},
	}

	s := getService()

	app := fiber.New()
	app.Post("/scans/:repoId", s.ScanRepository)

	for _, test := range tests {
		db := s.getServiceDB()
		db.resp = test.respRepo
		db.err = test.errRepo
		db.scanResp = test.respScan
		db.scanErr = test.errScan

		s.register = func(*models.Result) error {
			return test.err
		}

		req := httptest.NewRequest("POST", "/scans/1", nil)

		resp, err := app.Test(req, timeout)

		require.NoErrorf(t, err, test.name)
		assert.Equalf(t, test.code, resp.StatusCode, test.name)
	}
}
