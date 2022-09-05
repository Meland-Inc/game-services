package daprInvoke

import (
	"testing"
)

func Test_DaprClient(t *testing.T) {
	c, err := NewDaprClient("5770")
	t.Log(err)
	t.Log(c)
}
