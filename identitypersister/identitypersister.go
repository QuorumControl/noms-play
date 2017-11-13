package identitypersister

import (
	"fmt"
	"github.com/attic-labs/noms/go/types"
	"github.com/attic-labs/noms/go/marshal"
	"github.com/attic-labs/noms/go/chunks"
)

func newTestValueStore() *types.ValueStore {
	ts := &chunks.TestStorage{}
	return types.NewValueStore(ts.NewView())
}


func GetFields(i interface{}) error {
	val,err := marshal.Marshal(newTestValueStore(), i)

	fmt.Printf("val: %v", val)

	return err
}