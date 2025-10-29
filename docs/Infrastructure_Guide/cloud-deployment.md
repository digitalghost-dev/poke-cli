---
weight: 5
---

# 5 // Cloud Deployment
Once the services are created and configured in AWS, the virtual machine can be set up with the needed 
tools/libraries to run the data pipelines in Dagster.

## Installing Tools and Libraries
Connect to the virtual machine and run the following commands to get everything set up:

1. Install AWS CLI
    * Download via `curl`:
      ```shell
      curl "https://awscli.amazonaws.com/awscli-exe-linux-aarch64.zip" -o "awscliv2.zip"
      ```
    * Install `unzip` program:
      ```shell
      sudo apt install unzip
      ```
    * Unzip the compressed folder:
      ```shell
      unzip awscliv2.zip
      ```
    * Run the installer:
      ```shell
      sudo ./aws/install
      ```
    * Set the default region:
      ```shell
      export AWS_DEFAULT_REGION=us-west-2
      # or set it to ~./bashrc file
      echo 'export AWS_DEFAULT_REGION=us-east-1' >> ~/.bashrc
      source ~/.bashrc
      ```
    * Run an `aws` command such as `secretsmanager` to verify AWS connectivity:
      ```shell
      aws secretsmanager list-secrets
      ```
2. Clone Repository
    * Create a new directory: 
      ```shell
      git init <dir-name>
      ```
    * Change into new directory:
      ```shell
      cd <dir-name>
      ```
    * Add the remote repository:
      ```shell
      git remote add -f origin https://github.com/digitalghost-dev/poke-cli/
      ```
    * Edit the `git` config file to turn on sparse checkout:
      ```shell
      git config core.sparseCheckout true
      ```
    * Tell `git` which directory to check out. Then, pull that directory.
      ```shell
      echo "card_data/" >> .git/info/sparse-checkout
      ```
    * Pull the repo into the local directory
      ```shell
      git pull origin main
      ```
    * Verify that `card_data/` directory was created.
      ```shell
      ls
      ```
3. Install Tools
    * Install `uv` for Python: 
      ```shell
      curl -LsSf https://astral.sh/uv/0.7.21/install.sh | sh
      ```
    * Add to `PATH`:
      ```shell
      `source $HOME/.local/bin/env
      ```
    * Install libraries from `pyproject.toml` file: 
      ```shell
      uv sync
      ```
    * Activate virtual environment: 
      ```shell
      source .venv/bin/activate
      ```
    * Create `dagster.yaml` file:
      ```bash
      mkdir -p ~/.dagster && cat > ~/.dagster/dagster.yaml << 'EOF'
      storage:
        postgres:
          postgres_db:
            username: postgres
            password: "rds-password"
            hostname: "rds-hostname"
            db_name: postgres
            port: 5432
          params:
            sslmode: require
      EOF
      ```
    * Set environment variables:
        * `echo 'export DAGSTER_HOME="$HOME/.dagster"' >> ~/.bashrc`
        * `echo 'export SUPABASE_USER="<supabase_user>"' >> ~/.bashrc`
        * `echo 'export SUPABASE_PASSWORD="<supabase_password>"' >> ~/.bashrc`
    * `source ~/.bashrc` - to load variables in current session.
4. Verify Dagster and Connectivity
    * `dg dev --host 0.0.0.0 --port 3000`
    * In the browser, visit `http://<ip-address-of-vm>:3000`

---

## Automating Startup with `systemd`
_Optional_

In order to save on costs, the EC2 and RDS instances are scheduled to start and 
stop once each day with AWS EventBridge. To automate the starting of the Dagster webservice, 
`systemd`, along with a couple of shell scripts, will be used to create this automation.

### Service Files

The `card_data/infrastructure/` directory has the following files:

1. `dagster.service` - the main `systemd` file for defining the Dagster service and environment.
2. `wait-for-rds.sh` - stored as `ExecStartPre` in `dagster.service` to check if the RDS instance is available.
3. `start-dagster.sh` - If the RDS instance is ready, this will run and start the Dagster web service.

Although the files are included in this repository, they need to be moved or created in a specific directory on
the Linux virtual machine.

#### Move Files

```shell
/home/ubuntu/wait-for-rds.sh

/home/ubuntu/start-dagster.sh

/etc/systemd/system/dagster.service
```

#### Create Files
First, create `dagster.service`:

* Run `nano /etc/systemd/system/dagster.service` then copy and paste the following code in the editor:

??? note "dagster.service"

    ```bash
    [Unit]
    Description=Dagster Development Server
    After=network-online.target
    Wants=network-online.target
    
    [Service]
    Type=simple
    User=ubuntu
    WorkingDirectory=/home/ubuntu/card_data/card_data
    Environment="AWS_DEFAULT_REGION=us-west-2"
    Environment="PATH=/home/ubuntu/card_data/card_data/.venv/bin:/usr/local/bin:/usr/bin:/bin"
    ExecStartPre=/home/ubuntu/wait-for-rds.sh
    ExecStart=/home/ubuntu/start-dagster.sh
    Restart=on-failure
    RestartSec=10
    StandardOutput=journal
    StandardError=journal
    
    [Install]
    WantedBy=multi-user.target
    ```

### Start Service
