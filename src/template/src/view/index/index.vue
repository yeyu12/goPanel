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
                            icon-class="el-icon-folder-opened"
                            expand-on-click-node
                            draggable
                            @node-contextmenu="treeMenu"
                    >
                    </el-tree>
                </div>
            </el-aside>

            <el-container>
                <el-header id="panel-header" height="40px">Header</el-header>
                <el-main id="panel-main">Main</el-main>
            </el-container>
        </el-container>

        <!--<el-popover
                placement="right"
                width="200"
                trigger="manual"
                content="这里面是右键菜单"
                v-model="visible"
        >
        </el-popover>-->

        <!-- 挂一个遮罩-->
    </div>
</template>

<script>
    export default {
        name: "Index",
        data() {
            return {
                visible: true,
                search: '',
                data: [
                    {
                        id: 1,
                        label: '一级 1',
                        children: [{
                            id: 4,
                            label: '二级 1-1',
                            children: [{
                                id: 9,
                                label: '三级 1-1-1'
                            }, {
                                id: 10,
                                label: '三级 1-1-2'
                            }]
                        }]
                    }, {
                        id: 2,
                        label: '一级 2',
                        children: [{
                            id: 5,
                            label: '二级 2-1'
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
            treeMenu(event, nodeObj, node, data) {
                console.log(nodeObj, node, data)
            }
        },
    }
</script>

<style scoped>
    #panel {
        margin: 0;
        padding: 0;
    }

    #panel-left {
        border-right: 1px #d4d4d4 solid;
        height: 100vh;
        width: 200px;
        float: left;
    }

    #panel-header {
        border: 1px red solid;
        margin: 0;
        padding: 0;
    }

    #panel-main {
        border: 1px black solid;
        height: calc(100vh - 42px);
        width: calc(100vw - 200px);
        margin: 0;
        padding: 0;
    }
</style>