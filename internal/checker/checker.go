// Package checker contient la logique de vérification des cibles
// (lecture de targets.conf, vérification TCP/HTTP, mesure de latence).
package checker

import (
	"bufio"
	"net"
	"net/http"
	"os"
	"strings"
	"time"
)

// Result représente le résultat d'une vérification pour une cible.
type Result struct {
	Target  string
	Up      bool
	Slow    bool
	Latency time.Duration
	Err     error
}

// LoadTargets lit le fichier de configuration et retourne la liste des cibles.
// Les lignes vides et celles commençant par # sont ignorées.
func LoadTargets(path string) ([]string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var targets []string
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		targets = append(targets, line)
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return targets, nil
}

// Check vérifie une cible et retourne le résultat.
// Si la cible commence par http:// ou https://, une requête HEAD est envoyée.
// Sinon, la cible est traitée comme une IP/hôte et vérifiée en TCP sur le port 80.
// slowThreshold : si la latence dépasse ce seuil, le résultat est marqué Slow (0 = désactivé).
func Check(target string, timeout, slowThreshold time.Duration) Result {
	var r Result
	if strings.HasPrefix(target, "http://") || strings.HasPrefix(target, "https://") {
		r = checkHTTP(target, timeout)
	} else {
		r = checkTCP(target, timeout)
	}
	if r.Up && slowThreshold > 0 && r.Latency > slowThreshold {
		r.Slow = true
	}
	return r
}

func checkHTTP(target string, timeout time.Duration) Result {
	client := http.Client{Timeout: timeout}

	req, err := http.NewRequest(http.MethodHead, target, nil)
	if err != nil {
		return Result{Target: target, Up: false, Err: err}
	}

	start := time.Now()
	resp, err := client.Do(req)
	latency := time.Since(start)

	if err != nil {
		return Result{Target: target, Up: false, Latency: latency, Err: err}
	}
	defer resp.Body.Close()

	return Result{Target: target, Up: true, Latency: latency}
}

func checkTCP(target string, timeout time.Duration) Result {
	address := target
	if !strings.Contains(address, ":") {
		address = address + ":80"
	}

	start := time.Now()
	conn, err := net.DialTimeout("tcp", address, timeout)
	latency := time.Since(start)

	if err != nil {
		return Result{Target: target, Up: false, Latency: latency, Err: err}
	}
	defer conn.Close()

	return Result{Target: target, Up: true, Latency: latency}
}