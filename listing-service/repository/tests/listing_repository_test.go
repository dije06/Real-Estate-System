package tests

import (
	"real-estate-system/listing-service/models"
	"real-estate-system/listing-service/repository"
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

	dialector := postgres.New(postgres.Config{Conn: sqlDB})
	db, err := gorm.Open(dialector, &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	assert.NoError(t, err)

	return db, mock
}

func TestCreateListing(t *testing.T) {
	db, mock := setupMockDB(t)
	repo := repository.NewGormListingRepository(db)

	listing := &models.Listing{
		UserID:      1,
		Price:       500000,
		ListingType: "rent",
		CreatedAt:   123456789,
		UpdatedAt:   123456789,
	}

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(
		`INSERT INTO "listings" ("user_id","price","listing_type","created_at","updated_at") VALUES ($1,$2,$3,$4,$5) RETURNING "id"`)).
		WithArgs(listing.UserID, listing.Price, listing.ListingType, listing.CreatedAt, listing.UpdatedAt).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	mock.ExpectCommit()

	err := repo.CreateListing(listing)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetListings(t *testing.T) {
	db, mock := setupMockDB(t)
	repo := repository.NewGormListingRepository(db)

	rows := sqlmock.NewRows([]string{"id", "user_id", "price", "listing_type", "created_at", "updated_at"}).
		AddRow(1, 1, 100000, "sale", 123, 123).
		AddRow(2, 2, 200000, "rent", 123, 123)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "listings" ORDER BY created_at desc LIMIT $1`)).
		WithArgs(2).
		WillReturnRows(rows)

	listings, err := repo.GetListings(1, 2)
	assert.NoError(t, err)
	assert.Len(t, listings, 2)
	assert.Equal(t, "sale", listings[0].ListingType)
	assert.Equal(t, "rent", listings[1].ListingType)
	assert.NoError(t, mock.ExpectationsWereMet())
}
