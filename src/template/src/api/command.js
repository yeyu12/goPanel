import {post} from "@/utils/request"

export function addCommand(params) {
    return post("/command/add", params).then((data) => {
        return data.data;
    })
}