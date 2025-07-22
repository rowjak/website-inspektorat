package migrations

import (
	"github.com/goravel/framework/contracts/database/schema"
	"github.com/goravel/framework/facades"
)

type M20250720131432CreateFileEncryptTable struct {
}

// Signature The unique signature for the migration.
func (r *M20250720131432CreateFileEncryptTable) Signature() string {
	return "20250720131432_create_file_encrypt_table"
}

// Up Run the migrations.
func (r *M20250720131432CreateFileEncryptTable) Up() error {
	if !facades.Schema().HasTable("file_encrypt") {
		return facades.Schema().Create("file_encrypt", func(table schema.Blueprint) {
			table.ID()
			table.String("file_og_name")
			table.String("file_encrypt_name")
			table.String("file_path")
			table.String("file_mime")
			table.Integer("file_size").Unsigned()
			table.TimestampsTz()
		})
	}

	return nil
}

// Down Reverse the migrations.
func (r *M20250720131432CreateFileEncryptTable) Down() error {
	return facades.Schema().DropIfExists("file_encrypt")
}
