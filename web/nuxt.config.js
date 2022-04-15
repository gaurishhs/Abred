export default {
  head: {
    title: 'Abred',
    htmlAttrs: {
      lang: 'en',
    },
    meta: [
      { charset: 'utf-8' },
      { name: 'viewport', content: 'width=device-width, initial-scale=1' },
      { hid: 'description', name: 'description', content: 'Abred the only bot you shall need for managing temporary voice channels!' },
      { name: 'format-detection', content: 'telephone=no' },
    ],
    link: [
      { rel: 'icon', type: 'image/png', href: '/assets/transparent-logo.png' },
      {
        rel: "stylesheet",
        href: 'https://fonts.googleapis.com/icon?family=Material+Icons'
      }],
  },
  css: [
    "./styles/globals.scss"
  ],
  plugins: [],
  components: true,
  buildModules: [
    '@nuxt/typescript-build',
    '@nuxtjs/tailwindcss',
  ],
  modules: [
    '@nuxtjs/axios',
    '@nuxtjs/auth-next',
    '@nuxtjs/toast',
    '@nuxtjs/dotenv',
  ],
  auth: {
    strategies: {
      discord: {
        clientId: process.env.CLIENT_ID,
        clientSecret: process.env.CLIENT_SECRET,
      }
    },
    redirect: {
      callback: '/',
    }
  },
  axios: {
    baseURL: 'https://api.abred.bar',
  },
  build: {},
}
