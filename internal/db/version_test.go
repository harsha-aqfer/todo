package db

import (
	sqlMock "github.com/DATA-DOG/go-sqlmock"
	asserts "github.com/stretchr/testify/assert"
	"regexp"
	"testing"
)

func Test_Db_GetVersion(t *testing.T) {
	assert := asserts.New(t)
	db, mock, err := sqlMock.New()

	if assert.Nil(err) {
		defer func() {
			_ = db.Close()
		}()

		row := sqlMock.NewRows([]string{"version"}).AddRow("1.5")

		mock.ExpectQuery(regexp.QuoteMeta("SELECT VERSION()")).WillReturnRows(row)

		s := &DB{Sql: db, Version: NewVersionStore(db)}
		output, err := s.Version.GetVersion()

		if assert.Nil(err) {
			assert.Equal("1.5", output)
		}
	}
}
