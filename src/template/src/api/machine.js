import {post} from "@/utils/request"

export function add(params) {
    return post("/machine/add", params).then((data) => {
        return data.data;
    })
}