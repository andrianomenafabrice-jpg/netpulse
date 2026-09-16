# netpulse

Mini-moniteur de disponibilité réseau pour petits cybercafés, hôtels ou bornes WiFi. Vérifie une liste de cibles (routeur, DNS, site externe) et écrit un journal lisible.

## Stack
- Go (binaire unique, sans dépendance externe)
- Bash (scripts de lancement et de suivi)

## Statut du projet
🚧 En construction — étape 1 (structure du projet) terminée.

## Compiler et lancer
```bash
go build -o netpulse ./cmd/netpulse
./netpulse -config targets.conf -interval 30 -timeout 5
```

## Suivre le log
```bash
./scripts/tail.sh
```

## Exemple de targets.conf
```
192.168.1.1
https://google.com
8.8.8.8
```

## Pistes d'amélioration future
- Notification (email ou webhook Discord/Telegram) quand une cible passe DOWN.
- Endpoint HTTP `/status` pour consulter l'état depuis un navigateur local.