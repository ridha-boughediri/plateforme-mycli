
# Plateforme MyCLI

**Plateforme MyCLI** est une interface en ligne de commande (CLI) permettant de gérer des opérations sur des serveurs de stockage compatibles avec le protocole S3, tels que les serveurs locaux ou hébergés sur des services comme Contabo ou OVH.

## Fonctionnalités

- Créer, lister et supprimer des buckets.
- Uploader, télécharger et supprimer des fichiers dans des buckets.
- Gérer les authentifications pour interagir avec un serveur S3.
- Interface simple et intuitive basée sur `Cobra` et `Go`.

## Prérequis

Avant de commencer, assurez-vous d'avoir les éléments suivants installés sur votre machine :

- Go (version 1.18 ou plus récente)
- Un serveur compatible S3 (ex : MinIO, OVH, Contabo)

## Installation

1. Clonez le dépôt sur votre machine locale :
   ```bash
   git clone https://github.com/votre-utilisateur/plateforme-mycli.git
   ```

2. Accédez au répertoire du projet :
   ```bash
   cd plateforme-mycli
   ```

3. Installez les dépendances :
   ```bash
   go mod tidy
   ```

4. Compilez le projet :
   ```bash
   go build -o mycli main.go
   ```

## Configuration

Créez un fichier `.env` à la racine du projet et configurez-le avec les informations nécessaires. Voici un exemple :

```bash
S3_ENDPOINT=http://localhost:8080
S3_ACCESS_KEY=<votre-access-key>
S3_SECRET_KEY=<votre-secret-key>
```

## Utilisation des commandes

### Créer un bucket
```bash
./mycli create-bucket <bucket-name>
```
Crée un nouveau bucket avec le nom spécifié.

### Supprimer un bucket
```bash
./mycli delete-bucket <bucket-name>
```
Supprime le bucket spécifié.

### Lister les buckets
```bash
./mycli list-buckets
```
Affiche la liste des buckets existants.

### Uploader un fichier dans un bucket
```bash
./mycli upload-file <bucket-name> <file-path>
```
Téléverse un fichier dans le bucket spécifié.

### Télécharger un fichier depuis un bucket
```bash
./mycli download-file <bucket-name> <file-name>
```
Télécharge un fichier depuis le bucket.

### Supprimer un fichier d'un bucket
```bash
./mycli delete-file <bucket-name> <file-name>
```
Supprime un fichier d'un bucket.

### Lister les fichiers dans un bucket
```bash
./mycli list-files <bucket-name>
```
Liste les fichiers d'un bucket.

### Configurer l'authentification
```bash
./mycli configure-auth <access-key> <secret-key>
```
Configure les informations d'authentification pour accéder à votre serveur S3.

### Aide
```bash
./mycli --help
```
Affiche l'aide et la liste des commandes disponibles.

## Contribution

ridha , mathieu

## License

Ce projet est sous licence MIT. Voir le fichier [LICENSE](./LICENSE) pour plus d'informations.
