package httpServer

type ErrorResponse struct {
	Error string
}

type SuccessfulResponse struct {
	Status string
}

type ObfuscationResponse struct {
	SuccessfulResponse
	ProcessId string
}
