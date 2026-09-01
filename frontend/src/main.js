import { createApp } from 'vue'
import App from './App.vue'
import './styles/global.css';
import './styles/variables.css'
import { router } from './router/router.js';

createApp(App)
    .use(router)
    .mount('#app')