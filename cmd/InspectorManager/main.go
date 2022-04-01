package main

import (
	"InspectorManager/auditing"
	"InspectorManager/clients"
	"InspectorManager/inspector"
	"InspectorManager/security"
	"context"
	"flag"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/service/inspector2/types"
	"os"
)

func main() {
	awsAccount := flag.String("account", "", "Aws account")
	userName := flag.String("username", "", "Aws username")
	region := flag.String("region", "eu-west-1", "Default aws region")
	profile := flag.String("profile", "", "Profile for the aws account")
	action := flag.String("action", "NONE", "Filter action")
	filterName := flag.String("filter-name", "test", "name of filter")
	mfaToken := flag.String("mfa", "", "MFA Token")

	flag.Parse()

	logonCredentials := &clients.UserCredentials{
		AwsAccount:      *awsAccount,
		UserName:        *userName,
		MFAToken:        *mfaToken,
		Region:          *region,
		Profile:         *profile,
		FilterName:      *filterName,
		UserContext:     context.TODO(),
		SessionDuration: 3600,
		SessionName:     "inspector",
	}

	//Sets the default config
	logonCredentials.SetDefaultConfig()

	// Get Session Token
	stsFactory := clients.NewSTSClient()
	stsClient := stsFactory(logonCredentials.UserConfig)

	err := security.GetAWSSessionToken(logonCredentials, stsClient)

	if err != nil {
		auditing.Log(err.Error())
		os.Exit(1)
	}

	stsServiceFactory := clients.NewSTSClientSessionConfig()
	err = security.AssumeAccountRole(logonCredentials, stsServiceFactory, logonCredentials.AwsAccount)

	if err != nil {
		auditing.Log(err.Error())
		os.Exit(1)
	}

	inspectorFactory := clients.NewInspectorClientFactory()
	inspectorClient := inspectorFactory(logonCredentials.UserConfig, *logonCredentials)

	filterPipeline := inspector.InspectorFilterPipeline{
		AWSAccounts: []string{logonCredentials.AwsAccount},
		Action:      types.FilterAction(*action),
		FilterName:  logonCredentials.FilterName,
	}

	filterPipeline.PopulateAccountFilters("EQUALS").CreateFilterRequest().ProcessFilterRequest(inspectorClient, logonCredentials.UserContext)

	if filterPipeline.FilterError != nil {
		fmt.Printf("Error processing pipeline %s", filterPipeline.FilterError.Error())
		auditing.Log(filterPipeline.FilterError.Error())
	}

	fmt.Printf("Filter Output %s", *filterPipeline.FilterResponse.Arn)
}
