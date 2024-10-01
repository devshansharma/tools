package meta_test

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"github.com/devshansharma/tools/meta"
)

func TestPaginationParse(t *testing.T) {
	t.Run("default limit and page", func(t *testing.T) {
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = &http.Request{
			URL: &url.URL{
				RawQuery: "",
			},
		}

		pagination := &meta.Pagination{}
		err := pagination.Parse(ctx)
		assert.NoError(t, err)
		assert.Equal(t, uint(10), pagination.Limit)
		assert.Equal(t, uint(1), pagination.Page)
	})

	t.Run("parse limit and page", func(t *testing.T) {
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = &http.Request{
			URL: &url.URL{
				RawQuery: "limit=20&page=3",
			},
		}

		pagination := &meta.Pagination{}
		err := pagination.Parse(ctx)
		assert.NoError(t, err)
		assert.Equal(t, uint(20), pagination.Limit)
		assert.Equal(t, uint(3), pagination.Page)
	})

	t.Run("parse sort", func(t *testing.T) {
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = &http.Request{
			URL: &url.URL{
				RawQuery: "sort=created_at:desc,updated_at:asc",
			},
		}

		pagination := &meta.Pagination{}
		err := pagination.Parse(ctx)
		assert.NoError(t, err)
		assert.Equal(t, 2, len(pagination.Sort))
		assert.Equal(t, "created_at", pagination.Sort[0].Field)
		assert.Equal(t, "desc", pagination.Sort[0].Order)
		assert.Equal(t, "updated_at", pagination.Sort[1].Field)
		assert.Equal(t, "asc", pagination.Sort[1].Order)
	})

	t.Run("invalid sort field", func(t *testing.T) {
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = &http.Request{
			URL: &url.URL{
				RawQuery: "sort=created_at:desc,updated_at",
			},
		}

		pagination := &meta.Pagination{}
		err := pagination.Parse(ctx)
		assert.Equal(t, "invalid sort parameter: updated_at", err.Error())
	})
}
