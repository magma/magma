FROM debian:stretch AS %%PKG%%

COPY --from=%%DOCKER_REGISTRY%%/magma /cache/magma/install /cache/install
COPY --from=%%DOCKER_REGISTRY%%/magma /cache/magma/repo /cache/magma/repo
COPY --from=%%DOCKER_REGISTRY%%/magma /build/python/gen /build/python/gen
