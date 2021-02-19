FROM mcr.microsoft.com/mssql/server:2017-CU17-ubuntu

RUN mkdir -p /usr/src/app
WORKDIR /usr/src/app
COPY . /usr/src/app

RUN chmod +x /usr/src/app/init.sh

ENV MSSQL_PID Express

CMD /bin/bash ./entrypoint.sh
