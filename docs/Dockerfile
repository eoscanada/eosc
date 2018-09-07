FROM jojomi/hugo:0.48 as hugo-builder


WORKDIR /src
COPY . /src
RUN hugo config
RUN hugo -s /src -d /eosc-website



FROM nginx:stable-alpine-perl

WORKDIR /usr/share/nginx/html
COPY --from=hugo-builder /eosc-website /usr/share/nginx/html
