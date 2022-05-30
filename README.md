

Official Azure SDK [source](https://docs.microsoft.com/en-us/azure/go/).  


```
go build

export AZURE_STORAGE_ACCOUNT=
export AZURE_STORAGE_KEY=
export SA_CONTAINER_NAME=tfstate

./azure_blob -h


# create container

$ ./azure_blob --action createContainer
2022/05/26 19:13:50 Container tfstate created.

[optional]$ az storage container list --account-name $AZURE_STORAGE_ACCOUNT --account-key $AZURE_STORAGE_KEY --output table


# create files locally and upload to blobs (one file one blob)

$ ./azure_blob --action createUploadFile
2022/05/26 19:16:55 Creating a dummy file: tweety-c30b
2022/05/26 19:16:55 File tweety-c30b created.
2022/05/26 19:16:55 Uploading the file with blob name: tweety-c30b

$ ./azure_blob --action createUploadFile
2022/05/26 19:17:31 Creating a dummy file: tweety-8230
2022/05/26 19:17:31 File tweety-8230 created.
2022/05/26 19:17:31 Uploading the file with blob name: tweety-8230

$ ./azure_blob --action createUploadFile
2022/05/26 19:17:34 Creating a dummy file: tweety-35f7
2022/05/26 19:17:34 File tweety-35f7 created.
2022/05/26 19:17:34 Uploading the file with blob name: tweety-35f7


# list blobs

# >> 'size' is expressed in bytes

$ ./azure_blob --action list | jq
{
  "data": [
    {
      "id": 1,
      "file_name": "tweety-35f7",
      "size": 20,
      "creation": "2022-05-26T17:17:35Z",
      "content_type": "application/octet-stream"
    },
    {
      "id": 2,
      "file_name": "tweety-8230",
      "size": 20,
      "creation": "2022-05-26T17:17:32Z",
      "content_type": "application/octet-stream"
    },
    {
      "id": 3,
      "file_name": "tweety-c30b",
      "size": 20,
      "creation": "2022-05-26T17:16:56Z",
      "content_type": "application/octet-stream"
    }
  ]
}

$ ./azure_blob --action list | jq '.data[].file_name'
"tweety-35f7"
"tweety-8230"
"tweety-c30b"

[optional]$ az storage blob list --container-name $SA_CONTAINER_NAME --account-name $AZURE_STORAGE_ACCOUNT --account-key $AZURE_STORAGE_KEY --output table


# download data from blob


$ ./azure_blob --action download tweety-35f7
Blob tweety-b5c6 downloaded.

$ cat /tmp/tweety-35f7
Tweety vs Sylvester


# delete container

./azure_blob --action deleteContainer
2022/05/26 19:32:12 Container tfstate deleted.

[optional]
az storage container list --account-name $AZURE_STORAGE_ACCOUNT --account-key $AZURE_STORAGE_KEY --output table


# remove local files

$ ./azure_blob --action removeLocal ./tweety-35f7
2022/05/26 13:35:11 File ./tweety-35f7 removed.
```