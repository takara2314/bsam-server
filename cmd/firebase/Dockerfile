FROM node:20.16.0-slim

WORKDIR /firebase

RUN apt-get -y update && apt-get -y install default-jdk
RUN npm install -g firebase-tools

EXPOSE 4000

CMD ["firebase", "emulators:start"]
