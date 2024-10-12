import { defineConfig } from '@hey-api/openapi-ts'

export default defineConfig({
  client: '@hey-api/client-fetch',
  input: '../../backend/go/docs/swagger.json',
  output: {
    path: './src/api',
    lint: 'eslint',
    format: 'prettier'
  },
  types: {
    enums: 'typescript'
  },
  plugins: ['@tanstack/vue-query']
})