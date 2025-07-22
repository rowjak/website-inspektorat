package migrations

import (
	"github.com/goravel/framework/contracts/database/schema"
	"github.com/goravel/framework/facades"
)

type M20250720054341CreatePostImagesTable struct {
}

// Signature The unique signature for the migration.
func (r *M20250720054341CreatePostImagesTable) Signature() string {
	return "20250720054341_create_post_images_table"
}

// Up Run the migrations.
func (r *M20250720054341CreatePostImagesTable) Up() error {
	if !facades.Schema().HasTable("post_images") {
		return facades.Schema().Create("post_images", func(table schema.Blueprint) {
			table.ID()
			table.UnsignedBigInteger("post_id")
			table.String("image_sm")
			table.String("image_lg")
			table.TimestampsTz()

			table.Index("post_id")
		})
	}

	return nil
}

// Down Reverse the migrations.
func (r *M20250720054341CreatePostImagesTable) Down() error {
	return facades.Schema().DropIfExists("post_images")
}
