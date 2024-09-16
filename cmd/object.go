package cmd

import (
    "context"
    "fmt"
    "log"
    "os"
    "github.com/aws/aws-sdk-go-v2/aws"
    "github.com/aws/aws-sdk-go-v2/config"
    "github.com/aws/aws-sdk-go-v2/service/s3"
    "github.com/spf13/cobra"
)

func NewObjectCommand() *cobra.Command {
    var bucketName, objectName, filePath string

    var objectCmd = &cobra.Command{
        Use:   "object",
        Short: "Manage S3 Objects on Contabo",
    }

    // Sous-commande pour uploader un objet
    objectCmd.AddCommand(&cobra.Command{
        Use:   "upload",
        Short: "Upload an object to a bucket",
        Run: func(cmd *cobra.Command, args []string) {
            // Ouverture du fichier à uploader
            file, err := os.Open(filePath)
            if err != nil {
                log.Fatalf("Failed to open file %s, %v", filePath, err)
                return
            }
            defer file.Close()

            // Chargement de la configuration du SDK AWS
            cfg, err := config.LoadDefaultConfig(context.TODO(),
                config.WithRegion("eu-central-1"),  // Assurez-vous que la région est correcte pour votre bucket
                config.WithEndpointResolver(aws.EndpointResolverFunc(func(service, region string) (aws.Endpoint, error) {
                    return aws.Endpoint{
                        URL: "https://eu2.contabostorage.com/ridha-bucket",  // Remplacez ceci par votre endpoint Contabo S3
                    }, nil
                })),
            )
            if err != nil {
                log.Fatalf("Unable to load SDK config, %v", err)
                return
            }

            // Création d'un client S3
            s3Client := s3.NewFromConfig(cfg)

            // Upload de l'objet
            _, err = s3Client.PutObject(context.TODO(), &s3.PutObjectInput{
                Bucket: aws.String(bucketName),
                Key:    aws.String(objectName),
                Body:   file,
            })
            if err != nil {
                log.Fatalf("Unable to upload object, %v", err)
                return
            }

            fmt.Printf("Object %s uploaded successfully to bucket %s\n", objectName, bucketName)
        },
    })

    // Ajout de flags pour les opérations sur les objets
    objectCmd.PersistentFlags().StringVar(&bucketName, "bucket", "", "Bucket name (required)")
    objectCmd.PersistentFlags().StringVar(&objectName, "object", "", "Object name (required)")
    objectCmd.PersistentFlags().StringVar(&filePath, "file", "", "File path for upload (required for upload)")

    return objectCmd
}
