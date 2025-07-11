package tests

import (
	"real-estate-system/user-service/models"
	"real-estate-system/user-service/repository"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func setupMockDB(t *testing.T) (*gorm.DB, sqlmock.Sqlmock) {
	sqlDB, mock, err := sqlmock.New()
	assert.NoError(t, err)

	dialector := postgres.New(postgres.Config{
		Conn: sqlDB,
	})
	db, err := gorm.Open(dialector, &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	assert.NoError(t, err)

	return db, mock
}

func TestCreateUser(t *testing.T) {
	db, mock := setupMockDB(t)
	repo := repository.NewGormUserRepository(db)

	user := &models.User{
		Name:      "Alice",
		CreatedAt: 1752216806941602,
		UpdatedAt: 1752216806941602,
	}

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "users" ("name","created_at","updated_at") VALUES ($1,$2,$3) RETURNING "id"`)).
		WithArgs(user.Name, user.CreatedAt, user.UpdatedAt).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	mock.ExpectCommit()

	err := repo.CreateUser(user)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetUser(t *testing.T) {
	db, mock := setupMockDB(t)
	repo := repository.NewGormUserRepository(db)

	rows := sqlmock.NewRows([]string{"id", "name", "created_at", "updated_at"}).
		AddRow(1, "Charlie", 123456, 123456)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users" WHERE "users"."id" = $1 ORDER BY "users"."id" LIMIT $2`)).
		WithArgs(1, 1).
		WillReturnRows(rows)

	user, err := repo.GetUser(1)
	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, "Charlie", user.Name)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetUsers(t *testing.T) {
	db, mock := setupMockDB(t)
	repo := repository.NewGormUserRepository(db)

	rows := sqlmock.NewRows([]string{"id", "name", "created_at", "updated_at"}).
		AddRow(1, "User1", 123456, 123456).
		AddRow(2, "User2", 123456, 123456)

	// GORM may omit OFFSET if it's 0, so we test only the LIMIT
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users" ORDER BY created_at desc LIMIT $1`)).
		WithArgs(10).
		WillReturnRows(rows)

	users, err := repo.GetUsers(1, 10)
	assert.NoError(t, err)
	assert.Len(t, users, 2)
	assert.Equal(t, "User1", users[0].Name)
	assert.Equal(t, "User2", users[1].Name)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetUser_NotFound(t *testing.T) {
	db, mock := setupMockDB(t)
	repo := repository.NewGormUserRepository(db)

	// Simulate user not found
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users" WHERE "users"."id" = $1 ORDER BY "users"."id" LIMIT $2`)).
		WithArgs(999, 1).
		WillReturnError(gorm.ErrRecordNotFound)

	user, err := repo.GetUser(999)
	assert.Error(t, err)
	assert.Nil(t, user)
	assert.Equal(t, gorm.ErrRecordNotFound, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}
