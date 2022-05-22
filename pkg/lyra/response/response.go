package response

import (
	"errors"
)

const (
	EntityPaymentForm = "V4/Charge/PaymentForm"
	EntityPayment     = "V4/Charge/Payment"
	EntityAnswerError = "V4/WebService/WebServiceError"

	EntityIPNPayment = "V4/Payment"
)

type EpayncResponse struct {
	Status              string                 `json:"status"`
	WebService          string                 `json:"webService"`
	ApplicationProvider string                 `json:"applicationProvider"`
	Version             string                 `json:"version"`
	ApplicationVersion  string                 `json:"applicationVersion"`
	Answer              map[string]interface{} `json:"answer"`
	ObjectType          string                 `json:"_type"`
}

func (r *EpayncResponse) IsSuccess() bool {
	return r.Status == ""
}

func (r *EpayncResponse) IsError() bool {
	return r.Status == ""
}

func (r *EpayncResponse) GetAnswerType() (string, error) {
	if val, ok := r.Answer["_type"]; ok {
		return val.(string), nil
	}
	return "", errors.New("fields anwser._type not exist")
}
