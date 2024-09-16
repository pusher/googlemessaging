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
	InstanceIdApiUrl     = "https://iid.googleapis.com/iid/info"
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

// GetInstanceInfo returns app instance info. Reference: https://developers.google.com/instance-id/reference/server
func (c *fcmClient) GetInstanceInfo(token string) (*InstanceInformationResponse, error) {
	instanceInfoUrl := fmt.Sprintf("%s/%s", InstanceIdApiUrl, token)

	request, err := http.NewRequest(http.MethodGet, instanceInfoUrl, nil)
	if err != nil {
		return nil, fmt.Errorf("error build http request to %s", instanceInfoUrl)
	}

	request.Header.Set("access_token_auth", "true")
	response, err := c.httpClient.Do(request)
	if err != nil {
		return nil, fmt.Errorf("error making a http request to %s > %v", instanceInfoUrl, err)
	}
	defer response.Body.Close()

	switch response.StatusCode {
	case 400:
		errorResponseBytes, err := io.ReadAll(response.Body)
		if err != nil {
			return nil, fmt.Errorf("unable to read response from google: %w", err)
		}

		errorBody := ErrorResponse{}
		err = json.Unmarshal(errorResponseBytes, &errorBody)
		if err != nil {
			return nil, fmt.Errorf("unable to read response from google: %w", err)
		}

		if errorBody.Error == "InvalidToken" {
			return nil, ErrInvalidToken
		}

		return nil, ErrInvalidRequest

	case 401, 403:
		return nil, ErrInalidFCMServiceAccountFile
	case 404:
		return nil, ErrDeviceNotFound
	}

	if response.StatusCode > 299 {
		return nil, fmt.Errorf("unexpected status code: %d", response.StatusCode)
	}

	responseBodyBytes, err := io.ReadAll(response.Body)
	if err != nil {
		return &InstanceInformationResponse{}, fmt.Errorf("failed to read response body > %v", err)
	}

	responseBody := &InstanceInformationResponse{}
	err = json.Unmarshal(responseBodyBytes, &responseBody)
	if err != nil {
		return nil, err
	}

	return responseBody, nil
}

func (c *fcmClient) Send(m FcmMessageBody) (*FcmSendHttpResponse, error) {
	sendURL := fmt.Sprintf("%s/projects/%s/messages:send",
		BaseURL, c.ServiceAccountConfig.ProjectId,
	)

	body, err := json.Marshal(m)
	if err != nil {
		return nil, fmt.Errorf("error marshaling message>%v", err)
	}

	resp, err := c.httpClient.Post(sendURL, DefaultContentType, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("error sending request to HTTP connection server>%v", err)
	}

	defer resp.Body.Close()
	fcmResp := &FcmSendHttpResponse{Status: resp.StatusCode}
	responseBody, err := io.ReadAll(resp.Body)

	if err != nil {
		return nil, fmt.Errorf("failed to read response body > %v", err)
	}

	if fcmResp.Status != http.StatusOK {
		return fcmResp, fmt.Errorf("could not send a message as the server returned: %d", resp.StatusCode)
	}

	err = json.Unmarshal(responseBody, &fcmResp)
	if err != nil {
		return fcmResp, fmt.Errorf("error unmarshaling json from body: %v", err)
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
