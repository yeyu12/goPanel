import {post} from "@/utils/request"

export function login(params) {
    return post("/login", params).then((data) => {
        return data.data;
    })
}

export function register(params) {
    return post("/register", params).then((data) => {
        return data.data;
    })
}