package room

import (
	"errors"
)

var ErrNoAvatarURL = errors.New("chat: cannot get avatar url")

type Avatar interface {
	// return client's avatar url, 
	// if cannot get avatar_url from client, return ErrNoAvatarURL
	GetAvatarURL(c *client) (string, error)
}

type AuthAvatar struct{}
var UseAuthAvatar AuthAvatar
func (_ AuthAvatar) GetAvatarURL(c *client) (string, error) {
	if url, ok := c.userData["avatar_url"]; ok {
		if urlStr, ok := url.(string); ok {
			return urlStr, nil
		}
	}
	return "", ErrNoAvatarURL
}