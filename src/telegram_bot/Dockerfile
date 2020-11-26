FROM golang

RUN apt-get -y update
RUN apt-get install -y tree wget curl

RUN curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh

RUN mkdir go && mkdir go/src && mkdir go/bin && mkdir go/pkg

ENV GOPATH $HOME/go
ENV PATH $GOPATH/bin:/usr/local/go/bin:$PATH

USER root

WORKDIR $GOPATH/src/github.com/bmstu-iu8-g4-2020-project/todo_web_service/src/telegram_bot
ADD ./ $GOPATH/src/github.com/bmstu-iu8-g4-2020-project/todo_web_service/src/telegram_bot

RUN tree -L 4 ./

RUN chmod +x ./scripts/*
RUN ./scripts/build.sh

CMD ["./server.app"]
