package encryption

import (
	"github.com/hyperledger/fabric/protos/peer"
	"github.com/pkg/errors"
	"github.com/s7techlab/cckit/response"
	"github.com/s7techlab/cckit/router"
	"github.com/s7techlab/cckit/state"
)

// ArgsDecryptIfKeyProvided  - pre middleware, decrypts chaincode method arguments if key provided in transient map
func ArgsDecryptIfKeyProvided(next router.ContextHandlerFunc, pos ...int) router.ContextHandlerFunc {
	return argsDecryptor(next, false, nil)
}

// ArgsDecryptIfKeyProvided  - pre middleware, decrypts chaincode method arguments, key must be provided in transient map
func ArgsDecrypt(next router.ContextHandlerFunc, pos ...int) router.ContextHandlerFunc {
	return argsDecryptor(next, true, nil)
}

func ArgsDecryptExcept(exceptMethod ...string) router.ContextMiddlewareFunc {
	return func(next router.ContextHandlerFunc, pos ...int) router.ContextHandlerFunc {
		return argsDecryptor(next, true, exceptMethod)
	}
}

func decryptReplaceArgs(key []byte, c router.Context) error {
	args, err := DecryptArgs(key, c.Stub().GetArgs())
	if err != nil {
		return errors.Wrap(err, `args`)
	}
	c.ReplaceArgs(args)
	return nil
}

func argsDecryptor(next router.ContextHandlerFunc, keyShouldBeProvided bool, exceptMethod []string) router.ContextHandlerFunc {
	return func(c router.Context) peer.Response {
		if args := c.GetArgs(); len(exceptMethod) > 0 && len(args) > 0 {
			for _, m := range exceptMethod {
				if string(args[0]) == m {
					return next(c)
				}
			}
		}

		key, err := KeyFromTransient(c)
		// no key provided

		if err != nil {
			c.Logger().Debugf(`no decrypt key provided: %s`, err)
			if err == ErrKeyNotDefinedInTransientMap && keyShouldBeProvided {
				return response.Error(err)
			}
			return next(c)
		}

		if err = decryptReplaceArgs(key, c); err != nil {
			return response.Error(err)
		}

		return next(c)
	}
}

// EncStateContext replaces default state with encrypted state
func EncStateContext(next router.HandlerFunc, pos ...int) router.HandlerFunc {
	return func(c router.Context) (res interface{}, err error) {
		if err = replaceStateEncrypted(c); err != nil {
			return
		}
		return next(c)
	}
}

func replaceStateEncrypted(c router.Context) (err error) {
	var (
		s state.State
		e state.Event
	)

	if s, err = StateWithTransientKey(c); err != nil {
		return err
	}

	if e, err = EventWithTransientKey(c); err != nil {
		return err
	}

	c.UseState(s)
	c.UseEvent(e)
	return
}

// EncStateContext replaces default state with encrypted state
func EncStateContextIfKeyProvided(next router.HandlerFunc, pos ...int) router.HandlerFunc {
	return func(c router.Context) (res interface{}, err error) {
		if _, err = KeyFromTransient(c); err != nil {
			// skip state changing
			return next(c)
		}

		if err = replaceStateEncrypted(c); err != nil {
			return
		}
		return next(c)
	}
}
