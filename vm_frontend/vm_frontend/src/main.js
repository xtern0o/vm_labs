import { createApp } from 'vue'
import App from './App.vue'
import router from './router'
import Vue3Katex from 'vue3-katex'
import 'katex/dist/katex.min.css'
import './styles/style.scss'

createApp(App).use(router).use(Vue3Katex).mount('#app')
