package receivepack

import (
	"context"
	"os"

	"gitlab.com/gitlab-org/gitlab-shell/internal/command/commandargs"
	"gitlab.com/gitlab-org/gitlab-shell/internal/command/readwriter"
	"gitlab.com/gitlab-org/gitlab-shell/internal/command/shared/accessverifier"
	"gitlab.com/gitlab-org/gitlab-shell/internal/command/shared/customaction"
	"gitlab.com/gitlab-org/gitlab-shell/internal/command/shared/disallowedcommand"
	"gitlab.com/gitlab-org/gitlab-shell/internal/config"
)

type Command struct {
	Config     *config.Config
	Args       *commandargs.Shell
	ReadWriter *readwriter.ReadWriter
}

func (c *Command) Execute(ctx context.Context) error {
	args := c.Args.SshArgs
	if len(args) != 2 {
		return disallowedcommand.Error
	}

	repo := args[1]
	response, err := c.verifyAccess(ctx, repo)
	if err != nil {
		return err
	}

	if response.IsCustomAction() {
		customAction := customaction.Command{
			Config:     c.Config,
			ReadWriter: c.ReadWriter,
			EOFSent:    true,
		}
		return customAction.Execute(ctx, response)
	}

	var gitProtocolVersion string
	if c.Args.RemoteAddr != nil {
		gitProtocolVersion = c.Args.GitProtocolVersion
	} else {
		gitProtocolVersion = os.Getenv(commandargs.GitProtocolEnv)
	}

	return c.performGitalyCall(response, gitProtocolVersion)
}

func (c *Command) verifyAccess(ctx context.Context, repo string) (*accessverifier.Response, error) {
	cmd := accessverifier.Command{c.Config, c.Args, c.ReadWriter}

	return cmd.Verify(ctx, c.Args.CommandType, repo)
}
