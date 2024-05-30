### Help

listandfind -h
Usage of ./listandfind:
  -access-key string
    	S3 Access Key
  -bucket string
    	Select a specific bucket
  -endpoint string
    	S3 endpoint URL
  -insecure
    	Disable TLS verification
  -prefix string
    	Select an object/prefix
  -recursive
    	Enable recursive listing
  -secret-key string
    	S3 Secret Key
  -skiperror
    	Skip other errors

### Example :-

```sh
> listandfind --access-key minioadmin --secret-key minioadmin --bucket test --endpoint http://localhost:9000
```
