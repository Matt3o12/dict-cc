package plusTesting

import "os"

// ChangeEnv changes the environment temporarily until the return function
// is called. It is adviced to use defer when changing the environment
// like that:
// `defer ChangeEnv("hello", "world")()`
// this way, the environment will reset automatically as soon as
// the test finishes.
func ChangeEnv(key, value string) func() {
	old := os.Getenv(key)
	os.Setenv(key, value)

	return func() {
		os.Setenv(key, old)
	}
}
