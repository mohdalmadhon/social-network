import { fileURLToPath, URL } from 'node:url'

import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import vueDevTools from 'vite-plugin-vue-devtools'

const page = (path) => fileURLToPath(new URL(path, import.meta.url))

const pageRoutes = {
  '/login': '/login.html',
  '/home-feed': '/home-feed.html',
  '/profile': '/profile.html',
  '/groups': '/groups.html',
  '/chats': '/chats.html',
  '/notifications': '/notifications.html',
}

function cleanPageRoutes() {
  const routePages = (request, response, next) => {
    const requestUrl = new URL(request.url, 'http://localhost')

    if (requestUrl.pathname === '/') {
      response.writeHead(302, { Location: '/login' })
      response.end()
      return
    }

    const pageFile = pageRoutes[requestUrl.pathname]
    if (pageFile) {
      request.url = pageFile + requestUrl.search
    }

    next()
  }

  return {
    name: 'orbit-clean-page-routes',
    configureServer(server) {
      server.middlewares.use(routePages)
    },
    configurePreviewServer(server) {
      server.middlewares.use(routePages)
    },
  }
}

// https://vite.dev/config/
export default defineConfig({
  plugins: [
    vue(),
    vueDevTools(),
    cleanPageRoutes(),
  ],
  resolve: {
    alias: {
      '@': fileURLToPath(new URL('./src', import.meta.url)),
    },
  },
  build: {
    rollupOptions: {
      input: {
        login: page('./login.html'),
        homeFeed: page('./home-feed.html'),
        profile: page('./profile.html'),
        groups: page('./groups.html'),
        chats: page('./chats.html'),
        notifications: page('./notifications.html'),
      },
    },
  },
  server: {
    proxy: {
      '/api': {
        target: 'http://localhost:4033',
        changeOrigin: true,
      },
    },
  },
})
