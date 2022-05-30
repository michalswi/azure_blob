package main

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/url"
	"os"
	"time"

	"github.com/Azure/azure-storage-blob-go/azblob"
)

type azureBlobs struct {
	ID           int       `json:"id"`
	Filename     string    `json:"file_name"`
	Size         int64     `json:"size"`
	CreationTime time.Time `json:"creation"`
	ContentType  string    `json:"content_type"`
}

var datas []azureBlobs
var countID int
var finalJSON = make(map[string]interface{})

func main() {

	// take an action
	action := flag.String("action", "", "createContainer|createUploadFile\nlist|download + <file_name>\ndeleteContainer|removeLocal + <file_name>")
	flag.Parse()

	accountName, accountKey, containerName :=
		os.Getenv("AZURE_STORAGE_ACCOUNT"),
		os.Getenv("AZURE_STORAGE_KEY"),
		os.Getenv("SA_CONTAINER_NAME")
	if len(accountName) == 0 || len(accountKey) == 0 {
		log.Fatal("Either the AZURE_STORAGE_ACCOUNT or AZURE_STORAGE_KEY or SA_CONTAINER_NAME environment variable is not set")
	}

	// create a default request pipeline
	credential, err := azblob.NewSharedKeyCredential(accountName, accountKey)
	if err != nil {
		log.Fatalf("Invalid credentials with error: \n%v", err)
	}
	p := azblob.NewPipeline(credential, azblob.PipelineOptions{})

	// from the Azure portal, get your storage account blob service URL endpoint
	URL, _ := url.Parse(
		fmt.Sprintf("https://%s.blob.core.windows.net/%s", accountName, containerName))

	// create a ContainerURL object that wraps the container URL and a request pipeline to make requests
	containerURL := azblob.NewContainerURL(*URL, p)

	// never-expiring context
	ctx := context.Background()

	switch *action {
	case "list":
		// list blobs from a specific container
		listBlobs(ctx, containerURL, containerName)
	case "createContainer":
		_, err = containerURL.Create(ctx, azblob.Metadata{}, azblob.PublicAccessNone)
		if err != nil {
			log.Fatalf("Can't create container: %v", err)
		}
		log.Printf("Container %s created.", containerName)
	case "createUploadFile":
		// create and upload file
		createFile(ctx, containerURL)
	case "download":
		// download blobs
		fileName := flag.Args()[0]
		downloadFile(ctx, containerURL, fileName)
	case "deleteContainer":
		_, err = containerURL.Delete(ctx, azblob.ContainerAccessConditions{})
		if err != nil {
			log.Fatalf("Can't delete container: \n%v", err)
		}
		log.Printf("Container %s deleted.", containerName)
	case "removeLocal":
		// remove local file
		fileName := flag.Args()[0]
		os.Remove(fileName)
		log.Printf("File %s removed.\n", fileName)
	}
}

// randToken generate random string
func randToken(len int) string {
	b := make([]byte, len)
	rand.Read(b)
	return fmt.Sprintf("%x", b)
}

// createFile create a dummy file locally and upload to a container
func createFile(ctx context.Context, containerURL azblob.ContainerURL) {
	data := []byte("Tweety vs Sylvester\n")
	fileName := fmt.Sprintf("tweety-%s", randToken(2))
	log.Printf("Creating a dummy file: %s", fileName)
	err := ioutil.WriteFile(fileName, data, 0644)
	if err != nil {
		log.Fatalf("Can't create the file: %v", err)
	}

	blobURL := containerURL.NewBlockBlobURL(fileName)
	file, err := os.Open(fileName)
	if err != nil {
		log.Fatalf("Can't open the file: %v", err)
	}
	log.Printf("File %s created.\n", fileName)
	defer file.Close()

	// You can use the low-level PutBlob API to upload files. Low-level APIs are simple wrappers for the Azure Storage REST APIs.
	// Note that PutBlob can upload up to 256MB data in one shot.
	// Details: https://docs.microsoft.com/en-us/rest/api/storageservices/put-blob
	// Following is commented out intentionally because we will instead use UploadFileToBlockBlob API to upload the blob
	// _, err = blobURL.PutBlob(ctx, file, azblob.BlobHTTPHeaders{}, azblob.Metadata{}, azblob.BlobAccessConditions{})
	// handleErrors(err)

	// The high-level API UploadFileToBlockBlob function uploads blocks in parallel for optimal performance,
	// and can handle large files as well.
	// This function calls PutBlock/PutBlockList for files larger 256 MBs, and calls PutBlob for any file smaller
	log.Printf("Uploading the file with blob name: %s", fileName)
	_, err = azblob.UploadFileToBlockBlob(ctx, file, blobURL, azblob.UploadToBlockBlobOptions{
		BlockSize:   4 * 1024 * 1024,
		Parallelism: 16})
	if err != nil {
		log.Fatalf("Can't upload the file: \n%v", err)
	}
}

// downloadFile download blobs from a container
func downloadFile(ctx context.Context, containerURL azblob.ContainerURL, fileName string) {
	blobURL := containerURL.NewBlockBlobURL(fileName)
	downloadResponse, err := blobURL.Download(ctx, 0, azblob.CountToEnd, azblob.BlobAccessConditions{}, false, azblob.ClientProvidedKeyOptions{})
	if err != nil {
		log.Fatal(err)
	}
	// NOTE: automatically retries are performed if the connection fails
	bodyStream := downloadResponse.Body(azblob.RetryReaderOptions{MaxRetryRequests: 20})
	// read the body into a buffer
	downloadedData := bytes.Buffer{}
	_, err = downloadedData.ReadFrom(bodyStream)
	if err != nil {
		log.Fatal(err)
	}

	// downloaded blob data is in downloadData's buffer. :Let's print it
	// log.Printf("Downloaded the blob: " + downloadedData.String())

	// save downloaded file
	err = ioutil.WriteFile("/tmp/"+fileName, downloadedData.Bytes(), 0644)
	if err != nil {
		log.Fatalf("Can't download the file %s: %v", fileName, err)
	}
	log.Printf("Blob %s downloaded.", fileName)
}

// listBlobs list blobs from a container
func listBlobs(ctx context.Context, containerURL azblob.ContainerURL, containerName string) {
	// log.Printf("List blobs from %s container.", containerName)
	for marker := (azblob.Marker{}); marker.NotDone(); {
		// get a result segment starting with the blob indicated by the current Marker.
		listBlob, err := containerURL.ListBlobsFlatSegment(ctx, marker, azblob.ListBlobsSegmentOptions{})
		if err != nil {
			log.Fatalf("Can't list blobs from container: %v", err)
		}

		// listBlobs returns the start of the next segment; you MUST use this to get
		// the next segment (after processing the current result segment).
		marker = listBlob.NextMarker

		for _, blobInfo := range listBlob.Segment.BlobItems {
			countID++
			datas = append(datas, azureBlobs{
				ID:           countID,
				Filename:     blobInfo.Name,
				Size:         *blobInfo.Properties.ContentLength,
				CreationTime: blobInfo.Properties.LastModified,
				ContentType:  *blobInfo.Properties.ContentType,
			})
		}
	}
	finalJSON["data"] = datas
	res2B, _ := json.Marshal(finalJSON)
	fmt.Println(string(res2B))
}
