package api

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	apperrors "github.com/miguel/go-back-portfolo/pkg/errors"
	"github.com/miguel/go-back-portfolo/pkg/response"
	"github.com/miguel/go-back-portfolo/schema/db"
	"github.com/miguel/go-back-portfolo/service/visitorservice"
)

type handlerMockStore struct {
	db.Store
	createVisitorFn func(ctx context.Context, data []byte) (db.VisitorVisitor, error)
	getVisitorFn    func(ctx context.Context, id pgtype.UUID) (db.VisitorVisitor, error)
	listVisitorsFn  func(ctx context.Context, arg db.ListVisitorsParams) ([]db.VisitorVisitor, error)
	countVisitorsFn func(ctx context.Context) (int64, error)
	updateVisitorFn func(ctx context.Context, arg db.UpdateVisitorParams) (db.VisitorVisitor, error)
	deleteVisitorFn func(ctx context.Context, id pgtype.UUID) (pgtype.UUID, error)
}

func (m *handlerMockStore) CreateVisitor(ctx context.Context, data []byte) (db.VisitorVisitor, error) {
	return m.createVisitorFn(ctx, data)
}

func (m *handlerMockStore) GetVisitor(ctx context.Context, id pgtype.UUID) (db.VisitorVisitor, error) {
	return m.getVisitorFn(ctx, id)
}

func (m *handlerMockStore) ListVisitors(ctx context.Context, arg db.ListVisitorsParams) ([]db.VisitorVisitor, error) {
	return m.listVisitorsFn(ctx, arg)
}

func (m *handlerMockStore) CountVisitors(ctx context.Context) (int64, error) {
	return m.countVisitorsFn(ctx)
}

func (m *handlerMockStore) UpdateVisitor(ctx context.Context, arg db.UpdateVisitorParams) (db.VisitorVisitor, error) {
	return m.updateVisitorFn(ctx, arg)
}

func (m *handlerMockStore) DeleteVisitor(ctx context.Context, id pgtype.UUID) (pgtype.UUID, error) {
	return m.deleteVisitorFn(ctx, id)
}

func setupTestServer(mock *handlerMockStore) *Server {
	svc := visitorservice.NewService(mock)
	return &Server{
		visitorService: svc,
		store:          mock,
	}
}

func decodeResponse(t *testing.T, body []byte) response.APIResponse {
	t.Helper()
	var resp response.APIResponse
	err := json.Unmarshal(body, &resp)
	require.NoError(t, err)
	return resp
}

func TestCreateVisitor_Success(t *testing.T) {
	var id pgtype.UUID
	require.NoError(t, id.Scan("a1b2c3d4-e5f6-7890-abcd-ef1234567890"))

	mock := &handlerMockStore{
		createVisitorFn: func(ctx context.Context, data []byte) (db.VisitorVisitor, error) {
			return db.VisitorVisitor{ID: id, Data: data}, nil
		},
	}
	srv := setupTestServer(mock)
	router := SetupRouter(srv)

	body := `{"name":"Miguel","email":"miguel@test.com","country":"Mexico","city":"CDMX"}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/visitors", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	resp := decodeResponse(t, w.Body.Bytes())
	assert.True(t, resp.Success)
}

func TestCreateVisitor_ValidationError(t *testing.T) {
	mock := &handlerMockStore{}
	srv := setupTestServer(mock)
	router := SetupRouter(srv)

	body := `{"name":"","email":"not-an-email","country":"","city":""}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/visitors", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	resp := decodeResponse(t, w.Body.Bytes())
	assert.False(t, resp.Success)
}

