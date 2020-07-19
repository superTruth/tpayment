FROM alpine

ENV TPAYMENT_DB_CONFIG='root:123456@tcp(120.78.138.140)/tpayment?charset=utf8&parseTime=True&loc=Local'
ENV ENV_FILE_PATH='/conf/configfile.json'

COPY ./configfile.json /conf/configfile.json

ADD ./tpayment /tpayment
WORKDIR /
ENTRYPOINT [ "/tpayment" ]
