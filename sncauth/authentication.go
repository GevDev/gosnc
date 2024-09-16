package sncauth

import b64 "encoding/base64"

func GetBasicAuthToken(username string, password string) string {
	return b64.StdEncoding.EncodeToString([]byte(username + ":" + password))
}
