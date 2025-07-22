package tests

import (
	"github.com/goravel/framework/testing"

	"rowjak/website-inspektorat/bootstrap"
)

func init() {
	bootstrap.Boot()
}

type TestCase struct {
	testing.TestCase
}
