# Installing A Dev Email Server

### Prerequisites

- A Linux/Unix host

- Docker plus permissions to manage it

- The following line in `/etc/hosts`: `127.0.0.1 localhost.com`
  
  - `sudo` privileges are required

- Portainer or some GUI to manage Docker (optional)

- Thunderbird (optional)

### Steps

1. Install poste.io via the following command:
   
   ```
   docker run \
       -d \
       -p 1125:25 \
       -p 1180:80 \
       -p 1443:443 \
       -p 1110:110 \
       -p 1143:143 \
       -p 1465:465 \
       -p 1587:587 \
       -p 1993:993 \
       -p 1995:995 \
       -e TZ=America/Los_Angeles \
       -e HTTP_PORT=80 \
       -e HTTPS_PORT=443 \
       -e "DISABLE_CLAMAV=TRUE" \
       -e "DISABLE_RSPAMD=TRUE" \
       --restart unless-stopped \
       --name "PosteLocal" \
       -t analogic/poste.io
   ```

2. Navigate to `https://localhost:1443/admin/install/server` and configure the admin email and password. The admin email should end with `localhost.com`. **Note these down or lose access to the web panel.** 

3. Navigate to "Virtual Domains" > "localhost.com" > "Edit settings" > "Additional settings"
   
   1. Tick the box that says "Domain bin".
   
   2. Set the target email to be `admin@localhost.com` or the same one set during the setup phase if changed.
   
   3. Save changes

4. Navigate to "System settings" > "Advanced" and scroll down to "Connection blocking"
   
   1. Untick the box that says "Enabled" under that header. This prevents the email server from blacklisting your IP address due to an excess of connections.
   
   2. Save changes

5. Profit!

### Supplemental: Setup For Thunderbird

1. Navigate to "Settings" > "Account Settings" > "Account Actions" > "Add Mail Account".

2. Configure the basic account details with the following options:
   
   1. Your full name: <anything you want>
   
   2. Email address: `admin@localhost.com` (or the custom email you set earlier)
   
   3. Password: <the password set earlier>
   
   4. Tick "Remember password"

3. Click "Configure manually" and configure the connection with the following options:
   
   1. INCOMING SERVER
      
      1. Protocol: `POP3`
      
      2. Hostname: `localhost.com`
      
      3. Port: `1110`
      
      4. Connection security: `STARTTLS`
      
      5. Authentication method: `Autodetect`
      
      6. Username: `admin@localhost.com` (or the custom email you set earlier)
   
   2. OUTGOING SERVER
      
      1. Hostname: `localhost.com`
      
      2. Port: `1587`
      
      3. Connection security: `STARTTLS`
      
      4. Authentication method: `Autodetect`
      
      5. Username: `admin@localhost.com` (or the custom email you set earlier)
   
   3. Click the "Re-test" button to confirm that the connection was successful.

4. Click "Done" to save the settings.

5. At the "Add Security Exception" prompt, click "Confirm Security Exception" to add the self-signed certificate to the certificate store. If this comes up again while sending emails, simply repeat the process.

6. Click "Finish" to save the email account.

7. Profit!
