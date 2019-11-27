FROM fedora:latest AS formatter

RUN dnf update -y && dnf install -y clang findutils
WORKDIR /cache/devmand/repo/gateway
CMD find ./src \( -name "*.[hc]pp" -or -name "*.h" \) \
    -exec clang-format -i --style=file {} \;
