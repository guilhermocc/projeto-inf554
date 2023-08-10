package instance

import (
	"fmt"
	"log"

	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws/ec2"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func CreateWebServerInstance(ctx *pulumi.Context, sg pulumi.IDOutput, snetID pulumi.String, instanceNumber int) *ec2.Instance {
	ami := pulumi.String("ami-053b0d53c279acc90") // Ubuntu 22.04 free tier
	instanceType := pulumi.String("t2.micro")     // free tier
	keyPair := pulumi.String("legion-key")
	webServerSG := sg.ToStringOutput()
	instanceName := fmt.Sprintf("web-server-%d", instanceNumber)

	instance, err := ec2.NewInstance(ctx, instanceName, &ec2.InstanceArgs{
		SubnetId:     snetID,
		Tags:         pulumi.StringMap{"Name": pulumi.String(instanceName)},
		Ami:          ami,
		InstanceType: instanceType,
		KeyName:      keyPair,
		VpcSecurityGroupIds: pulumi.StringArray{
			webServerSG,
		},
		UserData: nginxWebServerScript(instanceNumber),
		//AssociatePublicIpAddress: pulumi.Bool(false),
	},
	)

	if err != nil {
		log.Fatalf("failed to create instance: %v", err)
	}

	ctx.Export(instanceName, instance.PublicIp)

	return instance
}

func nginxWebServerScript(instanceNumber int) pulumi.String {
	return pulumi.String(fmt.Sprintf(`
#cloud-config
package_upgrade: true
packages:
  - nginx

write_files:
  - path: /var/www/html/index.html
    content: |
      <html>
      	<head>
      		<title>Bem vindo ao servidor web NGINX</title>
      	</head>
      	<body>
      		<h1>Ola, essa eh a %d instancia</h1>
      	</body>
      </html>

runcmd:
  - systemctl restart nginx
`, instanceNumber))
}
