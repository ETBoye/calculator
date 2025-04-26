# Calculator
## Prerequisites for running the app
This app uses
 - Docker compose for easy deployment of the different services.
 - Dozzle for container observability
	 - For the real deployment on digital ocean, I want dozzle behind basic authentication. For testing, you can run 
	 `cp nginx/.htpasswd.unsafe-admin-admin nginx/.htpasswd` to create a simple user with username 'admin' and password 'admin'.
	

## Testing locally
You can test the whole application by using docker-compose by running
`POSTGRES_USER=user POSTGRES_PASSWORD=secret APP_ENV=test ./run.sh`

This exposes
 - The whole application as it will run on the website on  `http://localhost:8089`
	 - Note that for example `POST http://localhost:8089/api/compute` is proxied to the backend as `POST /compute`
 - Dozzle on `http://localhost:8089/dozzle`
 - The backend on `http://localhost:8080`


## Deployment on digital ocean
This application is running at [https://calculator.etboye.dk](https://calculator.etboye.dk) and deployed using a droplet on digital ocean. A few notes:
  - The droplet has nginx installed and uses a certificate from LetsEncrypt. Everything from nginx is proxied to `localhost:8089` where the nginx from the docker-compose app is exposed.
  - The droplet is behind a firewall configured directly from digital ocean.