package migrations

import (
	"github.com/goravel/framework/contracts/database/schema"
	"github.com/goravel/framework/facades"
)

type M20250721081947CreateLinkTerkaitTable struct{}

// Signature The unique signature for the migration.
func (r *M20250721081947CreateLinkTerkaitTable) Signature() string {
	return "20250721081947_create_link_terkait_table"
}

// Up Run the migrations.
func (r *M20250721081947CreateLinkTerkaitTable) Up() error {
	if !facades.Schema().HasTable("link_terkait") {
		return facades.Schema().Create("link_terkait", func(table schema.Blueprint) {
			table.ID()
			table.String("keterangan")
			table.String("link")
			table.Enum("status", []any{"Ditampilkan", "Disembunyikan"}).Default("Ditampilkan")
			table.String("image_sm")
			table.String("image_lg")
			table.TimestampsTz()
		})
	}

	return nil
}

// Down Reverse the migrations.
func (r *M20250721081947CreateLinkTerkaitTable) Down() error {
	return facades.Schema().DropIfExists("link_terkait")
}
