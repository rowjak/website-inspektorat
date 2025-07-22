package migrations

import (
	"github.com/goravel/framework/contracts/database/schema"
	"github.com/goravel/framework/facades"
)

type M20250721081409CreatePeraturanTable struct{}

// Signature The unique signature for the migration.
func (r *M20250721081409CreatePeraturanTable) Signature() string {
	return "20250721081409_create_peraturan_table"
}

// Up Run the migrations.
func (r *M20250721081409CreatePeraturanTable) Up() error {
	if !facades.Schema().HasTable("peraturan") {
		return facades.Schema().Create("peraturan", func(table schema.Blueprint) {
			table.ID()
			table.Enum("jenis", []any{
				"Peraturan Bupati", "Keputusan Bupati", "Undang-Undang", "Peraturan Presiden", "Peraturan Pemerintah"}).
				Default("Peraturan Bupati")
			table.String("nomor")
			table.String("tentang")
			table.Date("tanggal")
			table.UnsignedInteger("tahun")
			table.String("file_og_name")
			table.String("file_name")
			table.String("file_mime")
			table.String("file_size")
			table.String("link_jdih")
			table.TimestampsTz()
		})
	}

	return nil
}

// Down Reverse the migrations.
func (r *M20250721081409CreatePeraturanTable) Down() error {
	return facades.Schema().DropIfExists("peraturan")
}
