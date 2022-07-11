package drpc

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type Opcode = uint32

const (
	OpHandshake Opcode = iota
	OpFrame
	OpClose
	OpPing
	OpPong
)

// Message describes structure of a message.
type Message struct {
	Opcode
	Payload []byte
}

// Handshake is a handshake message.
type Handshake struct {
	V        string `json:"v"`
	ClientID string `json:"client_id"`
}

// FrameHeader is a header for a frame.
type FrameHeader struct {
	Command string `json:"cmd"`
	Event   string `json:"evt"`
	Nonce   string `json:"nonce"`
}

// NewFrameHeader creates a new frame header with the given command.
func NewFrameHeader(command string) FrameHeader {
	return FrameHeader{
		Command: command,
		Nonce:   uuid.NewString(),
	}
}

// NewFrameHeaderWithEvent creates a new frame header with the given command and event.
func NewFrameHeaderWithEvent(command, event string) FrameHeader {
	return FrameHeader{
		Command: command,
		Event:   event,
		Nonce:   uuid.NewString(),
	}
}

// FrameSetActivity is a frame of the SetActivity message.
type FrameSetActivity struct {
	FrameHeader
	Args FrameSetActivityArgs `json:"args"`
}

// FrameSetActivityArgs is a set of arguments for the SetActivity message.
type FrameSetActivityArgs struct {
	PID      int `json:"pid"`
	Activity `json:"activity"`
}

// Activity describes an user's current activity.
type Activity struct {
	// What the player is currently doing.
	Details string `json:"details,omitempty"`
	// The user's current state.
	State string `json:"state,omitempty"`
	// Timestamps for the start and end of the activity.
	Timestamps *Timestamps `json:"timestamps,omitempty"`
	// Profile artworks.
	Assets *Assets `json:"assets,omitempty"`
	// Party information.
	Party *Party `json:"party,omitempty"`
	// Secrets data. It cannot be used with the Buttons field.
	Secrets *Secrets `json:"secrets,omitempty"`
	// Buttons to display. It cannot be used with the Secrets field.
	Buttons []Button `json:"buttons,omitempty"`
}

// Timestamps contains timestamps for the start and end of the activity.
//
// Sending End timestamp will always have the time displayed as "remaining" until the given time.
// Sending Start timestamp will show "elapsed" as long as there is no End sent.
type Timestamps struct {
	// Epoch seconds for game start - including will show time as "elapsed".
	Start time.Time `json:"start,omitempty"`
	// Epoch seconds for game end - including will show time as "remaining".
	End time.Time `json:"end,omitempty"`
}

// Custom marshal method for Timestamps.
// It is required because protocol accepts only epoch seconds.
func (t Timestamps) MarshalJSON() ([]byte, error) {
	m := make(map[string]int64)
	if !t.Start.IsZero() {
		m["start"] = t.Start.Unix()
	}
	if !t.End.IsZero() {
		m["end"] = t.End.Unix()
	}
	return json.Marshal(m)
}

// Assets contains an user's profile artwork.
type Assets struct {
	// Name of the uploaded image for the large profile artwork.
	LargeImage string `json:"large_image,omitempty"`
	// Tooltip for the large profile artwork.
	LargeText string `json:"large_text,omitempty"`
	// Name of the uploaded image for the small profile artwork.
	SmallImage string `json:"small_image,omitempty"`
	// Tooltip for the small profile artwork.
	SmallText string `json:"small_text,omitempty"`
}

// Party contains party information.
type Party struct {
	// ID of the player's party, lobby, or group.
	ID string `json:"id,omitempty"`
	// The first element is the current party size, the second element is the max party size.
	Size [2]int `json:"size,omitempty"`
}

// Secrets contains secrets for joining and spectating.
type Secrets struct {
	// Unique hashed string for chat invitations and Ask to Join.
	Join string `json:"join,omitempty"`
	// Unique hashed string for Spectate button.
	Spectate string `json:"spectate,omitempty"`
	// (For future use) Unique hashed string for a player's match.
	Match string `json:"match,omitempty"`
}

// Button contains a button data.
// Both fields are required.
type Button struct {
	// Name of the button.
	Label string `json:"label,omitempty"`
	// URL to open when the button is clicked.
	URL string `json:"url,omitempty"`
}
