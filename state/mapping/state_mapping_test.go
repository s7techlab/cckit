package mapping_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"errors"

	"github.com/s7techlab/cckit/state"
	m "github.com/s7techlab/cckit/state/mapping"
	"github.com/s7techlab/cckit/state/mapping/testdata"
	"github.com/s7techlab/cckit/state/mapping/testdata/schema"
)

var _ = Describe(`State mappings`, func() {

	mappings := m.StateMappings{}.
		//key will be <`EntityWithComplexId`, {Id.IdPart1}, {Id.IdPart2} >
		Add(&schema.EntityWithComplexId{}, m.PKeyComplexId(&schema.EntityComplexId{}))

	It(`Got error if namespace not exists`, func() {
		_, err := mappings.GetByNamespace(state.Key{`this-namespace-not-exists`})
		Expect(errors.Is(err, m.ErrStateMappingNotFound)).To(BeTrue())

		_, err = mappings.Get([]string{`this-namespace-not-exists`})
		Expect(errors.Is(err, m.ErrStateMappingNotFound)).To(BeTrue())
	})

	It(`Allow to get mapping by namespace`, func() {
		mapping, err := mappings.GetByNamespace(state.Key{`EntityWithComplexId`})
		Expect(err).NotTo(HaveOccurred())

		Expect(mapping.Namespace()).To(Equal(state.Key{`EntityWithComplexId`}))

		mapping, err = mappings.Get([]string{`EntityWithComplexId`})
		Expect(err).NotTo(HaveOccurred())
		Expect(mapping.Namespace()).To(Equal(state.Key{`EntityWithComplexId`}))
	})

	It(`Allow to сreate primary key`, func() {
		mapping, _ := mappings.Get(&schema.EntityWithComplexId{})
		key, err := mapping.PrimaryKey(testdata.CreateEntityWithComplextId[0])

		Expect(err).NotTo(HaveOccurred())
		Expect(key).To(Equal(state.Key{`EntityWithComplexId`, `aaa`, `bb`, `ccc`, `2020-01-28`}))
	})
})
