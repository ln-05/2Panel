<template>
  <router-view v-slot="{ Component }">
    <transition mode="out-in" name="el-fade-in-linear">
      <keep-alive :include="routerStore.keepAliveRouters">
        <component :is="Component" />
      </keep-alive>
    </transition>
  </router-view>
</template>

<script setup>
import { useRouterStore } from '@/pinia/modules/router'
const routerStore = useRouterStore()
defineOptions({
  name: 'ContainerManagement'
})
</script>

