---
weight: 2
---

# 2. AWS

Amazon Web Services was the chosen cloud vendor for hosting this project's infrastructure. 

This page will describe how to create each resource manually first to get used to the console. Then, in [3. Terraform](terraform.md), 
IaC (Infrastructure as Code) files will be created so that all resources can better managed and easily destroyed or rebuilt.

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
2. 

---

## VPC
_Virtual Private Cloud_

Creating a VPC should be one of the initial services to configure so that it's available for selection when setting up other services later on.
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

1. Choose PostgreSQL
2. Choose dev/test
3. Single zone
4. Burstable class
5. t4g.micro  instance
6. Change storage to 20GB