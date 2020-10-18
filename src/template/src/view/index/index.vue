<template>
    <div id="panel">
        <el-container>
            <el-aside id="panel-left" width="200px">
                <div>
                    <el-input placeholder="过滤" v-model="search">
                        <el-button slot="prepend" icon="el-icon-plus"></el-button>
                    </el-input>

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
            <a class="menu-button">添加主机</a>
            <a class="menu-button" @click="openShell">打开终端</a>
            <a class="menu-button">打开桌面</a>
        </el-popover>
    </div>
</template>

<script>
    import '@/static/css/index.css';

    export default {
        name: "Index",
        data() {
            return {
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
            /*eslint no-unused-vars: ["error", { "args": "none" }]*/
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
            }
        },
    }
</script>

<style>

</style>