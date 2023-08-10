package secgroup

import (
	"log"

	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws/ec2"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func CreateWebServerSG(ctx *pulumi.Context, vpcID pulumi.ID) *ec2.SecurityGroup {
	sg, err := ec2.NewSecurityGroup(ctx, "web-server-sg-inf554", &ec2.SecurityGroupArgs{
		Name:        pulumi.String("web-server-sg-inf554"),
		VpcId:       vpcID,
		Description: pulumi.String("Allows SSH traffic to bastion hosts"),
		Ingress: ec2.SecurityGroupIngressArray{
			ec2.SecurityGroupIngressArgs{
				Protocol:    pulumi.String("tcp"),
				ToPort:      pulumi.Int(22),
				FromPort:    pulumi.Int(22),
				Description: pulumi.String("Allow inbound TCP 22"),
				CidrBlocks:  pulumi.StringArray{pulumi.String("0.0.0.0/0")},
			},
			ec2.SecurityGroupIngressArgs{
				Protocol:    pulumi.String("tcp"),
				ToPort:      pulumi.Int(80),
				FromPort:    pulumi.Int(80),
				Description: pulumi.String("Allow inbound HTTP 80"),
				CidrBlocks:  pulumi.StringArray{pulumi.String("0.0.0.0/0")},
			},
		},
		Egress: ec2.SecurityGroupEgressArray{
			ec2.SecurityGroupEgressArgs{
				Protocol:    pulumi.String("-1"),
				ToPort:      pulumi.Int(0),
				FromPort:    pulumi.Int(0),
				Description: pulumi.String("Allow all outbound traffic"),
				CidrBlocks:  pulumi.StringArray{pulumi.String("0.0.0.0/0")},
			},
		},
	})

	if err != nil {
		log.Fatalf("failed to create security group: %v", err)
	}

	return sg
}
