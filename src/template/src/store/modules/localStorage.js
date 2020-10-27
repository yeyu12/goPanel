const state = {};

const actions = {};

const mutations = {
    // 删除本地远程主机缓存的密码相关数据
    delComputer(state, data) {
        let del = data['host'] + ":" + data['port'];
        let computer = JSON.parse(window.localStorage.getItem('panel-computer'));
        if (computer) {
            delete computer[del];
            window.localStorage.setItem('panel-computer', JSON.stringify(computer));
        }
    }
};

export default {
    namespaced: true,
    state,
    actions,
    mutations
}