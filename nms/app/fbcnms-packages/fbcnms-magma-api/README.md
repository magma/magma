## How to re-generate MagmaAPIBindings.js
1. Place an up-to-date `swagger.yml` into this directory.
   You can get the full `swagger.yml` at `{orc8r domain}/apidocs/v1/swagger.yml`. 
   (e.g. `https://localhost:9443/apidocs/v1/swagger.yml`)
2. Run `bin/generateAPIFromSwagger.sh`
