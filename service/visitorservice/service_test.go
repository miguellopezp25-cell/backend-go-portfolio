package visitorservice

import (
	"context"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	apperrors "github.com/miguel/go-back-portfolo/pkg/errors"
	"github.com/miguel/go-back-portfolo/schema/db"
)

type mockStore struct {
	db.Store
	createVisitorFn func(ctx context.Context, data []byte) (db.VisitorVisitor, error)
	getVisitorFn    func(ctx context.Context, id pgtype.UUID) (db.VisitorVisitor, error)
	listVisitorsFn  func(ctx context.Context, arg db.ListVisitorsParams) ([]db.VisitorVisitor, error)
	countVisitorsFn func(ctx context.Context) (int64, error)
	updateVisitorFn func(ctx context.Context, arg db.UpdateVisitorParams) (db.VisitorVisitor, error)
	deleteVisitorFn func(ctx context.Context, id pgtype.UUID) (pgtype.UUID, error)
}

func (m *mockStore) CreateVisitor(ctx context.Context, data []byte) (db.VisitorVisitor, error) {
	return m.createVisitorFn(ctx, data)
}

func (m *mockStore) GetVisitor(ctx context.Context, id pgtype.UUID) (db.VisitorVisitor, error) {
	return m.getVisitorFn(ctx, id)
}

func (m *mockStore) ListVisitors(ctx context.Context, arg db.ListVisitorsParams) ([]db.VisitorVisitor, error) {
	return m.listVisitorsFn(ctx, arg)
}

func (m *mockStore) CountVisitors(ctx context.Context) (int64, error) {
	return m.countVisitorsFn(ctx)
}

func (m *mockStore) UpdateVisitor(ctx context.Context, arg db.UpdateVisitorParams) (db.VisitorVisitor, error) {
	return m.updateVisitorFn(ctx, arg)
}

func (m *mockStore) DeleteVisitor(ctx context.Context, id pgtype.UUID) (pgtype.UUID, error) {
	return m.deleteVisitorFn(ctx, id)
}

func TestCreateVisitor_Success(t *testing.T) {
	var id pgtype.UUID
	err := id.Scan("a1b2c3d4-e5f6-7890-abcd-ef1234567890")
	require.NoError(t, err)

	mock := &mockStore{
		createVisitorFn: func(ctx context.Context, data []byte) (db.VisitorVisitor, error) {
			return db.VisitorVisitor{ID: id, Data: data}, nil
		},
	}
	svc := NewService(mock)

	visitor, err := svc.Create(context.Background(), VisitorRequest{
		Name:    "Miguel",
		Email:   "miguel@test.com",
		Country: "Mexico",
		City:    "CDMX",
	})
	require.NoError(t, err)
	require.NotNil(t, visitor)
	assert.Equal(t, "a1b2c3d4-e5f6-7890-abcd-ef1234567890", visitor.ID)
}

func TestGetVisitor_NotFound(t *testing.T) {
	var emptyUID pgtype.UUID

	mock := &mockStore{
		getVisitorFn: func(ctx context.Context, id pgtype.UUID) (db.VisitorVisitor, error) {
			return db.VisitorVisitor{ID: emptyUID}, apperrors.ErrNotFound
		},
	}
	svc := NewService(mock)

	visitor, err := svc.GetByID(context.Background(), "00000000-0000-0000-0000-000000000001")
	require.Error(t, err)
	assert.ErrorIs(t, err, apperrors.ErrNotFound)
	assert.Nil(t, visitor)
}

func TestGetVisitor_Success(t *testing.T) {
	var id pgtype.UUID
	err := id.Scan("a1b2c3d4-e5f6-7890-abcd-ef1234567890")
	require.NoError(t, err)

	mock := &mockStore{
		getVisitorFn: func(ctx context.Context, uid pgtype.UUID) (db.VisitorVisitor, error) {
			return db.VisitorVisitor{
				ID:   uid,
				Data: []byte(`{"name":"Test"}`),
			}, nil
		},
	}
	svc := NewService(mock)

	visitor, err := svc.GetByID(context.Background(), "a1b2c3d4-e5f6-7890-abcd-ef1234567890")
	require.NoError(t, err)
	require.NotNil(t, visitor)
	assert.Equal(t, "a1b2c3d4-e5f6-7890-abcd-ef1234567890", visitor.ID)
	assert.Equal(t, "Test", visitor.Name)
}

