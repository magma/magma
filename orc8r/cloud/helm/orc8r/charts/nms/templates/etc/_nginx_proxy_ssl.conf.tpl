server {
  listen 443;
  ssl on;
  ssl_certificate /etc/nginx/conf.d/nms_nginx.pem;
  ssl_certificate_key /etc/nginx/conf.d/nms_nginx.key.pem;
  location / {
     proxy_pass http://magmalte:8081;
     proxy_set_header Host $http_host;
     proxy_set_header X-Forwarded-Proto $scheme;
  }
}
