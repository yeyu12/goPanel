import axios from 'axios';

let axiosObj = axios.create();

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
