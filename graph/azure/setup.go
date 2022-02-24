package azure

import (
	"context"
	"fmt"
	"github.com/Azure/azure-storage-blob-go/azblob"
	"net/url"
	"os"
)

func GetAzureAccountInfo() (string, string, string, string) {
	azrKey := os.Getenv("AZURE_BLOB_KEY")
	azrBlobAccountName := os.Getenv("AZURE_BLOB_ACCOUNT_NAME")
	azrPrimaryBlobServiceEndpoint := fmt.Sprintf("https://%s.blob.core.windows.net/", azrBlobAccountName)
	azrBlobContainer := os.Getenv("AZURE_BLOB_UPLOAD_PATH")
	return azrKey, azrBlobAccountName, azrPrimaryBlobServiceEndpoint, azrBlobContainer
}

func UploadBytesToBlob(b []byte, blobname string) (string, error) {
	azrKey, accountName, endPoint, container := GetAzureAccountInfo()
	u, _ := url.Parse(fmt.Sprint(endPoint, container, "/", blobname))
	credential, errC := azblob.NewSharedKeyCredential(accountName, azrKey)
	if errC != nil {
		return "", errC
	}

	blockBlobUrl := azblob.NewBlockBlobURL(*u, azblob.NewPipeline(credential, azblob.PipelineOptions{}))
	ctx := context.Background() // We create an empty context (https://golang.org/pkg/context/#Background)

	// Provide any needed options to UploadToBlockBlobOptions (https://godoc.org/github.com/Azure/azure-storage-blob-go/azblob#UploadToBlockBlobOptions)
	o := azblob.UploadToBlockBlobOptions{
		BlobHTTPHeaders: azblob.BlobHTTPHeaders{
			ContentType: "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet",
		},
		Parallelism: 16,
	}

	_, errU := azblob.UploadBufferToBlockBlob(ctx, b, blockBlobUrl, o)
	return blockBlobUrl.String(), errU
}

