import {get} from "@/utils/request"
import {post} from "../utils/request";

export function list() {
    return get("/machine/list").then((data) => {
        return data.data;
    })
}

export function save(data) {
    return post("/machine/save", data).then((ret) => {
        return ret.data
    })
}

export function reboot(id) {
    return get(`/machine/reboot?id=${id}`).then((ret) => {
        return ret.data
    })
}

export function restartService(id) {
    return get(`/machine/restartService?id=${id}`).then((ret) => {
        return ret.data
    })
}