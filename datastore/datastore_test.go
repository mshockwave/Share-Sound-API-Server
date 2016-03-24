package datastore

import (
	"testing"

	"google.golang.org/cloud/datastore"
	"time"
)

type testEntity struct {
	Name		string
	Description	string
	Created		time.Time
}
func (this testEntity) String() string{
	var ret = ""

	ret += ("Name: " + this.Name + " / ")
	ret += ("Description: " + this.Description + " / ")
	ret += ("TimeStamp: " + this.Created.String())

	return ret
}

func TestDataStoreClient(t *testing.T){

	t.Log("Start to test DataStore client...")

	if client, err := GetNewDataStoreClient(); err == nil {
		_, e := client.Client.RunInTransaction(client.Ctx, func(tx *datastore.Transaction) error{

			key := datastore.NewKey(client.Ctx, "TestKind", "sampleEntity", 0, nil)

			//Insert
			t.Log("Try to insert new entity...\n")
			_, txErr := tx.Put(key, &testEntity{
				Name: "Test Name",
				Description: "Test Description",
				Created: time.Now(),
			})
			if(txErr != nil){
				return txErr
			}

			return txErr
		})
		if e != nil {
			t.Error("Failed Commiting Insert Transaction")
			t.Error(e)
			t.FailNow()
		}

		_, e = client.Client.RunInTransaction(client.Ctx, func(tx *datastore.Transaction) error{

			key := datastore.NewKey(client.Ctx, "TestKind", "sampleEntity", 0, nil)

			//Get
			t.Log("Try to get entity...\n")
			result := testEntity{}
			txErr := tx.Get(key, &result)
			if(txErr != nil){
				return txErr
			}

			t.Logf("Get result: %s\n", result)

			//Delete
			t.Log("Try to delete entity...\n")
			txErr = tx.Delete(key)

			return txErr
		})
		if e != nil {
			t.Error("Failed Commiting Get/Delete Transaction")
			t.Error(e)
			t.FailNow()
		}
	}else{
		t.Error("Failed Getting DataStore Client")
		t.Error(err)
		t.FailNow()
	}
}
