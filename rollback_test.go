package migrations_test

import (
	"context"
	"errors"
	"os"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/reecerussell/migrations"
	"github.com/reecerussell/migrations/mock"
)

func TestRollback_GivenAppliedMigration_ReturnsNoError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	testContent := "My Migration Content"
	testCtx := context.Background()
	testMigration := &migrations.Migration{
		Name:     "MyMigration",
		DownFile: "MyFile",
	}

	mockProvider := mock.NewMockProvider(ctrl)
	mockProvider.EXPECT().GetAppliedMigrations(testCtx).Return([]*migrations.Migration{testMigration}, nil)
	mockProvider.EXPECT().Rollback(testCtx, "MyMigration", testContent).Return(nil)

	mockFileReader := mock.NewMockFileReader(ctrl)
	mockFileReader.EXPECT().Read("MyFile").Return(testContent, nil)

	err := migrations.Rollback(testCtx, []*migrations.Migration{testMigration}, mockProvider, mockFileReader, "")
	assert.NoError(t, err)
}

func TestRollback_GivenUnappliedMigration_SkipsMigration(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	testCtx := context.Background()
	testMigration := &migrations.Migration{}

	mockProvider := mock.NewMockProvider(ctrl)
	mockProvider.EXPECT().GetAppliedMigrations(testCtx).Return([]*migrations.Migration{}, nil)

	err := migrations.Rollback(testCtx, []*migrations.Migration{testMigration}, mockProvider, nil, "")
	assert.NoError(t, err)
}

func TestRollback_GivenMigrationWithMissingFile_ReturnsIsNotExists(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	testCtx := context.Background()
	testMigration := &migrations.Migration{
		Name:     "MyMigration",
		DownFile: "MyFile",
	}

	mockProvider := mock.NewMockProvider(ctrl)
	mockProvider.EXPECT().GetAppliedMigrations(testCtx).Return([]*migrations.Migration{testMigration}, nil)

	mockFileReader := mock.NewMockFileReader(ctrl)
	mockFileReader.EXPECT().Read("MyFile").Return("", os.ErrNotExist)

	err := migrations.Rollback(testCtx, []*migrations.Migration{testMigration}, mockProvider, mockFileReader, "")
	assert.True(t, os.IsNotExist(err))
}

func TestRollback_WhereRollbackFails_ReturnsError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	testContent := "My Migration Content"
	testCtx := context.Background()
	testMigration := &migrations.Migration{
		Name:     "MyMigration",
		DownFile: "MyFile",
	}
	testError := errors.New("an error occured")

	mockProvider := mock.NewMockProvider(ctrl)
	mockProvider.EXPECT().GetAppliedMigrations(testCtx).Return([]*migrations.Migration{testMigration}, nil)
	mockProvider.EXPECT().Rollback(testCtx, "MyMigration", testContent).Return(testError)

	mockFileReader := mock.NewMockFileReader(ctrl)
	mockFileReader.EXPECT().Read("MyFile").Return(testContent, nil)

	err := migrations.Rollback(testCtx, []*migrations.Migration{testMigration}, mockProvider, mockFileReader, "")
	assert.Equal(t, testError, err)
}

func TestRollback_RollbackTargetMigration_SkipsFurtherMigrations(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	testContent := "My Migration Content"
	testCtx := context.Background()
	testMigrations := []*migrations.Migration{
		&migrations.Migration{Name: "One"},
		&migrations.Migration{Name: "Two", DownFile: "MyFile"},
	}

	mockProvider := mock.NewMockProvider(ctrl)
	mockProvider.EXPECT().GetAppliedMigrations(testCtx).Return(testMigrations, nil)
	mockProvider.EXPECT().Rollback(testCtx, "Two", testContent).Return(nil)

	mockFileReader := mock.NewMockFileReader(ctrl)
	mockFileReader.EXPECT().Read("MyFile").Return(testContent, nil)

	err := migrations.Rollback(testCtx, testMigrations, mockProvider, mockFileReader, "Two")
	assert.NoError(t, err)
}

func TestRollback_FailsToGetAppliedMigrations_ReturnsError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	testCtx := context.Background()
	testError := errors.New("an error occured")

	mockProvider := mock.NewMockProvider(ctrl)
	mockProvider.EXPECT().GetAppliedMigrations(testCtx).Return(nil, testError)

	err := migrations.Rollback(testCtx, nil, mockProvider, nil, "")
	assert.Equal(t, testError, err)
}
