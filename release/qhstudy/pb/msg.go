package pb

type UserConfig struct {
	StreamingUri string `json:"streamingUri"`
	RtmpHost     string `json:"rtmpHost"`
	RtmpChannel  string `json:"rtmpChannel"`
}

// ==============================new==================================
type Sync_Hello struct {
	Ip      string //ip
	Port    int    //端口
	GinPort int    //http端口
}
