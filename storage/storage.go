package storage

import(
	"google.golang.org/cloud"
	"google.golang.org/cloud/storage"
	"golang.org/x/oauth2/google"
	"golang.org/x/oauth2"
	"golang.org/x/net/context"
	"github.com/mshockwave/share-sound-api-server/common"
)

var defaultTokenSource oauth2.TokenSource = nil
var ctx context.Context

func init(){

	ctx = context.Background()

	var err error
	defaultTokenSource, err = google.DefaultTokenSource( ctx,
		storage.ScopeFullControl,
	)
	if err != nil || defaultTokenSource == nil{
		common.LogE.Fatalf("Error getting storage token source: %s\n", err.Error())
	}
}

type StorageClient struct {
	Client *storage.Client
	Ctx	context.Context
}

func GetNewStorageClient() (*StorageClient, error) {
	client, err := storage.NewClient(ctx, cloud.WithTokenSource(defaultTokenSource))
	return &StorageClient{
		Client: client,
		Ctx: ctx,
	}, err
}

func (this *StorageClient) Close() { this.Client.Close() }

func (this *StorageClient) GetDefaultBucket() *storage.BucketHandle{ return this.Client.Bucket(common.MAIN_STORAGE_BUCKET) }
