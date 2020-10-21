import {post, get} from "@/utils/request"

export function list() {
    return get("/machine/list").then((data) => {
        return data.data;
    })
}

export function add(params) {
    return post("/machine/add", params).then((data) => {
        return data.data;
    })
}