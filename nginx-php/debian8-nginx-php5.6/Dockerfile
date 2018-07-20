FROM webdevops/php-nginx:debian-8

LABEL maintainer=slpcat@qq.com

#Container variables
ENV \
    TERM="xterm" \
    LANG="en_US.UTF-8" \
    LC_ALL="en_US.UTF-8" \
    TIMEZONE="Asia/Shanghai" \
    APP_ENV="dev" \
    APPLICATION_USER="www-data" \
    APPLICATION_GROUP="www-data" \
    APPLICATION_PATH="/data/www"

COPY etc /etc/
COPY www_root ${APPLICATION_PATH}/

# System update to latest
# Set timezone
RUN set -x \
    && unlink /etc/nginx/conf.d/10-docker.conf \
    && apt-get update \
    && apt-get dist-upgrade -y \
    && ln -sf /usr/share/zoneinfo/${TIMEZONE} /etc/localtime \
    && echo "${TIMEZONE}" > /etc/timezone \
    && apt-get install -y \
        nginx \
        ca-certificates \
        php5-xdebug \
        php5-dev \
        autoconf \
        make \
        gcc \
    && ln -sf /usr/bin/php5 /usr/bin/php \
    && pecl install yaf-2.3.5 \
    && echo "extension=yaf.so" > /etc/php5/mods-available/yaf.ini \
    && ln -sf /etc/php5/mods-available/yaf.ini /etc/php5/fpm/conf.d/20-yaf.ini \
    && ln -sf /etc/php5/mods-available/yaf.ini /etc/php5/cli/conf.d/20-yaf.ini \
    && docker-image-cleanup
