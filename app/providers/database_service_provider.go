package providers

import (
	"github.com/goravel/framework/contracts/foundation"
	"github.com/goravel/framework/facades"

	"rowjak/website-inspektorat/database"
)

type DatabaseServiceProvider struct {
}

func (receiver *DatabaseServiceProvider) Register(app foundation.Application) {

}

func (receiver *DatabaseServiceProvider) Boot(app foundation.Application) {
	kernel := database.Kernel{}
	facades.Schema().Register(kernel.Migrations())
	facades.Seeder().Register(kernel.Seeders())
}
