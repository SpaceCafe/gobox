package rest

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	problems "github.com/spacecafe/gobox/gin-problems"
	"github.com/spacecafe/gobox/gin-rest/types"
	"github.com/spacecafe/gobox/gin-rest/view/jsonapi"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type BookModel struct {
	types.Model `yaml:",inline"`
	Title       string `gorm:"not null;default:null"`
	AuthorID    string
	Author      AuthorModel
	CoAuthors   []AuthorModel `json:"co_authors" yaml:"co_authors" gorm:"many2many"`
	Category    CategoryModel
	CategoryID  int
}

type CategoryModel struct {
	ID   int
	Name string
}

type AuthorModel struct {
	types.Model `yaml:",inline"`
	Name        string
	Address     AddressModel
	AddressID   int
}

type AddressModel struct {
	ID      int
	Street  string
	City    string
	Country string
}

func setupDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	assert.NoError(t, err)
	err = db.AutoMigrate(&BookModel{}, &AuthorModel{})
	assert.NoError(t, err)
	return db
}

func setupREST(t *testing.T) (*gin.Engine, *gorm.DB) {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	r.Use(problems.New())
	rest := New(r.Group("/api/v1"))

	db := setupDB(t)

	api := &Resource[BookModel]{Repository: db}
	api.Apply(rest)

	return r, db
}

var sampleBook = BookModel{
	Model:     types.Model{ID: "4d6f5960-ff88-4dcf-b69c-54fce67ef28f"},
	Title:     "Sample Book",
	Author:    AuthorModel{Name: "Author One", Address: AddressModel{Street: "123 Main St", City: "Anytown", Country: "USA"}},
	Category:  CategoryModel{Name: "Fiction"},
	CoAuthors: []AuthorModel{{Name: "Author Two", Address: AddressModel{Street: "456 Elm St", City: "Othertown", Country: "USA"}}},
}

func TestREST_Create(t *testing.T) {
	testCases := []struct {
		name       string
		payload    string
		acceptType string
		wantStatus int
	}{
		{
			name:       "create first book",
			payload:    `{"id": "4d6f5960-ff88-4dcf-b69c-54fce67ef28f", "title":"test1","author":{"name":"Peter Parker","address":{"street":"Bakers Street"}},"co_authors":[{"name":"John Doe","address":{"street":"Numb Avenue"}}]}`,
			acceptType: "application/json",
			wantStatus: http.StatusCreated,
		},
		{
			name:       "create first book again",
			payload:    `{"id": "4d6f5960-ff88-4dcf-b69c-54fce67ef28f", "title":"test1","author":{"name":"Peter Parker","address":{"street":"Bakers Street"}},"co_authors":[{"name":"John Doe","address":{"street":"Numb Avenue"}}]}`,
			acceptType: "application/json",
			wantStatus: http.StatusConflict,
		},
		{
			name:       "create second book",
			payload:    `{"title":"test2","author":{"name":"Bruce Wayne","address":{"street":"Wayne Manor"}},"co_authors":[{"name":"Clark Kent","address":{"street":"Kent Farm"}}]}`,
			acceptType: "application/x-yaml",
			wantStatus: http.StatusCreated,
		},
		{
			name:       "create with empty payload",
			payload:    `{}`,
			acceptType: "application/json",
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "create with bad payload",
			payload:    `{"}`,
			acceptType: "application/json",
			wantStatus: http.StatusBadRequest,
		},
	}

	r, _ := setupREST(t)

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			recorder := httptest.NewRecorder()
			request := httptest.NewRequest("POST", "/api/v1/books/", strings.NewReader(tc.payload))
			request.Header.Set("Content-Type", "application/json")
			request.Header.Set("Accept", tc.acceptType)
			r.ServeHTTP(recorder, request)
			response := recorder.Result()
			_ = response.Body.Close()
			assert.Equal(t, tc.wantStatus, response.StatusCode)
		})
	}
}

func TestREST_Read(t *testing.T) {
	r, db := setupREST(t)
	db.Create(&sampleBook)

	tests := []struct {
		name       string
		id         string
		wantStatus int
	}{
		{
			name:       "read existing book",
			id:         sampleBook.ID,
			wantStatus: http.StatusOK,
		},
		{
			name:       "read non-existing book",
			id:         "999",
			wantStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			recorder := httptest.NewRecorder()
			request := httptest.NewRequest("GET", "/api/v1/books/"+tt.id, http.NoBody)
			request.Header.Set("Content-Type", "application/json")
			request.Header.Set("Accept", "application/vnd.api+json")
			r.ServeHTTP(recorder, request)
			response := recorder.Result()
			_ = response.Body.Close()

			assert.Equal(t, tt.wantStatus, response.StatusCode)
		})
	}
}

