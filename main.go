package main

import (
	"fmt"
	"os"

	"github.com/attic-labs/noms/go/spec"
	"github.com/attic-labs/noms/go/types"
)

func newPerson(givenName string, male bool) types.Struct {
	return types.NewStruct("Person", types.StructData{
		"given": types.String(givenName),
		"male":  types.Bool(male),
	})
}

func main() {
	sp, err := spec.ForDataset("nbs:tmp/noms::people")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not create dataset: %s\n", err)
		return
	}
	defer sp.Close()

	if headValue, ok := sp.GetDataset().MaybeHeadValue(); !ok {
		fmt.Fprintf(os.Stdout, "head is empty\n")
	} else {
		// type assertion to convert Head to List
		personList := headValue.(types.List)
		personEditor := personList.Edit()
		data := personEditor.Append(
			types.NewStruct("Person", types.StructData{
				"given":  types.String("Jon"),
				"family": types.String("Snow"),
				"male":   types.Bool(true),
			}),
		).List()

		fmt.Fprintf(os.Stdout, "data type: %v\n", types.TypeOf(data).Describe())

		_, err = sp.GetDatabase().CommitValue(sp.GetDataset(), data)
		if err != nil {
			fmt.Fprint(os.Stderr, "Error commiting: %s\n", err)
		}
	}
}