package scribo_test

import (
	"encoding/json"
	"time"

	. "github.com/bbengfort/scribo/scribo"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

// Note that the any model methods that require the database are in the
// db_test module that has the database connection runner, etc.
var _ = Describe("Models", func() {

	Describe("Nodes", func() {

		It("should not serialize the Key field", func() {

			node := Node{
				ID:      1,
				Name:    "apollo",
				Address: "108.51.64.223",
				DNS:     "bryant.bengfort.com",
				Key:     "werxhqb98rpaxn39848xrunpaw3489ruxnpa98w4rxn",
				Created: time.Now(),
				Updated: time.Now(),
			}

			data, err := json.Marshal(node)
			Ω(err).Should(BeNil())

			var obj map[string]*json.RawMessage
			err = json.Unmarshal(data, &obj)
			Ω(err).Should(BeNil())

			Ω(obj).ShouldNot(HaveKey("key"))
			Ω(obj).ShouldNot(HaveKey("Key"))
			Ω(obj).ShouldNot(HaveKey("KEY"))
		})

	})

})
