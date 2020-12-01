<template>
  <div class="content-container">
    <div class="account-wrapper">
      <div class="account-item-wrapper bind-action-wrapper" @click="auth">
        绑定账号
      </div>

      <div class="account-item-wrapper account-item" v-for="item in accounts" :key="item.AccountId">
        <div class="account-name-wrapper">{{ item.Name }}</div>
        <div class="sync-wrapper" @click="syncAdAccount(item.Token)">同步广告账号</div>
      </div>
    </div>

    <div class="search-wrapper">
      <Form :model="search" inline>
        <FormItem prop="广告账号">
          <Input type="text" v-model="search.name" placeholder="广告账号">
            <Icon type="logo-facebook" slot="prepend"></Icon>
          </Input>
        </FormItem>
        <FormItem prop="广告账号id">
          <Input type="text" v-model="search.id" placeholder="广告账号id">
            <Icon type="logo-facebook" slot="prepend"></Icon>
          </Input>
        </FormItem>
        <FormItem>
          <Button type="primary" @click="searchCommit">查询</Button>
          <Button @click="searchCancel">清空</Button>
        </FormItem>
      </Form>
    </div>

    <div class="upload-wrapper">
      <Button type="primary" @click="showUpload">受众上传</Button>
    </div>

    <Table border @on-select="changeAccount" @on-select-cancel="changeAccount" @on-select-all="changeAccountBatch" @on-select-all-cancel="changeAccountBatch" size="small" ref="selection" :columns="columns" :data="adAccounts"></Table>
    <div class="pagination-wrapper" style="margin-bottom: 20px;">
      <Page :total="total" :page-size="100" @on-change="changePage" />
    </div>

    <Modal v-model="show" :mask-closable="false" width="80%" @on-visible-change="modalVisible" @on-ok="createTask">
      <Steps :current="currentStep">
        <Step title="确认广告账号" content="移除不需要的广告账号"></Step>
        <Step title="选择文件" content="选择要进行上传的文件"></Step>
      </Steps>
      <div v-if="currentStep === 0">
        <Table class="pagination-wrapper" border size="small" :columns="selectedColumns" :data="Object.values(selectedAccount)">
          <template slot-scope="{ row, index }" slot="action">
            <Button type="error" size="small" @click="remove(row.AdAccountId)">移除</Button>
          </template>
        </Table>
        <div class="step-wrapper">
          <Button size="small" type="primary" class="pagination-wrapper" @click="currentStep = 1">下一步</Button>
        </div>
      </div>
      <div v-if="currentStep === 1">
        <Button class="pagination-wrapper" @click="chooseFiles">选择文件</Button>
        <div class="file-list-wrapper">
          <div class="file-item-wrapper" v-for="(file, i) in fileList" :key="i">
            <div class="file">
              <Icon type="ios-document" />
              <div class="file-title">{{ file }}</div>
            </div>
            <div class="remove-wrapper" @click="removeFile(i)">
              <Icon type="md-close" />
            </div>
          </div>
        </div>
        <div class="step-wrapper">
          <Button size="small" type="primary" class="pagination-wrapper" @click="currentStep = 0">上一步</Button>
        </div>
      </div>
    </Modal>
  </div>
</template>

<script>
const { ipcRenderer } = window.require('electron')

const columns = [
  {type: 'selection', width: 60, align: 'center'},
  {title: '广告账号', key: 'Name'},
  {title: '广告账号id', key: 'AdAccountId'},
  {title: 'FB账号', key: 'AccountName'}
]

const selectedColumns = [
  {title: '广告账号', key: 'Name'},
  {title: '广告账号id', key: 'AdAccountId'},
  {title: 'FB账号', key: 'AccountName'},
  {title: '操作', slot: 'action', width: 150, align: 'center'}
]

