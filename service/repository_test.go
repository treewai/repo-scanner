package service

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http/httptest"
	"secret-scanner/models"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func (m *mockDB) AddRepository(repo *models.Repository) error {
	m.req = repo
	if m.err != nil {
		return m.err
	}
	return nil
}

func (m *mockDB) GetRepositories() ([]*models.Repository, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.resp.([]*models.Repository), nil
}

func (m *mockDB) GetRepositoryById(repoId string) (*models.Repository, error) {
	m.req = repoId
	if m.err != nil {
		return nil, m.err
	}
	return m.resp.(*models.Repository), nil
}

func (m *mockDB) UpdateRepository(repo *models.Repository) error {
	m.req = repo
	return m.err
}

func (m *mockDB) DeleteRepository(repoId string) error {
	m.req = repoId
	return m.err
}

func TestServiceAddRepository(t *testing.T) {
	tests := []struct {
		name        string
		req         *models.Repository
		contentType string
		err         error
		code        int
	}{
		{
			name: "success",
			req: &models.Repository{
				ID:   "1",
				Name: "Repo-1",
				Url:  "https://example.com/test",
			},
			contentType: fiber.MIMEApplicationJSON,
			code:        fiber.StatusCreated,
		},
		{
			name: "unsupported media type error",
			req: &models.Repository{
				ID:   "1",
				Name: "Repo-1",
				Url:  "https://example.com/test",
			},
			code: fiber.StatusUnsupportedMediaType,
		},
		{
			name: "missing url error",
			req: &models.Repository{
				ID:   "1",
				Name: "Repo-1",
			},
			contentType: fiber.MIMEApplicationJSON,
			code:        fiber.StatusBadRequest,
		},
		{
			name: "missing name error",
			req: &models.Repository{
				ID:  "1",
				Url: "https://example.com/test",
			},
			contentType: fiber.MIMEApplicationJSON,
			code:        fiber.StatusBadRequest,
		},
		{
			name: "internal server error",
			req: &models.Repository{
				ID:   "1",
				Name: "Repo-1",
				Url:  "https://example.com/test",
			},
			contentType: fiber.MIMEApplicationJSON,
			err:         fiber.ErrInternalServerError,
			code:        fiber.StatusInternalServerError,
		},
	}

	s := getService()

	app := fiber.New()
	app.Post("/repositories", s.AddRepository)

	for _, test := range tests {
		s.getServiceDB().req = test.req
		s.getServiceDB().err = test.err

		b, err := json.Marshal(test.req)
		require.NoErrorf(t, err, test.name)

		req := httptest.NewRequest("POST", "/repositories", bytes.NewBuffer(b))
		req.Header.Set("Content-Type", test.contentType)

		resp, err := app.Test(req, timeout)

		require.NoErrorf(t, err, test.name)
		assert.Equalf(t, test.code, resp.StatusCode, test.name)

		switch test.code {
		case fiber.StatusCreated, fiber.StatusInternalServerError:
			assert.Equal(t, test.req, s.getServiceDB().req.(*models.Repository))
		}
	}
}

