package schema

import (
	"time"
	"github.com/mshockwave/share-sound-api-server/handlers/protos"
	"github.com/mshockwave/share-sound-api-server/storage"
	"github.com/mshockwave/share-sound-api-server/common"
	"mime"
)

const(
	USER_PROFILE_KIND = "UserProfile"
	STORY_KIND = "Story"
)

type HashId string
func (this HashId) String() string {
	return string(this)
}


type UserAuth struct {
	PasswordBcryptHash	string
	PasswordBcryptCost	int
}

type User struct {
	Username		string
	Email			string
	CreatedTimeStamp	time.Time

	Auth			UserAuth

	Thumbnail		string "" //Storage path
}

type Location struct {
	Longitude	float32
	Latitude	float32
}
func (this *Location) FromProtoBuf(location *protos.Location){
	this.Longitude = location.Longitude
	this.Latitude = location.Latitude
}

type StoryMeta struct {
	Id               HashId
	UploaderEmail    string

	Title            string

	Description      string

	AudioAttachments []AudioAttachmentMeta

	ImageAttachments []ImageAttachmentMeta
}

type AudioAttachmentMeta struct {
	Id		HashId

	Mime		string

	Name		string

	TimeStamp	time.Time

	Location	Location

	ContentPath	string "" //Storage path
}
func (this *AudioAttachmentMeta) FromProtoBuf(audio *protos.AudioAttachment, storageClient *storage.StorageClient) error{

	this.Id = HashId(common.GetDefaultSecureHash())

	if str, _, err := mime.ParseMediaType(audio.Mime); err == nil{
		this.Mime = str
	}else{
		return err
	}

	this.Name = audio.Name

	this.TimeStamp = time.Unix(int64(audio.TimeStamp), 0)

	(&this.Location).FromProtoBuf(audio.Location)

	bucket := storageClient.GetDefaultBucket()
	objName := common.PathJoin(common.STORAGE_AUDIO_FOLDER, common.GetDefaultSecureHash())
	objName += common.GetFileExtension(audio.Name)

	obj := bucket.Object(objName)

	writer := obj.NewWriter(storageClient.Ctx)
	defer writer.Close()

	_, e := writer.Write(audio.Content)

	return e
}

type ImageAttachmentMeta struct {
	Id		HashId

	Mime		string

	Name		string

	ContentPath	string "" //Storage path
}
func (this *ImageAttachmentMeta) FromProtoBuf(image *protos.ImageAttachment, storageClient *storage.StorageClient) error {
	this.Id = HashId(common.GetDefaultSecureHash())

	if str, _, err := mime.ParseMediaType(image.Mime); err == nil{
		this.Mime = str
	}else{
		return err
	}

	this.Name = image.Name

	bucket := storageClient.GetDefaultBucket()
	objName := common.PathJoin(common.STORAGE_IMAGE_FOLDER, common.GetDefaultSecureHash())
	objName += common.GetFileExtension(image.Name)

	obj := bucket.Object(objName)

	writer := obj.NewWriter(storageClient.Ctx)
	defer writer.Close()

	_, e := writer.Write(image.Content)

	return e
}
