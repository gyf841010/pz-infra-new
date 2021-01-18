package logging

import (
	"fmt"
	"testing"

	"github.com/pkg/errors"
)

func aTestError() error {
	err := errors.New("test_error")
	return err
}
func TestLogging(t *testing.T) {
	InitLogger("test")
	err := aTestError()
	Log.ErrorWithStack("test", WithError(err))
	Log.Error("test", WithError(err))

	fmt.Println(errors.WithMessage(errors.Wrap(err, "wrap info"), "with message").Error())
	fmt.Println(errors.Cause(
		errors.WithMessage(errors.Wrap(err, "wrap info"), "with message")))
}