func TestCreateVisitor_MissingFields(t *testing.T) {
	mock := &handlerMockStore{}
	srv := setupTestServer(mock)
	router := SetupRouter(srv)

	body := `{"name":"Miguel"}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/visitors", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestGetVisitor_Success(t *testing.T) {
	var id pgtype.UUID
	require.NoError(t, id.Scan("a1b2c3d4-e5f6-7890-abcd-ef1234567890"))

	mock := &handlerMockStore{
		getVisitorFn: func(ctx context.Context, uid pgtype.UUID) (db.VisitorVisitor, error) {
			return db.VisitorVisitor{
				ID:   uid,
				Data: []byte(`{"name":"Miguel","email":"miguel@test.com","country":"Mexico","city":"CDMX"}`),
			}, nil
		},
	}
	srv := setupTestServer(mock)
	router := SetupRouter(srv)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/visitors/a1b2c3d4-e5f6-7890-abcd-ef1234567890", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	resp := decodeResponse(t, w.Body.Bytes())
	assert.True(t, resp.Success)
}

func TestGetVisitor_NotFound(t *testing.T) {
	var emptyUID pgtype.UUID

	mock := &handlerMockStore{
		getVisitorFn: func(ctx context.Context, uid pgtype.UUID) (db.VisitorVisitor, error) {
			return db.VisitorVisitor{ID: emptyUID}, apperrors.ErrNotFound
		},
	}
	srv := setupTestServer(mock)
	router := SetupRouter(srv)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/visitors/00000000-0000-0000-0000-000000000001", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	resp := decodeResponse(t, w.Body.Bytes())
	assert.False(t, resp.Success)
}

func TestGetVisitor_InvalidUUID(t *testing.T) {
	mock := &handlerMockStore{}
	srv := setupTestServer(mock)
	router := SetupRouter(srv)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/visitors/not-a-uuid", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestListVisitors_Success(t *testing.T) {
	var id1, id2 pgtype.UUID
	require.NoError(t, id1.Scan("a1b2c3d4-e5f6-7890-abcd-ef1234567890"))
	require.NoError(t, id2.Scan("b2c3d4e5-f6a7-8901-bcde-f12345678901"))

	mock := &handlerMockStore{
		countVisitorsFn: func(ctx context.Context) (int64, error) { return 2, nil },
		listVisitorsFn: func(ctx context.Context, arg db.ListVisitorsParams) ([]db.VisitorVisitor, error) {
			return []db.VisitorVisitor{
				{ID: id1, Data: []byte(`{"name":"Alice","email":"alice@test.com","country":"US","city":"NYC"}`)},
				{ID: id2, Data: []byte(`{"name":"Bob","email":"bob@test.com","country":"UK","city":"London"}`)},
			}, nil
		},
	}
	srv := setupTestServer(mock)
	router := SetupRouter(srv)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/visitors?page=1&page_size=20", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestListVisitors_Empty(t *testing.T) {
	mock := &handlerMockStore{
		countVisitorsFn: func(ctx context.Context) (int64, error) { return 0, nil },
		listVisitorsFn: func(ctx context.Context, arg db.ListVisitorsParams) ([]db.VisitorVisitor, error) {
			return []db.VisitorVisitor{}, nil
		},
	}
	srv := setupTestServer(mock)
	router := SetupRouter(srv)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/visitors", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestListVisitors_InvalidPage(t *testing.T) {
	mock := &handlerMockStore{
		countVisitorsFn: func(ctx context.Context) (int64, error) { return 0, nil },
		listVisitorsFn: func(ctx context.Context, arg db.ListVisitorsParams) ([]db.VisitorVisitor, error) {
			assert.Equal(t, int32(20), arg.Limit)
			assert.Equal(t, int32(0), arg.Offset)
			return []db.VisitorVisitor{}, nil
		},
	}
	srv := setupTestServer(mock)
	router := SetupRouter(srv)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/visitors?page=-1&page_size=0", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestUpdateVisitor_Success(t *testing.T) {
	var id pgtype.UUID
	require.NoError(t, id.Scan("a1b2c3d4-e5f6-7890-abcd-ef1234567890"))

	mock := &handlerMockStore{
		updateVisitorFn: func(ctx context.Context, arg db.UpdateVisitorParams) (db.VisitorVisitor, error) {
			return db.VisitorVisitor{ID: arg.ID, Data: arg.Data}, nil
		},
	}
	srv := setupTestServer(mock)
	router := SetupRouter(srv)

	body := `{"name":"Updated","email":"updated@test.com","country":"US","city":"NYC"}`
	req := httptest.NewRequest(http.MethodPut, "/api/v1/visitors/a1b2c3d4-e5f6-7890-abcd-ef1234567890", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	resp := decodeResponse(t, w.Body.Bytes())
	assert.True(t, resp.Success)
}

func TestUpdateVisitor_NotFound(t *testing.T) {
	var id pgtype.UUID
	require.NoError(t, id.Scan("00000000-0000-0000-0000-000000000001"))

	mock := &handlerMockStore{
		updateVisitorFn: func(ctx context.Context, arg db.UpdateVisitorParams) (db.VisitorVisitor, error) {
			return db.VisitorVisitor{}, apperrors.ErrNotFound
		},
	}
	srv := setupTestServer(mock)
	router := SetupRouter(srv)

	body := `{"name":"Nope","email":"nope@test.com","country":"US","city":"NYC"}`
	req := httptest.NewRequest(http.MethodPut, "/api/v1/visitors/00000000-0000-0000-0000-000000000001", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestUpdateVisitor_InvalidUUID(t *testing.T) {
	mock := &handlerMockStore{}
	srv := setupTestServer(mock)
	router := SetupRouter(srv)

	body := `{"name":"Test","email":"test@test.com","country":"US","city":"NYC"}`
	req := httptest.NewRequest(http.MethodPut, "/api/v1/visitors/not-a-uuid", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestDeleteVisitor_Success(t *testing.T) {
	var id pgtype.UUID
	require.NoError(t, id.Scan("a1b2c3d4-e5f6-7890-abcd-ef1234567890"))

	mock := &handlerMockStore{
		deleteVisitorFn: func(ctx context.Context, uid pgtype.UUID) (pgtype.UUID, error) {
			return uid, nil
		},
	}
	srv := setupTestServer(mock)
	router := SetupRouter(srv)

	req := httptest.NewRequest(http.MethodDelete, "/api/v1/visitors/a1b2c3d4-e5f6-7890-abcd-ef1234567890", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNoContent, w.Code)
}

func TestDeleteVisitor_NotFound(t *testing.T) {
	var id pgtype.UUID
	require.NoError(t, id.Scan("00000000-0000-0000-0000-000000000001"))

	mock := &handlerMockStore{
		deleteVisitorFn: func(ctx context.Context, uid pgtype.UUID) (pgtype.UUID, error) {
			return pgtype.UUID{}, apperrors.ErrNotFound
		},
	}
	srv := setupTestServer(mock)
	router := SetupRouter(srv)

	req := httptest.NewRequest(http.MethodDelete, "/api/v1/visitors/00000000-0000-0000-0000-000000000001", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestDeleteVisitor_InvalidUUID(t *testing.T) {
	mock := &handlerMockStore{}
	srv := setupTestServer(mock)
	router := SetupRouter(srv)

	req := httptest.NewRequest(http.MethodDelete, "/api/v1/visitors/not-a-uuid", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestHealthCheck(t *testing.T) {
	srv := &Server{}
	router := SetupRouter(srv)

	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	resp := decodeResponse(t, w.Body.Bytes())
	assert.True(t, resp.Success)
}
