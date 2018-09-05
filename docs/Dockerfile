FROM jojomi/hugo:latest as hugo-builder

# RUN mkdir -p /src

WORKDIR /src
COPY ./docs/ /src
RUN hugo -s /src -d =/public --debug
RUN ls -l /src



FROM nginx:stable-alpine-perl

WORKDIR /usr/share/nginx/html
COPY --from=hugo-builder =/public /usr/share/nginx/html
