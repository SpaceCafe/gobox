package logger

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDefault(t *testing.T) {
	logger := New()
	assert.Equal(t, reflect.DeepEqual(Default(), logger), true, "expected Default() to be the same as New()")
	assert.Equal(t, reflect.ValueOf(SetFormat).Pointer(), reflect.ValueOf(logger.SetFormat).Pointer(), "expected SetFormat() to be the same as logger.SetFormat()")
	assert.Equal(t, reflect.ValueOf(SetLevel).Pointer(), reflect.ValueOf(logger.SetLevel).Pointer(), "expected SetLevel() to be the same as logger.SetLevel()")
	assert.Equal(t, reflect.ValueOf(SetOutput).Pointer(), reflect.ValueOf(logger.SetOutput).Pointer(), "expected SetOutput() to be the same as logger.SetOutput()")
	assert.Equal(t, reflect.ValueOf(Debug).Pointer(), reflect.ValueOf(logger.Debug).Pointer(), "expected Debug() to be the same as logger.Debug()")
	assert.Equal(t, reflect.ValueOf(Debugf).Pointer(), reflect.ValueOf(logger.Debugf).Pointer(), "expected Debugf() to be the same as logger.Debugf()")
	assert.Equal(t, reflect.ValueOf(Info).Pointer(), reflect.ValueOf(logger.Info).Pointer(), "expected Info() to be the same as logger.Info()")
	assert.Equal(t, reflect.ValueOf(Infof).Pointer(), reflect.ValueOf(logger.Infof).Pointer(), "expected Infof() to be the same as logger.Infof()")
	assert.Equal(t, reflect.ValueOf(Warn).Pointer(), reflect.ValueOf(logger.Warn).Pointer(), "expected Warn() to be the same as logger.Warn()")
	assert.Equal(t, reflect.ValueOf(Warnf).Pointer(), reflect.ValueOf(logger.Warnf).Pointer(), "expected Warnf() to be the same as logger.Warnf()")
	assert.Equal(t, reflect.ValueOf(Error).Pointer(), reflect.ValueOf(logger.Error).Pointer(), "expected Error() to be the same as logger.Error()")
	assert.Equal(t, reflect.ValueOf(Errorf).Pointer(), reflect.ValueOf(logger.Errorf).Pointer(), "expected Errorf() to be the same as logger.Errorf()")
	assert.Equal(t, reflect.ValueOf(Fatal).Pointer(), reflect.ValueOf(logger.Fatal).Pointer(), "expected Fatal() to be the same as logger.Fatal()")
	assert.Equal(t, reflect.ValueOf(Fatalf).Pointer(), reflect.ValueOf(logger.Fatalf).Pointer(), "expected Fatalf() to be the same as logger.Fatalf()")
}
