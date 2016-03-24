package datastore

import (
	"log"

	"golang.org/x/oauth2/google"
	"google.golang.org/cloud"
	"golang.org/x/net/context"
	"google.golang.org/cloud/datastore"
	"github.com/mshockwave/share-sound-api-server/common"
	"golang.org/x/oauth2"
)

var ctx context.Context
var defaultTokenSource oauth2.TokenSource = nil
func init(){

	ctx = context.Background()

	var err error
	defaultTokenSource, err = google.DefaultTokenSource(ctx,
		datastore.ScopeDatastore,
		datastore.ScopeUserEmail,
	)
	if err != nil || defaultTokenSource == nil{
		log.Fatalf("Error getting google default token source")
		return
	}
}

type DataStoreClient struct{
	Client *datastore.Client
	Ctx	context.Context
}

func GetNewDataStoreClient() (*DataStoreClient, error){
	client, err := datastore.NewClient(ctx, common.PROJECT_ID, cloud.WithTokenSource(defaultTokenSource))
	return &DataStoreClient{
		Client: client,
		Ctx: ctx,
	}, err
}
func GetContext() context.Context { return ctx }

func (this *DataStoreClient) NewKey(kind, name string, id int64, parent *datastore.Key) *datastore.Key {
	//Wrapper
	return datastore.NewKey(this.Ctx, kind, name, id, parent)
}
func (this *DataStoreClient) Run(query *datastore.Query) *datastore.Iterator {
	return this.Client.Run(this.Ctx, query)
}
func (this *DataStoreClient) RunInTransaction(f func(tx *datastore.Transaction) error, opts ...datastore.TransactionOption) (*datastore.Commit, error){
	return this.Client.RunInTransaction(this.Ctx, f, opts...)
}
