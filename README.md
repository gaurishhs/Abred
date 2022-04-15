# Abred

- Abred is a temporary voice channels discord bot built in [Golang](https://go.dev)

# Disclaimer

- I didn't have time and resources to continue hosting the bot, So i have discontinued this project.

# License

- This project is licensed under Attribution-NonCommercial-ShareAlike 4.0 International (CC BY-NC-SA 4.0)

# Self-Hosting

- Kindly abide by the license.
- All The Details here are provided assuming you have a ubuntu 20.04
  Build the application by running go build then you shall see a binary named abred, Create a service file for running the bot in background.
  Once the bot is ready, Install nginx and fill the nginx.conf file as follows:

```
server {
	root /var/www/html;
	index index.html index.htm index.nginx-debian.html;

	server_name api.abred.bar;

	location /manage {
		proxy_pass "http://localhost:8080/manage";
	        proxy_http_version 1.1;
    		proxy_set_header Upgrade $http_upgrade;
    		proxy_set_header Connection "upgrade";
	}

	location / {
		proxy_pass http://localhost:8080;
		proxy_http_version 1.1;
		proxy_set_header Upgrade $http_upgrade;
		proxy_set_header Connection 'upgrade';
		proxy_set_header Host $host;
		proxy_cache_bypass $http_upgrade;
	}

    listen [::]:443 ssl ipv6only=on; # managed by Certbot
    listen 443 ssl; # managed by Certbot
    ssl_certificate /etc/letsencrypt/live/api.abred.bar/fullchain.pem; # managed by Certbot
    ssl_certificate_key /etc/letsencrypt/live/api.abred.bar/privkey.pem; # managed by Certbot
    include /etc/letsencrypt/options-ssl-nginx.conf; # managed by Certbot
    ssl_dhparam /etc/letsencrypt/ssl-dhparams.pem; # managed by Certbot

}

map $sent_http_content_type $expires {
    "text/html"                 epoch;
    "text/html; charset=utf-8"  epoch;
    default                     off;
}

server {

    root /var/www/html;
    listen 443 ssl http2;
    listen [::]:443 ssl http2;
    ssl_certificate         /etc/ssl/cert.pem;
    ssl_certificate_key     /etc/ssl/key.pem;
    ssl_client_certificate /etc/ssl/cloudflare.crt;
    ssl_verify_client on;
    server_name abred.bar;
    gzip            on;
    gzip_types      text/plain application/xml text/css application/javascript;
    gzip_min_length 1000;


    location / {
	    expires $expires;
        proxy_redirect                      off;
        proxy_set_header Host               $host;
        proxy_set_header X-Real-IP          $remote_addr;
        proxy_set_header X-Forwarded-For    $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto  $scheme;
        proxy_read_timeout          1m;
        proxy_connect_timeout       1m;
	    proxy_pass http://localhost:3000;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection 'upgrade';
        proxy_set_header Host $host;
        proxy_cache_bypass $http_upgrade;
     }
}

server {
    if ($host = api.abred.bar) {
        return 301 https://$host$request_uri;
    } # managed by Certbot


	listen 80 default_server;
	listen [::]:80 default_server;

	server_name api.abred.bar;
    return 404; # managed by Certbot


}
```

- Don't forget to follow the ssl guide [here](https://www.digitalocean.com/community/tutorials/how-to-host-a-website-using-cloudflare-and-nginx-on-ubuntu-20-04)

- Once you've done all this, Install pm2 and follow nuxt's guide for deployment with pm2 from [here](https://nuxtjs.org/deployments/pm2/)
