package service

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http/httptest"
	"secret-scanner/models"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func (m *mockDB) GetResults() ([]*models.Result, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.resp.([]*models.Result), nil
}

func (m *mockDB) UpdateResult(res *models.Result) error {
	m.req = res
	return m.err
}

func TestServiceGetScanResults(t *testing.T) {
	tests := []struct {
		name string
		resp []*models.Result
		err  error
		code int
	}{
		{
			name: "success",
			resp: []*models.Result{
				{
					ID:       "1",
					Status:   models.StatusInProgress,
					RepoName: "Repop-1",
					RepoUrl:  "https://example.com/test1",
				},
				{
					ID:       "2",
					Status:   models.StatusSuccess,
					RepoName: "Repop-2",
					RepoUrl:  "https://example.com/test2",
				},
				{
					ID:       "3",
					Status:   models.StatusFailure,
					RepoName: "Repop-3",
					RepoUrl:  "https://example.com/test3",
				},
			},
			code: fiber.StatusOK,
		},
		{
			name: "internal server error",
			err:  fiber.ErrInternalServerError,
			code: fiber.StatusInternalServerError,
		},
	}

	s := getService()

	app := fiber.New()
	app.Get("/scans", s.GetScanResults)

	for _, test := range tests {
		s.getServiceDB().resp = test.resp
		s.getServiceDB().err = test.err

		req := httptest.NewRequest("GET", "/scans", nil)

		resp, err := app.Test(req, timeout)
		require.NoErrorf(t, err, test.name)
		assert.Equalf(t, test.code, resp.StatusCode, test.name)

		if test.code == fiber.StatusOK {
			body, err := ioutil.ReadAll(resp.Body)
			require.NoErrorf(t, err, test.name)

			defer resp.Body.Close()

			var respResult []*models.Result
			json.Unmarshal(body, &respResult)
			assert.Equalf(t, test.resp, respResult, test.name)
			fmt.Println(test.resp[0])
			fmt.Println(respResult[0])
		}
	}
}

func TestServicUpdateScanResult(t *testing.T) {
	tests := []struct {
		name string
		req  *models.Result
		err  error
	}{
		{
			name: "success",
			req: &models.Result{
				ID:     "1",
				Status: models.StatusQueued,
				RepoID: "100",
			},
		},
		{
			name: "internal server error",
			req: &models.Result{
				ID:     "1",
				Status: models.StatusQueued,
				RepoID: "100",
			},
			err: fiber.ErrInternalServerError,
		},
	}
	s := getService()

	for _, test := range tests {
		s.getServiceDB().req = test.req
		s.getServiceDB().err = test.err

		err := s.UpdateScanResult(test.req)
		assert.Equalf(t, test.err, err, test.name)
		assert.Equalf(t, test.req, s.getServiceDB().req.(*models.Result), test.name)
	}
}
