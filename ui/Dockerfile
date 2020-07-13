FROM node:14.5 AS staging

COPY ./package.json /srv/stable-coin-ui/package.json
COPY ./package-lock.json /srv/stable-coin-ui/package-lock.json

RUN cd /srv/stable-coin-ui && npm ci

FROM node:14.5 AS build

COPY ./ /srv/stable-coin-ui
COPY --from=staging /srv/stable-coin-ui/node_modules /srv/stable-coin-ui/node_modules

RUN cd /srv/stable-coin-ui && npm run build

FROM nginx:alpine

COPY ./nginx.conf /etc/nginx/conf.d/default.conf
RUN rm -rf /usr/share/nginx/html/*
COPY --from=build /srv/stable-coin-ui/dist /usr/share/nginx/html

CMD ["nginx", "-g", "daemon off;"]
