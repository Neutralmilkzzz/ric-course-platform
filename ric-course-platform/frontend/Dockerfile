# Simple dev Dockerfile for CRA
FROM node:20
WORKDIR /app
COPY package.json package-lock.json* ./
RUN npm install || true
COPY . .
EXPOSE 3000
CMD ["npm", "start"]
