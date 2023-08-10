package loadbalancer

import (
	"fmt"
	"log"

	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws/alb"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func CreateWebServerLoadBalancer(ctx *pulumi.Context, vpcID string, subnets []string, webServerSG pulumi.IDOutput, instanceIDs ...pulumi.IDOutput) {
	lb, err := alb.NewLoadBalancer(ctx, "web-server-lb-inf554", &alb.LoadBalancerArgs{
		Name:             pulumi.String("web-server-lb-inf554"),
		LoadBalancerType: pulumi.String("application"),
		SecurityGroups: pulumi.StringArray{
			webServerSG,
		},
		Subnets: pulumi.ToStringArray(subnets),
		Tags:    pulumi.StringMap{"Name": pulumi.String("web-server-lb-inf554")},
	})

	if err != nil {
		log.Fatalf("failed to create load balancer: %v", err)
	}

	tg, err := alb.NewTargetGroup(ctx, "lb-target-group-inf554", &alb.TargetGroupArgs{
		Port:       pulumi.Int(80),
		Protocol:   pulumi.String("HTTP"),
		TargetType: pulumi.String("instance"),
		VpcId:      pulumi.String(vpcID),
	})

	if err != nil {
		log.Fatalf("failed to create target group: %v", err)
	}

	for i, instanceID := range instanceIDs {
		attachmentName := fmt.Sprintf("lb-target-group-attachment-%d", i)
		_, err = alb.NewTargetGroupAttachment(ctx, attachmentName, &alb.TargetGroupAttachmentArgs{
			TargetGroupArn: tg.Arn,
			TargetId:       instanceID,
			Port:           pulumi.Int(80),
		})

		if err != nil {
			log.Fatalf("failed to create target group attachment: %v", err)
		}
	}

	_, err = alb.NewListener(ctx, "web-server-lb-listener-inf554", &alb.ListenerArgs{
		LoadBalancerArn: lb.Arn,
		Port:            pulumi.Int(80),
		Protocol:        pulumi.String("HTTP"),
		DefaultActions: alb.ListenerDefaultActionArray{
			alb.ListenerDefaultActionArgs{
				Type:           pulumi.String("forward"),
				TargetGroupArn: tg.Arn,
			},
		},
	})

	if err != nil {
		log.Fatalf("failed to create listener: %v", err)
	}

	ctx.Export("lb", lb.DnsName.ToStringOutput())

}
