<template>
    <div class="login-box" id="demo">
        <div class="input-content">
            <div class="login_tit">
                <div>
                    <i class="tit-bg left"></i>
                    GoPanel · 登录
                    <i class="tit-bg right"></i>
                </div>
                <p>Strive&nbsp;&nbsp;&nbsp;&nbsp;Everyday</p>
            </div>
            <div v-if="isLogin">
                <p class="p user_icon">
                    <input type="text" placeholder="账号" autocomplete="off" class="login_txtbx" v-model="username">
                </p>
                <p class="p pwd_icon">
                    <input type="password" placeholder="密码" autocomplete="off" class="login_txtbx" v-model="passwd"
                           @keypress="enter">
                </p>
            </div>
            <div v-else>
                <p class="p user_icon">
                    <input type="text" placeholder="账号" autocomplete="off" class="login_txtbx" v-model="username">
                </p>
                <p class="p pwd_icon">
                    <input type="password" placeholder="密码" autocomplete="off" class="login_txtbx" v-model="passwd">
                </p>
                <p class="p rpwd_icon">
                    <input type="password" placeholder="重复密码" autocomplete="off" class="login_txtbx"
                           v-model="repeat_passwd"
                           @keypress="enter">
                </p>
            </div>
            <!--<div class="p val_icon">
                <div class="checkcode">
                    <input type="text" id="J_codetext" placeholder="验证码" autocomplete="off" maxlength="4"
                           class="login_txtbx">
                    <canvas class="J_codeimg" id="myCanvas" onselectstart="return false">对不起，您的浏览器不支持canvas，请下载最新版浏览器!
                    </canvas>
                </div>
                <a class="ver_btn" onselectstart="return false">看不清，换一张</a>
            </div>-->
            <div class="signup">
                <a class="gv" href="javascript:" @click="login" v-if="isLogin">登&nbsp;&nbsp;录</a>
                <a class="gv" href="javascript:" v-else @click="register">注&nbsp;&nbsp;册</a>
            </div>
            <div>
                <a class="panel-register" v-if="!isLogin" @click="showLogin(true)">登&nbsp;录</a>
                <a class="panel-register" v-else @click="showLogin( false)">注&nbsp;册</a>
            </div>
        </div>
        <div class="canvaszz"></div>
        <canvas id="canvas"></canvas>
    </div>
</template>

<script>
    import '@/static/css/login.css';
    import {login, register} from '@/api/user';

    /*eslint no-unused-vars: ["error", { "args": "none" }]*/
    export default {
        name: "Login",
        data() {
            return {
                canvas: "",
                ctx: "",
                w: 0,
                h: 0,
                hue: 0,
                stars: [],
                count: 0,
                maxStars: 0,
                canvas2: "",
                ctx2: "",
                half: 0,
                gradient2: "",
                username: "",
                passwd: "",
                repeat_passwd: "",
                isLogin: true
            }
        },
        mounted() {
            this.canvas = document.getElementById('canvas');
            this.ctx = this.canvas.getContext('2d');
            this.w = this.canvas.width = window.innerWidth;
            this.h = this.canvas.height = window.innerHeight;
            this.hue = 217;
            this.stars = [];
            this.count = 0;
            this.maxStars = 2500; //星星数量

            this.canvas2 = document.createElement('canvas');
            this.ctx2 = this.canvas2.getContext('2d');
            this.canvas2.width = 100;
            this.canvas2.height = 100;
            this.half = this.canvas2.width / 2;
            this.gradient2 = this.ctx2.createRadialGradient(this.half, this.half, 0, this.half, this.half, this.half);
            this.gradient2.addColorStop(0.025, '#CCC');
            this.gradient2.addColorStop(0.1, 'hsl(' + this.hue + ', 61%, 33%)');
            this.gradient2.addColorStop(0.25, 'hsl(' + this.hue + ', 64%, 6%)');
            this.gradient2.addColorStop(1, 'transparent');

            this.ctx2.fillStyle = this.gradient2;
            this.ctx2.beginPath();
            this.ctx2.arc(this.half, this.half, this.half, 0, Math.PI * 2);
            this.ctx2.fill();

            let than = this;
            let Star = function () {
                this.orbitRadius = than.random(than.maxOrbit(than.w, than.h));
                this.radius = than.random(60, this.orbitRadius) / 18;
                //星星大小
                this.orbitX = than.w / 2;
                this.orbitY = than.h / 2;
                this.timePassed = than.random(0, than.maxStars);
                this.speed = than.random(this.orbitRadius) / 500000;
                //星星移动速度
                this.alpha = than.random(2, 10) / 10;

                than.count++;
                than.stars[than.count] = this;
            };

            than = this;
            Star.prototype.draw = function () {
                let x = Math.sin(this.timePassed) * this.orbitRadius + this.orbitX,
                    y = Math.cos(this.timePassed) * this.orbitRadius + this.orbitY,
                    twinkle = than.random(10);

                if (twinkle === 1 && this.alpha > 0) {
                    this.alpha -= 0.05;
                } else if (twinkle === 2 && this.alpha < 1) {
                    this.alpha += 0.05;
                }

                than.ctx.globalAlpha = this.alpha;
                than.ctx.drawImage(than.canvas2, x - this.radius / 2, y - this.radius / 2, this.radius, this.radius);
                this.timePassed += this.speed;
            };

            for (var i = 0; i < this.maxStars; i++) {
                new Star();
            }

            this.animation()
        },
        methods: {
            random(min, max) {
                if (arguments.length < 2) {
                    max = min;
                    min = 0;
                }

                if (min > max) {
                    let hold = max;
                    max = min;
                    min = hold;
                }

                return Math.floor(Math.random() * (max - min + 1)) + min;
            },
            maxOrbit(x, y) {
                let max = Math.max(x, y),
                    diameter = Math.round(Math.sqrt(max * max + max * max));
                return diameter / 2;
            },
            animation() {
                this.ctx.globalCompositeOperation = 'source-over';
                this.ctx.globalAlpha = 0.5; //尾巴
                this.ctx.fillStyle = 'hsla(' + this.hue + ', 64%, 6%, 2)';
                this.ctx.fillRect(0, 0, this.w, this.h);

                this.ctx.globalCompositeOperation = 'lighter';
                for (var i = 1, l = this.stars.length; i < l; i++) {
                    this.stars[i].draw();
                }

                window.requestAnimationFrame(this.animation);
            },
            login() {
                login({username: this.username, passwd: this.passwd}).then(data => {
                    if (data.code !== 200) {
                        this.$message.error(data.message);
                    } else {
                        localStorage.setItem('panel-token', data.data.token);
                        localStorage.setItem('panel-userinfo', this.$base64.encode(JSON.stringify(data.data)));

                        this.$router.push('/')
                    }
                }).catch(err => {
                    this.$message.error("服务器出小差！");
                })
            },
            register() {
                register({
                    username: this.username,
                    passwd: this.passwd,
                    repeat_passwd: this.repeat_passwd
                }).then(res => {
                    if (res.code !== 200) {
                        this.$message.error(res.message);
                    } else {
                        localStorage.setItem('panel-token', res.data.token);
                        localStorage.setItem('panel-userinfo', this.$base64.encode(JSON.stringify(res.data)));

                        this.$router.push('/')
                    }
                }).catch(err => {
                    this.$message.error("服务器出小差！");
                })
            },
            showLogin(val) {
                this.isLogin = val;
            },
            enter(e) {
                if (e.which === 13) {
                    if (this.isLogin) {
                        this.login();
                    } else {
                        this.register();
                    }
                }
            }
        }
    }
</script>

<style scoped>

</style>