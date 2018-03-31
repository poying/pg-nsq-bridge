package bridge

type Payload struct {
	Pid             int    `json:"pid"`
	PostgresChannel string `json:"postgres_channel"`
	Data            string `json:"data"`
}
