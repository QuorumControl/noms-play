package identitypersister

import (
	"fmt"
	"github.com/quorumcontrol/noms-play/marshal"
	"github.com/attic-labs/noms/go/datas"
	"github.com/attic-labs/noms/go/types"
	"github.com/quorumcontrol/qc/identity/identitypb"
	"reflect"
	"github.com/quorumcontrol/qc/simpcert"
)

func init() {
	simpcertType := reflect.TypeOf(simpcert.Certificate{})

	marshal.RegisterEncoder(simpcertType, CertificateEncoder)
	marshal.RegisterDecoder(simpcertType, CertificateDecoder)
}

func CertificateEncoder(v reflect.Value) types.Value {
	if !v.IsValid() {
		return types.String("")
	}

	cert := v.Interface().(simpcert.Certificate)
	pnter := &cert

	bytes,err := pnter.Marshal()

	if err != nil {
		panic("cannot marshal the simpcert")
	}

	return types.String(string(bytes))
}

func CertificateDecoder(v types.Value, rv reflect.Value) {
	if publicPem, ok := v.(types.String); ok {

		fmt.Printf("rv is: %+v, kind: %v, type: %v\n", rv, rv.Kind(), rv.Type())

		cert := &simpcert.Certificate{}
		cert.LoadFromString(string(publicPem))

		rv.Set(reflect.ValueOf(*cert))
	} else {
		panic("simpcert stored in non string format")
	}
}

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