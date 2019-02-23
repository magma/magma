# NMS X

# Building

    docker-compose -f docker/docker-compose-build.yml build magmalte

# Running

1. Create a .env file.  See [.env.example](.env.example)

2. We use [chamber](https://github.com/segmentio/chamber) for storing some environment variables.
See [docker/bin/start](docker/bin/start#L3) for more details.  You don't need to use this, but will
need to mount the keys into the docker volume and set `API_CERT_FILENAME` / `API_PRIVATE_KEY_FILENAME`.

3a. To run with chamber (the default):

    docker run --env-file=.env -t -i --rm magmalte

3b. Without chamber:

    docker run --env-file=.env -t -i --rm magmalte yarn run start
