package main

import (
	"github.com/pulumi/pulumi/sdk/v2/go/pulumi"
	"pulumi-for-aws-with-golang/components/resource"
)

func main() {
	pulumi.Run(resource.Setup)
}
