package db_test

import (
	"testing"

	db "github.com/fiatjaf/summadb/database"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestDatabase(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "CRUD Suite")
}

var _ = BeforeSuite(func() {
	Expect(db.Erase()).To(Succeed())
})

func value(v interface{}) map[string]interface{} {
	return map[string]interface{}{"_val": v}
}
