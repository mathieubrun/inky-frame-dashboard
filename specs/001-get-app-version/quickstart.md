# Quickstart: Get App Version

## Local CLI Verification
1. Build the application: `go build -o inky ./cmd/inky`
2. Run the version command: `./inky version`
3. Verify output matches the expected SemVer string.

## API Verification
1. Start the server (default port 8080): `./inky serve`
2. Query the version endpoint: `curl http://localhost:8080/version`
3. Verify the JSON response contains the correct version field.
