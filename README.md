# TALA (Type-Assisted Lambda Architecture)

A framework for building type-safe, AI-assistable serverless applications using Go and Cloudflare Workers(TBD).

## Overview

TALA solves the context problem for AI code assistants by providing:
- Type-driven workflows
- Stateless lambda functions
- Clear data flow patterns
- Edge-first architecture

## Architecture

```
├── lambdas/           # Individual lambda functions
│   ├── user_create/   # Example lambda
│   └── ...
├── orchestrator/      # Workflow orchestration
│   ├── executor.go    # Workflow execution engine
├── workflows/         # YAML workflow definitions
├── utils/            # Shared utilities
└── main.go           # Main application server
```

## Quick Start

1. **Prerequisites**
   ```bash
   # Install Go 1.21+
   brew install go

   # Install TinyGo for WASM compilation
   brew install tinygo
   ```

2. **Setup Database**
   ```bash
   # Start PostgreSQL
   docker run -d --name tala-db \
     -e POSTGRES_USER=user \
     -e POSTGRES_PASSWORD=password \
     -e POSTGRES_DB=tala \
     -p 5432:5432 \
     postgres:latest

   # Create tables
   psql -h localhost -U user -d tala -f scripts/schema.sql
   ```

3. **Configure Environment**
   ```bash
   # Copy example env file
   cp env.example .env

   # Edit .env with your settings
   DATABASE_URL=postgres://user:password@localhost:5432/tala?sslmode=disable
   ```

4. **Run Locally**
   ```bash
   # Start the main server (workflows + direct operations)
   go run main.go

   # Or start just the workflow orchestrator
   go run cmd/server/main.go
   ```

5. **Test the API**
   ```bash
   # Create a user directly
   curl -X POST http://localhost:8080/users \
     -H "Content-Type: application/json" \
     -d '{"email":"test@example.com","name":"Test User"}'

   # Run a workflow
   curl -X POST http://localhost:8080/run/user_signup_chain \
     -H "Content-Type: application/json" \
     -d '{"email":"test@example.com","name":"Test User"}'
   ```

## System Prompt for LLMs

When working with this codebase, use system prompts like this example to help LLMs understand the architecture:

```
You are assisting with a TALA (Type-Assisted Lambda Architecture) application. This system uses:

1. Type-Driven Workflows:
   - Workflows are defined in YAML with clear input/output types
   - Each step in a workflow is a stateless lambda function
   - Data flows between steps using typed templates
   - Error handling is explicit in the workflow definition

2. Lambda Functions:
   - Each lambda is a standalone Go function
   - Lambdas are compiled to WebAssembly using TinyGo
   - They communicate via HTTP with typed inputs/outputs
   - They are stateless and can be deployed to edge locations

3. Workflow Orchestration:
   - The orchestrator manages workflow execution
   - It maintains typed state between steps
   - It handles error propagation and recovery
   - It supports template-based data transformation

4. Type System:
   - WorkflowInput/WorkflowOutput for data flow
   - StepState for execution tracking
   - WorkflowError for error handling
   - All types are JSON-serializable

5. Edge-First Design:
   - Stateless design enables edge deployment
   - HTTP-based communication between components
   - Low-latency execution

When helping with this codebase:
1. Respect the type system and workflow patterns
2. Ensure lambdas remain stateless
3. Use the provided error handling patterns
4. Maintain clear data flow between components
5. Consider edge deployment implications
```

## Development

1. **Adding a New Lambda**
   ```bash
   # Create lambda directory
   mkdir -p lambdas/my_lambda

   # Create main.go
   touch lambdas/my_lambda/main.go

   # Build lambda
   ./scripts/build.sh
   ```

2. **Creating a Workflow**
   ```yaml
   # workflows/my_workflow.yaml
   name: my_workflow
   description: My workflow description
   steps:
     - name: step1
       lambda: my_lambda
       input_template: |
         {
           "data": {
             "input": "{{.input.data.input}}"
           }
         }
       pass_output_as: step1_output
   ```

3. **Testing**
   ```bash
   # Test workflow
   curl -X POST http://localhost:8080/run/my_workflow \
     -H "Content-Type: application/json" \
     -d '{"input":"test"}'
   ```

## Deployment

 **Build Lambdas**
   ```bash
   ./scripts/build.sh
   ```

## Contributing

1. Fork the repository
2. Create a feature branch
3. Commit your changes
4. Push to the branch
5. Create a Pull Request

## License

MIT 