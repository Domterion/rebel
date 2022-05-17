package rebel

import (
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/relvacode/iso8601"
)

/*
Library structs
*/

type Client struct {
	Token string

	onReadyFunction         func(*Client, *Ready)
	onMessageFunction       func(*Client, *Message)
	onMessageUpdateFunction func(*Client, *MessageUpdate)

	websocket *websocket.Conn
	http      *http.Client
}

type Context struct {
	Client *Client
}

/*
API related structs
*/

type ApiResponse struct {
	Type string `json:"type"`
}

type Ping struct {
	ApiResponse

	Data int `json:"data"`
}

type Ready struct {
	Users    []User    `json:"users"`
	Servers  []Server  `json:"servers"`
	Channels []Channel `json:"channels"`
}

type User struct {
	Id           string           `json:"_id"`
	Username     string           `json:"username"`
	Avatar       File             `json:"avatar"`
	Badges       int              `json:"badges"`
	Status       UserStatus       `json:"status"`
	Relationship UserRelationship `json:"relationship"`
	Online       bool             `json:"online"`
	Privileged   bool             `json:"privileged"`
	Flags        int              `json:"flags"`
	Bot          UserBot          `json:"bot"`
}

/*
This is a confusing one because the Type key is used to differeniate instead of separate objects so the fields are as follows:

SavedMessages : Id, User
DirectMessages: Id, Active, Recipients, LastMessageId
Group         : Id, Name, Owner, Description, Recipients, Icon, LastMessageId, Nsfw,
TextChannel   : Id, Server, Name, Description, Icon, LastMessageId, DefaultPermissions, RolePermissions, Nsfw
VoiceCHannel  : Id, Server, Name, Description, Icon, DefaultPermissions, RolePermissions, Nsfw

If anyone has ideas to improve this please feel free to make a PR
*/
type Channel struct {
	Id                 string                   `json:"_id"`
	ChannelType        string                   `json:"channel_type"`
	User               string                   `json:"user"`
	Name               string                   `json:"name"`
	Active             bool                     `json:"active"`
	Recipients         []string                 `json:"recipients"`
	LastMessageId      string                   `json:"last_message_id"`
	Owner              string                   `json:"owner"`
	Description        string                   `json:"description"`
	Icon               File                     `json:"icon"`
	Permissions        int                      `json:"permissions"`
	Nsfw               bool                     `json:"nsfw"`
	DefaultPermissions OverrideField            `json:"default_permissions"`
	RolePermissions    map[string]OverrideField `json:"role_permissions"`
}

type Server struct {
	Id                 string          `json:"_id"`
	Owner              string          `json:"owner"`
	Name               string          `json:"name"`
	Description        string          `json:"description"`
	Channels           []string        `json:"channels"`
	Categories         []Category      `json:"categories"`
	SystemMessages     SystemMessages  `json:"system_messages"`
	Roles              map[string]Role `json:"roles"`
	DefaultPermissions int             `json:"default_permissions"`
	Icon               File            `json:"file"`
	Banner             File            `json:"banner"`
	Nsfw               bool            `json:"nsfw"`
	Flags              int             `json:"flags"`
}

type Message struct {
	ID          string        `json:"_id"`
	Nonce       string        `json:"nonce"`
	Channel     string        `json:"channel"`
	Author      string        `json:"author"`
	Content     string        `json:"content"`
	Attachments []File        `json:"attachments"`
	Edited      *iso8601.Time `json:"edited"`
}

type SystemMessages struct {
	UserJoined string `json:"user_joined"`
	UserLeft   string `json:"user_left"`
	UserKicked string `json:"user_kicker"`
	UserBanned string `json:"user_banned"`
}

type Role struct {
	Name        string        `json:"name"`
	Permissions OverrideField `json:"permissions"`
	Colour      string        `json:"colour"`
	Hoist       bool          `json:"hoist"`
	Rank        int           `json:"rank"`
}

type OverrideField struct {
	Allowed    int `json:"a"`
	Disallowed int `json:"d"`
}

type Category struct {
	Id       string   `json:"id"`
	Title    string   `json:"title"`
	Channels []string `json:"channels"`
}

type File struct {
	Id          string       `json:"_id"`
	Tag         string       `json:"tag"`
	Name        string       `json:"filename"`
	Size        int          `json:"size"`
	Metadata    FileMetadata `json:"metadata"`
	ContentType string       `json:"content_type"`
}

type FileMetadata struct {
	Type   string `json:"type"`
	Height int    `json:"height"`
	Width  int    `json:"width"`
}

type UserStatus struct {
	Text     string `json:"text"`
	Presence string `json:"presence"`
}

type UserRelationship struct {
	Id     string `json:"_id"`
	Status string `json:"status"`
}

type UserBot struct {
	Owner string `json:"owner"`
}

type MessageUpdate struct {
	Id      string  `json:"id"`
	Channel string  `json:"channel"`
	Data    Message `json:"data"`
}
