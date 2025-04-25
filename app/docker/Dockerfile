FROM docker.io/python:3.13-bookworm

RUN pip install supervisor

#### This is for the workload, this is a bit hacky and not how it should be done
####
COPY ./ /data/www/apps

RUN cd /data/www/apps && \
    pip install -r requirements.txt
####
####

ENTRYPOINT [ "/usr/local/bin/supervisord" ]
