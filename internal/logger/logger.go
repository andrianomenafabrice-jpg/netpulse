// Package logger gère l'écriture des résultats de vérification
// (fichier de log horodaté + stdout, option -json).
package logger

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"netpulse/internal/checker"
)

// Logger écrit les résultats dans un fichier de log et sur stdout.
type Logger struct {
	file *os.File
}

// New ouvre (ou crée) le fichier de log en mode ajout.
func New(path string) (*Logger, error) {
	f, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}
	return &Logger{file: f}, nil
}

// Close ferme le fichier de log.
func (l *Logger) Close() error {
	return l.file.Close()
}

// jsonEntry est la structure utilisée pour la sortie JSON sur stdout.
type jsonEntry struct {
	Timestamp string `json:"timestamp"`
	Target    string `json:"target"`
	Status    string `json:"status"`
	LatencyMs int64  `json:"latency_ms"`
}

// Write écrit le résultat dans le fichier de log (toujours au format texte)
// et sur stdout (texte ou JSON selon jsonOutput).
func (l *Logger) Write(r checker.Result, jsonOutput bool) {
	ts := time.Now().UTC().Format(time.RFC3339)
	status := "OK"
	if !r.Up {
		status = "DOWN"
	}

	var line string
	if r.Up {
		line = fmt.Sprintf("%s [%s] %s %dms", ts, status, r.Target, r.Latency.Milliseconds())
	} else {
		line = fmt.Sprintf("%s [%s] %s timeout", ts, status, r.Target)
	}

	// Le fichier de log reste toujours au format texte, quel que soit -json.
	fmt.Fprintln(l.file, line)

	if jsonOutput {
		entry := jsonEntry{
			Timestamp: ts,
			Target:    r.Target,
			Status:    status,
			LatencyMs: r.Latency.Milliseconds(),
		}
		data, err := json.Marshal(entry)
		if err != nil {
			fmt.Fprintf(os.Stderr, "erreur JSON: %v\n", err)
			return
		}
		fmt.Println(string(data))
	} else {
		fmt.Println(line)
	}
}