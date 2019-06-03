/*
Copyright 2019 The Skaffold Authors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package runner

import (
	"context"
	"io"

	"github.com/GoogleContainerTools/skaffold/pkg/skaffold/build"
)

// Deploy deploys build artifacts.
func (r *SkaffoldRunner) Deploy(ctx context.Context, out io.Writer, artifacts []build.Artifact) error {
	if err := r.deploy(ctx, out, artifacts); err != nil {
		return err
	}
	if r.runCtx.Opts.Tail {
		images := make([]string, len(artifacts))
		for i, a := range artifacts {
			images[i] = a.ImageName
		}
		logger := r.newLoggerForImages(out, images)
		return r.TailLogs(ctx, out, logger)
	}
	return nil
}

// Deploy deploys the given artifacts and tail logs if tail present
func (r *SkaffoldRunner) deploy(ctx context.Context, out io.Writer, artifacts []build.Artifact) error {
	err := r.Deployer.Deploy(ctx, out, artifacts, r.labellers)
	r.hasDeployed = true
	if err != nil {
		return err
	}
	return r.statusCheck(ctx, out)
}

func (r *SkaffoldRunner) statusCheck(ctx context.Context, out io.Writer) error {
	// Check if we need to perform deploy status
	if !r.runCtx.Opts.StatusCheck {
		return nil
	}
	// Perform
	return deploy.GetDeployedManifests()
}
