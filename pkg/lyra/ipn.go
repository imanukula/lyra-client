package lyra

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
)

type IPNMessage struct {
	Hash          string `form:"kr-hash"`
	HashAlgorithm string `form:"kr-hash-algorithm"`
	HashKey       string `form:"kr-hash-key"`
	AnswerType    string `form:"kr-answer-type"`
	Answer        string `form:"kr-answer"`
}

// CheckHash - check kr-answer object signature
// docs : https://epaync.nc/doc/fr-FR/rest/V4.0/kb/payment_done.html#v%C3%A9rifier-la-signature-navigateur-hash
func CheckHash(params IPNMessage) error {
	// check if the hash algorithm is supported
	supported := []string{"sha256_hmac"}
	if !contains(supported, params.HashAlgorithm) {
		return errors.New(fmt.Sprint("hash algorithm not supported : %s", params.HashAlgorithm))
	}

	// if key is not defined, we use kr-hash-key POST parameter to choose it
	key := ""
	if len(params.HashKey) == 0 {
		return errors.New("invalid kr-hash-key POST parameter")
	} else {
		switch params.HashKey {
		case "sha256_hmac":
			key = DefaultHashKey
		case "password":
			key = DefaultPassword
		default:
			return errors.New("invalid kr-hash-key POST parameter")
		}
	}

	// return nil if calculated hash and sent hash are the same
	h := hmac.New(sha256.New, []byte(key))
	h.Write([]byte(params.Answer))

	if params.Hash != hex.EncodeToString(h.Sum(nil)) {
		log.Printf("%v", hex.EncodeToString(h.Sum(nil)))
		log.Printf("%v", params.Hash)
		return errors.New("hash answer and answer not match")
	}

	return nil
}

func contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}
	return false
}
