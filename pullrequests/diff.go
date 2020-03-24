package pullrequests

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/exec"
	"strconv"

	"github.com/gfleury/gobbs/common"
	"github.com/gfleury/gobbs/common/log"

	bitbucketv1 "github.com/gfleury/go-bitbucket-v1"
	"github.com/spf13/cobra"
)

const (
	WhiteColor = "\033[1;97m"
	GreenColor = "\033[0;32m"
	RedColor   = "\033[0;31m"
	CyanColor  = "\033[0;36m"
	ResetColor = "\033[0m"
)

var diffContextLines *int32

func init() {
	diffContextLines = Diff.Flags().Int32P("lines", "L", 3, "Diff context lines (how many lines before/after)")
}

// Digg is the cmd implementation for getting Raw diff of PullRequests
var Diff = &cobra.Command{
	Use:     "diff pullRequestID",
	Aliases: []string{"inf"},
	Short:   "Diff pull requests for repository",
	Args:    cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		prID, err := strconv.Atoi(args[0])
		if err != nil {
			log.Critical("Argument must be a pull request ID. Err: %s", err.Error())
			return err
		}

		apiClient, cancel, err := common.APIClient(cmd)
		defer cancel()

		if err != nil {
			log.Critical("Argument must be a pull request ID. Err: %s", err.Error())
			return err
		}

		stashInfo := cmd.Context().Value(common.StashInfoKey).(*common.StashInfo)
		err = mustHaveProjectRepo(stashInfo)
		if err != nil {
			return err
		}

		opts := map[string]interface{}{
			"contextLines": int32(3),
		}

		response, err := apiClient.DefaultApi.GetPullRequestDiff(*stashInfo.Project(), *stashInfo.Repo(), prID, opts)

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
			cmd.SilenceUsage = true
			return err
		}

		diffs, err := bitbucketv1.GetDiffResponse(response)
		if err != nil {
			return err
		}

		less := exec.Command("less", "-R")

		stdin, err := less.StdinPipe()
		if err != nil {
			log.Critical(err.Error())
			stdin = os.Stdin
		}

		less.Stdout = os.Stdout

		go func() {
			for _, diff := range diffs.Diffs {
				// --- a/pullrequests/pr.go
				// +++ b/pullrequests/pr.go
				fmt.Fprintf(stdin, "%s--- %s%s\n", WhiteColor, diff.Source.ToString, ResetColor)
				fmt.Fprintf(stdin, "%s+++ %s%s\n", WhiteColor, diff.Destination.ToString, ResetColor)
				for _, hunk := range diff.Hunks {
					// @@ -15,6 +15,7 @@
					fmt.Fprintf(stdin, "%s@@ -%d,%d +%d,%d @@%s\n", CyanColor, hunk.SourceLine, hunk.SourceSpan, hunk.DestinationLine, hunk.DestinationSpan, CyanColor)
					for _, segment := range hunk.Segments {
						lineStart := " "
						if segment.Type == "ADDED" {
							lineStart = GreenColor + "+"
						} else if segment.Type == "REMOVED" {
							lineStart = RedColor + "-"
						}
						for _, line := range segment.Lines {
							fmt.Fprintf(stdin, "%s%s%s\n", lineStart, line.Line, ResetColor)
						}
					}
				}
			}
			defer stdin.Close()
		}()

		err = less.Run()
		return err
	},
}
