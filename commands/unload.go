package commands

import (
	"fmt"

	"github.com/docker/model-cli/commands/completion"
	"github.com/docker/model-cli/desktop"
	"github.com/spf13/cobra"
)

func newUnloadCmd() *cobra.Command {
	var all bool
	var backend string

	const cmdArgs = "(MODEL [MODEL ...] [--backend BACKEND] | --all)"
	c := &cobra.Command{
		Use:   "unload " + cmdArgs,
		Short: "Unload running models",
		RunE: func(cmd *cobra.Command, models []string) error {
			unloadResp, err := desktopClient.Unload(desktop.UnloadRequest{All: all, Backend: backend, Models: models})
			if err != nil {
				err = handleClientError(err, "Failed to unload models")
				return handleNotRunningError(err)
			}
			unloaded := unloadResp.UnloadedRunners
			if unloaded == 0 {
				if all {
					cmd.Println("No models are running.")
				} else {
					cmd.Println("No such model(s) running.")
				}
			} else {
				cmd.Printf("Unloaded %d model(s).\n", unloaded)
			}
			return nil
		},
		ValidArgsFunction: completion.NoComplete,
	}
	c.Args = func(cmd *cobra.Command, args []string) error {
		if all {
			if len(args) > 0 {
				return fmt.Errorf(
					"'docker model unload' does not take MODEL when --all is specified.\n\n" +
						"Usage:  docker model unload " + cmdArgs + "\n\n" +
						"See 'docker model unload --help' for more information.",
				)
			}
			return nil
		}
		if len(args) < 1 {
			return fmt.Errorf(
				"'docker model unload' requires MODEL unless --all is specified.\n\n" +
					"Usage:  docker model unload " + cmdArgs + "\n\n" +
					"See 'docker model unload --help' for more information.",
			)
		}
		return nil
	}
	c.Flags().BoolVar(&all, "all", false, "Unload all running models")
	c.Flags().StringVar(&backend, "backend", "", "Optional backend to target")
	return c
}
