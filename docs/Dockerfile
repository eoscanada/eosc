FROM jojomi/hugo:0.48 as hugo-builder


WORKDIR /src
COPY ./docs/ /src
RUN mkdir /src/public
RUN hugo -s /src -d /src/public -b https://eosc.eoscanada.com/
RUN ls -l /src



FROM nginx:stable-alpine-perl

WORKDIR /usr/share/nginx/html
COPY --from=hugo-builder /src/public /usr/share/nginx/html
