package identitypersister

import (
	"testing"
)


type MySubType struct {
	ASubCoolField string
}

type MyType struct {
	CoolField string
	MySubType *MySubType
}

func TestGetFields(t *testing.T) {
	err := GetFields(&MyType{})

	if err != nil {
		t.Fatalf("error getting fields: %v", err)
	}
}