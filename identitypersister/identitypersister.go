package identitypersister

import (
	"fmt"
	"github.com/attic-labs/noms/go/datas"
	"github.com/attic-labs/noms/go/types"
	"time"
	"math/rand"
	"github.com/attic-labs/noms/go/marshal"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func RandString(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

const DefaultMemTableSize = 8 * (1 << 20) // 8MB


type IdentityLike struct {
	RootCertificate Certificate
	IntermediateCertificate Certificate
	Devices map[string]DeviceLike
	Metadata map[string]string
	UUID string
}

func NewIdentityLike() *IdentityLike {
	initialDevice := NewDeviceLike()
	return &IdentityLike{
		UUID: RandString(50),
		RootCertificate: NewCertificate(),
		IntermediateCertificate: NewCertificate(),
		Devices: map[string]DeviceLike{
			initialDevice.UUID: initialDevice,
		},
		Metadata: map[string]string{
			"EncryptedRootKey": "",
		},
	}
}

type DeviceLike struct {
	UUID string
	Certificate Certificate
	Description string
	Metadata    map[string]string
}

func NewDeviceLike() DeviceLike {
	return DeviceLike{
		UUID: RandString(100),
		Certificate: NewCertificate(),
		Description: RandString(200),
	}
}

type Certificate struct {
	Pem string
}

func NewCertificate() Certificate {
	return Certificate{Pem: RandString(2048)}
}



func getIdentities(ds datas.Dataset) types.Map {
	hv, ok := ds.MaybeHeadValue()
	if ok {
		fmt.Println("returning existing map")
		return hv.(types.Map)
	}
	fmt.Println("returning empty map")
	return types.NewMap(ds.Database())
}



func Save(ds datas.Dataset, id *IdentityLike) error {
	val,err := marshal.Marshal(ds.Database(), *id)
	if err != nil {
		return fmt.Errorf("error marshaling: %v", err)
	}

	_, err = ds.Database().CommitValue(ds, getIdentities(ds).Edit().Set(types.String(id.UUID), val).Map())

	if err != nil {
		return fmt.Errorf("error committing values: %v", err)
	}

	return nil
}