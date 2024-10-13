package logparser

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"github.com/eaardal/dig/config"
	"github.com/eaardal/dig/logentry"
	jsoniter "github.com/json-iterator/go"
	"io"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

type LogFile interface {
	Name() string
	Content() []byte
}

func ParseLogFile(ctx context.Context, sourceCh chan LogFile, sinkCh chan<- *logentry.LogEntry) error {
	defer close(sinkCh)

	for logFile := range sourceCh {
		reader := bufio.NewReader(bytes.NewReader(logFile.Content()))

		if err := readAndParseContent(ctx, reader, config.AppConfig.Keywords, sinkCh); err != nil {
			return err
		}
	}

	return nil
}

func readAndParseContent(ctx context.Context, reader *bufio.Reader, keywords config.KeywordConfig, logEntryCh chan<- *logentry.LogEntry) error {
	lineCount := 0

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			line, err := reader.ReadBytes('\n')
			if err != nil && err == io.EOF {
				return nil
			}
			if err != nil {
				return fmt.Errorf("failed to read log line: %w", err)
			}

			lineCount++
			logEntry := parseLogLine(line, lineCount, keywords)
			sendToSink(ctx, logEntryCh, logEntry)
		}
	}
}

//func readLine(reader *bufio.Reader) ([]byte, error) {
//	var fullLine []byte
//
//	for {
//		// Read a line of input, it may return only a part of the line (isPrefix = true)
//		line, isPrefix, err := reader.ReadLine()
//		if err == io.EOF {
//			// If we reach the end of file, return gracefully
//			return nil, io.EOF
//		}
//		if err != nil {
//			return nil, fmt.Errorf("failed to read line: %w", err)
//		}
//
//		// Append the current part to the full line
//		fullLine = append(fullLine, line...)
//
//		// If isPrefix is false, it means we have the full line. Otherwise, we need to read more in order to complete the line.
//		if !isPrefix {
//			break
//		}
//	}
//
//	return fullLine, nil
//}
//
//func readAndParseContent(ctx context.Context, reader *bufio.Reader, keywords config.KeywordConfig, logEntryCh chan<- *logentry.LogEntry) {
//	line, readErr := reader.ReadBytes('\n')
//	if readErr != nil {
//		log.Fatalf("failed to read log line: %v", readErr)
//		return
//	}
//
//	lineCount := 0
//
//	for readErr == nil {
//		lineCount++
//
//		logEntry := parseLogLine(line, lineCount, keywords)
//		sendToSink(ctx, logEntryCh, logEntry)
//
//		line, readErr = reader.ReadBytes('\n')
//	}
//
//	close(logEntryCh)
//}

func parseLogLine(line []byte, lineCount int, keywords config.KeywordConfig) *logentry.LogEntry {
	logEntry := &logentry.LogEntry{
		LineNumber:      lineCount,
		OriginalLogLine: line,
		Fields:          make(map[string]string),
	}

	parsedLogLine := make(map[string]interface{}, 0)

	if err := json.Unmarshal(line, &parsedLogLine); err != nil {
		logEntry.SetOriginalLogLine(line)
	} else {
		logEntry.SetFromJsonMap(parsedLogLine, keywords)
	}

	return logEntry
}

func sendToSink(ctx context.Context, logEntryCh chan<- *logentry.LogEntry, logEntry *logentry.LogEntry) {
	select {
	case <-ctx.Done():
		return
	case logEntryCh <- logEntry:
	}
}
