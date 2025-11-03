---
weight: 2
---

# 2 // Supabase

## Create an Account

Visit the Supabase [sign-up page](https://supabase.com/dashboard/sign-up) to create an account.
Signing in with GitHub is the easiest method.

## Retrieve PostgreSQL Connection String

1. Create an organization.
   * Provide an organization name.
   * Type of organization can be **Personal**.
   * Free plan is enough.
   * Click **Create**.
2. Create a new project.
   * Provide a project name.
   * Create or generate a database password.
   * Select the best region closest to users.
   * The **Security Options** and **Advanced Configuration** options can be left to default.
   * Click **Create**.
3. On the **Project Overview** page, in the top part of the page, click on **Connect**.
4. Under the **Connection String** section, change the **Method** dropdown to **Transaction Pooler**.
5. The connection string will be provided in the following format:
```shell
postgresql://postgres.[USERNAME]:[YOUR-PASSWORD]@aws-0-us-east-2.pooler.supabase.com:6543/postgres
```
6. Note the connection string for later instructions such as creating a secret of the string in AWS Secrets Manager[^1].

[^1]: Used in section: [3. AWS // Secrets Manager](aws.md#secrets-manager).