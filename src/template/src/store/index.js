import Vue from 'vue';
import Vuex from 'vuex';
import TopMenu from './modules/topMenu';
import LocalStorage from './modules/localStorage';

Vue.use(Vuex);

export default new Vuex.Store({
    modules: {
        TopMenu,
        LocalStorage,
    }
});
