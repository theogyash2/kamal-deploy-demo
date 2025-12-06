FROM node:18

WORKDIR /app_b
RUN echo "Current Build Context Root:" && pwd

COPY package.json package-lock.json* ./
RUN npm install

COPY . .

CMD ["node", "app.js"]
