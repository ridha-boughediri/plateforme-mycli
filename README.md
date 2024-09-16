
# Plateforme MyCLI

## Introduction

Ce projet vise à développer une interface en ligne de commande (CLI) en **Go**, capable d'interagir avec un serveur S3. Le client CLI permet d'exécuter des commandes basiques pour la gestion des buckets et des objets stockés sur un serveur S3 (comme celui précédemment développé en **Go**). Ce projet vous permettra de vous familiariser avec la création d'une interface utilisateur en ligne de commande et la manipulation de requêtes HTTP vers une API S3.

## Objectifs du Projet

Le but de ce projet est de créer un CLI capable d'exécuter les actions suivantes sur un serveur S3 :

- **Lister les buckets** : Affiche tous les buckets disponibles sur le serveur.
- **Créer un bucket** : Ajoute un nouveau bucket.
- **Supprimer un bucket** : Supprime un bucket existant.
- **Téléverser un fichier dans un bucket**.
- **Lister les objets dans un bucket**.
- **Télécharger un fichier à partir d'un bucket**.
- **Supprimer un fichier d'un bucket**.
- **Intégration de tests** : Tests fonctionnels pour vérifier le bon fonctionnement de l'API, ainsi que la bonne exécution des sauvegardes et restaurations.
- **Bonus** : Compatibilité avec un serveur S3 réel (Amazon S3 ou MinIO) via configuration des URL, clés d'accès, et signatures d'authentification.

## Compétences visées

Ce projet permet de développer les compétences suivantes :

- **Interface CLI** :
  - Le CLI doit permettre l'utilisation de commandes telles que `list-buckets`, `create-bucket`, `upload-file`, `delete-file`, etc.
  - Les commandes doivent accepter des options en ligne de commande pour spécifier des paramètres (comme le nom du bucket ou le chemin du fichier à téléverser).
  
- **Gestion des erreurs** :
  - Le CLI doit gérer proprement les erreurs telles que la non-existence d'un bucket ou d'un fichier.

- **Communication avec le serveur S3** :
  - Utilisation de requêtes HTTP (GET, PUT, DELETE) pour interagir avec l’API du serveur S3.
  - Authentification et signature des requêtes pour assurer la compatibilité avec un serveur S3 sécurisé.

- **Configuration** :
  - Permettre la configuration de l'URL du serveur S3 et des clés d’accès via des fichiers de configuration ou des variables d’environnement.

## Langage et Technologies

- **Langage utilisé** : Go
- **Technologies** : HTTP, API REST, S3, CLI

## Prérequis

- **Go** installé sur votre machine (version 1.18 ou supérieure).
- Un serveur S3 local ou distant pour tester les fonctionnalités (ex : [MinIO](https://min.io/)).

## Installation

1. Clonez le dépôt :

   ```bash
   git clone https://github.com/username/plateforme-mycli.git
   cd plateforme-mycli
   ```

2. Compilez et installez le projet :

   ```bash
   go build -o mycli
   ```

3. Configurez l'accès au serveur S3 via un fichier de configuration ou des variables d’environnement :

   ```bash
   export S3_URL="https://your-s3-server.com"
   export S3_ACCESS_KEY="your-access-key"
   export S3_SECRET_KEY="your-secret-key"
   ```

## Utilisation

Exemples d'utilisation du CLI :

- **Lister les buckets** :

  ```bash
  ./mycli list-buckets
  ```

- **Créer un bucket** :

  ```bash
  ./mycli create-bucket my-new-bucket
  ```

- **Téléverser un fichier dans un bucket** :

  ```bash
  ./mycli upload-file mybucket /chemin/vers/le/fichier.txt
  ```

- **Lister les objets dans un bucket** :

  ```bash
  ./mycli list-objects mybucket
  ```

- **Télécharger un fichier depuis un bucket** :

  ```bash
  ./mycli download-file mybucket fichier.txt /chemin/vers/destination
  ```

- **Supprimer un fichier dans un bucket** :

  ```bash
  ./mycli delete-file mybucket fichier.txt
  ```

- **Supprimer un bucket** :

  ```bash
  ./mycli delete-bucket mybucket
  ```

## Tests

Pour exécuter les tests fonctionnels afin de vérifier le bon fonctionnement de votre API et de votre CLI, lancez la commande suivante :

```bash
go test ./...
```

Ces tests s'assureront que chaque fonctionnalité du CLI interagit correctement avec le serveur S3.

