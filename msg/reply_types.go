package msg

type (
	// Success reply type for generic successful replies to commands
	Success struct{}
	// Failure reply type for generic failure replies to commands
	Failuer struct{}
)

// ReplyName implements core.Reply.ReplyName
func (Success) ReplyName() string { return "edat.msg.Success" }

// ReplyName implements core.Reply.ReplyName
func (Failuer) ReplyName() string { return "edat.msg.Failuer" }
