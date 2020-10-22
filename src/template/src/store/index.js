import Vue from 'vue'
import Vuex from 'vuex'
import TopMenu from './modules/topMenu'

Vue.use(Vuex);

export default new Vuex.Store({
    modules: {
        TopMenu
    }
});
