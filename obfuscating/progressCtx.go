package obfuscating

import (
	"github.com/google/uuid"
)

type ObfuscationProgress struct {
	ProcessId     string
	FinishedCount int
	TotalCount    int
	Error         string
}

var progressCtx = make(map[string]ObfuscationProgress)

func InitProcess(totalCount int) (string, error) {
	processUuid, err := uuid.NewRandom()
	if err != nil {
		return "", err
	}
	processId := processUuid.String()
	var progressEntry ObfuscationProgress
	progressEntry.ProcessId = processId
	progressEntry.FinishedCount = 0
	progressEntry.TotalCount = totalCount
	progressCtx[processId] = progressEntry
	return processId, err
}

func GetProcessCtx(processId string) (ObfuscationProgress, bool) {
	progressEntry, exists := progressCtx[processId]
	return progressEntry, exists
}

func EmptyProgressCtx() {
	progressCtx = make(map[string]ObfuscationProgress)
}

func increaseFinished(processId string) {
	entry := progressCtx[processId]
	entry.FinishedCount = entry.FinishedCount + 1
	progressCtx[processId] = entry
}

func writeError(processId string, err error) {
	entry := progressCtx[processId]
	entry.Error = err.Error()
	progressCtx[processId] = entry
}
