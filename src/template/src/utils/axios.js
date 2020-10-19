import axios from './request';

// 添加请求拦截器
axios.interceptors.request.use(
    function (config) {
        return config;
    },
    function (error) {
        requestAfter(error.config?.customParams);

        return Promise.reject(error);
    }
);

// 添加响应拦截器
// 统一在window unhandledrejection事件处理未捕获的promise事件
axios.interceptors.response.use(
    function (response) {
        return Promise.reject(response.data);
    },
    function (error) {
        // 对响应错误做点什么
        return Promise.reject({
            ...error,
        });
    }
);
