<template>
  <Menu mode="horizontal" theme="light" :active-name="activeIndex">
    <menu-item name="account" to="/">
      <Icon type="logo-facebook" />
      账号管理
    </menu-item>
    <menu-item name="task" to="/task">
      <Icon type="md-list-box" />
      任务列表
    </menu-item>
    <menu-item name="proxy" to="/proxy">
      <Icon type="ios-cog-outline" />
      设置
    </menu-item>
  </Menu>
</template>

<script>
export default {
  computed: {
    activeIndex () {
      let name = this.$route.path.split('/')
      if (name.length >= 2 && name[1] !== '') {
        return name[1]
      }

      return 'account'
    }
  }
}
</script>
