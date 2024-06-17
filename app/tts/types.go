package tts

type Req struct {
	Text      string `json:"text,omitempty" binding:"required"`
	Character string `json:"character" binding:"required"`
	Language  string `json:"language" binding:"required,oneof=中文 英文 日文"`
}

type TTSSendReq struct {
	Data        []interface{} `json:"data"`
	EventData   string        `json:"event_data"`
	FnIndex     int           `json:"fn_index"`
	SessionHash string        `json:"session_hash"`
}

type MetaData struct {
	Data string `json:"data"`
	Name string `json:"name"`
}

type TTSRecvResp struct {
	Msg    string     `json:"msg"`
	Output OutputData `json:"output"`
}

type OutputData struct {
	Data []struct {
		Name string `json:"name"`
	} `json:"data"`
}
