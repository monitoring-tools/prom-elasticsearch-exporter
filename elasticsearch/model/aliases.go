package model

// Aliases is a representation of ElasticSearch index aliases
type Aliases map[string]AliasInfo

// AliasInfo is a representation of ElasticSearch index alias
type AliasInfo struct {
	Aliases map[string]interface{} `json:"aliases"`
}
