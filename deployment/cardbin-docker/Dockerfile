FROM alpine

ENV ENV_FILE_PATH='/conf/configfile.json'

COPY ./configfile.json /conf/configfile.json

ADD ./cardbin /cardbin
WORKDIR /
ENTRYPOINT [ "/cardbin" ]
