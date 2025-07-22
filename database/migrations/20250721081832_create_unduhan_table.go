package migrations

import (
	"github.com/goravel/framework/contracts/database/schema"
	"github.com/goravel/framework/facades"
)

type M20250721081832CreateUnduhanTable struct{}

// Signature The unique signature for the migration.
func (r *M20250721081832CreateUnduhanTable) Signature() string {
	return "20250721081832_create_unduhan_table"
}

// Up Run the migrations.
func (r *M20250721081832CreateUnduhanTable) Up() error {
	if !facades.Schema().HasTable("unduhan") {
		return facades.Schema().Create("unduhan", func(table schema.Blueprint) {
			table.ID()
			table.Enum("jenis", []any{"Renstra", "LKjIP", "IKU", "Renja", "Perjanjian Kinerja"}).Default("Renstra")
			table.String("file_og_name")
			table.String("file_name")
			table.String("file_mime")
			table.String("file_size")
			table.Enum("status", []any{"Ditampilkan", "Disembunyikan"}).Default("Ditampilkan")
			table.UnsignedInteger("downloaded_count").Default(0)
			table.TimestampsTz()
		})
	}

	return nil
}

// Down Reverse the migrations.
func (r *M20250721081832CreateUnduhanTable) Down() error {
	return facades.Schema().DropIfExists("unduhan")
}
