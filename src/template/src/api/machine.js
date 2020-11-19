import {get} from "@/utils/request"

export function list() {
    return get("/machine/list").then((data) => {
        return data.data;
    })
}