package migrations

import (
	"github.com/goravel/framework/contracts/database/schema"
	"github.com/goravel/framework/facades"
)

type M20250721081228CreatePostCategoryTable struct{}

// Signature The unique signature for the migration.
func (r *M20250721081228CreatePostCategoryTable) Signature() string {
	return "20250721081228_create_post_category_table"
}

// Up Run the migrations.
func (r *M20250721081228CreatePostCategoryTable) Up() error {
	if !facades.Schema().HasTable("post_category") {
		return facades.Schema().Create("post_category", func(table schema.Blueprint) {
			table.ID()
			table.String("kategori")
			table.String("slug")
			table.TimestampsTz()
		})
	}

	return nil
}

// Down Reverse the migrations.
func (r *M20250721081228CreatePostCategoryTable) Down() error {
	return facades.Schema().DropIfExists("post_category")
}
