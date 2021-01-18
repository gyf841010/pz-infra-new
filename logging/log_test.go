package logging

import (
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
	Log.Error(err)
	Log.ErrorOld("test", WithError(err))
}
