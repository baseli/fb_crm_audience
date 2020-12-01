<template>
  <Form class="content-container proxy-wrapper" :model="form" :label-width="80">
    <FormItem label="代理地址">
      <Input v-model="form.input" placeholder="请输入代理地址，类似 http://127.0.0.1:7890"></Input>
    </FormItem>
    <FormItem>
      <Button type="primary" @click="save">保存</Button>
    </FormItem>
  </Form>
</template>

<script>
export default {
  data () {
    return {
      form: {
        input: ''
      }
    }
  },
  methods: {
    save () {
      this.$http({
        url: '/proxy',
        method: 'post',
        params: {proxy: this.form.input}
      }).then(res => {
        if (res.status === 200 && res.data.code === 0) {
          this.$Message.success('设置成功，请重启软件')
        }
      })
    }
  },
  mounted () {
    this.$http({
      url: '/proxy',
      method: 'get'
    }).then(res => {
      if (res.status === 200 && res.data.code === 0) {
        this.form.input = res.data.data
      }
    })
  }
}
</script>

<style>
.proxy-wrapper {
  padding: 20px 60vw 0 0;
}
</style>
