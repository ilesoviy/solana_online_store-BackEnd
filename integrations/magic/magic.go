package magic

import (
	"errors"

	"github.com/magiclabs/magic-admin-go"
	"github.com/magiclabs/magic-admin-go/client"
	"github.com/magiclabs/magic-admin-go/token"
)

var magicClient *client.API

func Init(apiSecret string) error {
	if magicClient != nil {
		return errors.New("magic service already initiated")
	}
	magicClient = client.New(apiSecret, magic.NewDefaultClient())
	return nil
}

func GetUserInfo(did string) (*magic.UserInfo, error) {
	return magicClient.User.GetMetadataByToken(did)
}

func ValidateDidToken(did string) (*token.Token, error) {
	tk, err := token.NewToken(did)
	if err != nil {
		return nil, err
	}
	if err := tk.Validate(); err != nil {
		return nil, err
	}
	return tk, nil
}
