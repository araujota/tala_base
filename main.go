package main

import (
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"tala_base/orchestrator"
	"tala_base/types"
	"tala_base/utils"
)

type Server struct {
	executor *orchestrator.ChainExecutor
}

func NewServer() *Server {
	executor := orchestrator.NewChainExecutor()

	// Load all workflows from the workflows directory
	workflowFiles, err := filepath.Glob("workflows/*.yaml")
	if err != nil {
		log.Printf("Warning: Failed to read workflows directory: %v", err)
	}

	for _, file := range workflowFiles {
		name := strings.TrimSuffix(filepath.Base(file), ".yaml")
		if err := executor.LoadWorkflow(name); err != nil {
			log.Printf("Warning: Failed to load workflow %s: %v", name, err)
		} else {
			log.Printf("Loaded workflow: %s", name)
		}
	}

	return &Server{executor: executor}
}

// handleLambda handles direct lambda invocations
func (s *Server) handleLambda(w http.ResponseWriter, r *http.Request) {
	utils.SetCORSHeaders(w)
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method != http.MethodPost {
		utils.RespondError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	// Extract lambda name from path
	path := strings.Trim(r.URL.Path, "/")
	parts := strings.Split(path, "/")
	if len(parts) != 2 {
		utils.RespondError(w, http.StatusBadRequest, "Invalid lambda path")
		return
	}
	lambdaName := parts[1]

	// Parse input
	var input map[string]interface{}
	if err := utils.DecodeJSONBody(w, r, &input); err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Create workflow input
	workflowInput := types.WorkflowInput{
		Data: input,
	}

	// Execute single step
	result, err := s.executor.ExecuteStep(types.Step{
		Name:   lambdaName,
		Lambda: lambdaName,
	}, &types.WorkflowState{
		Steps: map[string]types.StepState{
			lambdaName: {
				Input: workflowInput,
			},
		},
		CurrentStep: lambdaName,
	})

	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if result.Error != nil {
		utils.RespondError(w, http.StatusInternalServerError, result.Error.Message)
		return
	}

	utils.RespondJSON(w, http.StatusOK, result)
}

// handleWorkflow handles workflow executions
func (s *Server) handleWorkflow(w http.ResponseWriter, r *http.Request) {
	utils.SetCORSHeaders(w)
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method != http.MethodPost {
		utils.RespondError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	// Extract workflow name from path
	path := strings.Trim(r.URL.Path, "/")
	parts := strings.Split(path, "/")
	if len(parts) != 2 {
		utils.RespondError(w, http.StatusBadRequest, "Invalid workflow path")
		return
	}
	workflowName := parts[1]

	// Parse input
	var input map[string]interface{}
	if err := utils.DecodeJSONBody(w, r, &input); err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Create workflow input
	workflowInput := types.WorkflowInput{
		Data: input,
	}

	// Execute workflow
	result, err := s.executor.ExecuteChain(workflowName, workflowInput)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.RespondJSON(w, http.StatusOK, result)
}

func main() {
	server := NewServer()

	// Handle direct lambda invocations
	http.HandleFunc("/lambda/", server.handleLambda)

	// Handle workflow executions
	http.HandleFunc("/workflow/", server.handleWorkflow)

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Starting server on port %s", port)
	log.Printf("Available endpoints:")
	log.Printf("  Direct lambda:   POST /lambda/<lambda_name>")
	log.Printf("  Workflow:        POST /workflow/<workflow_name>")
	log.Printf("\nExample usage:")
	log.Printf("  # Call lambda directly")
	log.Printf("  curl -X POST http://localhost:%s/lambda/user_create -H \"Content-Type: application/json\" -d '{\"data\":{\"email\":\"test@example.com\",\"name\":\"Test User\"}}'", port)
	log.Printf("\n  # Execute workflow")
	log.Printf("  curl -X POST http://localhost:%s/workflow/user_signup_chain -H \"Content-Type: application/json\" -d '{\"data\":{\"email\":\"test@example.com\",\"name\":\"Test User\"}}'", port)

	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}
