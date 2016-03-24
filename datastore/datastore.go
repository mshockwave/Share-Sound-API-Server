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
