<template>
    <div id="panel">
        <el-container>
            <el-aside id="panel-left" width="200px">
                <div>
                    <el-input placeholder="过滤" v-model="search">
                        <el-button slot="prepend" icon="el-icon-plus" @click="showAddMenu"></el-button>
                    </el-input>

                    <transition name="el-fade-in-linear">
                        <div id="panel-add-menu" v-if="isAddMenu">
                            <el-button-group class="box-card">
                                <el-button @click="showAddDir">添加目录</el-button>
                                <el-button @click="showAddComputer">添加主机</el-button>
                            </el-button-group>
                        </div>
                    </transition>

                    <el-tree
                            class="filter-tree"
                            :data="machineData"
                            :props="defaultProps"
                            :filter-node-method="filterNode"
                            ref="tree"
                            expand-on-click-node
                            draggable
                            @node-contextmenu="treeRightMenu"
                            default-expand-all
                    >
                    </el-tree>
                </div>
            </el-aside>

            <el-container>
                <el-header id="panel-header" height="40px">
                    <el-tabs type="border-card" closable id="panel-header-menu" v-model="defaultTopTagMenu"
                             @tab-click="clickTagMenu"
                             @tab-remove="removeTagMenu"
                    >
                        <el-tab-pane v-for="(item, index) in topTagMenu" :label="item.name" :data-menu="item"
                                     :key="index" :name="index.toString()" lazy>
                            <component :is="item.menu_type" :menu="item" :tag-index="index"></component>
                        </el-tab-pane>
                    </el-tabs>
                    <el-dropdown id="panel-setting" @command="handleCommand">
                        <i class="el-icon-setting"></i>

                        <el-dropdown-menu slot="dropdown">
                            <!--                            <el-dropdown-item icon="el-icon-plus">黄金糕</el-dropdown-item>-->
                            <!--                            <el-dropdown-item icon="el-icon-circle-plus">狮子头</el-dropdown-item>-->
                            <el-dropdown-item icon="el-icon-check">批量执行命令</el-dropdown-item>
                            <el-dropdown-item icon="el-icon-circle-plus-outline">执行日志</el-dropdown-item>
                            <el-dropdown-item icon="el-icon-switch-button" command="loginout" divided>退出登录
                            </el-dropdown-item>
                        </el-dropdown-menu>
                    </el-dropdown>
                </el-header>
            </el-container>
        </el-container>

        <!-- 右键菜单-->
        <transition name="el-fade-in-linear">
            <el-popover
                    popper-class="menu"
                    placement="right"
                    width="150"
                    trigger="manual"
                    v-model="menuVisible">
                <a class="menu-button" @click="createComputer" v-if="isDir">添加</a>
                <a class="menu-button" @click="editTree">编辑</a>
                <a class="menu-button" @click="delTree">删除</a>
                <a class="menu-button" @click="openShell" v-if="!isDir">打开终端</a>
                <a class="menu-button" @click="showAddCommand" v-if="!isDir">执行命令</a>
                <!--                <a class="menu-button" v-if="!isDir">打开桌面</a>-->
            </el-popover>
        </transition>

        <!-- 添加目录-->
        <el-dialog
                title="添加/修改目录"
                :visible.sync="isAddDir"
                width="500px"
                center>
            <el-form :model="form.dir" label-width="100px" ref="dir" :rules="addDirVail">
                <el-form-item label="目录名称：" prop="name">
                    <el-input v-model="form.dir.name" placeholder="请输入目录名"></el-input>
                </el-form-item>
            </el-form>

            <span slot="footer" class="dialog-footer">
                <el-button @click="isAddDir = false">取 消</el-button>
                <el-button type="primary" @click="save(1, 'dir')">确 定</el-button>
            </span>
        </el-dialog>

        <!-- 添加主机-->
        <el-dialog
                title="添加/修改主机"
                :visible.sync="isAddComputer"
                width="500px"
                center>
            <el-form :model="form.computer" label-width="80px" ref="computer" :rules="addComputerVail">
                <el-form-item label="名称：" prop="name">
                    <el-input v-model="form.computer.name" placeholder="请输入主机名称"></el-input>
                </el-form-item>
                <el-form-item label="地址：" prop="host">
                    <el-input v-model="form.computer.host" placeholder="请输入主机地址"></el-input>
                </el-form-item>
                <el-form-item label="用户名：">
                    <el-input v-model="form.computer.user" placeholder="请输入主机用户名，默认root"></el-input>
                </el-form-item>
                <el-form-item label="密码：" prop="passwd">
                    <el-input v-model="form.computer.passwd" type="password" placeholder="请输入主机密码"></el-input>
                </el-form-item>
                <el-form-item label="端口：">
                    <el-input v-model="form.computer.port" placeholder="请输入主机端口，默认22"></el-input>
                </el-form-item>
            </el-form>

            <span slot="footer" class="dialog-footer">
                <el-button @click="isAddComputer = false">取 消</el-button>
                <el-button type="primary" @click="save(2, 'computer')">确 定</el-button>
            </span>
        </el-dialog>

        <!-- 执行命令-->
        <el-dialog title="执行命令" :visible.sync="isShowAddCommand" width="800px" center>
            <el-form :model="form.command" :inline="true" label-width="auto" ref="computer" :show-message="false">
                <el-form-item label="命令：" prop="command" required>
                    <el-input type="textarea" :rows="20" style="width: 600px" v-model="form.command.command"
                              placeholder="请输入要执行的命令" clearable></el-input>
                </el-form-item>
                <el-form-item label="执行方式：" required>
                    <el-radio-group v-model="form.command.flag">
                        <el-radio label="1">立即执行</el-radio>
                        <el-radio label="2">定时执行</el-radio>
                    </el-radio-group>
                </el-form-item>
                <el-form-item label="执行时间：" v-if="isShowAddCommandTime" required>
                    <el-date-picker
                            v-model="form.command.execTime"
                            type="datetime"
                            placeholder="选择执行日期时间">
                    </el-date-picker>
                </el-form-item>
            </el-form>
            <span slot="footer" class="dialog-footer">
                <el-button @click="isShowAddCommand = false">取 消</el-button>
                <el-button type="primary" @click="isShowAddCommand = false">确 定</el-button>
            </span>
        </el-dialog>
    </div>
