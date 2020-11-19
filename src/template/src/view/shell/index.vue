<template>
    <div class="container" style="padding: 0;margin: 0">
        <div :id="'terminal-container-'+tagIndex" style="height: calc(100vh - 40px);"></div>
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
                wsTimer: 0,
                isReconnection: true,
                passwd: '',
                url: '',
            };
        },
        props: [
            'menu',
            'tagIndex'
        ],
        created() {
            library.add(faBars, faClipboard, faDownload, faKey, faCog);
            dom.watch();

            this.url = 'ws' + '://' + window.location.hostname + ':10010';
        },
        /*eslint no-unused-vars: ["error", { "args": "none" }]*/
        mounted() {
            this.fitAddon = new FitAddon();
            let terminalContainer = document.getElementById("terminal-container-" + this.tagIndex);
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

            let computer = JSON.parse(window.localStorage.getItem('panel-computer'));
            if (computer) {
                this.passwd = computer[this.menu['host'] + ':' + this.menu['port']];
                this.connWebsocket();
            } else {
                console.log('在这里验证客户端没有密码的情况。');
            }
        },
        methods: {
            formatWs(event, data) {
                return new TextEncoder().encode(JSON.stringify({
                    event,
                    data,
                    type: 0
                }))
            },
            unFormatWs(data) {
                return JSON.parse(new TextDecoder().decode(data))
            },
            connWebsocket() {
                this.ws = new WebSocket(this.url + '/ws/ssh'); //地址
                this.ws.binaryType = "arraybuffer";
                //连接成功
                this.ws.onopen = (evt) => {
                    if (this.wsTimer) {
                        clearInterval(this.wsTimer);
                        this.wsTimer = 0;
                    }

                    this.ws.send(this.formatWs('init', {
                        token: window.sessionStorage.getItem('panel-token'),
                        cols: this.term.cols,
                        rows: this.term.rows,
                        id: this.menu.id,
                        passwd: this.passwd
                    }));

                    this.term.writeln("");
                };

                // 输入
                this.termDispose = this.term.onData(data => {
                    this.ws.send(this.formatWs('data', data));
                });

                // 返回
                this.ws.onmessage = (evt) => {
                    let wsData = this.unFormatWs(evt.data);
                    switch (wsData.event) {
                        case "data":
                            this.term.write(wsData.data);
                            break;
                        case "err":
                            this.isReconnection = false;

                            this.$notify.error({
                                title: '错误',
                                message: wsData.data,
                                duration: 10000,
                                showClose: false
                            });
                            break;
                    }
                };

                //关闭
                this.ws.onclose = (evt) => {
                    this.reconnect();
                };
                //错误
                this.ws.onerror = (evt) => {
                    this.reconnect();
                };
            },
            reconnect() {
                if (!this.wsTimer && this.isReconnection) {
                    this.wsTimer = setInterval(() => {
                        this.termDispose.dispose();
                        this.connWebsocket();
                        this.term.reset();
                    }, 10000);
                }
            }
        },
        destroyed() {
            window.clearInterval(this.wsTimer);
            this.isReconnection = false;
            this.ws.close();
        }
    };
</script>
