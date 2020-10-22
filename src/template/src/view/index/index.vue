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
                                     :key="index" :name="index.toString()"></el-tab-pane>
                    </el-tabs>
                    <div id="panel-setting">
                        <i class="el-icon-setting"></i>
                    </div>
                </el-header>
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
            <a class="menu-button" @click="createComputer" v-if="isDir">添加</a>
            <a class="menu-button">编辑</a>
            <a class="menu-button">删除</a>
            <a class="menu-button" @click="openShell" v-if="!isDir">打开终端</a>
            <a class="menu-button" v-if="!isDir">打开桌面</a>
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
                    <el-input v-model="form.computer.name" placeholder="请输入主机名称"></el-input>
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
    import {add, list} from "../../api/machine";

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
                        name: '',
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
                machineData: [],
                isDir: false,
                dirData: {},
                defaultTopTagMenu: '0',
                defaultProps: {
                    children: 'children',
                    label: 'name'
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
            }
        },
        mounted() {
            this.getMachineData()
        },
        methods: {
            filterNode(value, data) {
                if (!value) return true;
                return data.name.indexOf(value) !== -1;
            },
            treeRightMenu(MouseEvent, object, node, val) {
                localStorage.setItem('currentSelectTree', JSON.stringify(object));
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
                let jumpRoute = '/shell';
                if (this.$route.path !== jumpRoute) {
                    this.$router.push({
                        path: '/shell',
                    });
                }

                // 设置顶部菜单保存vuex中
                let menuData = JSON.parse(window.localStorage.getItem('currentSelectTree'));
                menuData['menu_type'] = this.$store.state.TopMenu.MENU_SHELL_TYPE;
                this.$store.commit("TopMenu/openTagMenuPush", menuData);
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
            getMachineData() {
                list().then(ret => {
                    if (ret.code === 200) {
                        this.machineData = ret.data;
                    }
                }).catch(err => {
                    this.$message.error('服务器出小差！');
                })
            },
            save(flag) {
                let saveData = {};
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

                add(saveData).then(res => {
                    if (res.code === 200) {
                        this.$message({
                            message: res.message,
                            type: 'success'
                        });

                        this.getMachineData()
                    } else {
                        this.$message.error(res.message);
                    }
                }).catch(err => {
                    this.$message.error('服务器出小差！');
                })
            },
            createComputer() {
                this.dirData = JSON.parse(window.localStorage.getItem('currentSelectTree'));
                this.isAddComputer = true;
            },
            clickTagMenu(tag, event) {
                this.$store.commit("TopMenu/upDefaultTagMenu", tag.name);
                console.log(tag, event)
            },
            removeTagMenu(index) {
                this.$store.commit("TopMenu/removeTagMenu", index);
            }
        },
        computed: {
            topTagMenu() {
                return this.$store.state.TopMenu.openTagMenu;
            },
        }
    }
</script>

<style>

</style>