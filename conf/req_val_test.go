package conf

import "testing"

func TestReqVal(t *testing.T) {
	newVal := "some_new_value"

	val := ReqVal{Name: "test", Default: "test"}
	val.Set(newVal)

	if *val.Get() != newVal || !val.IsDefined() {
		t.Errorf("ReqVal.Set func works incorrect, value wasn't set, or flag is steal false: %+v", val)
	}
}
