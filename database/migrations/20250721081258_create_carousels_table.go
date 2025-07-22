package migrations

import (
	"github.com/goravel/framework/contracts/database/schema"
	"github.com/goravel/framework/facades"
)

type M20250721081258CreateCarouselsTable struct{}

// Signature The unique signature for the migration.
func (r *M20250721081258CreateCarouselsTable) Signature() string {
	return "20250721081258_create_carousels_table"
}

// Up Run the migrations.
func (r *M20250721081258CreateCarouselsTable) Up() error {
	if !facades.Schema().HasTable("carousels") {
		return facades.Schema().Create("carousels", func(table schema.Blueprint) {
			table.ID()
			table.String("keterangan")
			table.String("image_sm")
			table.String("image_lg")
			table.String("link")
			table.Enum("status", []any{"Ditampilkan", "Disembunyikan"}).Default("Ditampilkan")
			table.TimestampsTz()
		})
	}

	return nil
}

// Down Reverse the migrations.
func (r *M20250721081258CreateCarouselsTable) Down() error {
	return facades.Schema().DropIfExists("carousels")
}
