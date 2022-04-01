package clients

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/inspector2"
	"github.com/aws/aws-sdk-go-v2/service/sts"
)

// UserCredentials is the initial logon credentials of the user/service
type UserCredentials struct {
	AwsAccount           string
	UserName             string
	MFAToken             string
	Region               string
	Profile              string
	FilterName           string
	UserContext          context.Context
	UserConfig           aws.Config
	SessionName          string
	SessionDuration      int32
	TemporaryCredentials Credentials
	ServiceCredentials   ReturnedCredentials
}

// Credentials holds the session credentials used for client creation
type Credentials struct {
	AccessKeyId       string
	SecretAccessKeyId string
	SessionToken      string
}

// ReturnedCredentials holds the client credentials for working with service clients
type ReturnedCredentials struct {
	AccessKeyId       string
	SecretAccessKeyId string
	SessionToken      string
	AssumedRole       string
}

const userAccount = "083917714948"

// SetDefaultConfig loads up the credentials from the .aws folder
func (u *UserCredentials) SetDefaultConfig() {
	u.UserConfig, _ = config.LoadDefaultConfig(u.UserContext,
		config.WithSharedConfigProfile(u.Profile), config.WithRegion(u.Region))
}

func (u *UserCredentials) SetRole(targetAccount string) {
	roleToAssume := fmt.Sprintf("arn:aws:iam::%s:role/hrk-role-inspector-reporter", targetAccount)
	u.ServiceCredentials.AssumedRole = roleToAssume
}

func (u *UserCredentials) GetSerialNumber() *string {
	serialNumber := fmt.Sprintf("arn:aws:iam::%s:mfa/%s", userAccount, u.UserName)
	fmt.Printf("***DEBUG*** %s", serialNumber)
	return &serialNumber
}

// NewSTSClient returns an STS client generated from default config
func NewSTSClient() func(cfg aws.Config) *sts.Client {
	return func(cfg aws.Config) *sts.Client {
		return sts.NewFromConfig(cfg)
	}
}

// NewSTSClientSessionConfig returns sts client generated from session config
func NewSTSClientSessionConfig() func(stsCredentials *UserCredentials) *sts.Client {
	return func(stsCredentials *UserCredentials) *sts.Client {
		return sts.New(sts.Options{
			Region: stsCredentials.Region,
			Credentials: aws.NewCredentialsCache(credentials.NewStaticCredentialsProvider(stsCredentials.TemporaryCredentials.AccessKeyId,
				stsCredentials.TemporaryCredentials.SecretAccessKeyId, stsCredentials.TemporaryCredentials.SessionToken)),
		})
	}
}

// NewInspectorClientFactory returns an Inspector client generated from default config
func NewInspectorClientFactory() func(cfg aws.Config, stsCredentials UserCredentials) *inspector2.Client {
	return func(cfg aws.Config, stsCredentials UserCredentials) *inspector2.Client {
		return inspector2.New(inspector2.Options{
			Region: stsCredentials.Region,
			Credentials: aws.NewCredentialsCache(credentials.NewStaticCredentialsProvider(stsCredentials.ServiceCredentials.AccessKeyId,
				stsCredentials.ServiceCredentials.SecretAccessKeyId, stsCredentials.ServiceCredentials.SessionToken)),
		})
	}
}
