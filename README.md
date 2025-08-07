# SSLPanel

Secure your website in just a few clicks — no technical expertise required.  
**SSLPanel** is a free and open-source tool that automates issuing, installing, and renewing SSL/TLS certificates to keep your website secure. You can self-host it or run it in a local development environment.

---

## 🚀 Self-Hosted Installation

Follow these steps to run SSLPanel on your own server:

1. **Download the latest release**  
   Get the latest version of SSLPanel from GitHub:  
   [Download SSLPanel](https://github.com/r2dtools/sslpanel/releases/latest/download/sslpanel.tar.gz)

2. **Extract the archive**  
   ```bash
   tar -xvzf sslpanel.tar.gz -C <your-directory>
   ```

3. **Configure SMTP for email notifications**  
   Edit the `.env.staging` file and set the following environment variables for your SMTP server:
   ```env
   CP_SMTP_HOST=<your-smtp-host>
   CP_SMTP_PORT=<your-smtp-port>
   CP_EMAIL_ADDRESS=<your-email-address>
   CP_EMAIL_PASSWORD=<your-email-password>
   ```
   These settings are required for account creation, password recovery, and other email notifications.

4. **Run SSLPanel**  
   Make the script executable and start the application:
   ```bash
   chmod +x <your-directory>/run.sh
   ./<your-directory>/run.sh
   ```

Once running, SSLPanel will be available at [http://localhost:5173](http://localhost:5173).

---

## 🛠 Development Setup

To run SSLPanel locally in development mode:

1. **Start the database**  
   ```bash
   docker compose -f docker-compose.dev.yml up -d
   ```

2. **Initialize the database**  
   ```bash
   make initdb
   ```

3. **Start the backend**  
   ```bash
   make start-back-dev
   ```

4. **Start the frontend**  
   ```bash
   make start-front-dev
   ```

Visit [http://localhost:5173](http://localhost:5173) to access the application.

> ✅ You can log in with the test account:  
> **Email:** `sslpanel@example.com`  
> **Password:** `1q2w3e4r`

---

## 📖 License

This project is open-source and available under the Apache License.  
View the source code on [GitHub](https://github.com/r2dtools/sslpanel).
