package service_test

import (
	"context"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/s7techlab/hlf-sdk-go/api"
	"github.com/s7techlab/hlf-sdk-go/client/chaincode"
	"github.com/s7techlab/hlf-sdk-go/client/chaincode/txwaiter"

	"github.com/s7techlab/cckit/gateway/service"
)

func TestContext(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Context Suite")
}

var (
	doOptsEmpty = []api.DoOption{}
	doOptsAll   = []api.DoOption{chaincode.WithTxWaiter(txwaiter.All)}
	doOptsSelf  = []api.DoOption{chaincode.WithTxWaiter(txwaiter.Self)}
)

var _ = Describe(`DoOption`, func() {

	It("Set Default DoOption", func() {
		ctx := service.ContextWithDefaultDoOption(context.Background(), doOptsSelf...)
		result := service.DoOptionFromContext(ctx)
		Expect(result).To(Equal(doOptsSelf))
	})

	It("Default dont allow update on ctx DoOption", func() {
		ctx := service.ContextWithDefaultDoOption(context.Background(), doOptsAll...)
		ctx = service.ContextWithDefaultDoOption(ctx, doOptsSelf...)
		result := service.DoOptionFromContext(ctx)
		Expect(result).To(Equal(doOptsAll))
	})

	It("Set DoOption to Context", func() {
		ctx := service.ContextWithDoOption(context.Background(), doOptsSelf...)
		result := service.DoOptionFromContext(ctx)
		Expect(result).To(Equal(doOptsSelf))
	})

	It("Update DoOption to Context", func() {
		ctx := service.ContextWithDoOption(context.Background(), doOptsAll...)
		ctx = service.ContextWithDoOption(ctx, doOptsSelf...)
		result := service.DoOptionFromContext(ctx)
		Expect(result).To(Equal(doOptsSelf))
	})

	It("Update DoOption to Context", func() {
		ctx := service.ContextWithDoOption(context.Background(), doOptsAll...)
		ctx = service.ContextWithDoOption(ctx, doOptsSelf...)
		result := service.DoOptionFromContext(ctx)
		Expect(result).To(Equal(doOptsSelf))
	})

	It("Allow get DoOptions from empty Context", func() {
		result := service.DoOptionFromContext(ctx)
		Expect(result).To(Equal(doOptsEmpty))
	})

	It("Allow get DoOptions from filled Context", func() {
		ctx := service.ContextWithDoOption(ctx, doOptsSelf...)
		result := service.DoOptionFromContext(ctx)
		Expect(result).To(Equal(doOptsSelf))
	})
})
