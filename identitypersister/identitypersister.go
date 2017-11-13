package identitypersister

import (
	"fmt"
	"github.com/quorumcontrol/noms-play/marshal"
	"github.com/attic-labs/noms/go/datas"
	"github.com/attic-labs/noms/go/types"
)

func getSet(ds datas.Dataset) types.Set {
	hv, ok := ds.MaybeHeadValue()
	if ok {
		return hv.(types.Set)
	}
	return types.NewSet(ds.Database())
}



func Save(ds datas.Dataset, i interface{}) error {
	val,err := marshal.Marshal(ds.Database(), i)
	if err != nil {
		return fmt.Errorf("error marshaling: %v", err)
	}

	_, err = ds.Database().CommitValue(ds, getSet(ds).Edit().Insert(val).Set())

	if err != nil {
		return fmt.Errorf("error committing values: %v", err)
	}

	return nil
}