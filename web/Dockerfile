FROM node:22-alpine

WORKDIR /app

# Copy lock files and install dependencies
COPY package*.json ./
RUN npm install

# Copy the rest of your application code
COPY . .

# Vite development port
EXPOSE 5173

CMD ["npm", "run", "dev"]