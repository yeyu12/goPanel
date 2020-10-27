import {post} from "@/utils/request"

export function addOne(params) {
    return post("/command/addOne", params).then((data) => {
        return data.data;
    })
}