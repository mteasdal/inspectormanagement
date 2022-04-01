package inspector

import (
	"github.com/aws/aws-sdk-go-v2/service/inspector2/types"
	"testing"
)

func TestInspectorFilterPipeline_PopulateAccountFilters(t *testing.T) {
	filterPipeline := InspectorFilterPipeline{
		AWSAccounts: []string{"12345348035"},
	}

	filterPipeline.PopulateAccountFilters("EQUALS")
	if filterPipeline.FilterError != nil {
		t.Errorf("Error populating account filters %s", filterPipeline.FilterError.Error())
	}

	if len(filterPipeline.AccountFilters) != 1 {
		t.Errorf("Error populating account filters expecting 1 gpt %d", len(filterPipeline.AWSAccounts))
	}

	if *filterPipeline.AccountFilters[0].Value != "12345348035" {
		t.Errorf("Error populating account filters account expecting 12345348035 got %s",
			*filterPipeline.AccountFilters[0].Value)
	}
}

func TestInspectorFilterPipeline_CreateFilterRequest(t *testing.T) {
	filterPipeline := InspectorFilterPipeline{
		AWSAccounts: []string{"12345348035"},
	}

	filterPipeline.Action = types.FilterAction("NONE")

	filterPipeline.PopulateAccountFilters("EQUALS")
	filterPipeline.CreateFilterRequest()
}
