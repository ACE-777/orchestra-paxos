package roles

type Stage string

const (
	REQUEST  Stage = "request"
	PREPARE  Stage = "prepare"
	PROMISE  Stage = "promise"
	ACCEPT   Stage = "accept"
	ACCEPTED Stage = "accepted"
	NACK     Stage = "nack"
)
