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
                ws: {},
                term: {},
                timeout: {},
                fitAddon: {},
                termDispose: {},
                wsTimer: 0
            };
        },
        created() {
            library.add(faBars, faClipboard, faDownload, faKey, faCog);
            dom.watch();
        },
        /*eslint no-unused-vars: ["error", { "args": "none" }]*/
        mounted() {
            this.fitAddon = new FitAddon();
            let terminalContainer = document.getElementById("terminal-container");
            this.term = new Terminal({
                cursorBlink: true,
                rendererType: "canvas", //渲染类型
                disableStdin: false, //是否应禁用输入。
                cursorStyle: "block", //光标样式
                tabStopWidth: 4
            });
            this.term.loadAddon(this.fitAddon);
            this.term.open(terminalContainer, true);
            this.term.focus();
            this.fitAddon.fit();

            this.connWebsocket();
        },
        methods: {
            formatWs(event, data) {
                return JSON.stringify({
                    event,
                    data: new TextEncoder().encode(data)
                })
            },
            connWebsocket() {
                this.ws = new WebSocket("ws://127.0.0.1:10010/ws/ssh/" + this.term.cols + "/" + this.term.rows + "/127.0.0.1"); //地址
                this.ws.binaryType = "arraybuffer";
                //连接成功
                this.ws.onopen = (evt) => {
                    if (this.wsTimer) {
                        clearInterval(this.wsTimer)
                        this.wsTimer = 0;
                    }

                    this.term.writeln("");
                }

                // 输入
                this.termDispose = this.term.onData(data => {
                    this.ws.send(data);
                });

                // 返回
                this.ws.onmessage = (evt) => {
                    this.term.write(evt.data)
                };

                //关闭
                this.ws.onclose = (evt) => {
                    this.reconnect()
                };
                //错误
                this.ws.onerror = (evt) => {
                    this.reconnect();
                };
            },
            reconnect() {
                if (!this.wsTimer) {
                    this.wsTimer = setInterval(() => {
                        this.termDispose.dispose()
                        this.connWebsocket();
                        this.term.reset();
                    }, 10000);
                }
            }
        },
        destroyed() {
            clearInterval(this.wsTimer);
            this.ws.close();
        }
    };
</script>
