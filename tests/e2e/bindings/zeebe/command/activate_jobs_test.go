// +build e2etests

// ------------------------------------------------------------
// Copyright (c) Microsoft Corporation.
// Licensed under the MIT License.
// ------------------------------------------------------------

package command

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/camunda-cloud/zeebe/clients/go/pkg/entities"
	"github.com/dapr/components-contrib/bindings"
	"github.com/dapr/components-contrib/bindings/zeebe/command"
	"github.com/dapr/components-contrib/tests/e2e/bindings/zeebe"
	"github.com/stretchr/testify/assert"
)

func TestActivateJobs(t *testing.T) {
	t.Parallel()

	id := zeebe.TestID()
	jobType := id + "-test"

	cmd, err := zeebe.Command()
	assert.NoError(t, err)

	deployment, err := zeebe.DeployProcess(
		cmd,
		zeebe.TestProcessFile,
		zeebe.ProcessIDModifier(id),
		zeebe.JobTypeModifier("test", jobType))
	assert.NoError(t, err)
	assert.Equal(t, id, deployment.BpmnProcessId)

	t.Run("activate a job", func(t *testing.T) {
		t.Parallel()

		data, err := json.Marshal(map[string]interface{}{
			"jobType":           jobType,
			"maxJobsToActivate": 100,
			"timeout":           "10m",
		})
		assert.NoError(t, err)

		_, err = zeebe.CreateProcessInstance(cmd, map[string]interface{}{
			"bpmnProcessId": id,
			"version":       1,
		})
		assert.NoError(t, err)
		time.Sleep(5 * time.Second)

		req := &bindings.InvokeRequest{Data: data, Operation: command.ActivateJobsOperation}
		res, err := cmd.Invoke(req)
		assert.NoError(t, err)
		assert.NotNil(t, res)

		jobs := &[]entities.Job{}
		err = json.Unmarshal(res.Data, jobs)
		assert.NoError(t, err)
		assert.Equal(t, 1, len(*jobs))
		assert.Nil(t, res.Metadata)
	})
}
