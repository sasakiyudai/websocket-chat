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

func TestGravatar(t *testing.T) {
	var gravatarAvatar GravatarAvatar
	client := new(client)
	client.userData =
		map[string]interface{}{"email": "MyEmailAddress@example.com"}
	url, err := gravatarAvatar.GetAvatarURL(client)
	if err != nil {
		t.Error("GravatarAvatar.GetAvatarURL MUST NOT return error")
	}
	if url != "//www.gravatar.com/avatar/0bc83cb571cd1c50ba6f3e8a78ef1346" {
		t.Errorf("GravatarAvatar.GetAvatarURL return wrong value")
	}
}
