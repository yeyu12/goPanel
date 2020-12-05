<template>
    <div>
        //虚拟化桌面将已canvas的形式被渲染到这里
        <div id="screen"></div>
    </div>

</template>

<script>
    import RFB from '@novnc/novnc/core/rfb';
    /*eslint no-unused-vars: ["error", { "args": "none" }]*/
    export default {
        name: "index",
        data() {
            return {
                rfb: null,
                url: undefined, //链接的url
                IsClean: false, //是否已断开并不可重新连接
                connectNum: 0 //重连次数
            }
        },
        methods: {
            // vnc连接断开的回调函数
            disconnectedFromServer(msg) {
                if (msg.detail.clean) {
                    // 根据 断开信息的msg.detail.clean 来判断是否可以重新连接
                    this.contentVnc()
                } else {
                    //这里做不可重新连接的一些操作
                }
            },
            // 连接成功的回调函数
            connectedToServer() {
                console.log('success')
            },

            //连接vnc的函数

            connectVnc() {
                const PASSWORD = '';
                let rfb = new RFB(document.getElementById('screen'), this.url, {
                    // 向vnc 传递的一些参数，比如说虚拟机的开机密码等
                    credentials: {password: PASSWORD}
                });
                rfb.addEventListener('connect', this.connectedToServer);
                rfb.addEventListener('disconnect', this.disconnectedFromServer);
                rfb.scaleViewport = true;  //scaleViewport指示是否应在本地扩展远程会话以使其适合其容器。禁用时，如果远程会话小于其容器，则它将居中，或者根据clipViewport它是否更大来处理。默认情况下禁用。
                rfb.resizeSession = true; //是一个boolean指示是否每当容器改变尺寸应被发送到调整远程会话的请求。默认情况下禁用
                this.rfb = rfb;

            }

        },
        // 在mounted周期里面连接vnc
        mounted() {
            //这时已经可以获取到dom元素
            this.connectVnc()
        }
    }
</script>

<style scoped>

</style>