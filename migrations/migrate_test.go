package migrations

import (
	"testing"

	"github.com/ferealqq/golang-trello-copy/server/pkg/database"
	"github.com/stretchr/testify/assert"
)

func TestMigrate(t *testing.T) {
	database.TestDBInit()
	MigrateSeedAfterwards(database.DBConn)
	defer database.Close()

	assert.True(t, true, "should not panic")
}
