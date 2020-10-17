import Vue from 'vue'
import App from './App.vue'
import router from './router'
import ElementUI from 'element-ui';
import md5 from 'js-md5';

import 'element-ui/lib/theme-chalk/index.css';
import './static/css/style.css';

Vue.use(ElementUI);

Vue.config.productionTip = false;

Vue.prototype.$md5 = md5;

new Vue({
    router,
    render: h => h(App),
}).$mount('#app');
