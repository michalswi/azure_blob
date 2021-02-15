

Official Azure SDK [source](https://docs.microsoft.com/en-us/azure/go/).  


```
go build

export AZURE_STORAGE_ACCOUNT=
export AZURE_STORAGE_KEY=
export TF_BACKEND_NAME=tfstate

./azure_blob -h


# CONTAINER

./azure_blob --action createContainer
Container tfstate created.

[optional]
$ az storage container list --account-name $AZURE_STORAGE_ACCOUNT --account-key $AZURE_STORAGE_KEY --output table
Name     Lease Status    Last Modified
-------  --------------  -------------------------
tfstate  unlocked        2020-02-07T10:59:10+00:00


# FILES

./azure_blob --action createUploadFile
Creating a dummy file: tweety-ce6d
File tweety-ce6d created.
Uploading the file with blob name: tweety-ce6d

./azure_blob --action createUploadFile
Creating a dummy file: tweety-b5c6
File tweety-b5c6 created.
Uploading the file with blob name: tweety-b5c6


# LIST blobs

# >> 'size' is expressed in bytes

./azure_blob --action list | jq
{
  "data": [
    {
      "id": 1,
      "file_name": "tweety-b5c6",
      "size": 20,
      "creation": "2020-02-07T11:04:14Z",
      "content_type": "application/octet-stream"
    },
    {
      "id": 2,
      "file_name": "tweety-ce6d",
      "size": 20,
      "creation": "2020-02-07T11:04:07Z",
      "content_type": "application/octet-stream"
    }
  ]
}

./azure_blob --action list | jq '.data[].file_name'
"tweety-b5c6"
"tweety-ce6d"


# DOWNLOAD blob

./azure_blob --action=download tweety-b5c6
Blob tweety-b5c6 downloaded.

cat /tmp/tweety-b5c6 
Tweety vs Sylvester


# DELETE container

./azure_blob --action deleteContainer
Container tfstate deleted.

[optional]
az storage container list --account-name $AZURE_STORAGE_ACCOUNT --account-key $AZURE_STORAGE_KEY --output table


# REMOVE local file (be aware which file you are removing!)

./azure_blob --action removeLocal /tmp/tweety-b5c6
File /tmp/tweety-b5c6 removed.
```