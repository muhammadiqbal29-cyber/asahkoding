# Stage 1: Install Dependencies dan Build
FROM node:alpine AS builder
WORKDIR /app
COPY package*.json ./
RUN npm install
COPY . .
RUN npm run build

# Stage 2: Produksi Ringan (Hanya membawa hasil build)
FROM node:alpine AS runner
WORKDIR /app
ENV NODE_ENV production

# Menyalin dari builder
COPY --from=builder /app/package.json ./package.json
COPY --from=builder /app/.next ./.next
COPY --from=builder /app/public ./public
COPY --from=builder /app/node_modules ./node_modules

EXPOSE 3000
CMD ["npm", "run", "start"]
