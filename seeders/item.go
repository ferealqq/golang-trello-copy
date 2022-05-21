package seeders

import (
	"log"

	faker "github.com/bxcodec/faker/v3"
	"github.com/ferealqq/wienerlist/boardapi/models"
	"github.com/ferealqq/wienerlist/pkg/seed"
	"gorm.io/gorm"
)

func CreateItem(db *gorm.DB, title string, desc string, workspace uint, section uint) *seed.SeedOut[models.Item] {
	return seed.SeedModel(db, models.Item{
		Title:       title,
		Description: desc,
		WorkspaceId: workspace,
		SectionId:   section,
	})
}

func CreateItemFaker(db *gorm.DB) *seed.SeedOut[models.Item] {
	section := CreateSectionFaker(db).Model.ID
	var workspace models.Workspace
	if err := db.First(&workspace).Error; err != nil {
		workspace = CreateWorkspaceFaker(db).Model
	}

	return CreateItem(
		db,
		faker.Word(),
		faker.Sentence(),
		workspace.ID,
		section,
	)
}

// funky ass util function for seeding items
func fiveItemsFor(db *gorm.DB, sectionId uint, workspaceId uint) ([]models.Item, []error) {
	var section models.Section
	var workspace models.Workspace
	var items []models.Item
	var errors []error

	if err := db.First(&workspace, workspaceId).Error; err != nil {
		workspace = CreateWorkspaceFaker(db).Model
	}
	if err := db.First(&section, sectionId).Error; err != nil {
		section = CreateSectionFaker(db).Model
	}
	for i := 0; i < 5; i++ {
		r := CreateItem(db, faker.Word(), faker.Sentence(), workspace.ID, section.ID)
		items = append(items, r.Model)
		errors = append(errors, r.Error)
	}

	return items, errors
}

func ItemAll() []seed.Seed[[]models.Item] {
	return []seed.Seed[[]models.Item]{
		{
			Name: "Five items for section 1",
			Run: func(db *gorm.DB) *seed.SeedOut[[]models.Item] {
				items, errs := fiveItemsFor(db, 1, 1)
				var err error
				for _, e := range errs {
					if e != nil {
						err = e
					}
				}
				return &seed.SeedOut[[]models.Item]{
					Model: items,
					Error: err,
				}
			},
		},
		{
			Name: "Five items for section 2",
			Run: func(db *gorm.DB) *seed.SeedOut[[]models.Item] {
				items, errs := fiveItemsFor(db, 2, 1)
				var err error
				for _, e := range errs {
					if e != nil {
						err = e
					}
				}
				return &seed.SeedOut[[]models.Item]{
					Model: items,
					Error: err,
				}
			},
		},
		{
			Name: "Five items for section 3",
			Run: func(db *gorm.DB) *seed.SeedOut[[]models.Item] {
				items, errs := fiveItemsFor(db, 3, 1)
				var err error
				for _, e := range errs {
					if e != nil {
						err = e
					}
				}
				return &seed.SeedOut[[]models.Item]{
					Model: items,
					Error: err,
				}
			},
		},
	}
}

func SeedItems(db *gorm.DB) {
	for _, seed := range ItemAll() {
		if result := seed.Run(db); result.Error != nil {
			log.Fatalf("Running seed '%s', failed with error: %s", seed.Name, result.Error)
		}
	}
}
