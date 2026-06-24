package visitorservice

import (
	"context"
	"testing"

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
}

func (m *mockStore) CreateVisitor(ctx context.Context, data []byte) (db.VisitorVisitor, error) {
	return m.createVisitorFn(ctx, data)
}

func (m *mockStore) GetVisitor(ctx context.Context, id pgtype.UUID) (db.VisitorVisitor, error) {
	return m.getVisitorFn(ctx, id)
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
