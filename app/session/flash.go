package session

type FlashType string

const (
	FlashType_Info    FlashType = "info"
	FlashType_Success FlashType = "success"
	FlashType_Error   FlashType = "error"
)

type FlashMessage struct {
	Type    FlashType
	Content string
}

func (m FlashMessage) IsInfo() bool {
	return m.Type == FlashType_Info
}

func (m FlashMessage) IsSuccess() bool {
	return m.Type == FlashType_Success
}

func (m FlashMessage) IsError() bool {
	return m.Type == FlashType_Error
}
