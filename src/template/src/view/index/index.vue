<template>
    <div id="panel">
        <el-container>
            <el-aside id="panel-left" width="200px">
                <div>
                    <el-input placeholder="搜索主机" v-model="search"></el-input>
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
                            @node-click="treeDoubleClick"
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
                            <el-dropdown-item icon="el-icon-c-scale-to-original" command="showAddCommands">批量执行命令
                            </el-dropdown-item>
                            <el-dropdown-item icon="el-icon-coin">执行日志</el-dropdown-item>
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
                <!--                <a class="menu-button" @click="createComputer" v-if="isDir">添加</a>-->
                <!--                <a class="menu-button" @click="editTree">编辑</a>-->
                <!--                <a class="menu-button" @click="delTree">删除</a>-->
                <a class="menu-button" @click="openShell" v-if="!isDir">打开终端</a>
                <a class="menu-button" @click="showAddCommand(1)" v-if="!isDir">执行命令</a>
                <!--                <a class="menu-button" v-if="!isDir">打开桌面</a>-->
            </el-popover>
        </transition>

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
                            v-model="form.command.plan_exec_time"
                            type="datetime"
                            value-format="yyyy-MM-dd HH:mm:ss"
                            placeholder="选择执行日期时间">
                    </el-date-picker>
                </el-form-item>
            </el-form>
            <span slot="footer" class="dialog-footer">
                <el-button @click="isShowAddCommand = false">取 消</el-button>
                <el-button type="primary" @click="saveCommand">确 定</el-button>
            </span>
        </el-dialog>

        <!--        批量执行命令-->
        <!--        抽屉选择主机-->
        <el-drawer
                direction="rtl"
                size="500px"
                :visible.sync="isShowAddCommandComputerList"
                :show-close="false"
                id="panel-commands-drawer">
            <span slot="title">批量执行命令，选择主机</span>
            <el-divider></el-divider>

            <el-table :data="machineData" stripe height="calc(100vh - 150px)" ref="multipleTable"
                      @selection-change="handleSelectionTableChange">
                <el-table-column type="selection"></el-table-column>
                <el-table-column property="id" label="ID"></el-table-column>
                <el-table-column property="name" label="名称"></el-table-column>
            </el-table>

            <el-button-group style="position:relative; right: 10px;float: right;top: 10px">
                <el-button @click="isShowAddCommands = false">取 消</el-button>
                <el-button type="primary" @click="showAddCommand(2)">提 交</el-button>
            </el-button-group>
        </el-drawer>
    </div>
</template>

<script>
    import '@/static/css/index.css';
    import {list} from '../../api/machine';
    import {addCommand} from '../../api/command';
    import shell from '../shell/index';

    /*eslint no-unused-vars: ["error", { "args": "none" }]*/
    export default {
        name: "Index",
        data() {
            return {
                form: {
                    command: {
                        command: '',
                        flag: '1', // 是否定时执行（1立即执行，2定时执行
                        plan_exec_time: '',
                        is_type: 1, // 执行类型（1单个执行，2批量执行
                        ids: [],
                    }
                },
                menuVisible: false,
                search: '',
                machineData: [],
                isDir: false,
                dirData: {},
                defaultTopTagMenu: '0',
                isShowAddCommand: false,
                isShowAddCommandComputerList: false,
                isShowAddCommandTime: false,
                defaultProps: {
                    children: 'children',
                    label: 'name'
                },
                treeClickCount: 0,
                timer: {},
                multipleTableSelection: []
            }
        },
        created() {
            if (!window.sessionStorage.getItem('panel-token')) {
                this.$router.push('/login');
            }

            // 初始化操作
            this.$store.commit('LocalStorage/init');
            this.$store.commit('TopMenu/init');

            this.defaultTopTagMenu = this.$store.state.TopMenu.defaultTagMenu;
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
            treeRightMenu(MouseEvent, data, node, val) {
                this.$store.commit('TopMenu/currentSelectMenuEdit', data);
                this.menuVisible = false;
                this.menuVisible = true;
                var menu = document.querySelector('.menu');
                menu.style.left = MouseEvent.clientX + 'px';
                document.addEventListener('click', this.clearEventRightMenu);
                menu.style.top = MouseEvent.clientY - 10 + 'px';

                if (data.is_dir) {
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
            getMachineData() {
                list().then(ret => {
                    if (ret.code === 200) {
                        this.machineData = ret.data;
                    }
                }).catch(err => {
                    this.$message.error('服务器出小差！');
                })
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
                sessionStorage.removeItem('panel-token');
                sessionStorage.removeItem('panel-userinfo');

                this.$router.push('/login')
            },
            showAddCommand(isType) {
                this.form.command.is_type = isType;
                this.form.command = {
                    command: '',
                    flag: '1',
                    plan_exec_time: '',
                    is_type: isType,
                    ids: [],
                    passwd: {}
                }

                this.isShowAddCommand = true;
            },
            showAddCommands() {
                this.isShowAddCommandComputerList = true;
            },
            treeDoubleClick(data, node) {
                if (data.is_dir) {
                    return;
                }

                this.treeClickCount++;
                if (this.treeClickCount > 2) {
                    return;
                }

                this.$store.commit('TopMenu/currentSelectMenuEdit', data);

                //计时器,计算300毫秒为单位,可自行修改
                this.timer = window.setTimeout(() => {
                    if (this.treeClickCount === 1) {
                        //把次数归零
                        this.treeClickCount = 0;
                    } else if (this.treeClickCount > 1) {
                        //把次数归零
                        this.treeClickCount = 0;
                        //双击事件
                        this.openShell()
                    }
                }, 300);
            },
            handleSelectionTableChange(val) {
                this.multipleTableSelection = val;
            },
            saveCommand() {
                if ((this.form.command.is_type === 2) && this.multipleTableSelection.length) {
                    for (let i in this.multipleTableSelection) {
                        this.form.command.ids.push(this.multipleTableSelection[i]['id'])
                    }
                } else {
                    this.form.command.ids = [this.dirData.id];
                }

                if (this.form.command.ids.length > 0) {
                    let idsMap = {};
                    for (let i in this.form.command.ids) {
                        let idStr = this.form.command.ids[i];
                        idsMap[idStr] = true;
                    }

                    for (let i in this.computerData) {
                        let computerId = this.computerData[i]['id'];
                        if (idsMap[computerId]) {
                            this.form.command.passwd[computerId] = this.$store.state.LocalStorage.computerPasswdData[this.computerData[i].host + ':' + this.computerData[i].port]
                        }
                    }
                }

                addCommand(this.form.command).then(res => {
                    if (res.code === 200) {
                        this.$message({
                            message: res.message,
                            type: 'success'
                        });
                    } else {
                        this.$message.error(res.message);
                    }

                    this.isShowAddCommand = false;
                }).catch(err => {
                    this.$message.error('服务器出小差！');
                })
            }
        },
        computed: {
            topTagMenu() {
                return this.$store.state.TopMenu.openTagMenu;
            }
        },
        components: {
            shell
        },
        destroyed() {
            clearTimeout(this.timer);
        }
    }
</script>

<style>

</style>