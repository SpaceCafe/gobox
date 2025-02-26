package rest

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	authorization "github.com/spacecafe/gobox/gin-authorization"
	jwt "github.com/spacecafe/gobox/gin-jwt"
	problems "github.com/spacecafe/gobox/gin-problems"
	"github.com/spacecafe/gobox/gin-rest/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	_ types.IModel           = (*Book)(nil)
	_ types.IModelReadable   = (*Book)(nil)
	_ types.IModelUpdatable  = (*Book)(nil)
	_ types.IModelFilterable = (*Book)(nil)
	_ types.IModelSortable   = (*Book)(nil)
)

type Book struct {
	Model      `yaml:",inline"`
	Title      string   `json:"title" yaml:"title" gorm:"not null;default:null"`
	AuthorID   string   `json:"-" yaml:"-"`
	Author     Author   `json:"author" yaml:"author" binding:"omitempty"`
	CoAuthors  []Author `json:"co_authors" yaml:"co_authors" gorm:"many2many" binding:"omitempty"`
	Category   Category `json:"category" yaml:"category" binding:"omitempty"`
	CategoryID int      `json:"-" yaml:"-"`
}

func (r *Book) BeforeRender(_ *gin.Context) (_ error) {
	r.CreatedAt = nil
	r.UpdatedAt = nil
	return
}

func (r *Book) Readable(_ *gin.Context) []string {
	return []string{"id", "title", "author_id", "category_id"}
}

func (r *Book) Updatable(_ *gin.Context) []string {
	return []string{"title"}
}

func (r *Book) Filterable(_ *gin.Context) map[string]struct{} {
	return map[string]struct{}{"id": {}, "title": {}}
}

func (r *Book) Sortable(_ *gin.Context) map[string]struct{} {
	return map[string]struct{}{"created_at": {}}
}

type Category struct {
	ID   int    `json:"-" yaml:"-"`
	Name string `json:"name" yaml:"name"`
}

type Author struct {
	Model     `json:"-" yaml:"-"`
	Name      string  `json:"name" yaml:"name"`
	Address   Address `json:"-" yaml:"-"`
	AddressID int     `json:"-" yaml:"-"`
}

type Address struct {
	ID      int
	Street  string
	City    string
	Country string
}

func setupDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	assert.NoError(t, err)
	err = db.AutoMigrate(&Book{}, &Author{})
	assert.NoError(t, err)
	return db
}

func setupREST(t *testing.T) (*gin.Engine, *gorm.DB, *jwt.Config) {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	r.Use(problems.New())

	db := setupDB(t)
	jwtConfig := jwt.NewConfig(nil)
	jwtConfig.SecretKey = "c2VjcmV0"
	require.NoError(t, jwtConfig.Validate())

	api := New(r.Group("/api/v1"), jwtConfig, authorization.NewConfig())
	api.Register(NewDatabaseResource[Book](nil, nil, nil, db))

	return r, db, jwtConfig
}

func token(jwtConfig *jwt.Config, entitlements []string) string {
	claims := jwt.NewClaims(jwtConfig, "test")
	claims.Entitlements = entitlements
	token, _ := jwt.NewToken(jwtConfig, claims)
	return token.String()
}

