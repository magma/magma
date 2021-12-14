server {
  listen 443;
  ssl on;
  ssl_certificate /var/opt/magma/certs/nms/tls.crt;
  ssl_certificate_key /var/opt/magma/certs/nms/tls.key;
  location / {
     proxy_pass http://magmalte:8081;
     proxy_set_header Host $http_host;
     proxy_set_header X-Forwarded-Proto $scheme;
  }
}
