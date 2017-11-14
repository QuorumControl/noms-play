package identitypersister

import (
	"testing"
	"github.com/spf13/afero"
	"github.com/quorumcontrol/qc/identity"
	"github.com/attic-labs/noms/go/types"
	"github.com/quorumcontrol/qc/identity/identitypb"
	"github.com/quorumcontrol/noms-play/marshal"
	"github.com/attic-labs/noms/go/datas"
	"github.com/attic-labs/noms/go/nbs"
)

const DefaultMemTableSize = 8 * (1 << 20) // 8MB


func TestGetFields(t *testing.T) {
	fs := afero.NewOsFs()
	fs.RemoveAll("tmp/noms")
	fs.MkdirAll("tmp/noms", 0755)

	alice := identity.GenerateIdentity("alice", "insaasity")

	sp := datas.NewDatabase(nbs.NewLocalStore("tmp/noms", DefaultMemTableSize))

	err := Save(sp.GetDataset("identities"), alice)

	if err != nil {
		t.Fatalf("error saving: %v", err)
	}

	//update alice

	alice.Metadata = map[string]string{"myUpdate": "another thing"}

	err = Save(sp.GetDataset("identities"), alice)

	if err != nil {
		t.Fatalf("error getting fields: %v", err)
	}

	hv, ok := sp.GetDataset("identities").MaybeHeadValue()
	if ok {
		people := hv.(types.Map)

		dbAliceMarshaled := people.Get(types.String(alice.Uid()))

		dbAlice := &identitypb.Identity{}

		err = marshal.Unmarshal(dbAliceMarshaled, dbAlice)

		if err != nil {
			t.Fatalf("Error unmarshaling: %v", err)
		}

		if !alice.Equal(dbAlice) {
			t.Errorf("alices were not equal\n\n alice:\n %v \n\ndbAlice:\n %v", alice, dbAlice)
		}

		//for name,equalPairs := range map[string][]string {
		//	"Name": {alice.Name, dbAlice.Name},
		//	"Organization": {alice.Organization, dbAlice.Organization},
		//	"CurrentDeviceFingerprint": {alice.CurrentDevice().Certificate.Pem.PublicKeyFingerprint(), dbAlice.CurrentDevice().Certificate.Pem.PublicKeyFingerprint()},
		//	"CurrentDeviceCreatedAt": {strconv.FormatInt(alice.CurrentDevice().CreatedAt.UnixNano(), 10), strconv.FormatInt(dbAlice.CurrentDevice().CreatedAt.UnixNano(), 10)},
		//	} {
		//	if equalPairs[0] != equalPairs[1] {
		//		t.Errorf("%s values did not match: %s != %s", name, equalPairs[0], equalPairs[1])
		//	}
		//}



	} else {
		t.Fatalf("no head value")
	}

}