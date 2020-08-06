FROM alpine

ENV ENV_FILE_PATH='/conf/configfile.json'

COPY ./configfile.json /conf/configfile.json

ADD ./tpayment /tpayment
WORKDIR /
ENTRYPOINT [ "/tpayment" ]
