package service_test

import (
	"context"
	"reflect"
	"testing"

	"github.com/s7techlab/hlf-sdk-go/api"
	"github.com/s7techlab/hlf-sdk-go/client/chaincode"
	"github.com/s7techlab/hlf-sdk-go/client/chaincode/txwaiter"

	"github.com/s7techlab/cckit/gateway/service"
)

func TestContextWithDefaultDoOption(tt *testing.T) {
	var (
		doOptsAll  = []api.DoOption{chaincode.WithTxWaiter(txwaiter.All)}
		doOptsSelf = []api.DoOption{chaincode.WithTxWaiter(txwaiter.Self)}
	)

	tt.Run("Context empty", func(t *testing.T) {
		ctx := service.ContextWithDefaultDoOption(context.Background(), doOptsSelf...)
		result := service.DoOptionFromContext(ctx)
		if !reflect.DeepEqual(result, doOptsSelf) {
			t.Errorf("Expected find doOption on ctx %p, got -%p", doOptsSelf, result)
		}
	})

	tt.Run("Dont allow update on ctx DoOption", func(t *testing.T) {

		ctx := service.ContextWithDefaultDoOption(context.Background(), doOptsAll...)

		ctx = service.ContextWithDefaultDoOption(ctx, doOptsSelf...)
		result := service.DoOptionFromContext(ctx)
		if !reflect.DeepEqual(result, doOptsAll) {
			t.Errorf("Expected find first option on ctx %p, got -%p", doOptsAll, result)
		}
	})
}

func TestContextWithDoOption(tt *testing.T) {
	var (
		doOptsAll  = []api.DoOption{chaincode.WithTxWaiter(txwaiter.All)}
		doOptsSelf = []api.DoOption{chaincode.WithTxWaiter(txwaiter.Self)}
	)

	tt.Run("Context empty", func(t *testing.T) {
		ctx := service.ContextWithDoOption(context.Background(), doOptsSelf...)
		result := service.DoOptionFromContext(ctx)
		if !reflect.DeepEqual(result, doOptsSelf) {
			t.Errorf("Expected find doOption on ctx %p, got -%p", doOptsSelf, result)
		}
	})

	tt.Run("Allow update on ctx DoOption", func(t *testing.T) {
		ctx := service.ContextWithDoOption(context.Background(), doOptsAll...)
		ctx = service.ContextWithDoOption(ctx, doOptsSelf...)
		result := service.DoOptionFromContext(ctx)
		if !reflect.DeepEqual(result, doOptsSelf) {
			t.Errorf("Expected find second option on ctx %p, got -%p", doOptsSelf, result)
		}
	})
}

func TestDoOptionFromContext(tt *testing.T) {
	var (
		doOptsEmpty = []api.DoOption{}
		doOptsSelf  = []api.DoOption{chaincode.WithTxWaiter(txwaiter.Self)}
	)

	tt.Run("Context empty", func(t *testing.T) {
		result := service.DoOptionFromContext(ctx)
		if !reflect.DeepEqual(result, doOptsEmpty) {
			t.Errorf("Expected find doOption on ctx %p, got -%p", doOptsEmpty, result)
		}
	})

	tt.Run("Context with doOption", func(t *testing.T) {
		ctx = service.ContextWithDoOption(ctx, doOptsSelf...)
		result := service.DoOptionFromContext(ctx)
		if !reflect.DeepEqual(result, doOptsSelf) {
			t.Errorf("Expected find doOption on ctx %p, got -%p", doOptsSelf, result)
		}
	})
}
