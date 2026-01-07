// Configuration for your app
// https://v2.quasar.dev/quasar-cli-vite/quasar-config-file

import { defineConfig } from '#q-app/wrappers'

export default defineConfig((/* ctx */) => {
  return {

    boot: [
      'api',
      'pinia',
      'lang'
    ],

    css: ['app.scss'],

    extras: [
      'roboto-font',
      'material-icons',
    ],

    build: {
      target: {
        browser: ['es2022', 'firefox115', 'chrome115', 'safari14'],
        node: 'node20',
      },

      vueRouterMode: 'hash',

      vitePlugins: [
        [
          'vite-plugin-checker',
          {
            eslint: {
              lintCommand: 'eslint -c ./eslint.config.js "./src*/**/*.{js,mjs,cjs,vue}"',
              useFlatConfig: true,
            },
          },
          { server: false },
        ],
      ],
    },

    devServer: {
      open: true,
    },

    framework: {
      lang: 'pt-BR',
      config: {
        brand: {
          primary: '#610659',
          secondary: '#3a054e',
          accent: '#FFD200',
          positive: '#21BA45',
          negative: '#C10015',
          info: '#1C77C3',
          warning: '#F2C037'
        },
        notify: {
          position: 'top-right',
          timeout: 3000,
          textColor: 'white'
        },
      },
      plugins: ['Notify', 'Dialog']
    },

    animations: [],

    ssr: {
      prodPort: 3000,
      middlewares: ['render'],
      pwa: false,
    },

    // ✅ PWA CONFIG — NOME + ÍCONES + COR DO APP
    pwa: {
      workboxMode: 'GenerateSW',

      manifest: {
        name: 'Midia IEQ Maria Goretti',
        short_name: 'Midia IEQ',
        description: 'Plataforma de mídia da IEQ Maria Goretti',
        display: 'standalone',
        orientation: 'portrait',
        background_color: '#ffffff',
        theme_color: '#610659',
        start_url: '.',
        scope: '.',
        icons: [
          { src: 'icons/icon-128x128.png', sizes: '128x128', type: 'image/png' },
          { src: 'icons/icon-192x192.png', sizes: '192x192', type: 'image/png' },
          { src: 'icons/icon-256x256.png', sizes: '256x256', type: 'image/png' },
          { src: 'icons/icon-384x384.png', sizes: '384x384', type: 'image/png' },
          { src: 'icons/icon-512x512.png', sizes: '512x512', type: 'image/png' }
        ]
      }
    },

    cordova: {},

    capacitor: {
      hideSplashscreen: true,
    },

    electron: {
      preloadScripts: ['electron-preload'],
      inspectPort: 5858,
      bundler: 'packager',
      builder: {
        appId: 'midia.ieq.mariagoretti',
      },
    },

    bex: {
      extraScripts: [],
    },
  }
})
