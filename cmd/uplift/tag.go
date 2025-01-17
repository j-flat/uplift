/*
Copyright (c) 2021 Gemba Advantage

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
*/

package main

import (
	"github.com/gembaadvantage/uplift/internal/context"
	"github.com/gembaadvantage/uplift/internal/middleware/logging"
	"github.com/gembaadvantage/uplift/internal/middleware/skip"
	"github.com/gembaadvantage/uplift/internal/semver"
	"github.com/gembaadvantage/uplift/internal/task"
	"github.com/gembaadvantage/uplift/internal/task/currentversion"
	"github.com/gembaadvantage/uplift/internal/task/fetchtag"
	"github.com/gembaadvantage/uplift/internal/task/gittag"
	"github.com/gembaadvantage/uplift/internal/task/lastcommit"
	"github.com/gembaadvantage/uplift/internal/task/nextcommit"
	"github.com/gembaadvantage/uplift/internal/task/nextversion"
	"github.com/spf13/cobra"
)

const (
	tagDesc = `Tags a git repository with the next semantic version. The tag
is based on the conventional commit message from the last commit.`
)

func newTagCmd(ctx *context.Context) *cobra.Command {
	var pre string

	cmd := &cobra.Command{
		Use:   "tag",
		Short: "Tag a git repository with the next semantic version",
		Long:  tagDesc,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			// Handle prerelease suffix if one is provided
			if pre != "" {
				var err error
				if ctx.Prerelease, ctx.Metadata, err = semver.ParsePrerelease(pre); err != nil {
					return err
				}
			}

			return tagRepo(ctx)
		},
	}

	f := cmd.Flags()
	f.BoolVar(&ctx.FetchTags, "fetch-all", false, "fetch all tags from the remote repository")
	f.BoolVar(&ctx.NextTagOnly, "next", false, "output the next tag only")
	f.StringVar(&pre, "prerelease", "", "append a prerelease suffix to next calculated semantic version")

	return cmd
}

func tagRepo(ctx *context.Context) error {
	tsks := []task.Runner{
		fetchtag.Task{},
		lastcommit.Task{},
		currentversion.Task{},
		nextversion.Task{},
		nextcommit.Task{},
		gittag.Task{},
	}

	for _, tsk := range tsks {
		if err := skip.Running(tsk.Skip, logging.Log(tsk.String(), tsk.Run))(ctx); err != nil {
			return err
		}
	}

	return nil
}
