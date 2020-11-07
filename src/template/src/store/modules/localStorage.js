const state = {
    computerPasswdData: {}
};

const actions = {};

const mutations = {
    init(state) {
        let computer = window.localStorage.getItem('gps-computer');
        computer && (state.computerPasswdData = JSON.parse(computer));
    },
    pushComputerPasswdData(state, data) {
        let key = data.host + ':' + data.port.toString();
        state.computerPasswdData[key] = data.passwd
        window.localStorage.setItem('gps-computer', JSON.stringify(state.computerPasswdData))
    },
    // 删除本地远程主机缓存的密码相关数据
    delComputerPasswd(state, data) {
        let del = data['host'] + ":" + data['port'];
        let computer = JSON.parse(window.localStorage.getItem('gps-computer'));
        if (computer) {
            delete computer[del];
            window.localStorage.setItem('gps-computer', JSON.stringify(computer));
        }
    },
};

export default {
    namespaced: true,
    state,
    actions,
    mutations
}