package types

// Step represents a single step in a workflow
type Step struct {
	Name          string `yaml:"name"`
	Lambda        string `yaml:"lambda"`
	InputTemplate string `yaml:"input_template"`
	PassOutputAs  string `yaml:"pass_output_as"`
	ErrorHandler  string `yaml:"error_handler,omitempty"`
}

// Workflow represents a complete workflow definition
type Workflow struct {
	Name        string `yaml:"name"`
	Description string `yaml:"description"`
	Steps       []Step `yaml:"steps"`
}

// WorkflowState represents the state of a workflow execution
type WorkflowState struct {
	Steps       map[string]StepState `json:"steps"`
	CurrentStep string               `json:"current_step"`
	Completed   bool                 `json:"completed"`
}

// StepState represents the state of a single step execution
type StepState struct {
	Input  WorkflowInput  `json:"input"`
	Output WorkflowOutput `json:"output"`
}

// WorkflowInput represents the input to a workflow
type WorkflowInput struct {
	Data    map[string]interface{} `json:"data"`
	Context map[string]interface{} `json:"context"`
}

// WorkflowOutput represents the output of a workflow
type WorkflowOutput struct {
	Data    map[string]interface{} `json:"data"`
	Context map[string]interface{} `json:"context"`
	Error   *WorkflowError         `json:"error,omitempty"`
}

// WorkflowError represents an error in workflow execution
type WorkflowError struct {
	Step    string `json:"step"`
	Message string `json:"message"`
	Code    string `json:"code"`
}

// StepResult represents the result of a single step execution
type StepResult struct {
	Data  map[string]interface{} `json:"data"`
	Error *WorkflowError         `json:"error,omitempty"`
}
