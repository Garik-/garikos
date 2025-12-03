import './assets/main.css'

import { createApp } from 'vue'
import App from './App.vue'

import * as constants from '@/config/constants'

console.log('App constants:', constants)

const app = createApp(App)
app.mount('#app')
