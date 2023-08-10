package main

import (
	"projeto-inf0554-aws/src/instance"
	"projeto-inf0554-aws/src/loadbalancer"
	"projeto-inf0554-aws/src/secgroup"
	vpc2 "projeto-inf0554-aws/src/vpc"

	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

const numberOfInstances = 2

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		// Get subnets from default VPC
		subnets := vpc2.GetSubnetIDs(ctx)

		// Create web server security group
		webServerSg := secgroup.CreateWebServerSG(ctx, vpc2.DefaultVPCID)

		// Create web server instances
		instanceIDs := make([]pulumi.IDOutput, numberOfInstances)
		for i := 0; i < numberOfInstances; i++ {
			inst := instance.CreateWebServerInstance(ctx, webServerSg.ID(), pulumi.String(subnets[0]), i+1)
			instanceIDs[i] = inst.ID()
		}

		// Create web server load balancer
		loadbalancer.CreateWebServerLoadBalancer(ctx, vpc2.DefaultVPCID, subnets, webServerSg.ID(), instanceIDs...)

		return nil
	})
}
