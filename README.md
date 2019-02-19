# Bigdata on Google Cloud Plateform


![WeSpeakCloud](img/wespeakcloud.png)

Please follow this link to get the SpeakerDeck presentation

## Objectif

This project is used to:

- get data from Twitter and Mastodon
- convert to a common object
- improve data
- store data
- create statistics
- get data by GSheet and REST API

## Prerequisite

### To be installer

- [Gcloud CLI](https://cloud.google.com/sdk/install)
- [Terraform](https://www.terraform.io/) d'Hashicorp
- [Helm](https://docs.helm.sh/using_helm/)

### Configuration

Clone the current project.
Create config file in the root path of the project.


```language-bash
#/bin/bash

export CONSUMER_KEY="***"
export CONSUMER_SECRET="***"
export ACCESS_TOKEN="***"
export ACCESS_SECRET="***"


export MASTODON_SERVER="https://linuxrocks.online"
export MASTODON_CLIENT_ID="***"
export MASTODON_CLIENT_SECRET="***"
export MASTODON_LOGIN="***"
export MASTODON_PASSWORD="***"

export CLIENT_ID="***"
export CLIENT_SECRET="***"
```

CONSUMER_KEY, CONSUMER_SECRET, ACCESS_TOKEN and ACCESS_SECRET are used to connect to Twitter.

MASTODON_SERVER, MASTODON_CLIENT_ID, MASTODON_CLIENT_SECRET, MASTODON_LOGIN and MASTODON_PASSWORD are used to connect to Mastodon.

## Création de la plateform

### Création de l'infrastructure

Go to "iac". And launch "create-all.sh".

## Tests

### Flux temps réel

Go to:

- [Public site: "https://public.gcp-wescale.slavayssiere.fr"](https://public.gcp-wescale.slavayssiere.fr).
- [Admin site: "https://iap.gcp-wescale.slavayssiere.fr"](https://iap.gcp-wescale.slavayssiere.fr).

### API

#### Génération d'aggregas

```language-bash
curl -X POST https://public.gcp-wescale.slavayssiere.fr/aggregator/stats | jq .
```

```language-bash
curl -X GET https://public.gcp-wescale.slavayssiere.fr/aggregator/stats/1 | jq .
curl -X GET https://us-central1-slavayssiere-sandbox.cloudfunctions.net/laststat | jq .
curl -X GET https://us-central1-slavayssiere-sandbox.cloudfunctions.net/getstat?id=1 | jq .
```

```language-bash
curl -X GET http://public.gcp-wescale.slavayssiere.fr/aggregator/top10 | jq .
```

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
