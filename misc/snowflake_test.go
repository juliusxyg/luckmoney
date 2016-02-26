package misc

import (
	"testing"
)

func TestCheckUniqueId(t *testing.T) {
	err := CheckUniqueId()
	if err!=nil {
		t.Error("init snowflake service error")
	}
}

func TestUniqueId(t *testing.T) {
	id := UniqueId()
	if id == 0 {
		t.Error("generate id error")
	}

	id_2 := UniqueId()
	if id_2 == 0 {
		t.Error("generate id2 error")
	}

	t.Logf("id %v, next id %v", id, id_2)

	if id >= id_2 {
		t.Error("compare ids error")
	}
}