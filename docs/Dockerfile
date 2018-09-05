FROM jojomi/hugo:latest as hugo-builder

# RUN mkdir -p /src

WORKDIR /src
COPY ./docs/ /src
RUN ls -l
RUN hugo -s /src -d /src/public
RUN ls -l /src
RUN ls -l /



FROM nginx:stable-alpine-perl

WORKDIR /usr/share/nginx/html
COPY --from=hugo-builder /src/public /usr/share/nginx/html
