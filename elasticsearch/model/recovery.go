package model

// Recovery is a ES /_recovery representation
type Recovery map[string]IndexRecovery

// IndexRecovery is a index recovery representation
type IndexRecovery struct {
	Shards []ShardRecovery `json:"shards"`
}

// IndexRecoveryState is an index recovery statement representation
type IndexRecoveryState struct {
	Size                       IndexRecoverySize  `json:"size"`
	Files                      IndexRecoveryFiles `json:"files"`
	TotalTimeInMillis          int64              `json:"total_time_in_millis"`
	SourceThrottleTimeInMillis int64              `json:"source_throttle_time_in_millis"`
	TargetThrottleTimeInMillis int64              `json:"target_throttle_time_in_millis"`
}

// IndexRecoverySize is an index recovery size representation
type IndexRecoverySize struct {
	TotalInBytes     int64  `json:"total_in_bytes"`
	ReusedInBytes    int64  `json:"reused_in_bytes"`
	RecoveredInBytes int64  `json:"recovered_in_bytes"`
	Percent          string `json:"percent"`
}

// IndexRecoveryFiles is an index recovery files representation
type IndexRecoveryFiles struct {
	Total     int64  `json:"total"`
	Reused    int64  `json:"reused"`
	Recovered int64  `json:"recovered"`
	Percent   string `json:"percent"`
}

// IndexRecoveryTranslog is an index recovery translog representation
type IndexRecoveryTranslog struct {
	Recovered         int64  `json:"recovered"`
	Total             int64  `json:"total"`
	Percent           string `json:"percent"`
	TotalOnStart      int64  `json:"total_on_start"`
	TotalTimeInMillis int64  `json:"total_time_in_millis"`
}

// ShardRecovery is an index shard recovery representation
type ShardRecovery struct {
	ID                int64                 `json:"id"`
	Type              string                `json:"type"`
	Stage             string                `json:"stage"`
	Primary           bool                  `json:"primary"`
	StartTimeInMillis int64                 `json:"start_time_in_millis"`
	TotalTimeInMillis int64                 `json:"total_time_in_millis"`
	Source            RecoveryDestination   `json:"source"`
	Target            RecoveryDestination   `json:"target"`
	Index             IndexRecoveryState    `json:"index"`
	Translog          IndexRecoveryTranslog `json:"translog"`
	VerifyIndex       struct {
		CheckIndexTimeInMillis int64 `json:"check_index_time_in_millis"`
		TotalTimeInMillis      int64 `json:"total_time_in_millis"`
	} `json:"verify_index"`
}

// RecoveryDestination is a shard recovery source / target
type RecoveryDestination struct {
	ID               string `json:"id"`
	Host             string `json:"host"`
	TransportAddress string `json:"transport_address"`
	IP               string `json:"ip"`
	Name             string `json:"name"`
}
