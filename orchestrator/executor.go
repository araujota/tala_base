package orchestrator

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"text/template"

	"tala_base/types"

	"gopkg.in/yaml.v3"
)

type ChainExecutor struct {
	workflows map[string]types.Workflow
}

func NewChainExecutor() *ChainExecutor {
	return &ChainExecutor{
		workflows: make(map[string]types.Workflow),
	}
}

func (e *ChainExecutor) LoadWorkflow(name string) error {
	file, err := os.ReadFile(fmt.Sprintf("workflows/%s.yaml", name))
	if err != nil {
		return fmt.Errorf("failed to read workflow file: %w", err)
	}

	var workflow types.Workflow
	if err := yaml.Unmarshal(file, &workflow); err != nil {
		return fmt.Errorf("failed to parse workflow: %w", err)
	}

	e.workflows[name] = workflow
	return nil
}

func (e *ChainExecutor) ExecuteStep(step types.Step, state *types.WorkflowState) (*types.StepResult, error) {
	// Parse input template
	tmpl, err := template.New("input").Parse(step.InputTemplate)
	if err != nil {
		return nil, fmt.Errorf("failed to parse input template: %w", err)
	}

	// Execute template with current state
	var inputBuf bytes.Buffer
	if err := tmpl.Execute(&inputBuf, state); err != nil {
		return nil, fmt.Errorf("failed to execute template: %w", err)
	}

	// Call lambda
	lambdaURL := fmt.Sprintf("http://%s:8080", step.Lambda)
	resp, err := http.Post(lambdaURL, "application/json", &inputBuf)
	if err != nil {
		return nil, fmt.Errorf("failed to call lambda: %w", err)
	}
	defer resp.Body.Close()

	// Read response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read lambda response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return &types.StepResult{
			Error: &types.WorkflowError{
				Step:    step.Name,
				Message: fmt.Sprintf("lambda returned error: %s", string(body)),
				Code:    "LAMBDA_ERROR",
			},
		}, nil
	}

	// Parse response
	var result types.StepResult
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to parse lambda response: %w", err)
	}

	return &result, nil
}

func (e *ChainExecutor) ExecuteChain(name string, input types.WorkflowInput) (*types.WorkflowOutput, error) {
	workflow, exists := e.workflows[name]
	if !exists {
		return nil, fmt.Errorf("workflow %s not found", name)
	}

	state := &types.WorkflowState{
		Steps:       make(map[string]types.StepState),
		CurrentStep: workflow.Steps[0].Name,
	}

	// Initialize first step
	state.Steps[workflow.Steps[0].Name] = types.StepState{
		Input: input,
	}

	for i, step := range workflow.Steps {
		// Execute step
		result, err := e.ExecuteStep(step, state)
		if err != nil {
			return nil, fmt.Errorf("step %s failed: %w", step.Name, err)
		}

		// Update state
		stepState := state.Steps[step.Name]
		stepState.Output = types.WorkflowOutput{
			Data:  result.Data,
			Error: result.Error,
		}
		state.Steps[step.Name] = stepState

		// Handle error if any
		if result.Error != nil {
			if step.ErrorHandler != "" {
				// Execute error handler
				errorStep := workflow.Steps[i+1]
				errorResult, err := e.ExecuteStep(errorStep, state)
				if err != nil {
					return nil, fmt.Errorf("error handler %s failed: %w", errorStep.Name, err)
				}
				state.Steps[errorStep.Name] = types.StepState{
					Input: stepState.Input,
					Output: types.WorkflowOutput{
						Data:  errorResult.Data,
						Error: errorResult.Error,
					},
				}
			}
			return &types.WorkflowOutput{
				Error: result.Error,
			}, nil
		}

		// Move to next step
		if i < len(workflow.Steps)-1 {
			nextStep := workflow.Steps[i+1]
			state.CurrentStep = nextStep.Name
			state.Steps[nextStep.Name] = types.StepState{
				Input: types.WorkflowInput{
					Data:    result.Data,
					Context: stepState.Input.Context,
				},
			}
		}
	}

	// Workflow completed successfully
	state.Completed = true
	lastStep := workflow.Steps[len(workflow.Steps)-1]
	lastState := state.Steps[lastStep.Name]

	return &types.WorkflowOutput{
		Data:    lastState.Output.Data,
		Context: lastState.Input.Context,
	}, nil
}
