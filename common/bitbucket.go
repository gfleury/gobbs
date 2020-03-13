package common

import (
	"context"
	"fmt"
	"time"

	bitbucketv1 "github.com/gfleury/go-bitbucket-v1"
	"github.com/spf13/cobra"
)

// APIClientContext initialize all configuration for CLI
func APIClientContext(stashInfo *StashInfo) (context.Context, context.CancelFunc) {

	ctx, cancel := context.WithTimeout(context.Background(), 15000*time.Millisecond)
	ctx = context.WithValue(ctx, StashInfoKey, stashInfo)

	return ctx, cancel
}

// APIClient return Stash APIClient
func APIClient(cmd *cobra.Command) *bitbucketv1.APIClient {
	ctx := cmd.Context()
	stashInfo := cmd.Context().Value(StashInfoKey).(*StashInfo)

	if stashInfo.host == "" && stashInfo.project == "" && stashInfo.repo == "" {
		stashInfo.host, stashInfo.project, stashInfo.repo = GitInfo()
		stashInfo.host = fmt.Sprintf("https://%s", stashInfo.host)
	}

	basicAuth := bitbucketv1.BasicAuth{UserName: stashInfo.Credential().GetUser(), Password: stashInfo.Credential().GetPasswd()}
	ctx = context.WithValue(ctx, bitbucketv1.ContextBasicAuth, basicAuth)

	client := bitbucketv1.NewAPIClient(
		ctx,
		bitbucketv1.NewConfiguration(fmt.Sprintf("%s/rest", stashInfo.host)),
	)
	return client
}
