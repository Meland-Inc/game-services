package configModule

import (
	"testing"
)

func Test_TaskCnf(t *testing.T) {
	t.Log("-----------begin----------")
	err := Init()
	t.Log("InitTaskCnfDB err : ", err)
	t.Log(AllTaskCnfs())
	t.Log(AllTaskListCnfs())
}
