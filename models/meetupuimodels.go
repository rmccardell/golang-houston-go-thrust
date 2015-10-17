package models

//An object that encapsulates a meet up oauth2 token response.
//MeetUpToken...
type MeetUpToken struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int64  `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
}

//Used for messaging pasing
//godoc comment for RPCMessage
type RPCMessage struct {
	MessageID string
	Content   string
}
