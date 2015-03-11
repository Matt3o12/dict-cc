package plusTesting

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

const testEnvKey = "dict_cc_test_foo"

func TestChangeEnv(t *testing.T) {
	os.Setenv(testEnvKey, "bar")
	deferFunc := ChangeEnv(testEnvKey, "foobar")
	assert.Equal(t, "foobar", os.Getenv(testEnvKey))
	deferFunc()
	assert.Equal(t, "bar", os.Getenv(testEnvKey))
}
