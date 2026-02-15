# Official Node.js runtime as a parent image
FROM node:22-alpine

# working directory in the container
WORKDIR /app

ARG PROJECT_ID
ENV VITE_PROJECT_ID=$PROJECT_ID

# Copy package.json and package-lock.json to the working directory
COPY package*.json index.html *.js ./

# Install packages
RUN npm install

# Bundle app source
COPY . .

# Build the Vue.js application
RUN npm run build

# Make port 8080 available
EXPOSE 8080

ENV OTEL_SERVICE_NAME=auto-claims-frontend
ENV OTEL_TRACES_EXPORTER=otlp
ENV OTEL_EXPORTER_OTLP_ENDPOINT=https://telemetry.googleapis.com
ENV VITE_OTLP_EXPORTER_URL=https://telemetry.googleapis.com

# Run the app when the container launches
CMD ["npm", "start"]