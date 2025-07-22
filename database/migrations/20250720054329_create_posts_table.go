package migrations

import (
	"github.com/goravel/framework/contracts/database/schema"
	"github.com/goravel/framework/facades"
)

type M20250720054329CreatePostsTable struct {
}

// Signature The unique signature for the migration.
func (r *M20250720054329CreatePostsTable) Signature() string {
	return "20250720054329_create_posts_table"
}

// Up Run the migrations.
func (r *M20250720054329CreatePostsTable) Up() error {
	if !facades.Schema().HasTable("posts") {
		return facades.Schema().Create("posts", func(table schema.Blueprint) {
			table.ID()
			table.String("judul")
			table.String("slug")
			table.String("readmore")
			table.LongText("isi")
			table.MediumText("tags")
			table.UnsignedInteger("kategori_id")
			table.UnsignedInteger("user_id")
			table.UnsignedInteger("views").Default(0)
			table.Enum("status", []any{"Ditampilkan", "Disembunyikan"}).Default("Ditampilkan")
			table.String("thumbnail_sm")
			table.String("thumbnail_lg")
			table.String("attachment_og_name")
			table.String("attachment_name")
			table.String("attachment_mime")
			table.String("attachment_size")
			table.TimestampsTz()

			table.Unique("slug")
			table.Index("slug", "kategori_id", "user_id")
		})
	}

	return nil
}

// Down Reverse the migrations.
func (r *M20250720054329CreatePostsTable) Down() error {
	return facades.Schema().DropIfExists("posts")
}
