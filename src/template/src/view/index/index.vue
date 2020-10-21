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
                            :data="data"
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
                <el-header id="panel-header" height="40px">Header</el-header>
                <el-main id="panel-main">
                    <router-view></router-view>
                </el-main>
            </el-container>
        </el-container>

        <el-popover
                popper-class="menu"
                placement="right"
                width="150"
                trigger="manual"
                v-model="menuVisible">
            <a class="menu-button">添加</a>
            <a class="menu-button">编辑</a>
            <a class="menu-button">删除</a>
            <a class="menu-button" @click="openShell">打开终端</a>
            <a class="menu-button">打开桌面</a>
        </el-popover>

        <el-dialog
                title="添加目录"
                :visible.sync="isAddDir"
                width="500px"
                center>
            <el-form :model="form.dir" label-width="80px">
                <el-form-item label="目录名称">
                    <el-input v-model="form.dir.name" placeholder="请输入目录名"></el-input>
                </el-form-item>
            </el-form>

            <span slot="footer" class="dialog-footer">
                <el-button @click="isAddDir = false">取 消</el-button>
                <el-button type="primary" @click="save(1)">确 定</el-button>
            </span>
        </el-dialog>

        <el-dialog
                title="添加主机"
                :visible.sync="isAddComputer"
                width="500px"
                center>
            <el-form :model="form.computer" label-width="80px">
                <el-form-item label="名称">
                    <el-input v-model="form.computer.alias" placeholder="请输入主机名称"></el-input>
                </el-form-item>
                <el-form-item label="地址">
                    <el-input v-model="form.computer.host" placeholder="请输入主机地址"></el-input>
                </el-form-item>
                <el-form-item label="用户名">
                    <el-input v-model="form.computer.user" placeholder="请输入主机用户名，默认root"></el-input>
                </el-form-item>
                <el-form-item label="端口">
                    <el-input v-model="form.computer.port" placeholder="请输入主机端口，默认22"></el-input>
                </el-form-item>
            </el-form>

            <span slot="footer" class="dialog-footer">
                <el-button @click="isAddComputer = false">取 消</el-button>
                <el-button type="primary" @click="save(2)">确 定</el-button>
            </span>
        </el-dialog>
    </div>
</template>

<script>
    import '@/static/css/index.css';
    import {add} from "../../api/machine";

    const ADD_MACHINE_DIR = 1;
    const ADD_MACHINE_COMPUTER = 2;

    /*eslint no-unused-vars: ["error", { "args": "none" }]*/
    export default {
        name: "Index",
        data() {
            return {
                form: {
                    dir: {
                        name: '',
                        flag: ADD_MACHINE_DIR
                    },
                    computer: {
                        alias: '',
                        host: '',
                        user: '',
                        port: '',
                        machine_group_id: 0,
                        flag: ADD_MACHINE_COMPUTER
                    }
                },
                isAddComputer: false,
                isAddMenu: false,
                isAddDir: false,
                menuVisible: false,
                search: '',
                data: [
                    {
                        id: 2,
                        label: '一级 2',
                        children: [{
                            id: 5,
                            label: '二级 2-1',
                            host: '127.0.0.1'
                        }, {
                            id: 6,
                            label: '二级 2-2'
                        }]
                    }, {
                        id: 3,
                        label: '一级 3',
                        children: [{
                            id: 7,
                            label: '二级 3-1'
                        }, {
                            id: 8,
                            label: '二级 3-2'
                        }]
                    }],
                defaultProps: {
                    children: 'children',
                    label: 'label'
                }
            }
        },
        created() {
            if (!localStorage.getItem('panel-token')) {
                this.$router.push('/login')
            }
        },
        watch: {
            search(val) {
                this.$refs.tree.filter(val);
            }
        },
        methods: {
            filterNode(value, data) {
                if (!value) return true;
                return data.label.indexOf(value) !== -1;
            },
            treeRightMenu(MouseEvent, object, node, val) {
                localStorage.setItem('currentSelectTree', JSON.stringify(object));
                this.menuVisible = false;
                this.menuVisible = true;
                var menu = document.querySelector('.menu');
                menu.style.left = MouseEvent.clientX + 'px';
                document.addEventListener('click', this.clearEventRightMenu);
                menu.style.top = MouseEvent.clientY - 10 + 'px';
            },
            clearEventRightMenu() { // 取消鼠标监听事件 菜单栏
                this.menuVisible = false;
                document.removeEventListener('click', this.foo);
            },
            openShell() {
                this.$router.push({
                    path: '/shell/' + this.$md5((new Date()).getTime().toString()),
                });
            },
            showAddMenu() {
                this.isAddMenu = !this.isAddMenu;
            },
            showAddDir() {
                this.isAddDir = true;
                this.isAddMenu = false;
            },
            showAddComputer() {
                this.isAddComputer = true;
                this.isAddMenu = false;
            },
            save(flag) {
                let saveData = {};
                if (flag === ADD_MACHINE_DIR) {
                    saveData = this.form.dir;
                    this.isAddDir = false;
                } else if (flag === ADD_MACHINE_COMPUTER) {
                    saveData = this.form.computer;
                    saveData.machine_group_id = saveData.machine_group_id ? saveData.machine_group_id : 0;
                    saveData.port = saveData.port ? parseInt(saveData.port) : 22;
                    saveData.user = saveData.user ? saveData.user : 'root';
                    this.isAddComputer = false;
                }

                add(saveData).then(res => {
                    if (res.code === 200) {
                        this.$message({
                            message: res.message,
                            type: 'success'
                        });
                    } else {
                        this.$message.error(res.message);
                    }
                }).catch(err => {
                    this.$message.error('服务器出小差！');
                })
            }
        },
    }
</script>

<style>
    #panel-add-menu {
        position: absolute;
        z-index: 888;
    }

    .box-card {
        width: 200px;
        height: 40px;
        margin: 0 0 0 2px !important;
        padding: 0 !important;
    }
</style>