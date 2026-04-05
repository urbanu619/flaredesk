# 构建阶段
ARG s3Image=654654279104.dkr.ecr.ap-southeast-1.amazonaws.com/jcwork-test:node-18.20-slim-aws-git-lokalise
FROM ${s3Image} AS build

WORKDIR /web
COPY . .

ARG DOCKER_ENV

ENV DEPLOY_ENV=${DOCKER_ENV}

RUN npm install

# 如果是测试环境，使用测试配置构建，否则使用正式配置构建
RUN if [ "$DEPLOY_ENV" = "coinvex-eco-test" ]; then npm run build-test; else npm run build; fi

# 最终阶段
FROM nginx:latest

ENV TZ=UTC

# 拷贝构建产物到 nginx 默认目录
COPY --from=build /web/dist /usr/share/nginx/html

# 可选：复制自定义配置
COPY nginx.conf /etc/nginx/nginx.conf

EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]