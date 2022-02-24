package azure

import (
	"bytes"
	"context"
	"log"
	url2 "net/url"

	"github.com/Azure/azure-storage-blob-go/azblob"
	"github.com/eztrade/kpi/graph/logengine"
)

func DownloadFileFromBlobURL(url string) (*bytes.Buffer, error) {
	// From the Azure portal, get your Storage account blob service URL endpoint.
	accountKey, accountName, _, _ := GetAzureAccountInfo()

	u, err := url2.Parse(url)
	if err != nil {
		logengine.GetTelemetryClient().TrackException(err.Error())
	}
	credential, err := azblob.NewSharedKeyCredential(accountName, accountKey)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	blobURL := azblob.NewBlobURL(*u, azblob.NewPipeline(credential, azblob.PipelineOptions{}))
	ctx := context.Background()

	response, err := blobURL.Download(ctx, 0, 0, azblob.BlobAccessConditions{}, false)
	if err != nil {
		return nil, err
	}

	blobData := &bytes.Buffer{}
	reader := response.Body(azblob.RetryReaderOptions{})
	blobData.ReadFrom(reader)
	defer reader.Close()
	return blobData, nil
}
