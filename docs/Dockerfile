FROM jojomi/hugo:latest as hugo-builder

# RUN mkdir -p /src

WORKDIR /src
COPY ./docs/ /src
RUN mkdir /src/public
RUN hugo -s /src -d /src/public
RUN ls -l /src



FROM nginx:stable-alpine-perl

WORKDIR /usr/share/nginx/html
RUN ls -l /src
RUN ls -l /app
RUN ls -l /docs/src
COPY --from=hugo-builder /src/public /usr/share/nginx/html
