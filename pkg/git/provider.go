package git

import "strings"

const (
	ProviderGitHub      = "github.com"
	ProviderAzureDevOps = "dev.azure.com"
	ProviderUnknown     = "unknown Git provider"

	Github      = "github"
	AzureDevOps = "azure-devops"
)

var (
	Providers = map[string]string{
		ProviderGitHub:      Github,
		ProviderAzureDevOps: AzureDevOps,
	}
)

func Provider(uri string) string {
	if strings.Contains(uri, ProviderGitHub) {
		return ProviderGitHub
	}

	if strings.Contains(uri, ProviderAzureDevOps) {
		return ProviderAzureDevOps
	}

	return ProviderUnknown
}
