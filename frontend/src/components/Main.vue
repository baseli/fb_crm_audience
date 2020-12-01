<template>
  <div class="container">
    <navigation/>
    <div class="main-container">
      <router-view/>
    </div>
  </div>
</template>

<script>
import Navigation from './Navigation'

export default {
  components: {
    Navigation
  }
}
</script>

<style>
.container {
  background-color: #eee;
}

.main-container {
  margin-top: 20px;
}

.content-container {
  background-color: #fff;
  padding: 0 20px;
}
</style>