func TestListVisitors_Success(t *testing.T) {
	var id1, id2 pgtype.UUID
	_ = id1.Scan("a1b2c3d4-e5f6-7890-abcd-ef1234567890")
	_ = id2.Scan("b2c3d4e5-f6a7-8901-bcde-f12345678901")

	mock := &mockStore{
		countVisitorsFn: func(ctx context.Context) (int64, error) { return 2, nil },
		listVisitorsFn: func(ctx context.Context, arg db.ListVisitorsParams) ([]db.VisitorVisitor, error) {
			return []db.VisitorVisitor{
				{ID: id1, Data: []byte(`{"name":"Alice","email":"alice@test.com","country":"US","city":"NYC"}`)},
				{ID: id2, Data: []byte(`{"name":"Bob","email":"bob@test.com","country":"UK","city":"London"}`)},
			}, nil
		},
	}
	svc := NewService(mock)

	visitors, total, err := svc.List(context.Background(), 1, 20)
	require.NoError(t, err)
	assert.Equal(t, int64(2), total)
	require.Len(t, visitors, 2)
	assert.Equal(t, "a1b2c3d4-e5f6-7890-abcd-ef1234567890", visitors[0].ID)
	assert.Equal(t, "Alice", visitors[0].Name)
	assert.Equal(t, "b2c3d4e5-f6a7-8901-bcde-f12345678901", visitors[1].ID)
	assert.Equal(t, "Bob", visitors[1].Name)
}

func TestListVisitors_Empty(t *testing.T) {
	mock := &mockStore{
		countVisitorsFn: func(ctx context.Context) (int64, error) { return 0, nil },
		listVisitorsFn: func(ctx context.Context, arg db.ListVisitorsParams) ([]db.VisitorVisitor, error) {
			return []db.VisitorVisitor{}, nil
		},
	}
	svc := NewService(mock)

	visitors, total, err := svc.List(context.Background(), 1, 20)
	require.NoError(t, err)
	assert.Equal(t, int64(0), total)
	require.Empty(t, visitors)
}

func TestListVisitors_PaginationDefaults(t *testing.T) {
	mock := &mockStore{
		countVisitorsFn: func(ctx context.Context) (int64, error) { return 0, nil },
		listVisitorsFn: func(ctx context.Context, arg db.ListVisitorsParams) ([]db.VisitorVisitor, error) {
			assert.Equal(t, int32(20), arg.Limit)
			assert.Equal(t, int32(0), arg.Offset)
			return []db.VisitorVisitor{}, nil
		},
	}
	svc := NewService(mock)

	_, _, err := svc.List(context.Background(), 0, 0)
	require.NoError(t, err)
}

func TestUpdateVisitor_Success(t *testing.T) {
	var id pgtype.UUID
	require.NoError(t, id.Scan("a1b2c3d4-e5f6-7890-abcd-ef1234567890"))

	mock := &mockStore{
		updateVisitorFn: func(ctx context.Context, arg db.UpdateVisitorParams) (db.VisitorVisitor, error) {
			return db.VisitorVisitor{ID: arg.ID, Data: arg.Data}, nil
		},
	}
	svc := NewService(mock)

	visitor, err := svc.Update(context.Background(), "a1b2c3d4-e5f6-7890-abcd-ef1234567890", VisitorRequest{
		Name:    "Updated",
		Email:   "updated@test.com",
		Country: "US",
		City:    "NYC",
	})
	require.NoError(t, err)
	require.NotNil(t, visitor)
	assert.Equal(t, "a1b2c3d4-e5f6-7890-abcd-ef1234567890", visitor.ID)
	assert.Equal(t, "Updated", visitor.Name)
}

func TestUpdateVisitor_NotFound(t *testing.T) {
	var id pgtype.UUID
	require.NoError(t, id.Scan("00000000-0000-0000-0000-000000000001"))

	mock := &mockStore{
		updateVisitorFn: func(ctx context.Context, arg db.UpdateVisitorParams) (db.VisitorVisitor, error) {
			return db.VisitorVisitor{}, pgx.ErrNoRows
		},
	}
	svc := NewService(mock)

	visitor, err := svc.Update(context.Background(), "00000000-0000-0000-0000-000000000001", VisitorRequest{
		Name: "Nope", Email: "nope@test.com", Country: "US", City: "NYC",
	})
	require.Error(t, err)
	assert.ErrorIs(t, err, apperrors.ErrNotFound)
	assert.Nil(t, visitor)
}

func TestDeleteVisitor_Success(t *testing.T) {
	var id pgtype.UUID
	require.NoError(t, id.Scan("a1b2c3d4-e5f6-7890-abcd-ef1234567890"))

	mock := &mockStore{
		deleteVisitorFn: func(ctx context.Context, uid pgtype.UUID) (pgtype.UUID, error) {
			return uid, nil
		},
	}
	svc := NewService(mock)

	err := svc.Delete(context.Background(), "a1b2c3d4-e5f6-7890-abcd-ef1234567890")
	require.NoError(t, err)
}

func TestDeleteVisitor_NotFound(t *testing.T) {
	var id pgtype.UUID
	require.NoError(t, id.Scan("00000000-0000-0000-0000-000000000001"))

	mock := &mockStore{
		deleteVisitorFn: func(ctx context.Context, uid pgtype.UUID) (pgtype.UUID, error) {
			return pgtype.UUID{}, pgx.ErrNoRows
		},
	}
	svc := NewService(mock)

	err := svc.Delete(context.Background(), "00000000-0000-0000-0000-000000000001")
	require.Error(t, err)
	assert.ErrorIs(t, err, apperrors.ErrNotFound)
}
