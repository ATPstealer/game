import { defineConfig } from '@kubb/core'
import { pluginOas } from '@kubb/plugin-oas'
import { pluginTs } from '@kubb/plugin-ts'
import { pluginVueQuery } from '@kubb/plugin-vue-query'

const baseURL = 'http://staging.game.k8s.atpstealer.com/api/v2'

export default defineConfig(() => {
  return [
    {
      root: '.',
      input: {
        path: '../../backend/go/docs/swagger.yaml'
      },
      output: {
        path: './src/gen'
      },
      plugins: [
        pluginOas(),
        pluginTs({
          syntaxType: 'interface',
          unknownType: 'any'
        }),
        pluginVueQuery({
          output: {
            path: './hooks'
          },
          client: {
            baseURL,
            importPath: '@/api/customClientAxios'
          },
          query: {
            importPath: '@tanstack/vue-query'
          },
          mutation: {
            methods: [ 'post', 'put', 'delete' ]
          },
          infinite: {
            queryParam: 'next_page',
            initialPageParam: 0,
            cursorParam: 'nextCursor'
          }
        })
      ]
    }
  ]
})