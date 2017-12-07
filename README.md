# ssm-env-all
Get ALL environment variables from Amazon Parameter Store
Used in production (mic drop)

## SSM Doc

https://aws.amazon.com/blogs/compute/managing-secrets-for-amazon-ecs-applications-using-parameter-store-and-iam-roles-for-tasks/

## Usage:

AWS Go SDK will load credentials that were previously set using `awscli`, the only
thing that has to be specified is `AWS_REGION`.

`AWS_REGION=us-east-1 ./ssm-env-all ./test.sh`

### Multiple paths and script arguments

`AWS_REGION=us-east-1 ./ssm-env-all -path /common,/develop ./test.sh arg1 arg2 arg3`

### Ignoring authentication errors (not recommended)

`AWS_REGION=us-east-1 SSM_IGNORE_ERRORS=1 ./ssm-env-all env`

or

`AWS_REGION=us-east-1 ./ssm-env-all -ignore-errors env`



## Dockerfile (sample)


```
FROM alpine:latest

# Prepare entrypoint to populate env variables from AWS Parameter store
# ssm-env-all should get comma separated list of paths via the 'SSM_PATH' variable,
# e.g. SSM_PATH="/develop/generic,/develop/be"
RUN apk update && apk upgrade && apk add --update curl
RUN curl -L https://<>/ssm-env-all.linux.amd64 > /usr/local/bin/ssm-env-all && chmod +x /usr/local/bin/ssm-env-all
ENTRYPOINT ["/usr/local/bin/ssm-env-all"]

# # Execute service
CMD [ "yourapp", "start" ]
```

## Building


`make`


# Thanks

https://github.com/remind101/ssm-env
