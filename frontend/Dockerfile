# build image
FROM node:10.3 as node

WORKDIR /app
COPY package.json /app/
COPY package-lock.json /app/
RUN npm ci --quiet

COPY ./ /app/
RUN npm run build -- --progress=false

# production image
FROM nginx:1.13
COPY --from=node /app/dist/my-app /usr/share/nginx/html
COPY ./nginx-custom.conf /etc/nginx/conf.d/default.conf
