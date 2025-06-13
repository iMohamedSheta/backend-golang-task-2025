package requests

type Request struct {
	RequestSentFields map[string]any `json:"-"`
}

func (r *Request) GetRequestSentFields() map[string]any {
	return r.RequestSentFields
}

func (r *Request) SetRequestSentFields(requestSentFields map[string]any) {
	r.RequestSentFields = requestSentFields
}
