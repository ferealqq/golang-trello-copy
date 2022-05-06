package controllers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"

	. "github.com/ferealqq/golang-trello-copy/server/boardapi/models"
	app "github.com/ferealqq/golang-trello-copy/server/pkg/appenv"
	ctrl "github.com/ferealqq/golang-trello-copy/server/pkg/controller"
	"github.com/ferealqq/golang-trello-copy/server/pkg/database"
	. "github.com/ferealqq/golang-trello-copy/server/pkg/testing"
	. "github.com/ferealqq/golang-trello-copy/server/seeders"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestListItemsHandler(t *testing.T) {
	action := HttpTestAction[Board]{
		Method: http.MethodGet,
		RouterFunc: func(e *gin.Engine, ae app.AppEnv) {
			e.GET("/items", ctrl.MakeHandler(ae, ListItemsHandler))
		},
		ReqPath: "/items",
		Seeders: []func(db *gorm.DB){SeedItems},
		Tables:  []string{"items"},
	}
	response := action.Run()
	assert.Equal(t, http.StatusOK, response.Code, "they should be equal")
	var d map[string]interface{}

	if err := json.Unmarshal(response.Body.Bytes(), &d); err != nil {
		assert.Fail(t, "Unmarshal should not fail")
		return
	}
	var count int64
	err := database.DBConn.Model(&Item{}).Count(&count)
	if err.Error != nil {
		assert.Fail(t, "Fetching count should not fail")
	}
	assert.Equal(t, int(count), int(d["count"].(float64)), "they should be equal")
}

/**
func TestListSectionsHandlerFromBoard(t *testing.T) {
	action := HttpTestAction[Board]{
		Method: http.MethodGet,
		RouterFunc: func(e *gin.Engine, ae app.AppEnv) {
			e.GET("/sections", ctrl.MakeHandler(ae, ListSectionsHandler))
		},
		Seeders: []func(db *gorm.DB){SeedSections},
		Tables:  []string{},
	}
	w := CreateWorkspaceFaker(database.DBConn).Model
	b1 := Board{
		Title:       "Only this boards sections will be listed",
		Description: "This is a test board",
		WorkspaceId: w.ID,
	}
	database.DBConn.Create(&b1)
	CreateSection(database.DBConn, "Section one of results wanted", "this is a test section", b1.ID)
	CreateSection(database.DBConn, "Section two of results wanted", "this is a test section", b1.ID)
	b2 := Board{
		Title:       "Only this boards sections will be listed",
		Description: "This is a test board",
		WorkspaceId: w.ID,
	}
	database.DBConn.Create(&b2)
	CreateSection(database.DBConn, "Section one of results wanted", "this is a test section", b2.ID)

	// Query with multiple board ids
	action.ReqPath = "/sections?BoardId=" +
		strconv.FormatUint(uint64(b1.ID), 10) + "&BoardId=" +
		strconv.FormatUint(uint64(b2.ID), 10)

	response := action.Run()
	assert.Equal(t, http.StatusOK, response.Code, "they should be equal")
	var d map[string]interface{}

	if err := json.Unmarshal(response.Body.Bytes(), &d); err != nil {
		assert.Fail(t, "Unmarshal should not fail")
		return
	}
	assert.Equal(t, 3, int(d["count"].(float64)), "they should be equal")
}
*/

func TestPatchItemHandler(t *testing.T) {
	// BoardId 1 is created by SeedSections function, by default if there is no board in the database SeedSections will create a new board
	i := Item{
		SectionId: 7,
	}
	b, _ := json.Marshal(i)

	action := HttpTestAction[Section]{
		Method: http.MethodPatch,
		RouterFunc: func(e *gin.Engine, ae app.AppEnv) {
			e.PATCH("/items/:id", ctrl.MakeHandler(ae, UpdateItemHandler))
		},
		ReqPath: "/items/1",
		Body:    bytes.NewReader(b),
		Seeders: []func(db *gorm.DB){SeedWorkspaces, SeedSections, SeedItems, func(db *gorm.DB) {
			// generate lots of random sections
			for i := 0; i < 5; i++ {
				CreateSectionFaker(db)
			}
		}},
		Tables: []string{"workspaces", "sections", "items"},
	}

	response := action.Run()

	assert.Equal(t, http.StatusOK, response.Code, "they should be equal")

	var rItem map[string]interface{}
	if err := json.Unmarshal(response.Body.Bytes(), &rItem); err != nil {
		assert.Fail(t, "Unmarshal should not fail")
		return
	}
	var item Item
	assert.Nil(t, database.DBConn.First(&item, 1).Error, "item should be found")
	assert.Equal(t, item.SectionId, uint(7), "they should be equal")
	assert.NotNil(t, item.Title, "should not be nil")
	assert.NotNil(t, item.Description, "should not be nil")
	assert.NotNil(t, item.WorkspaceId, "should not be nil")
	assert.NotNil(t, item.UpdatedAt, "should not be nil")
}
