import {post, get} from "@/utils/request"

export function list() {
    return get("/machine/list").then((data) => {
        return data.data;
    })
}

export function save(params) {
    return post("/machine/save", params).then((data) => {
        return data.data;
    })
}