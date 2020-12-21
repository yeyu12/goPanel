import axios from 'axios';
import router from '@/router/index'

let axiosObj = axios.create();
let loginStatus = [3000, 4003, 4004];

// 添加请求拦截器
axiosObj.interceptors.request.use(
    config => {
        config.url = "http://" + window.location.hostname + ":10000" + config.url;
        config.headers = {
            ...config.headers,
            'Account-Token': window.sessionStorage.getItem('panel-token') || '',
        };

        return config;
    },
    error => {
        return Promise.reject(error);
    }
);

// 添加响应拦截器
// 统一在window unhandledrejection事件处理未捕获的promise事件
axiosObj.interceptors.response.use(
    response => {
        if (loginStatus.indexOf(response.data.code) > -1) {
            window.sessionStorage.clear()
            router.app.$router.push('/login')
        }

        return Promise.resolve(response);
    },
    error => {
        // 对响应错误做点什么
        return Promise.reject({
            ...error,
        });
    }
);

export const get = (url, val, config = {}) => {
    return axiosObj.get(url, {
        params: val,
        ...config,
        data: {customParams: config.customParams},
    });
};
// 删除公用方法
export const del = (url, data, config = {}) => {
    return axiosObj.delete(url, {
        data: config.customParams
            ? {...data, customParams: config.customParams}
            : data,
        ...config,
    });
};

export const post = (url, val, config = {}) => {
    let contentType;

    return axiosObj.request({
        url,
        data: config.customParams
            ? {...val, customParams: config.customParams}
            : val,
        method: 'post',
        headers: {
            'Content-type': contentType,
        },
        ...config,
    });
};

// 修改数据公用方法
export const put = (url, val, config = {}) => {
    let contentType;

    return axiosObj.request({
        url,
        data: val,
        method: 'put',
        headers: {
            'Content-type': contentType,
        },
        ...config,
    });
};

// formadata post 提交数据
export const postFormData = (url, val, config = {}) => {
    return axiosObj.request({
        url,
        data: val,
        method: 'post',
        headers: {
            'Content-type': 'application/x-www-form-urlencoded',
        },
        ...config,
    });
};

export default axiosObj;
