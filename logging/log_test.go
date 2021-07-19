package logging

import (
	"fmt"
	"testing"

	"github.com/gyf841010/pz-infra-new/errorUtil"
	"github.com/pkg/errors"
)

func aTestError() error {
	err := errors.New("test_error")
	return err
}

func TestLogging(t *testing.T) {
	InitLogger("test")
	err := aTestError()
	Log.ErrorWithStack("my", err, "", With("with msg", "test"), With("struct", struct {
		A string `json:"a"`
	}{A: "some msg"}))

	Log.Error("test", WithError(err))
	switch t := err.(type) {
	case *errorUtil.HError:
		fmt.Println("custom error", t.Code(), t.Error(), t.Message, t.ResCode)
	default:
		fmt.Println("default", errors.Cause(err))
	}
	Log.Info("-------------------\n")
	s := struct{ str string }{"test_error_f"}
	Log.ErrorWithStackf(err, "结构体%+v,数字:%d", s, 1111111111111111)
}
