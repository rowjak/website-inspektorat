package providers

import (
	"rowjak/website-inspektorat/app/helpers"

	"github.com/goravel/framework/contracts/foundation"
	"github.com/goravel/framework/facades"
)

type AppServiceProvider struct {
}

func (receiver *AppServiceProvider) Register(app foundation.Application) {

}

func (receiver *AppServiceProvider) Boot(app foundation.Application) {
	facades.View().Share("h", helpers.ViewHelpers)
}
