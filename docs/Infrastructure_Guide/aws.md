---
weight: 4
---

# 4 // AWS
Amazon Web Services was the chosen cloud vendor for hosting this project's infrastructure. This takes the 
local deployment and moves it into the cloud. 


!!! question "What is AWS?"

    AWS (Amazon Web Services) is a cloud platform that gives you on-demand access to things like computing 
    power, storage, and databases, along with a wide variety of other services. Instead of setting up and
    maintaining physical servers, you can use AWS to quickly build, deploy, and scale applications of all
    sizes. From hosting websites to running machine learning models, AWS provides flexible tools to support
    different kinds of projects, with built-in options for security, monitoring, and global reach.
    
    View more [about AWS](https://aws.amazon.com/what-is-aws/)

## Services Used
* [IAM](#iam)
* [VPC](#vpc)
* [RDS](#rds)
* [EC2](#ec2)
* [Elastic IPs](#elastic-ips)
* [EventBridge](#eventbridge)
* [Secrets Manager](#secrets-manager)

!!! note
    
    The instructions below are all focused on creating AWS resources through the web console (can be helpful if new to AWS to learn how to
    navigate the console) . Since this project uses Terraform, all resources can be created and destroyed through IaC. Refer to the 
    [Terraform](terraform.md) page to create the resources through Terraform.

---

## IAM
[Identity and Access Management](https://docs.aws.amazon.com/IAM/latest/UserGuide/introduction.html) (IAM) is a service that helps control 
access to resources on AWS. With IAM, you can manage permissions that control which AWS resources users can access.

In the case of being a solo developer, different "users" will be treated as service accounts. One example for this usage is creating a service 
account that can only access [Elastic Container Registry](https://docs.aws.amazon.com/AmazonECR/latest/userguide/what-is-ecr.html) (ECR) in
a CI/CD pipeline that pushes a new image to ECR.

By setting up the service account to only access to ECR, 
the [Principle of Least Privilege](https://www.paloaltonetworks.com/cyberpedia/what-is-the-principle-of-least-privilege) is ensured.

Currently, the project has two service accounts:

1. `elastic-container-registry-user`
2. `terraform-user`

Service account names should make it obvious which resources they can access. AWS recommends adding services accounts to a group and then 
assigning permissions to that group. However, again with being a solo developer on a single project, an IAM group is not used.

This would make more sense if there were several users or different projects under the same account.

### Setup Instructions

1. Visit the IAM Console.
2. On the left, under **Access Management**, click on **Users**.
3. Click on the **Create User** button in the upper-right. 
4. Provide a name for the user. Ideally, the name should reflect the role or service it'll work with.
5. Click **next**.
6. Choose to **attack policies directly**
7. In the **Permission Policy** section, the option to attach an existing AWS managed policy or create a custom one exists.
   * AWS Managed Policies
     * Depending on what the user account is being created for, an existing AWS managed policy could suffice.
     * For example, this project's `elastic-container-registry-user` account has the AWS managed `SecretsManagerReadWrite` policy that
      allows it to read and write secrets from/to [Secret Manager](https://docs.aws.amazon.com/aws-managed-policy/latest/reference/SecretsManagerReadWrite.html).
   * Custom Policies
     * For even more fine-grain control and granting [least-privilege](https://docs.aws.amazon.com/IAM/latest/UserGuide/best-practices.html#grant-least-privilege), [custom](https://docs.aws.amazon.com/IAM/latest/UserGuide/access_policies_create.html) or _customer managed_ policies can be created.
     * For example, this project's `terraform-user` account has a policy that grants access to describe resources with in the EC2 service:
     ```json
     {
       "Version": "2012-10-17",
       "Statement": [
         {
           "Sid": "TerraformerRDSPermissions",
           "Effect": "Allow",
           "Action": [
             "ec2:DescribeVpcAttribute",
             "ec2:DescribeVpcs",
             "ec2:DescribeRouteTables",
             "ec2:DescribeSubnets",
             "ec2:DescribeInternetGateways",
             "ec2:DescribeSecurityGroups",
             "ec2:DescribeNatGateways",
             "ec2:DescribeVpcEndpoints"
          ],
          "Resource": "*"
         }
       ]
     }
     ```
8. ...
---

## VPC
_Virtual Private Cloud_

Creating a custom VPC instead of using the default one provides full control over network configuration, security, and isolation tailored to specific application requirements.
At first, the VPC will have public subnets to test the local version of Dagster to make sure everything is working correctly.
The VPC will then be modified to only have private subnet groups.

AWS creates a default VPC, but learning to create one can be invaluable when needing to trouble connection issues.

### Setup Instructions

1. Visit the VPC Console.
2. Choose to delete or keep the default VPC(s).
3. Click on **Create VPC**.
4. Under _Resources to create_, choose **VPC and more**.
5. For _Name tag auto-generation_, enter a name. Such as the project's name.
6. This project uses a CIDR block of `10.0.0.0/20` but a different can be chosen if needed. Learn more about [CIDR](https://aws.amazon.com/what-is/cidr/).
7. Choose 2 public subnets. (Only for testing Dagster locally).
8. Choose 2 private subnets.
9. Choose 0 NAT gateways since there is a cost to use them.
10. Create tags if wanted to organize resources on AWS.
11. Click **Create VPC**.

---

## RDS
_Relational Database Service_

Amazon RDS is a managed service that simplifies the setup, operation, and scaling of relational databases in the cloud.
In this project, [PostgreSQL](https://www.postgresql.org/) is the database engine of choice for storing the metadata of Dagster.

The cost to maintain the database with the project's configuration options come out to ~$15.00 USD.

### Setup Instructions
_**Note:** these are the configuration options that were chosen for this project. Costs are the main driver behind these options.
Feel free to choose any other options that could be more suitable._

1. Visit the [RDS console](https://console.aws.amazon.com/rds/home).
2. On the **dashboard**, there should be an option **Create a Database**. If not, click on **Databases** on the left menu.
   Then click **Create Database** in the upper-right.
3. Under **Engine Options**, choose PostgreSQL.
4. Under **Templates**, choose Dev/Test. If eligible, use Free tier.
5. Under **Availability and Durability**, choose Single Instance Deployment.
6. Under **Settings**, give the database a name and let AWS manage the credentials.
7. Under **Instance Configuration**, choose **Burstable Classes** and then select `t4g.micro` instance.
8. Change storage to the minimum of 20GB.
9. Under **Connectivity**, choose to not connect to an EC2 instance. This can be done later.
   * Choose the VPC that was created in the [previous step](#vpc).
   * The subnet from the VPC should be already selected.
   * Choose **no** for Public Access
   * Keep the **default** VPC security group.
   * This project does not have a preference on **Availability Zones** and uses the auto-generated **Certificate Authority**.
10. Under **Tags**, create a new tag if desired for resource organization.
11. Under **Database Authentication**, choose _password authentication_.
12. Under **Monitoring**, choose the standard version of _Database Insights_. All other options in this section can be left as default.
13. Review the **Estimated Monthly Costs**, make any changes if necessary, then click on _create database_.

---

## EC2
_Elastic Compute Cloud_

AWS EC2 (Elastic Compute Cloud) is a cloud service that provides resizable virtual servers to run applications and workloads on demand.

### Launch Instance
1. Visit the [EC2 console](https://console.aws.amazon.com/ec2).
2. Click on **Launch Instance**.
3. Provide a **name** for the virtual machine.
4. Under **Application and OS Images**, choose _Ubuntu 24.04 (HVM), SSD Volume Type 64-bit ARM_ or a different image if preferred.
5. Under **Instance Type**, choose `t2.small`.
6. Under **Key Pair (login)**, select a key pair or create a new one. If a new one is created, check for the `.pem` file in the downloads folder.
7. Under **Network Settings**:
   * Select the VPC created earlier.
   * Switch to a _public_ subnet to allow connection to the virtual machine.
   * Enable _Auto-assign public IP_.
   * For the _Firewall_, select the default security group that should've been created when setting up the VPC.
8. Under **Configure Storage**, leave as default.
9. Under **Advanced Details**, lease as default.

### Verify Connection to Instance
1. First, configure a trusted connection to the previously created RDS instance.
   * Visit the [RDS console](https://console.aws.amazon.com/rds/home).
   * Click on the RDS instance previously created.
   * Scroll down to the **Connected Compute Resources** section, in the **Actions** drop-down, click **Set up EC2 Connection**.
   * On the next screen, select the created EC2 instance from the drop-down. Then, select **Continue**.
   * On the **Review and Confirm** screen, review all information then click **Continue**.
2. SSH into machine.
   * Back in the [EC2 console](https://console.aws.amazon.com/ec2), click on the created EC2 instance.
   * In the top-right of the **Summary** section, click on the **Connect** button.
   * On the next page, click on the **SSH Client** tab.
   * Instructions on how to connect will be provided and `ssh` command will be provided. For example:
     * `ssh -i "dagster-vm-key-pair.pem" ubuntu@ec2-<ip-address-of-vm>.<region>.compute.amazonaws.com`
       * **Note:** Run this command in the directory of the `.pem` file.
       * **Note:** Since the virtual machine was created with the default VPC security group, 
       make sure the **Inbound Rules** of the security allows your IP address to connect.
   * The terminal should show an Ubuntu welcome screen once connected.

---

## Elastic IPs
An Elastic IP is a static public IPv4 address in AWS that can be assigned to an EC2 instance.

1. Visit the [EC2 console](https://console.aws.amazon.com/ec2).
2. On the left, under **Network & Security**, click on **Elastic IPs**.
3. In the upper-right, click on **Allocate Elastic IP Address**.
4. On the configuration screen, leave all as default, provide a tag if needed, then click on **Allocate**.
5. After the static IP address is allocated, select the IP address, then on the **Actions** dropdown, and finally, select **Associate Elastic IP Address**.
6. On the next page, choose instance as the **Resource Types**.
7. Select the previously created instance and its private IP.
8. Then, click on **Associate**.

---

## EventBridge
_Optional_

AWS [EventBridge](https://aws.amazon.com/eventbridge/) is a serverless event bus that can use rules to automatically 
trigger actions—such as starting or stopping RDS and EC2 instances—based on scheduled times or specific events.

EventBridge helps save on costs by stopping and starting the RDS and EC2 instance during off-peak hours. In the case 
of this project, it would be between the hours of 11:00PM and 6:00AM PST. This is optional but recommended.

### Create Policy
Once the role is created, a custom policy needs to be created to allow the role access to start/stop the targeted service.

1. Visit the [IAM console home](https://us-east-1.console.aws.amazon.com/iam/).
2.  On the left, click on **Policies**.
3. In the upper-right, click on **Create Policy**.
4. In the **Policy Editor** section, choose the **JSON** editor and paste in the following (Update `Action` to match the targeted service):
```json
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Effect": "Allow",
            "Action": [
                "ec2:StartInstances",
                "ec2:StopInstances"
            ],
            "Resource": "arn:aws:ec2:us-west-2:*:instance/*"
        }
    ]
}
```
5. Click **Next**. On the following screen, enter a name for the policy. For example, `ec2-start-stop`.
6. Add a description and/or tags. Click on **Create Policy**.

### Create Roles
EventBridge needs a role to connect and access other AWS resources. This project uses separate roles for each service.
The instructions below discuss the creation of a single role and a single attached policy.

1. Visit the [IAM console home](https://us-east-1.console.aws.amazon.com/iam/).
2. On the left, click on **Roles**.
3. In the upper-right, click on **Create Role**.
4. For **Trusted Entity Type**, choose **AWS Service**.
5. For **Use Case**, search and choose **EC2**. Leave at default selected option. Click **Next**.
6. For **Permissions Policies**, search and choose the custom policy created in the previous step.
7. Enter descriptive role name. For example, `EC2-StartStop-Scheduler-Role`.
8. Enter a description and/or tags. Then, click **Create Role**.

### Create Schedule

1. Visit the [EventBridge console home](https://console.aws.amazon.com/events).
   * Ensure AWS is in the correct region as the rest of the project's services.
   * If an EventBridge schedule is created in a different region than, say, an EC2 instance, it will not find it.
2. On the right, under **Scheduler**, click on **Schedules**.
3. Click on **Create Schedule**.

#### Specify Schedule Details

1. Specify Schedule Detail Page
   * Provide a name and description. 
   * For **Schedule Pattern**, choose **Recurring Schedule**. 
   * Choose the correct timezone and select **Cron-Base Schedule**. 
   * The `cron` expression only needs to run once a day. So, a schedule for starting the instance could be `0 6 * * ? *` - (run at 6am). 
   * Add a **Flexible Time Window** if needed.
   * The **Timeframe** section can be left blank. Click **Next**.
2. Select Target Page
   * Under **Target Detail**, select **All APIs** and search for and select EC2. 
   * In the API search box, search for `StartInstances` or `StopInstances`. 
   * In the next window, grab the **Instance ID** for the EC2 instance created earlier and paste it in to the JSON input.
    ```json
    {
      "InstanceIds": [
        "i-01dwf82r21c9da71f"
      ]
    }
    ```
   * Click **Next**.
3. Settings Page
   * Enable schedule state.
   * Choose `NONE` for **Action after Schedule Completion**.
   * Choose `None` for **Retry Policy and DLQ**
   * Customize encryption if needed.
   * In the **Permissions** section, choose **Use Existing Role**, then select the previously created role associated with this schedule’s service.
   * Click **Next**.
4. Review and Create Schedule Page
    * Review all the configuration changes then click **Create Schedule**.

---

## Secrets Manager
AWS Secrets Manager is a service for storing and managing sensitive information like database credentials, API keys, and passwords. 
It provides encryption, automatic rotation, IAM-based access control, and integrates with RDS and other AWS services to manage 
credentials programmatically.

### Supabase Secrets
The EC2 instance for Dagster needs to authenticate with Supabase to write the card data. 

Create a secret for Supabase credentials[^1]:

```shell
aws secretsmanager create-secret \
--name supabase \
--secret-string '{"password":"supabase-password","user":"supabase-user","database_uri":"postgresql://postgres.<user>:<password>@aws-0-us-east-2.pooler.supabase.com:6543/postgres"}' \
--region us-west-2
```

!!! note

    The connection string from Supabase has the username and password in it as such:
    ```shell
    postgresql://<postgres.rest-of-user-name>:<password>@aws-0-us-east-2.pooler.supabase.com:6543/postgres
    ```

    The username includes the `postgres` prefix like `postgres.rest-of-user-name`.

!!! note "Reminder"

    The EC2 instance has the `EC2-SecretsManager-Role` IAM role attached, which grants
    permissions to create and retrieve secrets from AWS Secrets Manager. This uses an
    IAM role (not an IAM user), so the instance receives temporary credentials automatically.

[^1]: Can be retrieved by creating a project as shown in [2 // Supabase](supabase.md).