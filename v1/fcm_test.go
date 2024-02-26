package googlemessaging

import (
	"fmt"
	"os"
	"testing"
)

const (
	ServiceAccountFilePath = "fixtures/service_account.json"
)

func assertEqual(t *testing.T, v, e interface{}) {
	if v != e {
		t.Fatalf("%#v != %#v", v, e)
	}
}

func readFixture(path string) string {
	data, err := os.ReadFile(path)
	if err != nil {
		panic(fmt.Sprintf("error while reading %v", path))
	}

	return string(data)
}

func TestGetValidServiceAccount(t *testing.T) {
	serviceAccountFileContent := readFixture(ServiceAccountFilePath)

	acc, err := getServiceAccountFromContent(serviceAccountFileContent)

	assertEqual(t, err, nil)
	assertEqual(t, acc.ProjectId, "pusher-project-id")
	assertEqual(t, acc.PrivateKeyId, "12345")
	assertEqual(t, acc.PrivateKey, "-----BEGIN PRIVATE KEY-----\nprivate_key\n-----END PRIVATE KEY-----\n")
	assertEqual(t, acc.ClientEmail, "firebase-adminsdk-g680q@pusher-project-id.iam.gserviceaccount.com")
	assertEqual(t, acc.ClientId, "12345")
	assertEqual(t, acc.AuthUri, "https://accounts.google.com/o/oauth2/auth")
	assertEqual(t, acc.TokenUri, "https://oauth2.googleapis.com/token")
}
