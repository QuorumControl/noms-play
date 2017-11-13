package identitypersister

import (
	"testing"
	"github.com/attic-labs/noms/go/spec"
	"github.com/spf13/afero"
	"fmt"
	"github.com/quorumcontrol/qc/identity/identitypb"
)


type MySubType struct {
	ASubCoolField string
}

type MyType struct {
	CoolField string
	MySubType *MySubType
}

func TestGetFields(t *testing.T) {
	fs := afero.NewOsFs()
	fs.RemoveAll("tmp/noms")
	fs.MkdirAll("tmp/noms", 0755)


	sp,err := spec.ForDataset("nbs:tmp/noms::people")

	if err != nil {
		t.Fatalf("error getting dataset: %v", err)
	}

	//err = Save(sp.GetDataset(), *identity.GenerateIdentity("alice", "insaasity"))
	err = Save(sp.GetDataset(), identitypb.Identity{})

	if err != nil {
		t.Fatalf("error getting fields: %v", err)
	}

	val := sp.GetDataset().Head()
	fmt.Printf("head: %v", val)
}