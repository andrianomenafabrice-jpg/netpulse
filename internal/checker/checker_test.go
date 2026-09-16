package checker

import (
	"net"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestLoadTargets(t *testing.T) {
	content := `# commentaire à ignorer
192.168.1.1

https://google.com
   
8.8.8.8
`
	dir := t.TempDir()
	path := filepath.Join(dir, "targets.conf")
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatalf("écriture du fichier de test : %v", err)
	}

	targets, err := LoadTargets(path)
	if err != nil {
		t.Fatalf("LoadTargets a retourné une erreur : %v", err)
	}

	want := []string{"192.168.1.1", "https://google.com", "8.8.8.8"}
	if len(targets) != len(want) {
		t.Fatalf("nombre de cibles = %d, attendu %d (%v)", len(targets), len(want), targets)
	}
	for i, w := range want {
		if targets[i] != w {
			t.Errorf("cible[%d] = %q, attendu %q", i, targets[i], w)
		}
	}
}

func TestLoadTargets_FichierInexistant(t *testing.T) {
	_, err := LoadTargets("chemin/qui/n/existe/pas.conf")
	if err == nil {
		t.Fatal("attendu une erreur pour un fichier inexistant, obtenu nil")
	}
}

func TestCheck_TCPUp(t *testing.T) {
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("impossible de démarrer le listener de test : %v", err)
	}
	defer ln.Close()

	go func() {
		for {
			conn, err := ln.Accept()
			if err != nil {
				return
			}
			conn.Close()
		}
	}()

	target := ln.Addr().String() // ex: 127.0.0.1:54321, inclut déjà le port
	r := Check(target, 2*time.Second, 0)

	if !r.Up {
		t.Errorf("attendu Up=true pour %s, obtenu Up=false (err: %v)", target, r.Err)
	}
	if r.Slow {
		t.Errorf("attendu Slow=false avec seuil désactivé, obtenu Slow=true")
	}
}

func TestCheck_TCPDown(t *testing.T) {
	// Port improbable qu'aucun service n'écoute, avec un timeout court.
	r := Check("127.0.0.1:1", 200*time.Millisecond, 0)

	if r.Up {
		t.Errorf("attendu Up=false pour une cible fermée, obtenu Up=true")
	}
}

func TestCheck_HTTPUp(t *testing.T) {
	srv := httptest.NewServer(nil)
	defer srv.Close()

	r := Check(srv.URL, 2*time.Second, 0)

	if !r.Up {
		t.Errorf("attendu Up=true pour %s, obtenu Up=false (err: %v)", srv.URL, r.Err)
	}
}

// TestApplySlowThreshold teste la logique du seuil SLOW de façon isolée,
// avec des Result fabriqués à la main plutôt qu'une vraie latence réseau
// (le timing réel d'une connexion TCP locale est trop instable pour un test).
func TestApplySlowThreshold(t *testing.T) {
	cases := []struct {
		name      string
		r         Result
		threshold time.Duration
		wantSlow  bool
	}{
		{
			name:      "latence au-dessus du seuil",
			r:         Result{Up: true, Latency: 100 * time.Millisecond},
			threshold: 50 * time.Millisecond,
			wantSlow:  true,
		},
		{
			name:      "latence en dessous du seuil",
			r:         Result{Up: true, Latency: 10 * time.Millisecond},
			threshold: 50 * time.Millisecond,
			wantSlow:  false,
		},
		{
			name:      "seuil désactivé (0)",
			r:         Result{Up: true, Latency: time.Hour},
			threshold: 0,
			wantSlow:  false,
		},
		{
			name:      "cible down, jamais marquée slow même si latence haute",
			r:         Result{Up: false, Latency: time.Hour},
			threshold: 1 * time.Millisecond,
			wantSlow:  false,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := applySlowThreshold(c.r, c.threshold)
			if got.Slow != c.wantSlow {
				t.Errorf("Slow = %v, attendu %v (latence=%v, seuil=%v)", got.Slow, c.wantSlow, c.r.Latency, c.threshold)
			}
		})
	}
}