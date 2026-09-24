import { createApp } from 'vue'
import App from './App.vue'
import './styles/global.css';
import './styles/variables.css'
import { router } from './router/router.js';
import { initTheme } from './helpers/theme.js';

initTheme()

createApp(App)
    .use(router)
    .mount('#app')