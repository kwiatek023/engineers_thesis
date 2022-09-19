package simulation

// JsonStatsStructure - structure for saving statistics to file
type JsonStatsStructure struct {
	Size               int       `json:"size"`
	Result             float64   `json:"result"`
	NofRounds          int       `json:"nof_rounds"`
	MaxReceivedMsgs    int       `json:"max_received_msgs"`
	MinReceivedMsgs    int       `json:"min_received_msgs"`
	AllReceivedMsgs    int       `json:"all_received_msgs"`
	AvgReceivedMsgs    float64   `json:"avg_received_msgs"`
	StddevReceivedMsgs float64   `json:"stddev_received_msgs"`
	MaxSentMsgs        int       `json:"max_sent_msgs"`
	MinSentMsgs        int       `json:"min_sent_msgs"`
	AllSentMsgs        int       `json:"all_sent_msgs"`
	AvgSentMsgs        float64   `json:"avg_sent_msgs"`
	StddevSentMsgs     float64   `json:"stddev_sent_msgs"`
	AllMemory          int       `json:"all_memory"`
	MaxMemory          int       `json:"max_memory"`
	MinMemory          int       `json:"min_memory"`
	AvgMemory          float64   `json:"avg_memory"`
	StddevMemory       float64   `json:"stddev_memory"`
	Stations           []Station `json:"stations"`
}
