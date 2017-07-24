# mesos2iam

Provide IAM credentials to containers running inside a mesos cluster based on environment variable of the container.



## Usage

#### IAM roles

It is necessary to create an IAM role for mesos agents. This role should be able to assume other roles.



#### Running mesos2iam

mesos2iam is supposed to run as a daemon inside mesos agents.

Environment variables to be set before running mesos2iam and there default values.

```
MESOS2IAM_LISTENING_IP             		= "0.0.0.0"
MESOS2IAM_HOST_IP                  		= ""
MESOS2IAM_SERVER_PORT			   		= 51679
MESOS2IAM_AWS_CONTAINER_CREDENTIALS_IP 	= "169.254.170.2"
MESOS2IAM_CREDENTIALS_URL 				= "http://127.0.0.1:8080"
MESOS2IAM_PREFIX 						= "TARDIS_SCHID="
```

**Build**

```
make build
```

##### **Run**

```
build/mesos2iam
```

##### Example

Run mesos2iam with a custom container identifier

```
MESOS2IAM_PREFIX="MY_CONTAINER_ID=" build/mesos2iam
```

Having a mesos task running with next environment variables:

```
MY_CONTAINER_ID=1234 
AWS_CONTAINER_CREDENTIALS_RELATIVE_URI=/v2/credentials
```

Once any aws-sdk try to retrieve credentails from amazon 169.254.170.2/v2/credentials/

The mesos2iam retrieves the container id and do the call to 

```
${MESOS2IAM_CREDENTIALS_URL}/credentials/1234
```

and returns the credentials to the aws-sdk

```
{
    "RoleArn": "arn:aws:iam::ACCOUNTID:role/configured-aws-credentials-for-1234",
    "AccessKeyId": "...",
    "SecretAccessKey": "...",
    "Token": "...",
    "Expiration": "2017-08-01T21:06:06Z"
}
```

So its up to every user how they implement the service that returns the aws credentials.