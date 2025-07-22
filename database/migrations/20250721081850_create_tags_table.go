package migrations

import (
	"github.com/goravel/framework/contracts/database/schema"
	"github.com/goravel/framework/facades"
)

type M20250721081850CreateTagsTable struct{}

// Signature The unique signature for the migration.
func (r *M20250721081850CreateTagsTable) Signature() string {
	return "20250721081850_create_tags_table"
}

// Up Run the migrations.
func (r *M20250721081850CreateTagsTable) Up() error {
	if !facades.Schema().HasTable("tags") {
		return facades.Schema().Create("tags", func(table schema.Blueprint) {
			table.ID()
			table.String("slug")
			table.String("nama")
			table.TimestampsTz()
		})
	}

	return nil
}

// Down Reverse the migrations.
func (r *M20250721081850CreateTagsTable) Down() error {
	return facades.Schema().DropIfExists("tags")
}
