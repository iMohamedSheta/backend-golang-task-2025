package contracts

type Validatable interface {
	Messages() map[string]string
	GetRequestSentFields() map[string]any
	SetRequestSentFields(map[string]any)
}
