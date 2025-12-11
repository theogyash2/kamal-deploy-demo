gem install kamal

kamal init

1. Create the EC2 instance with following user data:

```
#cloud-config
runcmd:
  # Ensure root .ssh directory
  - sudo mkdir -p /root/.ssh
  - sudo chmod 700 /root/.ssh
  # Copy ubuntu's authorized_keys
  -  sudo cp /home/ubuntu/.ssh/authorized_keys /root/.ssh/authorized_keys.tmp
  # Remove AWS forced command restrictions from the key
  -  sudo sed 's/^no-port-forwarding[^"]*exit 142" //' /root/.ssh/authorized_keys.tmp > /root/.ssh/authorized_keys
  # Fix permissions
  -  sudo chmod 600 /root/.ssh/authorized_keys
  -  sudo chown root:root /root/.ssh/authorized_keys
  # Clean up temporary file
  -  sudo rm -f /root/.ssh/authorized_keys.tmp
```


2. Enable traffic: ()
To make the instance available over the ssh and load balancer enable the port: (EC2 -> Security Groups -> Select Name -> Inbound Rules -> Edit Inbound Rules -> Add rule for 0.0.0.0 -> Save)
22 
80 
443 


3. Dockerfile

```
#exampledockerfile
# Use a production-ready Node.js image
FROM node:18-alpine 

# Set the working directory
WORKDIR /app

# Copy package files and install dependencies
COPY package*.json ./
RUN npm install --production

# Copy the rest of the application code
COPY . .

# Expose the application port
EXPOSE 3000

# Define the command to start the app
CMD ["node", "app.js"]
Or 
CMD ["sh", "-c", "node ${MAIN}.js"] #changes according to env variable

```

4. Initialize Configuration Files
```
cd project-folder-name
```
```
kamal init
```

This generates: config/deploy.yml and a secret file for environment variables (.env).
Setup the AWS account credentials:

``` aws configure --profile new_profile_name ```

3.3 Edit config/deploy.yml

Update this file with your server details, registry credentials, and application domain.
YAML
# config/deploy.yml


service: my-node-app #Unique name for all containers and services                
image: prod-node #repo name
deploy_timeout: 120 #default is 30
builder:
  arch: amd64




servers:
 Web:
   hosts:
     - 54.12.34.56 (#replace with your instance public ip)
     - 54.34.87.87 (2nd instance ip if need multiple same containers)
    env:
  	# The 'clear' section is for non-sensitive env data
  	Clear:
    		- RAILS_ENV: production
		- MAIN: web (the variable that decides which service to run, one dockerfile and 2 services based on the env variable)
  	Secrets:
    	- DB_USERNAME  # sensitive environment variable fetched from aws secret manager 


 Worker:
   hosts:
     - 54.12.34.56 (#replace with your instance public ip)
     - 54.34.87.87 (2nd instance ip if need multiple same containers)
   env:
     MAIN: worker
   # cmd: sh -c node worker.js


 
NOTE: For multiple containers of the same application or service, we need to use multiple instances, each instance will have one container of that particular service or application.


ssh:
 user: root               (always go for the root user)
 keys:
      - /path/to/your/ssh/key.pem


proxy:
   app_port: 3000
   ssl: true
   forward_headers: true
   host: foo.example.com (domain name if one)
#Or (if multiple use hosts instead of host)
 #hosts:
    #- foo.example.com
    #- bar.example.comdomain/subdomain
  healthcheck:
   path: /up              
   interval: 3
   timeout: 3


registry:
 server: 123456789012.dkr.ecr.us-east-1.amazonaws.com (link to the repo)
 username: AWS   (#for ecr it remains constant)
 Password: <%= `aws ecr get-login-password --region us-east-1 --profile your-aws-profile` %> (the profile should be set in our local shell always)













Secrets file (auto generated):
#!/usr/bin/env bash
set -e
# Fetch the entire JSON secret payload from AWS Secrets Manager
# This command requires the AWS CLI to be set up locally.
SECRETS_PAYLOAD=$(kamal secrets fetch --adapter aws_secrets_manager --account default --from myapp/production/env)
# --- Extract and Export Individual Environment Variables ---
# The left side (e.g., DB_USERNAME) is the final environment variable name
# that your application code will use inside the container.
# Export the value of the 'DB_USER' key from the JSON payload 
DB_USERNAME=$(kamal secrets extract DB_USER $SECRETS_PAYLOAD)



4. 🚀 Deployment Steps
4.1 Commit Changes
Kamal uses the latest Git commit hash to tag your Docker image. Always ensure your changes are committed before deploying.
Bash
git add .
git commit -m "Initial commit for Kamal deployment"

4.2 Initial Server Setup
Run the setup command once. This connects to your server, installs Docker, and sets up the Kamal Proxy (Traefik).
Bash
kamal setup

This command may take a few minutes as it provisions the server.
4.3 DNS Configuration (Mandatory)
Ensure the domain name you listed in deploy.yml (myapp.yourdomain.com) has an A Record pointing to the public IP of your remote server (54.12.34.56).
4.4 Deploy the Application
Run the main deployment command. Kamal will build the image, push it to your registry, pull it to the server, run the new container, wait for the /up health check, and finally switch the proxy.
Bash
kamal deploy

4.5 Verify Status
You can check the status and logs from your local machine:
Bash
# Check container status
kamal status

# View live container logs
kamal app logs -f

Navigate to http://myapp.yourdomain.com in your browser. If successful, you should see: "Hello from Kamal Deployed App!"

5. 🔁 Ongoing Operations
Command
Purpose
kamal deploy
Deploys a new version (builds, pushes, pulls, switches proxy).
kamal rollback
Rolls back to the previous healthy deployed image (zero-downtime).
kamal env push
Pushes updated environment variables to the remote server.
kamal accessory add redis
Adds an accessory service like Redis or a database to the server.
kamal app exec -i bash
Opens an interactive bash shell inside the running container.



