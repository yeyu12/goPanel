<template>
    <div class="container" style="padding: 0;margin: 0">
        <div id="terminal-container" style="height: calc(100vh - 40px);"></div>
    </div>
</template>

<script>
    import "xterm/css/xterm.css"
    import {Terminal} from "xterm";
    import {FitAddon} from 'xterm-addon-fit'
    import {dom, library} from '@fortawesome/fontawesome-svg-core'
    import {faBars, faClipboard, faCog, faDownload, faKey} from '@fortawesome/free-solid-svg-icons'

    export default {
        data() {
            return {
                copy: "",
                isBackspace: false,
                lock: true
            };
        },
        created() {
            library.add(faBars, faClipboard, faDownload, faKey, faCog);
            dom.watch();
        },
        /*eslint no-unused-vars: ["error", { "args": "none" }]*/
        mounted() {
            let fitAddon = new FitAddon();
            let terminalContainer = document.getElementById("terminal-container");
            let term = new Terminal({
                // 光标闪烁
                cursorBlink: true,
                rendererType: "canvas", //渲染类型
                disableStdin: false, //是否应禁用输入。
                cursorStyle: "block", //光标样式
                tabStopWidth: 4
            });
            term.loadAddon(fitAddon);
            term.open(terminalContainer, true);
            term.focus();
            fitAddon.fit();

            let websocket = new WebSocket("ws://127.0.0.1:10010/ws/ssh/" + term.cols + "/" + term.rows + "/127.0.0.1"); //地址
            websocket.binaryType = "arraybuffer";
            //连接成功
            websocket.onopen = (evt) => {
                term.writeln("");
            }

            // 输入
            term.onData(data => {
                websocket.send(data);
            });

            // 返回
            websocket.onmessage = function (evt) {
                term.write(evt.data)
            };

            //关闭
            websocket.onclose = function (evt) {
                console.log("close", evt);
            };
            //错误
            websocket.onerror = function (evt) {
                console.log("error", evt);
            };
        }
        ,
        methods: {
            formatWs(event, data) {
                return JSON.stringify({
                    event,
                    data: new TextEncoder().encode(data)
                })
            }
        }
    }
    ;
</script>
