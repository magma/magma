set -x
# build container

docker build . -f services/build/Dockerfile


# # mobilityd
# docker build . -f services/mobilityd/Dockerfile -t mobilityd
# # enodedb
# docker build . -f services/enodedb/Dockerfile -t enodedb
# # health
# docker build . -f services/health/Dockerfile -t health
# # monitord
# docker build . -f services/monitord/Dockerfile -t monitord
# # pipelined
# docker build . -f services/pipelined/Dockerfile -t pipelined
# # pkt_tester
# docker build . -f services/pkt_tester/Dockerfile -t pkt_tester
# # policydb
# docker build . -f services/policydb/Dockerfile -t policydb
# # redirectd
# docker build . -f services/redirectd/Dockerfile -t redirectd
# # smsd
# docker build . -f services/smsd/Dockerfile -t smsd
# # subscriberd
# docker build . -f services/subscriberd/Dockerfile -t subscriberd
# # tests
# docker build . -f services/tests/Dockerfile -t tests
