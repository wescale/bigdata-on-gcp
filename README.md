# sandbox-gcp

## Objectif

Voir la présentation disponible [ici](perdu.com)

## Pré-requis

### A installer

- gcloud
- terraform

### Configuration

Créer le fichier de configuration dans le répertoire racine du projet.

```language-bash
#/bin/bash

export GCP_PROJECT="***"
export CONSUMER_KEY="***"
export CONSUMER_SECRET="***"
export ACCESS_TOKEN="***"
export ACCESS_SECRET="***"


export MASTODON_SERVER="https://linuxrocks.online"
export MASTODON_CLIENT_ID="***"
export MASTODON_CLIENT_SECRET="***"
export MASTODON_LOGIN="***"
export MASTODON_PASSWORD="***"
```

## Création de la plateform

### Création de l'infrastructure

Aller dans le répertoire "iac". Puis lancer le script "create-all.sh".

### Déploiement

Aller dans le répertoire "deployment". Puis lancer le script "deploy.sh".

## Tests

### Flux temps réel

Aller sur le [site: "http://public.gcp-wescale.slavayssiere.fr"](http://public.gcp-wescale.slavayssiere.fr).

### API

#### Génération d'aggregas

```language-bash
curl -X POST http://public.gcp-wescale.slavayssiere.fr/aggregator/stats | jq .
```

```language-bash
curl -X GET http://public.gcp-wescale.slavayssiere.fr/aggregator/stats/1 | jq .
curl -X GET https://us-central1-slavayssiere-sandbox.cloudfunctions.net/laststat | jq .
curl -X GET https://us-central1-slavayssiere-sandbox.cloudfunctions.net/getstat?id=1 | jq .
```

```language-bash
curl -X GET http://public.gcp-wescale.slavayssiere.fr/aggregator/top10 | jq .
```

curl -vvv --header "Authorization: Bearer $TOKEN" https://iap.gcp-wescale.slavayssiere.fr/

#### To test datavisualisation

Aller sur Google Drive

kubectl create ns test
kubectl -n test run -i --tty busybox --image=busybox --restart=Never -- sh

## Observability

Connect to bastion and create local SSH tunnel:

```language-bash
sudo ssh -i /Users/slavayssiere/.ssh/id_rsa admin@bastion.gcp-wescale.slavayssiere.fr -L 80:admin.gcp.wescale:80
```

To view Istio observability stack, please follow: 

- [Service Graph](http://servicegraph.localhost/force/forcegraph.html)
- [Jaeger](http://jaeger.localhost)
- [Prometheus](http://prometheus.localhost)
- [Grafana](http://grafana.localhost)
- [Kiali](http://kiali.localhost)
