<template>
  <div class="content-container">
    <div class="search-wrapper">
      <Form :model="search" inline>
        <FormItem prop="广告账号id">
          <Input type="text" v-model="search.id" placeholder="广告账号id">
            <Icon type="logo-facebook" slot="prepend"></Icon>
          </Input>
        </FormItem>
        <FormItem prop="状态">
          <Select v-model="search.status">
            <Option :value="0">上传中</Option>
            <Option :value="1">成功</Option>
            <Option :value="2">失败</Option>
          </Select>
        </FormItem>
        <FormItem>
          <Button type="primary" @click="searchCommit">查询</Button>
          <Button @click="searchCancel">清空</Button>
        </FormItem>
      </Form>
    </div>

    <div>
      <Button type="primary" @click="retry">重试所有</Button>
      <Button class="margin-left-16" type="primary" @click="remove">清空记录</Button>
    </div>

    <Table class="margin-top-16" border size="small" ref="selection" :columns="column" :data="taskList">
      <template slot-scope="{ row, index }" slot="status">
        <div>
          {{ row.Status == 0 ? '待处理' : row.Status == 1 ? '成功' : row.Status == 2 ? '失败' + row.Msg : '' }}
        </div>
      </template>
    </Table>
    <div class="margin-top-16" style="margin-bottom: 20px;">
      <Page :total="total" :page-size="100" @on-change="changePage" />
    </div>
  </div>
</template>

<script>
const column = [
  {title: '广告账号', key: 'AdAccountName'},
  {title: '广告账号名', key: 'AdAccountId'},
  {title: '文件', key: 'File'},
  {title: '状态', slot: 'status'},
  {title: '时间', key: 'CreatedAt'}
]

export default {
  data () {
    return {
      column,
      search: {
        id: '',
        status: ''
      },
      total: 0,
      taskList: []
    }
  },
  methods: {
    retry () {
      this.$Modal.confirm({
        content: '确认重试所有错误任务？',
        onOk: () => {
          this.$http({
            method: 'post',
            url: '/task/retry'
          }).then(res => {
            if (res.status === 200 && res.data.code === 0) {
              this.$Message.success('重试成功')
              this.getTaskList()
            }
          })
        }
      })
    },
    remove () {
      this.$Modal.confirm({
        content: '确认清空上传列表？',
        onOk: () => {
          this.$http({
            method: 'delete',
            url: '/task'
          }).then(res => {
            if (res.status === 200 && res.data.code === 0) {
              this.$Message.success('删除成功')
              this.getTaskList()
            }
          })
        }
      })
    },
    changePage (page) {
      this.getTaskList({...this.search, page})
    },
    searchCommit () {
      this.getTaskList(this.search)
    },
    searchCancel () {
      this.search = {
        id: '',
        status: ''
      }
      this.getTaskList()
    },
    getTaskList (params = {}) {
      this.$http({
        method: 'get',
        url: '/task',
        params
      }).then(res => {
        if (res.status === 200 && res.data.code === 0) {
          this.taskList = res.data.data
          this.total = res.data.size
        }
      })
    }
  },
  mounted () {
    this.getTaskList()
  }
}
</script>

<style scoped>
.search-wrapper {
  padding-top: 20px;
}

.margin-top-16 {
  margin-top: 16px;
}

.margin-left-16 {
  margin-left: 16px;
}
</style>
