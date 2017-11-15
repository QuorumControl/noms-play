package identitypersister

import (
	"testing"
	"github.com/spf13/afero"
	"github.com/attic-labs/noms/go/types"
	"github.com/attic-labs/noms/go/datas"
	"github.com/attic-labs/noms/go/nbs"
	"github.com/attic-labs/noms/go/marshal"
)

func TestSaveAndUpdate(t *testing.T) {
	fs := afero.NewOsFs()
	fs.RemoveAll("tmp/noms")
	fs.MkdirAll("tmp/noms", 0755)

	alice := NewIdentityLike()

	sp := datas.NewDatabase(nbs.NewLocalStore("tmp/noms", DefaultMemTableSize))

	err := Save(sp.GetDataset("identities"), alice)

	if err != nil {
		t.Fatalf("error saving: %v", err)
	}

	//update alice

	alice.Metadata = map[string]string{"myUpdate": "another thing"}

	newDevice := NewDeviceLike()

	alice.Devices[newDevice.UUID] = newDevice

	err = Save(sp.GetDataset("identities"), alice)

	if err != nil {
		t.Fatalf("error getting fields: %v", err)
	}

	hv, ok := sp.GetDataset("identities").MaybeHeadValue()
	if ok {
		people := hv.(types.Map)

		dbAliceMarshaled := people.Get(types.String(alice.UUID))

		dbAlice := &IdentityLike{}

		err = marshal.Unmarshal(dbAliceMarshaled, dbAlice)

		if err != nil {
			t.Fatalf("Error unmarshaling: %v", err)
		}

		if alice.UUID != dbAlice.UUID {
			t.Errorf("alices were not equal\n\n alice:\n %v \n\ndbAlice:\n %v", alice, dbAlice)
		}

	} else {
		t.Fatalf("no head value")
	}

}