import request from "@/utils/request"

export function login(params) {
    return request.post("http://localhost:10010/login", params).then((data) => {
        return data.data;
    })
}