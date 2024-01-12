package msg

type (
	// Success reply type for generic successful replies to commands
	Success struct{}
	// Failure reply type for generic failure replies to commands
	Failure struct{}
)

// ReplyName implements core.Reply.ReplyName
func (Success) ReplyName() string { return "edat.msg.Success" }

// ReplyName implements core.Reply.ReplyName
func (Failure) ReplyName() string { return "edat.msg.Failure" }
