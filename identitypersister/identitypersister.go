package identitypersister

import (
	"fmt"
	"github.com/quorumcontrol/noms-play/marshal"
	"github.com/attic-labs/noms/go/datas"
	"github.com/attic-labs/noms/go/types"
	"github.com/quorumcontrol/qc/identity/identitypb"
)

func getIdentities(ds datas.Dataset) types.Map {
	hv, ok := ds.MaybeHeadValue()
	if ok {
		return hv.(types.Map)
	}
	return types.NewMap(ds.Database())
}



func Save(ds datas.Dataset, id *identitypb.Identity) error {
	val,err := marshal.Marshal(ds.Database(), *id)
	if err != nil {
		return fmt.Errorf("error marshaling: %v", err)
	}

	_, err = ds.Database().CommitValue(ds, getIdentities(ds).Edit().Set(types.String(id.Uid()), val).Map())

	if err != nil {
		return fmt.Errorf("error committing values: %v", err)
	}

	return nil
}