func TestREST_CRUD(t *testing.T) {
	r, _, jwtConfig := setupREST(t)

	type args struct {
		method       string
		path         string
		payload      string
		accept       string
		entitlements []string
	}
	type wants struct {
		status int
		body   string
	}
	tests := []struct {
		name  string
		args  args
		wants wants
	}{
		// Create
		{
			name: "create first book",
			args: args{
				method:       "POST",
				payload:      `{"id": "4d6f5960-ff88-4dcf-b69c-54fce67ef28e", "title":"test1","author":{"name":"Peter Parker","address":{"street":"Bakers Street"}},"co_authors":[{"name":"John Doe","address":{"street":"Numb Avenue"}}],"category":{"name":"Fiction"}}`,
				accept:       "application/json",
				entitlements: []string{"create_books"},
			},
			wants: wants{
				status: http.StatusCreated,
			},
		},
		{
			name: "create first book again",
			args: args{
				method:       "POST",
				payload:      `{"id": "4d6f5960-ff88-4dcf-b69c-54fce67ef28e", "title":"test1","author":{"name":"Peter Parker","address":{"street":"Bakers Street"}},"co_authors":[{"name":"John Doe","address":{"street":"Numb Avenue"}}]}`,
				accept:       "application/json",
				entitlements: []string{"create_books"},
			},
			wants: wants{
				status: http.StatusConflict,
			},
		},
		{
			name: "create second book",
			args: args{
				method:       "POST",
				payload:      `{"title":"test2","author":{"name":"Bruce Wayne","address":{"street":"Wayne Manor"}},"co_authors":[{"name":"Clark Kent","address":{"street":"Kent Farm"}}]}`,
				accept:       "application/x-yaml",
				entitlements: []string{"create_books"},
			},
			wants: wants{
				status: http.StatusCreated,
			},
		},
		{
			name: "create with empty payload",
			args: args{
				method:       "POST",
				payload:      `{}`,
				accept:       "application/json",
				entitlements: []string{"create_books"},
			},
			wants: wants{
				status: http.StatusBadRequest,
			},
		},
		{
			name: "create with bad payload",
			args: args{
				method:       "POST",
				payload:      `{"}`,
				accept:       "application/json",
				entitlements: []string{"create_books"},
			},
			wants: wants{
				status: http.StatusBadRequest,
			},
		},
		{
			name: "create without entitlements",
			args: args{
				method:       "POST",
				payload:      `{}`,
				accept:       "application/json",
				entitlements: []string{},
			},
			wants: wants{
				status: http.StatusForbidden,
			},
		},

		// Read
		{
			name: "read existing book as json",
			args: args{
				method:       "GET",
				path:         "4d6f5960-ff88-4dcf-b69c-54fce67ef28e",
				accept:       "application/json",
				entitlements: []string{"read_books"},
			},
			wants: wants{
				status: http.StatusOK,
				body:   `{"page":1,"page_size":10,"total":1,"total_pages":1,"data":{"id":"4d6f5960-ff88-4dcf-b69c-54fce67ef28e","title":"test1","author":{"name":"Peter Parker"},"co_authors":[{"name":"John Doe"}],"category":{"name":"Fiction"}}}`,
			},
		},
		{
			name: "read existing book as yaml",
			args: args{
				method:       "GET",
				path:         "4d6f5960-ff88-4dcf-b69c-54fce67ef28e",
				accept:       "application/x-yaml",
				entitlements: []string{"read_books"},
			},
			wants: wants{
				status: http.StatusOK,
				body:   `{"page":1,"page_size":10,"total":1,"total_pages":1,"data":{"id":"4d6f5960-ff88-4dcf-b69c-54fce67ef28e","title":"test1","author":{"name":"Peter Parker"},"co_authors":[{"name":"John Doe"}],"category":{"name":"Fiction"}}}`,
			},
		},
		{
			name: "read existing book with invalid accept",
			args: args{
				method:       "GET",
				path:         "4d6f5960-ff88-4dcf-b69c-54fce67ef28e",
				accept:       "text/plain",
				entitlements: []string{"read_books"},
			},
			wants: wants{
				status: http.StatusUnsupportedMediaType,
				body:   "",
			},
		},
		{
			name: "read non-existing book",
			args: args{
				method:       "GET",
				path:         "999",
				accept:       "*/*",
				entitlements: []string{"read_books"},
			},
			wants: wants{
				status: http.StatusNotFound,
				body:   `{"detail":"The specified key does not exist.","instance":"/api/v1/books/999","status":404,"title":"No such key","type":"/errors/no-such-key"}`,
			},
		},
		{
			name: "read without entitlements",
			args: args{
				method:       "GET",
				path:         "4d6f5960-ff88-4dcf-b69c-54fce67ef28e",
				accept:       "application/json",
				entitlements: []string{},
			},
			wants: wants{
				status: http.StatusForbidden,
				body:   `{"detail":"Insufficient permission","instance":"/api/v1/books/4d6f5960-ff88-4dcf-b69c-54fce67ef28e","status":403,"title":"Insufficient permission","type":"/errors/insufficient-permission"}`,
			},
		},

		// List
		{
			name: "list one books as json",
			args: args{
				method:       "GET",
				path:         "?page[number]=1&page[size]=1&sort=+created_at",
				accept:       "application/json",
				entitlements: []string{"list_books"},
			},
			wants: wants{
				status: http.StatusOK,
				body:   `{"page":1,"page_size":1,"total":2,"total_pages":2,"data":[{"id":"4d6f5960-ff88-4dcf-b69c-54fce67ef28e","title":"test1","author":{"name":"Peter Parker"},"co_authors":[{"name":"John Doe"}],"category":{"name":"Fiction"}}]}`,
			},
		},
		{
			name: "list one book as yaml",
			args: args{
				method:       "GET",
				path:         "?page[number]=2&page[size]=1&sort=-created_at",
				accept:       "application/x-yaml",
				entitlements: []string{"list_books"},
			},
			wants: wants{
				status: http.StatusOK,
				body:   `{"page":2,"page_size":1,"total":2,"total_pages":2,"data":[{"id":"4d6f5960-ff88-4dcf-b69c-54fce67ef28e","title":"test1","author":{"name":"Peter Parker"},"co_authors":[{"name":"John Doe"}],"category":{"name":"Fiction"}}]}`,
			},
		},
		{
			name: "list second book",
			args: args{
				method:       "GET",
				path:         "?page[size]=1&filter[title][ne]=test2",
				accept:       "application/json",
				entitlements: []string{"list_books"},
			},
			wants: wants{
				status: http.StatusOK,
				body:   `{"page":1,"page_size":1,"total":1,"total_pages":1,"data":[{"id":"4d6f5960-ff88-4dcf-b69c-54fce67ef28e","title":"test1","author":{"name":"Peter Parker"},"co_authors":[{"name":"John Doe"}],"category":{"name":"Fiction"}}]}`,
			},
		},
		{
			name: "list without entitlements",
			args: args{
				method:       "GET",
				accept:       "application/json",
				entitlements: []string{},
			},
			wants: wants{
				status: http.StatusForbidden,
				body:   `{"detail":"Insufficient permission","instance":"/api/v1/books/","status":403,"title":"Insufficient permission","type":"/errors/insufficient-permission"}`,
			},
		},

		// Update
		{
			name: "update existing book",
			args: args{
				method:       "PUT",
				path:         "4d6f5960-ff88-4dcf-b69c-54fce67ef28e",
				payload:      `{"title": "Updated Title"}`,
				accept:       "application/json",
				entitlements: []string{"update_books"},
			},
			wants: wants{
				status: http.StatusNoContent,
			},
		},
		{
			name: "update non-existing book",
			args: args{
				method:       "PUT",
				path:         "999",
				payload:      `{"title": "Non-Existent Book"}`,
				accept:       "application/x-yaml",
				entitlements: []string{"update_books"},
			},
			wants: wants{
				status: http.StatusNotFound,
			},
		},
		{
			name: "update existing book without entitlement",
			args: args{
				method:       "PUT",
				path:         "4d6f5960-ff88-4dcf-b69c-54fce67ef28e",
				payload:      `{"title": "Updated Title"}`,
				accept:       "application/json",
				entitlements: []string{},
			},
			wants: wants{
				status: http.StatusForbidden,
			},
		},

		// Delete
		{
			name: "delete existing book",
			args: args{
				method:       "DELETE",
				path:         "4d6f5960-ff88-4dcf-b69c-54fce67ef28e",
				accept:       "application/json",
				entitlements: []string{"delete_books"},
			},
			wants: wants{
				status: http.StatusNoContent,
			},
		},
		{
			name: "delete non-existing book",
			args: args{
				method:       "DELETE",
				path:         "999",
				accept:       "application/x-yaml",
				entitlements: []string{"delete_books"},
			},
			wants: wants{
				status: http.StatusNotFound,
			},
		},
		{
			name: "delete existing book without entitlement",
			args: args{
				method:       "DELETE",
				path:         "4d6f5960-ff88-4dcf-b69c-54fce67ef28e",
				accept:       "application/json",
				entitlements: []string{},
			},
			wants: wants{
				status: http.StatusForbidden,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			recorder := httptest.NewRecorder()
			request := httptest.NewRequest(tt.args.method, "/api/v1/books/"+tt.args.path, strings.NewReader(tt.args.payload))
			request.Header.Set("Content-Type", "application/json")
			request.Header.Set("Authorization", "Bearer "+token(jwtConfig, tt.args.entitlements))
			request.Header.Set("Accept", tt.args.accept)

			r.ServeHTTP(recorder, request)
			response := recorder.Result()
			_ = response.Body.Close()

			assert.Equal(t, tt.wants.status, response.StatusCode)
			fmt.Printf("%+v", response.Header)
			if tt.wants.body != "" {
				switch tt.args.accept {
				case "application/json":
					assert.JSONEq(t, tt.wants.body, recorder.Body.String())
				case "application/x-yaml":
					assert.YAMLEq(t, tt.wants.body, recorder.Body.String())
				default:
				}
			}
		})
	}
}
