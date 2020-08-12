#!/bin/bash
docker images -f label=commit_id --format "{{.ID}}" | xargs docker inspect --format '{{ index .Config.Labels "commit_id" }}'
