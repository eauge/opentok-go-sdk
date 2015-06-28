package opentok

import "encoding/xml"

type xmlSessions struct {
	XMLName  xml.Name     `xml:"sessions"`
	Sessions []xmlSession `xml:"Session"`
}

type xmlSession struct {
	XMLName    xml.Name `xml:"Session"`
	PartnerID  int      `xml:"partner_id"`
	SessionID  string   `xml:"session_id"`
	CreateDate string   `xml:"create_dt"`
}

// MediaMode specifies how the streams will be managed
// by the OpenTok platform
type MediaMode string

const (

	// Routed the streams are handled by our routing component, Mantis.
	// This mode is useful for sessions with more than 3 participants
	// or sessions that need archiving
	Routed MediaMode = "disabled"

	// Relayed the streams are either send P2P, if possible. If that's not
	// possible, they are routed to a TURN server. When using this mode,
	// archiving will not be available.
	// Useful for sessions with up to three people that do not need
	// archiving
	Relayed = "enabled"
)

// ArchiveMode is the archiving mode that is used by the session
type ArchiveMode string

const (
	// Always archiveMode will archive the session always, and
	// ArchiveStart and ArchiveStop cannot be used
	Always ArchiveMode = "always"

	// Manual is the default archiveMode. The client can use
	// ArchiveStart and ArchiveStop to start and stop archiving
	Manual ArchiveMode = "manual"
)

// SessionProps contains the different options when creating a
// new session object
type SessionProps struct {

	// Location an IP address that enables the user to specifiy a preferred
	// location to be considered when creating the session in
	// the platform
	Location    string
	MediaMode   MediaMode
	ArchiveMode ArchiveMode
}

// Role for a client connected to an OpenTok Session
type Role string

const (

	// Moderator has the highest level of control in the session.
	// Can subscribe and pubslish to a session and can also force
	// the disconnections of other participants in the session
	Moderator Role = "moderator"

	// Publisher is the default role. It allows the participants to
	// subscribe to the streams sent by other participants as well as
	// pubslish its own stream to the session
	Publisher Role = "publisher"

	// Subscriber is the role with the lowest control. It's only allowed
	// to subscribe to the streams in the session published by publishers
	// and moderators and cannot publish its stream
	Subscriber Role = "subscriber"
)

// Token used by the opentok clients to connect to a session
type Token string

func (t *Token) String() string {
	return string(*t)
}

// TokenProps are the properties that can be set when creating
// a token
type TokenProps struct {
	// Role that the client using that token will have in the opentok
	// session
	Role Role

	// ExpireTime is the time that the token can be used to access a
	// token before it expires
	ExpireTime int64

	// Data is any extra data that needs to be added to the token. It
	// can be used to store information about the client that will use
	// it to connect to the OpenTok Session
	Data string
}

// Session struct that represents an OpenTok Session
type Session struct {
	ID string
}
