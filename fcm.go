package googlemessaging

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"golang.org/x/oauth2/google"
	"golang.org/x/oauth2/jwt"
)

const (
	BaseURL              = "https://fcm.googleapis.com/v1"
	FirebaseRequestScope = "https://www.googleapis.com/auth/firebase.messaging"
	DefaultContentType   = "application/json"
)

type ServiceAccount struct {
	ServiceAccountType string `json:"type"`
	ProjectId          string `json:"project_id"`
	PrivateKeyId       string `json:"private_key_id"`
	PrivateKey         string `json:"private_key"`
	ClientEmail        string `json:"client_email"`
	ClientId           string `json:"client_id"`
	AuthUri            string `json:"auth_uri"`
	TokenUri           string `json:"token_uri"`
}

type fcmClient struct {
	httpClient           *http.Client
	ServiceAccountConfig *ServiceAccount
}

func getServiceAccountFromContent(content string) (*ServiceAccount, error) {
	serviceAccount := &ServiceAccount{}
	err := json.Unmarshal([]byte(content), &serviceAccount)
	if err != nil {
		return nil, fmt.Errorf("error umarshalling service account file>%v", err)
	}

	return serviceAccount, nil
}

func NewFcmClient(serviceAccountFileContent string, ctx context.Context) (*fcmClient, error) {
	serviceAccountConfig, err := getServiceAccountFromContent(serviceAccountFileContent)
	if err != nil {
		return nil, err
	}

	conf := &jwt.Config{
		Email:        serviceAccountConfig.ClientEmail,
		PrivateKeyID: serviceAccountConfig.PrivateKeyId,
		PrivateKey:   []byte(serviceAccountConfig.PrivateKey),
		Scopes:       []string{FirebaseRequestScope},
		TokenURL:     google.JWTTokenURL,
	}

	return &fcmClient{
		httpClient:           conf.Client(ctx),
		ServiceAccountConfig: serviceAccountConfig,
	}, nil
}

func (c *fcmClient) Send(m FcmMessageBody) (*FcmSendHttpResponse, error) {
	sendURL := fmt.Sprintf("%s/projects/%s/messages:send",
		BaseURL, c.ServiceAccountConfig.ProjectId,
	)

	body, err := json.Marshal(m)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling message>%v", err)
	}

	resp, err := c.httpClient.Post(sendURL, DefaultContentType, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("error sending request to HTTP connection server>%v", err)
	}

	responseBody, _ := io.ReadAll(resp.Body)
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("could not send a message as the server returned: %s", resp.StatusCode)
	}

	fcmResp := &FcmSendHttpResponse{}
	err = json.Unmarshal(responseBody, &fcmResp)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling json from body: %v", err)
	}

	return fcmResp, nil
}

func SendPush(ctx context.Context, serviceAccountConfig string, m FcmMessageBody) (*FcmSendHttpResponse, error) {
	c, err := NewFcmClient(serviceAccountConfig, ctx)
	if err != nil {
		return nil, err
	}

	return c.Send(m)
}
