import Vue from 'vue'
import Router from 'vue-router'

import Index from '@/view/index/index'
import Login from '@/view/login/index'

Vue.use(Router);

export default new Router({
    routes: [{
        path: '/',
        name: 'index',
        component: Index,
    }, {
        path: '/login',
        name: 'login',
        component: Login,
    }]
});