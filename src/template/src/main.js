import Vue from 'vue';
import App from './App.vue';
import router from './router';
import store from './store';
import ElementUI from 'element-ui';
import Md5 from 'js-md5';
// import Sha1 from 'js-sha1';
import {Base64} from 'js-base64';

import "./static/icon/iconfont";

import 'element-ui/lib/theme-chalk/index.css';
import './static/css/style.css';

Vue.use(ElementUI);

Vue.config.productionTip = false;

Vue.prototype.$md5 = Md5;
// Vue.prototype.$sha1 = Sha1;
Vue.prototype.$base64 = Base64;

new Vue({
    router,
    store,
    render: h => h(App),
}).$mount('#app');
