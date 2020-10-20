import {post} from "@/utils/request"

export function login(params) {
    return post("/login", params).then((data) => {
        return data.data;
    })
}