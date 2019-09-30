FROM debian:stable AS validator

RUN apt-get -q -y update && apt-get -q -y install \
  clang \
  cmake \
  libpcre3-dev \
  git

RUN git clone https://github.com/CESNET/libyang
RUN mkdir -p libyang/build
WORKDIR /libyang/build
RUN cmake -DCMAKE_INSTALL_PREFIX:PATH=/usr ..
RUN make
RUN make install

WORKDIR /
RUN mkdir -p /usr/share/yang/models

RUN git clone https://github.com/YangModels/yang/
RUN find yang/standard -name \*.yang -exec cp {} /usr/share/yang/models \;
RUN cp yang/experimental/ietf-extracted-YANG-modules/ietf-geo-location@2019-02-17.yang \
       /usr/share/yang/models/ietf-geo-location.yang

RUN git clone https://github.com/openconfig/public
COPY validator/standard_compliance.sh /
RUN find public -name \*.yang -exec /standard_compliance.sh {} \;
RUN find public -name \*.yang -exec cp {} /usr/share/yang/models \;

COPY yang/fbc-symphony-device.yang /usr/share/yang/models

COPY validator/validate_json.sh /bin/
COPY validator/validate_yang.sh /bin/

CMD validate_json.sh
