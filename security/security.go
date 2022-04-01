package security

import (
	"InspectorManager/clients"
	"github.com/aws/aws-sdk-go-v2/service/sts"
)

// AssumeAccountRole returns an assumed role
func AssumeAccountRole(userCredentials *clients.UserCredentials, factory func(stsCredentials *clients.UserCredentials) *sts.Client,
	targetAccount string) error {
	stsClient := factory(userCredentials)
	userCredentials.SetRole(targetAccount)
	assumeRoleResult, err := stsClient.
		AssumeRole(userCredentials.UserContext, &sts.AssumeRoleInput{
			DurationSeconds: &userCredentials.SessionDuration,
			RoleArn:         &userCredentials.ServiceCredentials.AssumedRole,
			RoleSessionName: &userCredentials.SessionName,
		})

	if err != nil {
		return err
	}
	userCredentials.ServiceCredentials.AccessKeyId = *assumeRoleResult.Credentials.AccessKeyId
	userCredentials.ServiceCredentials.SecretAccessKeyId = *assumeRoleResult.Credentials.SecretAccessKey
	userCredentials.ServiceCredentials.SessionToken = *assumeRoleResult.Credentials.SessionToken

	return nil
}

// GetAWSSessionToken gets a valid aws session token.
func GetAWSSessionToken(userCredentials *clients.UserCredentials, stsClient *sts.Client) error {
	sessionTokenResult, err := stsClient.GetSessionToken(userCredentials.UserContext, &sts.GetSessionTokenInput{
		TokenCode:       &userCredentials.MFAToken,
		DurationSeconds: &userCredentials.SessionDuration,
		SerialNumber:    userCredentials.GetSerialNumber(),
	}, func(options *sts.Options) {
		options.Region = userCredentials.Region
	})

	if err != nil {
		return err
	}

	userCredentials.TemporaryCredentials.AccessKeyId = *sessionTokenResult.Credentials.AccessKeyId
	userCredentials.TemporaryCredentials.SecretAccessKeyId = *sessionTokenResult.Credentials.SecretAccessKey
	userCredentials.TemporaryCredentials.SessionToken = *sessionTokenResult.Credentials.SessionToken

	return nil
}
