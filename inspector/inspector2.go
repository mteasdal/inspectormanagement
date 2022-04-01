package inspector

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/service/inspector2"
	"github.com/aws/aws-sdk-go-v2/service/inspector2/types"
)

type InspectorFilterPipeline struct {
	AWSAccounts     []string
	CVETitles       []string
	AccountFilters  []types.StringFilter
	CVETitleFilters []types.StringFilter
	FilterRequest   *inspector2.CreateFilterInput
	FilterResponse  *inspector2.CreateFilterOutput
	Action          types.FilterAction
	FilterName      string
	FilterError     error
}

// GetFilterOnCVETitle creates a filter for inspector
func GetFilterOnCVETitle(cveTitle string, cveComparisonOperator string) types.StringFilter {

	filterType := types.StringFilter{
		Comparison: types.StringComparison(cveComparisonOperator),
		Value:      &cveTitle,
	}

	return filterType
}

// getFilterOnAWSAccount creates a filter on account
func getFilterOnAWSAccount(awsAccount string, accountComparisonOperator string) types.StringFilter {
	filterType := types.StringFilter{
		Comparison: types.StringComparison(accountComparisonOperator),
		Value:      &awsAccount,
	}
	return filterType
}

func (i *InspectorFilterPipeline) PopulateAccountFilters(comparisonOperator string) *InspectorFilterPipeline {
	var accountFilters []types.StringFilter
	if len(i.AWSAccounts) > 0 {
		for _, awsAccount := range i.AWSAccounts {
			accountFilter := getFilterOnAWSAccount(awsAccount, comparisonOperator)
			accountFilters = append(accountFilters, accountFilter)
		}
	}
	i.AccountFilters = accountFilters
	return i
}

func (i *InspectorFilterPipeline) PopulateTitleFilters(comparisonOperator string) *InspectorFilterPipeline {
	var titleFilters []types.StringFilter
	if len(i.CVETitles) > 0 {
		for _, cveTitle := range i.CVETitles {
			accountFilter := getFilterOnAWSAccount(cveTitle, comparisonOperator)
			titleFilters = append(titleFilters, accountFilter)
		}
	}
	i.CVETitleFilters = titleFilters
	return i
}

// CreateFilterRequest adds the filters to Inspector
func (i *InspectorFilterPipeline) CreateFilterRequest() *InspectorFilterPipeline {
	filterRequest := inspector2.CreateFilterInput{}
	filterRequest.Action = i.Action
	filterRequest.Name = &i.FilterName
	if len(i.AccountFilters) > 0 {
		f := types.FilterCriteria{AwsAccountId: i.AccountFilters}
		filterRequest.FilterCriteria = &f
	}
	i.FilterRequest = &filterRequest
	return i
}

// ProcessFilterRequest processes the filter request to inspector
func (i *InspectorFilterPipeline) ProcessFilterRequest(client *inspector2.Client,
	ctx context.Context) *InspectorFilterPipeline {
	fmt.Printf("Processing Filter Request")

	response, err := client.CreateFilter(ctx, i.FilterRequest)
	if err != nil {
		i.FilterError = err
	} else {
		i.FilterResponse = response
	}
	return i
}
