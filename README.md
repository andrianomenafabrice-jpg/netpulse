# netpulse

Mini-moniteur de disponibilité réseau pour petits cybercafés, hôtels ou bornes WiFi. Vérifie une liste de cibles (routeur local, DNS, site externe) à intervalle régulier et écrit un journal horodaté, avec un état clair OK/DOWN dans le terminal.

## Pourquoi ce projet

À Madagascar, beaucoup de petits cybercafés, hôtels ou boutiques avec WiFi libre n'ont aucun outil pour savoir si leur connexion ou leur box est tombée. Ils s'en rendent compte seulement quand un client se plaint. netpulse vérifie en continu plusieurs niveaux du réseau (routeur local, DNS, site externe) pour distinguer une panne locale d'une panne côté fournisseur, et garde une trace exploitable dans un log.

## Stack

- **Go** : binaire unique, sans dépendance externe, tourne sur un vieux PC ou un Raspberry Pi.
- **Bash** : scripts de lancement en arrière-plan et de suivi du log.

## Installation

Prérequis : [Go 1.22+](https://go.dev/dl/) installé.

```bash
git clone <URL_DU_REPO_GITHUB>
cd netpulse
go build -o netpulse ./cmd/netpulse
```

Sous Windows (Git Bash), le binaire produit est `netpulse.exe` si vous êtes sous PowerShell/CMD, ou `netpulse` sous Git Bash — utilisez `./netpulse` dans les deux cas depuis Git Bash.

## Configuration

Éditez `targets.conf` : une cible par ligne, IP/hôte simple ou URL complète.

```
# une cible par ligne
192.168.1.1
https://google.com
8.8.8.8
```

- Une IP ou un nom d'hôte simple (ex. `192.168.1.1`) est vérifié par une connexion TCP sur le port 80.
- Une URL `http://` ou `https://` est vérifiée par une requête HEAD.

## Utilisation

### Lancement direct (au premier plan)

```bash
./netpulse -config targets.conf -interval 30 -timeout 5
```

Options disponibles :

| Flag        | Défaut          | Description                                  |
|-------------|-----------------|-----------------------------------------------|
| `-config`   | `targets.conf`  | Chemin vers le fichier de cibles              |
| `-interval` | `30`            | Intervalle entre les vérifications (secondes) |
| `-timeout`  | `5`             | Timeout par cible (secondes)                  |
| `-json`     | `false`         | Sortie stdout au format JSON                  |
| `-log`      | `netpulse.log`  | Chemin du fichier de log                      |

Exemple de sortie terminal :
```
[OK] 192.168.1.1 — 4ms
[DOWN] 8.8.8.8 — timeout
```

Le programme retourne un code de sortie **1** si au moins une cible était DOWN lors du dernier cycle (utile pour scripts), **0** sinon.

### Lancement en arrière-plan

```bash
./scripts/run.sh -config targets.conf
```

Lance netpulse en tâche de fond et enregistre son PID dans `netpulse.pid`.

### Suivre le log en direct

```bash
./scripts/tail.sh
```

### Arrêter le processus en arrière-plan

```bash
kill $(cat netpulse.pid)
rm netpulse.pid netpulse.out
```

## Format du log

Chaque ligne de `netpulse.log` est horodatée en UTC (RFC3339) :

```
2026-09-15T10:32:11Z [OK] 192.168.1.1 4ms
2026-09-15T10:32:11Z [DOWN] 8.8.8.8 timeout
```

Avec `-json`, la sortie **stdout** (pas le fichier log, qui reste toujours en texte) devient :

```json
{"timestamp":"2026-09-15T10:32:11Z","target":"192.168.1.1","status":"OK","latency_ms":4}
```

## Structure du projet

```
netpulse/
├── cmd/netpulse/main.go       — point d'entrée, CLI, boucle principale
├── internal/checker/checker.go — vérification TCP/HTTP, mesure de latence
├── internal/logger/logger.go   — écriture du log (texte + JSON)
├── scripts/run.sh              — lancement en arrière-plan
├── scripts/tail.sh             — suivi du log
├── targets.conf                — liste des cibles (exemple)
├── go.mod
└── .gitignore
```

## Pistes d'amélioration future

- Notification (email ou webhook Discord/Telegram) quand une cible passe DOWN.
- Endpoint HTTP `/status` pour consulter l'état depuis un navigateur local.
- Seuil d'alerte `[SLOW]` configurable pour les latences élevées.