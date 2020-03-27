package search

import (
	"context"
	"errors"
	"fmt"
	"html"
	"net"
	"net/http"
	"strings"

	"github.com/gfleury/gobbs/common"
	"github.com/gfleury/gobbs/common/log"

	bitbucketv1 "github.com/gfleury/go-bitbucket-v1"
	"github.com/spf13/cobra"
)

var (
	noColor *bool
	limit   *int
)

func init() {
	noColor = Search.Flags().BoolP("nocolor", "N", false, "Disable color output (getting a working diff)")
	limit = Search.Flags().IntP("limit", "L", 5, "Limit of result items")
}

// Search is the cmd implementation for Searching Users
var Search = &cobra.Command{
	Use:     "search",
	Aliases: []string{"s"},
	Short:   "Search stash for any string",
	Args:    cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		var results []bitbucketv1.SearchResult

		queryString := args[0]

		if len(args) > 1 {
			queryString = strings.Join(args, " ")
		}

		limits := bitbucketv1.Limits{}

		if *limit >= 25 {
			limits.Primary = 25
			limits.Secondary = 25
		} else {
			limits.Primary = *limit
			limits.Secondary = *limit
		}

		searchQuery := bitbucketv1.SearchQuery{
			Query:  queryString,
			Limits: limits,
		}

		for {
			var hasNext bool
			apiClient, cancel, err := common.APIClient(cmd)
			defer cancel()

			if err != nil {
				cmd.SilenceUsage = true
				return err
			}

			response, err := apiClient.DefaultApi.SearchCode(searchQuery)

			if netError, ok := err.(net.Error); (!ok || (ok && !netError.Timeout())) &&
				!errors.Is(err, context.Canceled) &&
				!errors.Is(err, context.DeadlineExceeded) &&
				response != nil && response.Response != nil &&
				response.Response.StatusCode >= http.StatusMultipleChoices {
				common.PrintApiError(response.Values)
				cmd.SilenceUsage = true
				log.Debugf(err.Error())
				return fmt.Errorf("Unable to process request, API Error")
			} else if err != nil {
				// Print what we already collect
				if len(results) > 0 {
					break
				}
				cmd.SilenceUsage = true
				return err
			}

			pagedSearchResult, err := bitbucketv1.GetSearchResultResponse(response)
			if err != nil {
				cmd.SilenceUsage = true
				return err
			}
			results = append(results, pagedSearchResult)

			hasNext = !pagedSearchResult.Code.IsLastPage
			if _, ok := response.Values["code"]; !hasNext || !ok {
				break
			}

			*limit -= pagedSearchResult.Code.Count

			if *limit >= 25 {
				searchQuery.Entities.Code.Limit = 25
			} else if *limit <= 0 {
				break
			} else {
				searchQuery.Entities.Code.Limit = *limit
			}

			searchQuery.Entities.Code.Start = pagedSearchResult.Code.NextStart
		}

		less, stdin := common.Pager()

		go func() {
			for idx := range results {
				for _, value := range results[idx].Code.Values {
					fmt.Fprintf(stdin, "// Repository\n")
					fmt.Fprintf(stdin, "// %s/%s: %s\n", value.Repository.Project.Key, value.Repository.Slug, value.File)
					for i := range value.HitContexts {
						for j := range value.HitContexts[i] {
							fmt.Fprintf(stdin, "\t%d  %s\n", value.HitContexts[i][j].Line, unHtml(value.HitContexts[i][j].Text, !*noColor))
						}
					}
					fmt.Fprintf(stdin, "\n")
				}
			}
			defer stdin.Close()
		}()

		err := less.Run()
		return err
	},
}

func unHtml(line string, color bool) (ret string) {
	green := "\033[0;32m"
	reset := "\033[0m"

	if !color {
		green = ""
		reset = ""
	}

	ret = html.UnescapeString(line)
	ret = strings.ReplaceAll(ret, "<em>", green)
	ret = strings.ReplaceAll(ret, "</em>", reset)
	return ret
}
