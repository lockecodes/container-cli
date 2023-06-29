version: '3.9'
services:
  __NAME__:
    build:
      context: ./src
      dockerfile: Dockerfile
    environment:
      - CONTEXT_DIR
      - VERSION
      - IN_DOCKER
      - DISPLAY
    volumes:
      - $USER_HOME:/opt/usr/home
      - $CONTEXT_DIR:/opt/context
