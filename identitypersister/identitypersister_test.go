package identitypersister

import (
	"testing"
	"github.com/attic-labs/noms/go/spec"
	"github.com/spf13/afero"
	"github.com/quorumcontrol/qc/identity"
	"github.com/attic-labs/noms/go/types"
	"github.com/quorumcontrol/qc/identity/identitypb"
	"github.com/quorumcontrol/noms-play/marshal"
)

func TestGetFields(t *testing.T) {
	fs := afero.NewOsFs()
	fs.RemoveAll("tmp/noms")
	fs.MkdirAll("tmp/noms", 0755)


	sp,err := spec.ForDataset("nbs:tmp/noms::people")

	if err != nil {
		t.Fatalf("error getting dataset: %v", err)
	}

	alice := identity.GenerateIdentity("alice", "insaasity")

	err = Save(sp.GetDataset(), alice)

	if err != nil {
		t.Fatalf("error getting fields: %v", err)
	}

	hv, ok := sp.GetDataset().MaybeHeadValue()
	if ok {
		people := hv.(types.Map)

		dbAliceMarshaled := people.Get(types.String(alice.Uid()))

		dbAlice := &identitypb.Identity{}

		err = marshal.Unmarshal(dbAliceMarshaled, dbAlice)

		if err != nil {
			t.Fatalf("Error unmarshaling: %v", err)
		}

		if !dbAlice.Equal(alice) {
			t.Errorf("db Alice %+v does not equal alice: %+v", dbAlice, alice)
		}


	} else {
		t.Fatalf("no head value")
	}

}