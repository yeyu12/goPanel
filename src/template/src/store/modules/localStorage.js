const state = {
    computerData: {}
};

const actions = {};

const mutations = {
    init(state) {
        let computer = window.localStorage.getItem('panel-computer');
        computer && (state.computerData = JSON.parse(computer));
    },
    pushComputerData(state, data) {
        let key = data.host + ':' + data.port.toString();
        state.computerData[key] = data.passwd
        window.localStorage.setItem('panel-computer', JSON.stringify(state.computerData))
    },
    // 删除本地远程主机缓存的密码相关数据
    delComputer(state, data) {
        let del = data['host'] + ":" + data['port'];
        let computer = JSON.parse(window.localStorage.getItem('panel-computer'));
        if (computer) {
            delete computer[del];
            window.localStorage.setItem('panel-computer', JSON.stringify(computer));
        }
    },
};

export default {
    namespaced: true,
    state,
    actions,
    mutations
}