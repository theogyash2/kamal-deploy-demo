FROM node:18

WORKDIR /app

COPY package.json package-lock.json* ./
RUN npm install

COPY . .

# MAIN decides what runs: web OR worker
CMD ["sh", "-c", "node ${MAIN}.js"]
