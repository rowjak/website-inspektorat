package database

import (
	"rowjak/website-inspektorat/database/migrations"
	"rowjak/website-inspektorat/database/seeders"

	"github.com/goravel/framework/contracts/database/schema"
	"github.com/goravel/framework/contracts/database/seeder"
)

type Kernel struct {
}

func (kernel Kernel) Migrations() []schema.Migration {
	return []schema.Migration{
		&migrations.M20210101000001CreateUsersTable{},
		&migrations.M20210101000002CreateJobsTable{},
		&migrations.M20250720054329CreatePostsTable{},
		&migrations.M20250720054341CreatePostImagesTable{},
		&migrations.M20250720131432CreateFileEncryptTable{},
		&migrations.M20250721081228CreatePostCategoryTable{},
		&migrations.M20250721081258CreateCarouselsTable{},
		&migrations.M20250721081409CreatePeraturanTable{},
		&migrations.M20250721081832CreateUnduhanTable{},
		&migrations.M20250721081850CreateTagsTable{},
		&migrations.M20250721081947CreateLinkTerkaitTable{},
	}
}
func (kernel Kernel) Seeders() []seeder.Seeder {
	return []seeder.Seeder{
		&seeders.DatabaseSeeder{},
		&seeders.UserSeeder{},
	}
}
