package common

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"time"

	bitbucketv1 "github.com/gfleury/go-bitbucket-v1"
	"github.com/spf13/cobra"
)

// APIClientContext initialize all configuration for CLI
func APIClientContext(stashInfo *StashInfo) context.Context {

	ctx := context.Background()
	ctx = context.WithValue(ctx, StashInfoKey, stashInfo)

	return ctx
}

// APIClient return Stash APIClient
func APIClient(cmd *cobra.Command) (client *bitbucketv1.APIClient, cancel context.CancelFunc, err error) {
	ctx := cmd.Context()
	stashInfo := ctx.Value(StashInfoKey).(*StashInfo)

	ctx, cancel = context.WithTimeout(ctx, time.Duration(*stashInfo.Timeout())*time.Second)

	if stashInfo.host == "" || stashInfo.project == "" || stashInfo.repo == "" {
		host, project, repo, _ := GitInfo()
		if stashInfo.host == "" && host == "" {
			keys := Config().AllKeys()
			sort.Strings(keys)
			for _, key := range keys {
				if !strings.HasPrefix(key, "password_method") {
					stashInfo.host = strings.Split(key, "::")[0]
					break
				}
			}
		} else if stashInfo.host == "" {
			stashInfo.host = host
		}
		if stashInfo.project == "" {
			stashInfo.project = project
		}
		if stashInfo.repo == "" {
			stashInfo.repo = repo
		}

		if !strings.HasPrefix(stashInfo.host, "https://") &&
			!strings.HasPrefix(stashInfo.host, "http://") {
			stashInfo.host = fmt.Sprintf("https://%s", stashInfo.host)
		}
		if stashInfo.host == "" {
			err = fmt.Errorf("Unable to autodetect Stash host/url. Please inform with -H flag")
		}
	}

	basicAuth := bitbucketv1.BasicAuth{UserName: stashInfo.Credential().GetUser(stashInfo.host), Password: stashInfo.Credential().GetPasswd()}
	ctx = context.WithValue(ctx, bitbucketv1.ContextBasicAuth, basicAuth)

	client = bitbucketv1.NewAPIClient(
		ctx,
		bitbucketv1.NewConfiguration(fmt.Sprintf("%s/rest", stashInfo.host)),
	)
	return client, cancel, err
}
