FROM node:10.5.0-stretch

RUN npm install -g serverless --unsafe-perm=true

WORKDIR /app