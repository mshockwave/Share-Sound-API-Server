package schema

import "time"

const(
	USER_PROFILE_KIND = "UserProfile"
)

type UserAuth struct {
	PasswordBcryptHash	string
	PasswordBcryptCost	int
}

type User struct {
	Username		string
	Email			string
	CreatedTimeStamp	time.Time

	Auth			UserAuth

	Thumbnail		string ""//Storage path
}
