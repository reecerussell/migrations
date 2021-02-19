package migrations_test

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/reecerussell/migrations"
	"github.com/reecerussell/migrations/mock"
)

func TestRollback_GivenAppliedMigration_ReturnsNoError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	testCtx := context.Background()
	testMigration := &migrations.Migration{}

	mockProvider := mock.NewMockProvider(ctrl)
	mockProvider.EXPECT().GetAppliedMigrations(testCtx).Return([]*migrations.Migration{testMigration}, nil)
	mockProvider.EXPECT().Rollback(testCtx, testMigration).Return(nil)

	err := migrations.Rollback(testCtx, []*migrations.Migration{testMigration}, mockProvider, "")
	assert.NoError(t, err)
}

func TestRollback_GivenUnappliedMigration_SkipsMigration(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	testCtx := context.Background()
	testMigration := &migrations.Migration{}

	mockProvider := mock.NewMockProvider(ctrl)
	mockProvider.EXPECT().GetAppliedMigrations(testCtx).Return([]*migrations.Migration{}, nil)

	err := migrations.Rollback(testCtx, []*migrations.Migration{testMigration}, mockProvider, "")
	assert.NoError(t, err)
}

func TestRollback_WhereRollbackFails_ReturnsError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	testCtx := context.Background()
	testMigration := &migrations.Migration{}
	testError := errors.New("an error occured")

	mockProvider := mock.NewMockProvider(ctrl)
	mockProvider.EXPECT().GetAppliedMigrations(testCtx).Return([]*migrations.Migration{testMigration}, nil)
	mockProvider.EXPECT().Rollback(testCtx, testMigration).Return(testError)

	err := migrations.Rollback(testCtx, []*migrations.Migration{testMigration}, mockProvider, "")
	assert.Equal(t, testError, err)
}

func TestRollback_RollbackTargetMigration_SkipsFurtherMigrations(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	testCtx := context.Background()
	testMigrations := []*migrations.Migration{
		&migrations.Migration{Name: "One"},
		&migrations.Migration{Name: "Two"},
	}

	mockProvider := mock.NewMockProvider(ctrl)
	mockProvider.EXPECT().GetAppliedMigrations(testCtx).Return(testMigrations, nil)
	mockProvider.EXPECT().Rollback(testCtx, testMigrations[1]).Return(nil)

	err := migrations.Rollback(testCtx, testMigrations, mockProvider, "Two")
	assert.NoError(t, err)
}

func TestRollback_FailsToGetAppliedMigrations_ReturnsError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	testCtx := context.Background()
	testError := errors.New("an error occured")

	mockProvider := mock.NewMockProvider(ctrl)
	mockProvider.EXPECT().GetAppliedMigrations(testCtx).Return(nil, testError)

	err := migrations.Rollback(testCtx, nil, mockProvider, "")
	assert.Equal(t, testError, err)
}