func TestREST_List(t *testing.T) {
	r, db := setupREST(t)
	db.Create(&sampleBook)
	db.Create(&BookModel{Title: "Example Book"})

	t.Run("List books", func(t *testing.T) {
		recorder := httptest.NewRecorder()
		request := httptest.NewRequest("GET", "/api/v1/books/?sort=-created_at,title,%2Bupdated_at&filter[title][like]=Sample%25", http.NoBody)
		request.Header.Set("Accept", "application/vnd.api+json")
		r.ServeHTTP(recorder, request)
		response := recorder.Result()
		body, _ := io.ReadAll(response.Body)
		_ = response.Body.Close()
		fmt.Printf("%s", body)

		jsonAPI := jsonapi.NewResponse[BookModel]()
		err := json.Unmarshal(body, &jsonAPI)
		assert.NoError(t, err)

		// Assert Meta
		assert.Equal(t, 1, jsonAPI.Meta.Page)
		assert.Equal(t, 10, jsonAPI.Meta.PageSize)
		assert.Equal(t, 1, jsonAPI.Meta.Total)
		assert.Equal(t, 1, jsonAPI.Meta.TotalPages)

		// Assert JSONAPI
		assert.Equal(t, "1.1", jsonAPI.JSONAPI.Version)

		// Assert Data
		assert.Len(t, *jsonAPI.Data, 1)
	})
}

func TestREST_PartialUpdate(t *testing.T) {
	r, db := setupREST(t)
	db.Create(&sampleBook)

	tests := []struct {
		name       string
		id         string
		payload    string
		wantStatus int
	}{
		{
			name:       "update existing book",
			id:         sampleBook.ID,
			payload:    `{"title": "Updated Title"}`,
			wantStatus: http.StatusNoContent,
		},
		{
			name:       "update non-existing book",
			id:         "999",
			payload:    `{"title": "Non-Existent Book"}`,
			wantStatus: http.StatusNotFound,
		},
		{
			name:       "update with empty payload",
			id:         sampleBook.ID,
			payload:    `{}`,
			wantStatus: http.StatusNoContent,
		},
		{
			name:       "update with bad payload",
			id:         sampleBook.ID,
			payload:    `{"title": "Non-Existent}`,
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			recorder := httptest.NewRecorder()
			request := httptest.NewRequest("PATCH", "/api/v1/books/"+tt.id, strings.NewReader(tt.payload))
			request.Header.Set("Content-Type", "application/json")
			request.Header.Set("Accept", "application/vnd.api+json")
			r.ServeHTTP(recorder, request)
			response := recorder.Result()
			_ = response.Body.Close()

			assert.Equal(t, tt.wantStatus, response.StatusCode)
		})
	}
}

func TestREST_Update(t *testing.T) {
	r, db := setupREST(t)
	db.Create(&sampleBook)

	tests := []struct {
		name       string
		id         string
		payload    string
		wantStatus int
	}{
		{
			name:       "update existing book",
			id:         sampleBook.ID,
			payload:    `{"title": "Updated Title"}`,
			wantStatus: http.StatusNoContent,
		},
		{
			name:       "update non-existing book",
			id:         "999",
			payload:    `{"title": "Non-Existent Book"}`,
			wantStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			recorder := httptest.NewRecorder()
			request := httptest.NewRequest("PUT", "/api/v1/books/"+tt.id, strings.NewReader(tt.payload))
			request.Header.Set("Content-Type", "application/json")
			request.Header.Set("Accept", "application/vnd.api+json")
			r.ServeHTTP(recorder, request)
			response := recorder.Result()
			_ = response.Body.Close()

			assert.Equal(t, tt.wantStatus, response.StatusCode)
		})
	}
}

func TestREST_Delete(t *testing.T) {
	r, db := setupREST(t)
	db.Create(&sampleBook)

	tests := []struct {
		name       string
		id         string
		wantStatus int
	}{
		{
			name:       "delete existing book",
			id:         sampleBook.ID,
			wantStatus: http.StatusNoContent,
		},
		{
			name:       "delete non-existing book",
			id:         "999",
			wantStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			recorder := httptest.NewRecorder()
			request := httptest.NewRequest("DELETE", "/api/v1/books/"+tt.id, http.NoBody)
			request.Header.Set("Accept", "application/vnd.api+json")
			r.ServeHTTP(recorder, request)
			response := recorder.Result()
			_ = response.Body.Close()

			assert.Equal(t, tt.wantStatus, response.StatusCode)
		})
	}
}
