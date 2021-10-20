package room

import (
	"testing"
)

func TestAuthAvatar(t *testing.T) {
	var authAvatar AuthAvatar
	client := new(client)
	url, err := authAvatar.GetAvatarURL(client)
	if err != ErrNoAvatarURL {
		t.Error("if there is no url, AuthAvatar.GetAvatarURL must return ErrNoAvatarURL")
	}
	testUrl := "http://url-to-avatar/"
	client.userData = map[string]interface{}{"avatar_url": testUrl}
	url, err = authAvatar.GetAvatarURL(client)
	if err != nil {
		t.Error("if there is url, AuthAvatar.GetAvatarURL MUST NOT return error")
	} else {
		if url != testUrl {
			t.Error("AuthAvatar.GetAvatarURL must return correct url")
		}
	}
}
