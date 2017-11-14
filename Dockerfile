FROM golang:1.9.2

COPY . .

RUN make

ENTRYPOINT [ "./bin/pressure" ]

CMD []

