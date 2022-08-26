package debug_test

import (
	"strconv"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/s7techlab/cckit/extensions/debug"
	testcc "github.com/s7techlab/cckit/testing"
	. "github.com/s7techlab/cckit/testing/gomega"
)

var _ = Describe(`State`, func() {

	var (
		dbg = debug.NewStateService()

		cc, ctx = testcc.NewTxHandler(`Debug`)
	)

	Context(`Empty`, func() {

		It("Allow to get empty key list", func() {
			cc.From(Owner).Tx(func() {
				keys, err := dbg.ListKeys(ctx, nil) // get all keys
				Expect(err).NotTo(HaveOccurred())

				Expect(keys.Keys).To(HaveLen(0))
			})
		})

		It("Allow to clean", func() {
			cc.From(Owner).Tx(func() {
				deleted, err := dbg.DeleteStates(ctx, &debug.Prefixes{Prefixes: []*debug.Prefix{{
					Key: []string{`some`},
				}, {
					Key: []string{`not exists prefix`},
				}}}) // get all keys
				Expect(err).NotTo(HaveOccurred())

				Expect(len(deleted.Matches)).To(Equal(2))
				Expect(deleted.Matches[`some`]).To(Equal(uint32(0)))
				Expect(deleted.Matches[`not exists prefix`]).To(Equal(uint32(0)))
			})

		})

	})

	Context(`Non Empty`, func() {

		It("Allow put value in state", func() {
			for i := 0; i < 5; i++ {
				cc.From(Owner).Tx(func() {
					val := &debug.Value{
						Key:   []string{`prefixA`, `key` + strconv.Itoa(i)},
						Value: []byte(`value` + strconv.Itoa(i)),
					}
					valReceived, err := dbg.PutState(ctx, val)
					Expect(err).NotTo(HaveOccurred())
					Expect(valReceived).To(StringerEqual(val))
				})

			}

			for i := 0; i < 7; i++ {
				cc.From(Owner).Tx(func() {
					_, err := dbg.PutState(ctx, &debug.Value{
						Key:   []string{`prefixB`, `subprefixA`, `key` + strconv.Itoa(i)},
						Value: []byte(`value` + strconv.Itoa(i)),
					})
					Expect(err).NotTo(HaveOccurred())
				})

				cc.From(Owner).Tx(func() {
					_, err := dbg.PutState(ctx, &debug.Value{
						Key:   []string{`prefixB`, `subprefixB`, `key` + strconv.Itoa(i)},
						Value: []byte(`value` + strconv.Itoa(i)),
					})
					Expect(err).NotTo(HaveOccurred())
				})
			}

			// put proto message
			cc.From(Owner).Tx(func() {
				_, err := dbg.PutState(ctx, &debug.Value{
					Key:   []string{`keyA`},
					Value: testcc.MustProtoMarshal(&debug.Prefix{Key: []string{`keyA`}}),
				})
				Expect(err).NotTo(HaveOccurred())
			})
			//cc.From(Owner).Invoke(`debugStatePut`, []string{`keyB`}, []byte(`valueKeyB`))
			//cc.From(Owner).Invoke(`debugStatePut`, []string{`keyC`}, []byte(`valueKeyC`))
		})

		It("Allow to get value from state", func() {
			cc.From(Owner).Tx(func() {
				val, err := dbg.GetState(ctx, &debug.CompositeKey{
					Key: []string{`prefixB`, `subprefixA`, `key` + strconv.Itoa(1)},
				})
				Expect(err).NotTo(HaveOccurred())
				Expect(val.Value).To(Equal([]byte(`value1`)))
			})

			cc.From(Owner).Tx(func() {
				val, err := dbg.GetState(ctx, &debug.CompositeKey{
					Key: []string{`keyA`},
				})
				Expect(err).NotTo(HaveOccurred())
				Expect(val.Value).To(Equal(testcc.MustProtoMarshal(&debug.Prefix{Key: []string{`keyA`}})))
				Expect(val.Json).To(Equal(``)) // we have no mapping in state for this entry
			})
		})

		It("Allow to get keys", func() {
			cc.From(Owner).Tx(func() {
				keys, err := dbg.ListKeys(ctx, &debug.Prefix{
					Key: []string{`prefixA`},
				})
				Expect(err).NotTo(HaveOccurred())
				Expect(keys.Keys).To(HaveLen(5))

				for i := range keys.Keys {
					Expect(keys.Keys[i].Key[0]).To(Equal(`prefixA`))
				}
			})

			cc.From(Owner).Tx(func() {
				keys, err := dbg.ListKeys(ctx, &debug.Prefix{
					Key: []string{`prefixB`, `subprefixB`},
				})
				Expect(err).NotTo(HaveOccurred())
				Expect(keys.Keys).To(HaveLen(7))

				for i := range keys.Keys {
					Expect(keys.Keys[i].Key[0]).To(Equal(`prefixB`))
					Expect(keys.Keys[i].Key[1]).To(Equal(`subprefixB`))
				}
			})
		})

		It("Allow to get ALL keys", func() {
			cc.From(Owner).Tx(func() {
				keys, err := dbg.ListKeys(ctx, nil)
				Expect(err).NotTo(HaveOccurred())
				Expect(keys.Keys).To(HaveLen(20)) //total state entries inserted
			})
		})

		It("Allow to delete state entry", func() {
			cc.From(Owner).Tx(func() {
				value, err := dbg.DeleteState(ctx, &debug.CompositeKey{
					Key: []string{`prefixA`, `key0`},
				})
				Expect(err).NotTo(HaveOccurred())
				Expect(value).To(StringerEqual(&debug.Value{
					Key:   []string{`prefixA`, `key0`},
					Value: []byte(`value0`),
				}))
			})

			cc.From(Owner).Tx(func() {
				keys, _ := dbg.ListKeys(ctx, &debug.Prefix{
					Key: []string{`prefixA`},
				})
				Expect(keys.Keys).To(HaveLen(4))
			})

			cc.From(Owner).Tx(func() {
				_, _ = dbg.DeleteState(ctx, &debug.CompositeKey{
					Key: []string{`prefixA`, `key4`},
				})
			})

			cc.From(Owner).Tx(func() {
				keys, _ := dbg.ListKeys(ctx, &debug.Prefix{
					Key: []string{`prefixA`},
				})
				Expect(keys.Keys).To(HaveLen(3))
			})
		})

		It("Allow to clean state entries with key prefix", func() {
			cc.From(Owner).Tx(func() {
				deleted, err := dbg.DeleteStates(ctx, &debug.Prefixes{Prefixes: []*debug.Prefix{{
					Key: []string{`prefixA`},
				}}})
				Expect(err).NotTo(HaveOccurred())

				Expect(len(deleted.Matches)).To(Equal(1))
				Expect(deleted.Matches[`prefixA`]).To(Equal(uint32(3)))
			})

			cc.From(Owner).Tx(func() {
				keys, _ := dbg.ListKeys(ctx, nil)
				Expect(keys.Keys).To(HaveLen(15)) //total state after last clean and delete in previous test
			})
		})

		It("Allow to clean ALL state", func() {
			cc.From(Owner).Tx(func() {
				deleted, err := dbg.DeleteStates(ctx, &debug.Prefixes{Prefixes: []*debug.Prefix{{
					Key: nil,
				}}})
				Expect(err).NotTo(HaveOccurred())

				Expect(len(deleted.Matches)).To(Equal(1))
				Expect(deleted.Matches[``]).To(Equal(uint32(15)))
			})

			cc.From(Owner).Tx(func() {
				keys, _ := dbg.ListKeys(ctx, nil)
				Expect(keys.Keys).To(HaveLen(0)) //total state after last clean and delete in previous test
			})
		})
	})

})
