package documents

import "gopkg.in/mgo.v2/bson"

type ChannelDocument struct {
	Id              bson.ObjectId   `bson:"_id,omitempty"`
	ChannelName     string          `bson:"channel_name"`
	IsPrivate       bool            `bson:"is_private"`
	ChannelPassword string          `bson:"channel_password"`
	UsersIds        []bson.ObjectId `bson:"users_ids"`
}

type UserDocument struct {
	Id           bson.ObjectId `bson:"_id,omitempty"`
	Username     string
	PasswordHash string          `bson:"password_hash"`
	SessionToken string          `bson:"session_token"`
	ChannelsIds  []bson.ObjectId `bson:"channels"`
}

func (u *UserDocument) ToMap() map[string]string {
	return map[string]string{
		"id":            u.Id.Hex(),
		"username":      u.Username,
		"session_token": u.SessionToken,
	}
}

func (d *UserDocument) String() string {
	return d.Username
}

type Message struct {
	Id     bson.ObjectId `bson:"_id,omitempty"`
	Author string        `json:"author"`
	Body   string        `json:"body"`
}

func (m *Message) String() string {
	return m.Author + " says " + m.Body
}
