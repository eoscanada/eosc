FROM jojomi/hugo:latest as hugo-builder

# RUN mkdir -p /src

WORKDIR /src
COPY . /src
RUN ls -l
RUN hugo -s /src



FROM nginx:stable-alpine-perl

WORKDIR /usr/share/nginx/html
COPY --from=hugo-builder /src/public /usr/share/nginx/html
