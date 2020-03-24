package repos

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"os"

	"github.com/gfleury/gobbs/common"
	"github.com/gfleury/gobbs/common/log"

	bitbucketv1 "github.com/gfleury/go-bitbucket-v1"
	"github.com/spf13/cobra"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing/transport/ssh"
)

// Clone is the cmd implementation for Cloning repositories locally
var Clone = &cobra.Command{
	Use:     "clone repoName",
	Aliases: []string{"cl"},
	Short:   "Clone repository",
	Args:    cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		apiClient, cancel, err := common.APIClient(cmd)
		defer cancel()

		if err != nil {
			cmd.SilenceUsage = true
			return err
		}

		stashInfo := cmd.Context().Value(common.StashInfoKey).(*common.StashInfo)
		err = mustHaveProject(stashInfo)
		if err != nil {
			return err
		}

		response, err := apiClient.DefaultApi.GetRepository(*stashInfo.Project(), args[0])

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

		repo, err := bitbucketv1.GetRepositoryResponse(response)
		if err != nil {
			cmd.SilenceUsage = true
			return err
		}

		if len(repo.Links.Clone) < 1 {
			cmd.SilenceUsage = true
			return fmt.Errorf("Can't find a clonable link for the repository.")
		}

		homeDir, err := os.UserHomeDir()
		if err != nil {
			return err
		}

		authSSH, err := ssh.NewPublicKeysFromFile("git", fmt.Sprintf("%s/.ssh/id_rsa", homeDir), "")
		if err != nil {
			return err
		}

		for _, link := range repo.Links.Clone {
			if link.Name == "ssh" {
				_, err = git.PlainClone(args[0], false, &git.CloneOptions{
					URL:               link.Href,
					RecurseSubmodules: git.DefaultSubmoduleRecursionDepth,
					Auth:              authSSH,
				})
				break
			}
		}

		if err != nil {
			return err
		}

		log.Noticef("Repository successfully cloned into %s", args[0])
		return err
	},
}
