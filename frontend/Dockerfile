FROM ubuntu:latest AS build-env

# Install required dependencies
RUN apt-get update && \
    apt-get install -y \
    curl \
    git \
    wget \
    unzip \
    xz-utils \
    libglu1-mesa \
    && rm -rf /var/lib/apt/lists/*

# Set up Flutter
ENV FLUTTER_HOME=/flutter
ENV PATH=$FLUTTER_HOME/bin:$PATH
ENV FLUTTER_WEB_PORT=3000

# Clone Flutter and switch to stable
RUN git clone https://github.com/flutter/flutter.git $FLUTTER_HOME && \
    cd $FLUTTER_HOME && \
    git checkout stable

# Run basic check and enable web
RUN flutter doctor -v && \
    flutter config --enable-web

# Set up project
WORKDIR /app

# Copy pubspec files
COPY pubspec.* ./

# Get packages
RUN flutter pub get

# Copy the rest of the project files
COPY . .

# Clean and get packages again to ensure everything is set up
RUN flutter clean && \
    flutter pub get && \
    flutter config --enable-web

# Build web with specific settings
RUN flutter build web \
    --release \
    --web-renderer html \
    --dart-define=FLUTTER_WEB_USE_SKIA=true

# Stage 2 - Create the run-time image
FROM nginx:stable-alpine

# Copy the build output
COPY --from=build-env /app/build/web /usr/share/nginx/html

# Copy nginx config
COPY nginx.conf /etc/nginx/conf.d/default.conf

# Expose port
EXPOSE 3000 