</template>

<script>
    import '@/static/css/index.css';
    import {del, list, save} from '../../api/machine';
    import shell from '../shell/index';

    const ADD_MACHINE_DIR = 1;
    const ADD_MACHINE_COMPUTER = 2;

    /*eslint no-unused-vars: ["error", { "args": "none" }]*/
    export default {
        name: "Index",
        data() {
            var validip = (rule, value, callback) => { // eslint-disable-line no-unused-vars
                const reg = /^(\d{1,2}|1\d\d|2[0-4]\d|25[0-5])\.(\d{1,2}|1\d\d|2[0-4]\d|25[0-5])\.(\d{1,2}|1\d\d|2[0-4]\d|25[0-5])\.(\d{1,2}|1\d\d|2[0-4]\d|25[0-5])$/;
                if (reg.test(value)) {
                    callback();
                } else {
                    return callback(new Error('输入格式不合法！'));
                }
            };

            return {
                form: {
                    dir: {
                        id: 0,
                        name: '',
                        flag: ADD_MACHINE_DIR
                    },
                    computer: {
                        id: 0,
                        name: '',
                        host: '',
                        user: '',
                        port: '',
                        machine_group_id: 0,
                        flag: ADD_MACHINE_COMPUTER
                    },
                    command: {
                        command: '',
                        flag: '1',
                        execTime: '',
                    }
                },
                isAddComputer: false,
                isAddMenu: false,
                isAddDir: false,
                menuVisible: false,
                search: '',
                machineData: [],
                isDir: false,
                dirData: {},
                defaultTopTagMenu: '0',
                isShowAddCommand: false,
                isShowAddCommandTime: false,
                defaultProps: {
                    children: 'children',
                    label: 'name'
                },
                addDirVail: {
                    name: [{required: true, message: '名称不能为空'}]
                },
                addComputerVail: {
                    host: [
                        {required: true, message: '地址不能为空', trigger: 'blur'},
                        {validator: validip, trigger: 'blur'}
                    ],
                    passwd: [{required: true, message: '密码不能为空'}]
                }
            }
        },
        created() {
            if (!localStorage.getItem('panel-token')) {
                this.$router.push('/login')
            }

            let menuData = JSON.parse(window.localStorage.getItem('panel-tag-menu'));
            menuData && this.$store.commit("TopMenu/openTagMenu", menuData);

            let defaultMenuIndex = window.localStorage.getItem('panel-default-tag-menu');

            if (defaultMenuIndex !== undefined) {
                this.$store.commit("TopMenu/upDefaultTagMenu", defaultMenuIndex);
                this.defaultTopTagMenu = defaultMenuIndex;
            }
        },
        watch: {
            search(val) {
                this.$refs.tree.filter(val);
            },
            "$store.state.TopMenu.defaultTagMenu"(val) {
                this.defaultTopTagMenu = this.$store.state.TopMenu.defaultTagMenu;
            },
            "$store.state.TopMenu.currentSelectMenu"(val) {
                this.dirData = val;
            },
            "form.command.flag"(val) {
                if (val === '2') {
                    this.isShowAddCommandTime = true
                } else {
                    this.isShowAddCommandTime = false;
                }
            }
        },
        mounted() {
            this.getMachineData();
        },
        methods: {
            filterNode(value, data) {
                if (!value) return true;
                return data.name.indexOf(value) !== -1;
            },
            treeRightMenu(MouseEvent, object, node, val) {
                // window.localStorage.setItem('currentSelectTree', JSON.stringify(object));
                this.$store.commit('TopMenu/currentSelectMenuEdit', object);
                this.menuVisible = false;
                this.menuVisible = true;
                var menu = document.querySelector('.menu');
                menu.style.left = MouseEvent.clientX + 'px';
                document.addEventListener('click', this.clearEventRightMenu);
                menu.style.top = MouseEvent.clientY - 10 + 'px';

                if (object.is_dir) {
                    this.isDir = true;
                } else {
                    this.isDir = false;
                }
            },
            clearEventRightMenu() { // 取消鼠标监听事件 菜单栏
                this.menuVisible = false;
            },
            openShell() {
                // 设置顶部菜单保存vuex中
                let menuData = this.dirData;
                menuData['menu_type'] = this.$store.state.TopMenu.MENU_SHELL_TYPE;
                this.$store.commit('TopMenu/openTagMenuPush', menuData);
            },
            showAddMenu() {
                this.isAddMenu = !this.isAddMenu;
            },
            showAddDir() {
                this.form.dir = {
                    id: 0,
                    name: '',
                    flag: ADD_MACHINE_DIR
                };
                this.isAddDir = true;
                this.isAddMenu = false;
            },
            showAddComputer() {
                this.form.computer = {
                    id: 0,
                    name: '',
                    host: '',
                    user: '',
                    port: '',
                    machine_group_id: 0,
                    flag: ADD_MACHINE_COMPUTER
                };
                this.isAddComputer = true;
                this.isAddMenu = false;
            },
            getMachineData() {
                list().then(ret => {
                    if (ret.code === 200) {
                        this.machineData = ret.data;
                    }
                }).catch(err => {
                    this.$message.error('服务器出小差！');
                })
            },
            save(flag, formName) {
                let saveData = {};

                this.$refs[formName].validate((valid) => {
                    if (valid) {
                        if (flag === ADD_MACHINE_DIR) {
                            saveData = this.form.dir;
                            this.isAddDir = false;
                        } else if (flag === ADD_MACHINE_COMPUTER) {
                            saveData = this.form.computer;
                            saveData.machine_group_id = this.dirData.is_dir ? this.dirData.is_dir : 0;
                            saveData.port = saveData.port ? parseInt(saveData.port) : 22;
                            saveData.user = saveData.user ? saveData.user : 'root';
                            this.isAddComputer = false;
                        }

                        this.dirData = {};

                        save(saveData).then(res => {
                            if (res.code === 200) {
                                this.$message({
                                    message: res.message,
                                    type: 'success'
                                });

                                if (res.data.is_dir === undefined) {
                                    let computerCache = {};
                                    let computer = window.localStorage.getItem('panel-computer');
                                    if (computer) {
                                        computerCache = JSON.parse(computer);
                                    }

                                    computerCache[res.data.host + ':' + res.data.port.toString()] = res.data.passwd;
                                    window.localStorage.setItem('panel-computer', JSON.stringify(computerCache))
                                }

                                this.getMachineData()
                            } else {
                                this.$message.error(res.message);
                            }
                        }).catch(err => {
                            this.$message.error('服务器出小差！');
                        })
                    } else {
                        return false;
                    }
                });
            },
            createComputer() {
                this.isAddComputer = true;
            },
            editTree() {
                if (this.dirData.is_dir) {
                    this.form.dir = {
                        id: this.dirData.id,
                        name: this.dirData.name,
                        flag: ADD_MACHINE_DIR
                    };

                    this.isAddDir = true;
                    this.isAddMenu = false;
                } else {
                    this.form.computer = {
                        id: this.dirData.id,
                        name: this.dirData.name,
                        host: this.dirData.host,
                        user: this.dirData.user,
                        port: this.dirData.port,
                        machine_group_id: this.dirData.machine_group_id,
                        flag: ADD_MACHINE_COMPUTER
                    };

                    this.isAddComputer = true;
                    this.isAddMenu = false;
                }
            },
            delTree() {
                this.$confirm('此操作将永久删除, 是否继续?', '提示', {
                    confirmButtonText: '确定',
                    cancelButtonText: '取消',
                    type: 'warning'
                }).then(() => {
                    let req = {};
                    if (this.dirData.is_dir) {
                        req = {
                            id: this.dirData.id,
                            flag: ADD_MACHINE_DIR
                        }
                    } else {
                        req = {
                            id: this.dirData.id,
                            flag: ADD_MACHINE_COMPUTER
                        }
                    }

                    del(req).then(res => {
                        if (res.code === 200) {
                            this.$message({
                                message: res.message,
                                type: 'success'
                            });

                            if (!this.dirData.is_dir) {
                                // 删除已打开的tag标签
                                this.$store.commit("TopMenu/openTagMenuDel", this.dirData);
                                this.$store.commit("LocalStorage/delComputer", this.dirData);
                            }
                            this.getMachineData()
                        } else {
                            this.$message.error(res.message);
                        }
                    }).catch(err => {
                        console.log(err);
                        this.$message.error('服务器出小差！');
                    })
                }).catch(() => {
                });
            },
            clickTagMenu(tag, event) {
                this.$store.commit("TopMenu/upDefaultTagMenu", tag.name);
            },
            removeTagMenu(index) {
                this.$store.commit("TopMenu/removeTagMenu", index);
            },
            handleCommand(command) {
                eval("this." + command + "()")
            },
            loginout() {
                localStorage.removeItem('panel-token');
                localStorage.removeItem('panel-userinfo');

                this.$router.push('/login')
            },
            showAddCommand() {
                this.isShowAddCommand = true;
            },
            addCommand() {
                console.log(this.dirData);
                // this.isAddComputer = true;
            }
        },
        computed: {
            topTagMenu() {
                return this.$store.state.TopMenu.openTagMenu;
            }
        },
        components: {
            shell
        }
    }
</script>

<style>

</style>