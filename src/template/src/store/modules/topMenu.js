const state = {
    MENU_SHELL_TYPE: 'shell',
    openTagMenu: [],
    defaultTagMenu: '0'
};

const actions = {};

const mutations = {
    openTagMenuPush(state, data) {
        state.openTagMenu.push(data);
        mutations.upDefaultTagMenu(state, (state.openTagMenu.length - 1).toString())
        window.localStorage.setItem('panel-tag-menu', JSON.stringify(state.openTagMenu));
    },
    openTagMenuClear(state) {
        state.openTagMenu = [];
        window.localStorage.removeItem('panel-tag-menu');
    },
    openTagMenu(state, data) {
        state.openTagMenu = data;
    },
    removeTagMenu(state, index) {
        state.openTagMenu.splice(parseInt(index), 1);
        let intDefaultTagMenu = parseInt(state.defaultTagMenu);
        intDefaultTagMenu -= 1;
        if (intDefaultTagMenu < 0) {
            intDefaultTagMenu = 0;
        }

        mutations.upDefaultTagMenu(state, intDefaultTagMenu.toString())
        window.localStorage.setItem('panel-tag-menu', JSON.stringify(state.openTagMenu));
    },
    upDefaultTagMenu(state, index) {
        state.defaultTagMenu = index;
        window.localStorage.setItem('panel-default-tag-menu', index);
    }
};

export default {
    namespaced: true,
    state,
    actions,
    mutations
}