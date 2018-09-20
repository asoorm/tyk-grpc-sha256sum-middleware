package main

import (
	"crypto/sha256"
	"fmt"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/TykTechnologies/tyk-protobuf/bindings/go"
)

const sharedSecret = "mySharedSecret"

func Sha256SumHook(object *coprocess.Object) (*coprocess.Object, error) {

	authKey := object.Request.Headers["Authorization"]
	authSignature := object.Request.Headers["X-Signature"]

	currentTime := time.Now().Unix()

	signatureValid := false
	for i := int64(0); i < 300; i++ {
		raw := fmt.Sprintf("%s%s%d", authKey, sharedSecret, currentTime+i)
		if checkShaSum32(raw, authSignature) {
			signatureValid = true
		}

		raw = fmt.Sprintf("%s%s%d", authKey, sharedSecret, currentTime-i)
		if checkShaSum32(raw, authSignature) {
			signatureValid = true
		}
	}

	if !signatureValid {
		log.Println("Bad auth on Sha256SumHook")
		return object, nil
	}

	log.Println("Successful authentication on Sha256SumHook, setting session token")

	// Set the ID extractor deadline, useful for caching valid keys:
	extractorDeadline := time.Now().Add(time.Hour * 1).Unix()

	object.Session = &coprocess.SessionState{
		Rate:                1000.0,
		Per:                 1.0,
		QuotaMax:            int64(1000),
		QuotaRenews:         time.Now().Unix(),
		IdExtractorDeadline: extractorDeadline,
	}

	object.Metadata = map[string]string{
		"token":     authKey,
		"signature": authSignature,
	}

	return object, nil
}

func checkShaSum32(raw string, authSignature string) bool {
	b := sha256.Sum256([]byte(raw))
	return string(b[:]) == authSignature
}
