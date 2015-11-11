package db_test

import (
	db "github.com/fiatjaf/summadb/database"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Array", func() {
	Context("put or delete entire arrays", func() {
		It("should save a simple array", func() {
			Expect(db.SaveArrayAt("/objects/home/wood", []interface{}{
				`"chair"`,
				`"table"`,
				`"true"`,
			})).To(Succeed())
			Expect(db.GetValueAt("/objects/home/wood/_array[0]")).To(BeEquivalentTo(`"chair"`))
			Expect(db.GetValueAt("/objects/home/wood/_array[1]")).To(BeEquivalentTo(`"table"`))
			Expect(db.GetValueAt("/objects/home/wood/_array[2]")).To(BeEquivalentTo("true"))
		})

		It("should save an array in a tree operation", func() {
			Expect(db.SaveTreeAt("/objects/outside", map[string]interface{}{
				"steel": []map[string]interface{}{
					map[string]interface{}{
						"name":        "litterbin",
						"state-owned": true,
					},
					map[string]interface{}{
						"name":     "lamppost",
						"quantity": 10,
					},
				},
			})).To(Succeed())
			Expect(db.GetValueAt("/objects/outside/steel/_array[0]/name")).To(BeEquivalentTo(`"litterbin"`))
			Expect(db.GetTreeAt("/objects/outside")).To(Equal(map[string]interface{}{
				"steel": []map[string]interface{}{
					map[string]interface{}{
						"name":        value("litterbin"),
						"state-owned": value(true),
					},
					map[string]interface{}{
						"name":     value("lamppost"),
						"quantity": value(10),
					},
				},
			}))
		})
	})
})
