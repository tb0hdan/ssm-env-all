# ssm-env-all
Get ALL environment variables from Amazon Parameter Store

## SSM Doc

https://aws.amazon.com/blogs/compute/managing-secrets-for-amazon-ecs-applications-using-parameter-store-and-iam-roles-for-tasks/

## Usage:

AWS Go SDK will load credentials that were previously set using `awscli`, the only
thing that has to be specified is `AWS_REGION`.

`AWS_REGION=us-east-1 ./ssm-env-all ./test.sh`

### Multiple paths and script arguments

`AWS_REGION=us-east-1 ./ssm-env-all -path /common,/develop ./test.sh arg1 arg2 arg3`


## Building


`make`


# Thanks

https://github.com/remind101/ssm-env