func TestServiceGetRepositories(t *testing.T) {
	tests := []struct {
		name string
		resp []*models.Repository
		err  error
		code int
	}{
		{
			name: "success",
			resp: []*models.Repository{
				{
					ID:   "1",
					Name: "Repo-1",
					Url:  "https://example.com/test1",
				},
				{
					ID:   "2",
					Name: "Repo-2",
					Url:  "https://example.com/test2",
				},
				{
					ID:   "3",
					Name: "Repo-3",
					Url:  "https://example.com/test3",
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
	app.Get("/repositories", s.GetRepositories)

	for _, test := range tests {
		s.getServiceDB().resp = test.resp
		s.getServiceDB().err = test.err

		req := httptest.NewRequest("GET", "/repositories", nil)

		resp, err := app.Test(req, timeout)
		require.NoErrorf(t, err, test.name)
		assert.Equalf(t, test.code, resp.StatusCode, test.name)

		if test.code == fiber.StatusOK {
			body, err := ioutil.ReadAll(resp.Body)
			require.NoErrorf(t, err, test.name)

			defer resp.Body.Close()

			var respRepos []*models.Repository
			json.Unmarshal(body, &respRepos)
			assert.Equalf(t, test.resp, respRepos, test.name)
		}
	}
}

func TestServiceGetRositoryById(t *testing.T) {
	tests := []struct {
		name string
		req  *models.Repository
		err  error
		code int
	}{
		{
			name: "success",
			req: &models.Repository{
				ID:   "1",
				Name: "Repo-1",
				Url:  "https://example.com/test1",
			},
			code: fiber.StatusOK,
		},
		{
			name: "internal server error",
			req: &models.Repository{
				ID: "1",
			},
			err:  fiber.ErrInternalServerError,
			code: fiber.StatusInternalServerError,
		},
	}

	s := getService()

	app := fiber.New()
	app.Get("/repositories/:repoId", s.GetRepositoryById)

	for _, test := range tests {
		s.getServiceDB().resp = test.req
		s.getServiceDB().err = test.err

		req := httptest.NewRequest("GET", "/repositories/1", nil)

		resp, err := app.Test(req, timeout)
		require.NoErrorf(t, err, test.name)
		assert.Equalf(t, test.req.ID, s.db.(*mockDB).req.(string), test.name)
		assert.Equalf(t, test.code, resp.StatusCode, test.name)

		if test.code == fiber.StatusOK {
			body, err := ioutil.ReadAll(resp.Body)
			require.NoError(t, err)

			defer resp.Body.Close()

			var respRepo *models.Repository
			json.Unmarshal(body, &respRepo)
			assert.Equal(t, test.req, respRepo)
		}
	}
}

func TestServiceUpdateRepository(t *testing.T) {
	tests := []struct {
		name        string
		req         *models.Repository
		contentType string
		err         error
		code        int
	}{
		{
			name: "success",
			req: &models.Repository{
				ID:   "1",
				Name: "Repo-1",
				Url:  "https://example.com/test1",
			},
			contentType: fiber.MIMEApplicationJSON,
			code:        fiber.StatusOK,
		},
		{
			name: "unsupport media type error",
			req: &models.Repository{
				ID:   "1",
				Name: "Repo-1",
				Url:  "https://example.com/test1",
			},
			code: fiber.StatusUnsupportedMediaType,
		},
		{
			name: "missing both name and url error",
			req: &models.Repository{
				ID: "1",
			},
			contentType: fiber.MIMEApplicationJSON,
			code:        fiber.StatusBadRequest,
		},
		{
			name: "internal error",
			req: &models.Repository{
				ID:   "1",
				Name: "Repo-1",
				Url:  "https://example.com/test1",
			},
			contentType: fiber.MIMEApplicationJSON,
			err:         fiber.ErrInternalServerError,
			code:        fiber.StatusInternalServerError,
		},
	}

	s := getService()

	app := fiber.New()
	app.Patch("/repositories/:repoId", s.UpdateRepository)

	for _, test := range tests {
		s.getServiceDB().err = test.err

		b, err := json.Marshal(test.req)
		require.NoError(t, err)

		req := httptest.NewRequest("PATCH", "/repositories/1", bytes.NewBuffer(b))
		req.Header.Set("Content-Type", test.contentType)

		resp, err := app.Test(req, timeout)
		require.NoErrorf(t, err, test.name)
		assert.Equalf(t, test.code, resp.StatusCode, test.name)

		switch test.code {
		case fiber.StatusOK, fiber.StatusInternalServerError:
			assert.Equalf(t, test.req, s.getServiceDB().req.(*models.Repository), test.name)
		}
	}
}

func TestServicDeleteRepository(t *testing.T) {
	tests := []struct {
		name   string
		repoId string
		err    error
		code   int
	}{
		{
			name:   "success",
			repoId: "1",
			code:   fiber.StatusOK,
		},
		{
			name:   "internal server error",
			repoId: "2",
			err:    fiber.ErrInternalServerError,
			code:   fiber.StatusInternalServerError,
		},
	}
	s := getService()

	app := fiber.New()
	app.Delete("/repositories/:repoId", s.DeleteRepository)

	for _, test := range tests {
		s.getServiceDB().err = test.err

		req := httptest.NewRequest("DELETE", "/repositories/"+test.repoId, nil)

		resp, err := app.Test(req, timeout)
		require.NoErrorf(t, err, test.name)
		assert.Equalf(t, test.code, resp.StatusCode, test.name)
		assert.Equalf(t, test.repoId, s.getServiceDB().req.(string), test.name)
	}
}
