const state = {
    MENU_SHELL_TYPE: 'shell',
    openTagMenu: [],
    defaultTagMenu: '0',
    currentSelectMenu: {}
};

const actions = {};

const mutations = {
    init(state) {
        let menuData = JSON.parse(window.localStorage.getItem('gps-tag-menu'));
        menuData && mutations.openTagMenu(state, menuData)

        let defaultMenuIndex = window.localStorage.getItem('gps-default-tag-menu');
        defaultMenuIndex && mutations.upDefaultTagMenu(state, defaultMenuIndex)
    },
    openTagMenuPush(state, data) {
        state.openTagMenu.push(data);
        mutations.upDefaultTagMenu(state, (state.openTagMenu.length - 1).toString());
        window.localStorage.setItem('gps-tag-menu', JSON.stringify(state.openTagMenu));
    },
    openTagMenuDel(state, data) {
        let defaultTagMenuInt = parseInt(state.defaultTagMenu);
        let del = data['host'] + ":" + data['port'];
        let delIndexObj = {};
        let delIndexArr = [];

        // 找出要删除的位置
        for (let i = 0; i < state.openTagMenu.length; i++) {
            let tag = state.openTagMenu[i]['host'] + ":" + state.openTagMenu[i]['port'];

            if (tag === del) {
                delIndexObj[i] = true;
                delIndexArr.push(i);
            }
        }

        // 开始执行删除  没有要删除的tag标签，直接执行删除
        if (!delIndexArr.length) {
            return;
        }
        if (delIndexArr[0] === defaultTagMenuInt) {
            defaultTagMenuInt -= 1;
        } else if (delIndexArr[0] < defaultTagMenuInt) {
            let tmp = 0;
            for (let i in delIndexObj) {
                if (i <= defaultTagMenuInt) {
                    tmp++;
                }
            }

            defaultTagMenuInt -= tmp;
        } else {
            // 小于的情况不需要处理
            return;
        }

        state.openTagMenu = state.openTagMenu.filter((item, index) => {
            if (!delIndexObj[index]) {
                return item;
            }
        });

        mutations.upDefaultTagMenu(state, defaultTagMenuInt.toString());
        window.localStorage.setItem('gps-tag-menu', JSON.stringify(state.openTagMenu));
    },
    openTagMenuClear(state) {
        state.openTagMenu = [];
        window.localStorage.removeItem('gps-tag-menu');
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

        mutations.upDefaultTagMenu(state, intDefaultTagMenu.toString());
        window.localStorage.setItem('gps-tag-menu', JSON.stringify(state.openTagMenu));
    },
    upDefaultTagMenu(state, index) {
        state.defaultTagMenu = index;
        window.localStorage.setItem('gps-default-tag-menu', index);
    },
    currentSelectMenuEdit(state, data) {
        state.currentSelectMenu = data
    }
};

export default {
    namespaced: true,
    state,
    actions,
    mutations
}