export default {
  data () {
    return {
      accounts: [],
      adAccounts: [],
      columns,
      selectedColumns,
      search: {
        name: '',
        id: ''
      },
      selectedAccount: {},
      total: 0,
      show: false,
      fileList: [],
      currentStep: 0
    }
  },
  methods: {
    syncAdAccount (token) {
      this.$http({
        method: 'post',
        url: '/bind',
        params: {token}
      }).then(res => {
        this.$Message.success('成功，请稍后查询是否同步完成')
      })
    },
    createTask () {
      if (Object.keys(this.selectedAccount).length === 0) {
        this.$Message.error('请选择广告账号')
        return
      }

      if (this.fileList.length === 0) {
        this.$Message.error('请选择文件')
        return
      }

      this.$http({
        method: 'post',
        url: '/task',
        data: {
          files: this.fileList,
          adAccounts: Object.keys(this.selectedAccount)
        }
      }).then(res => {
        if (res.status === 200 && res.data.code === 0) {
          this.$Message.success('任务创建成功')
          this.$router.push({path: '/task'})
        }
      })
    },
    removeFile (i) {
      this.fileList.splice(i, 1)
    },
    modalVisible (visible) {
      if (!visible) {
        this.currentStep = 0
      }
    },
    remove (adAccountId) {
      delete this.selectedAccount[adAccountId]
      this.$forceUpdate()
    },
    changeAccountBatch (selector) {
      if (selector.length === 0) {
        this.selectedAccount = {}
      } else {
        selector.forEach(row => {
          this.selectedAccount[row.AdAccountId] = row
        })
      }
    },
    chooseFiles () {
      ipcRenderer.send('chooseFile')
    },
    changePage (page) {
      this.getAdAccounts({...this.search, page})
    },
    searchCommit () {
      this.getAdAccounts(this.search)
    },
    showUpload () {
      this.show = true
    },
    searchCancel () {
      this.search = {
        name: '',
        id: ''
      }
      this.getAdAccounts()
    },
    auth () {
      ipcRenderer.send('auth')
    },
    getAccounts () {
      this.$http({
        method: 'get',
        url: '/accounts'
      }).then(res => {
        if (res.status === 200 && res.data.code === 0) {
          this.accounts = res.data.data
        }
      })
    },
    getAdAccounts (params = {}) {
      this.$http({
        method: 'get',
        url: '/adAccounts',
        params
      }).then(res => {
        if (res.status === 200 && res.data.code === 0) {
          let data = res.data.data
          for (let i = 0; i < data.length; i++) {
            if (this.selectedAccount.hasOwnProperty(data[i]['AdAccountId'])) {
              data[i]['_checked'] = true
            } else {
              data[i]['_checked'] = false
            }
          }
          this.adAccounts = data
          this.total = res.data.size
        }
      })
    },
    changeAccount (selector, row) {
      if (this.selectedAccount.hasOwnProperty(row.AdAccountId)) {
        delete this.selectedAccount[row.AdAccountId]
      } else {
        this.selectedAccount[row.AdAccountId] = row
      }
    }
  },
  mounted () {
    this.getAccounts()
    this.getAdAccounts()

    ipcRenderer.on('chooseFile', (event, args) => {
      if (args) {
        let files = new Set()
        this.fileList.forEach(file => files.add(file))
        args.forEach(file => files.add(file))
        this.fileList = Array.from(files)
      }
    })

    ipcRenderer.on('auth', (event, args) => {
      if (args) {
        this.$http({
          method: 'post',
          url: '/bind',
          params: {token: args}
        }).then(res => {
          this.$Message.success('绑定成功，请稍后查询是否同步完成')
          setTimeout(() => {
            this.getAccounts()
            this.getAdAccounts()
          }, 5000)
        })
      }
    })
  }
}
</script>

<style lang="css" scoped>
.account-wrapper {
  padding-top: 20px;
  display: flex;
}

.account-item-wrapper {
  margin-right: 12px;
  width: 160px;
  height: 54px;
  border: 1px solid #e9e9e9;
  border-radius: 4px;
  box-sizing: content-box;
}

.account-item-wrapper.account-item {
  display: flex;
  flex-direction: column;
  justify-content: center;
}

.account-name-wrapper {
  padding: 0 12px;
  flex: 0 0 27px;
  line-height: 27px;
}

.sync-wrapper {
  text-align: center;
  flex: 0 0 27px;
  line-height: 27px;
  font-size: 12px;
  color: #008cff;
  background-color: #f7f9fa;
  cursor: pointer;
}

.account-item-wrapper.bind-action-wrapper {
  line-height: 54px;
  color: #fff;
  text-align: center;
  background-color: #008cff;
  cursor: pointer;
}

.search-wrapper {
  margin-top: 20px;
}

.upload-wrapper {
  margin-bottom: 16px;
}

.pagination-wrapper {
  margin-top: 16px;
}

.step-wrapper {
  display: flex;
  justify-content: flex-end;
}

.file-list-wrapper {
  margin-top: 16px;
}

.file-item-wrapper {
  display: flex;
  justify-content: space-between;
  align-items: center;
  height: 32px;
  line-height: 32px;
  border-bottom: 1px solid #dddddd;
}

.file {
  display: flex;
  align-items: center;
}

.file-title {
  margin-left: 12px;
}

.remove-wrapper {
  flex: 0 0 80px;
  cursor: pointer;
  color: darkred;
  text-align: center;
}
</style>
