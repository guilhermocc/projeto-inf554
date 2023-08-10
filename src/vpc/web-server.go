package vpc

import (
	"log"

	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws/ec2"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

const DefaultVPCID = "vpc-03e02ff18de66ce75"

func GetSubnetIDs(ctx *pulumi.Context) []string {
	subnets, err := ec2.GetSubnetIds(ctx, &ec2.GetSubnetIdsArgs{
		VpcId: DefaultVPCID,
	})
	if err != nil {
		log.Fatalf("failed to get subnet ids: %v", err)
	}

	return subnets.Ids
}
