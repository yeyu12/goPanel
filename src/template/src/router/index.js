import Vue from 'vue';
import Router from 'vue-router';

import Index from '@/view/index/index';
import Login from '@/view/login/index';
// import Shell from '@/view/shell/index';

Vue.use(Router);

export default new Router({
    routes: [{
        path: '/',
        name: 'index',
        component: Index,
        children: [
            // {
            //     path: 'shell',
            //     name: 'Shell',
            //     component: Shell
            // }
        ]
    }, {
        path: '/login',
        name: 'login',
        component: Login,
    }]
});