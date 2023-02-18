# aws-lambda-rpc

## Installing
`go install github.com/chay-man/aws-lambda-rpc@latest`

### How to use
Once the tool is installed all you need to do is provide a valid json file to invoke your lambda with.

Example:
```shell
aws-lambda-rpc -f path/to/my/file.json
```

All other arguments are optional. --help shows all arguments. Example:
```
aws-lambda-rpc --help
Usage of ./aws-lambda-rpc:
  -debug
    	Enables debug logging
  -f string
    	The file to send in the InvokeRequest. Must be valid json
  -h string
    	The ip the lambda is running on. This must match the environment variable AWS_LAMBDA_RUNTIME_API (default "localhost")
  -p int
    	The port the lambda is listening to. This must match the environment variable _LAMBDA_SERVER_PORT (default 8080)
  -t int
    	The time to wait in seconds until sending the done signal to the lambda. Default is 1 day (default 86400)
```